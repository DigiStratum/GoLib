// DigiStratum GoLib - LRU Cache
package golib

/*

A Least-Recently-Used (LRU) Cache for clients to improve inexpensive cache hits before incurring
expensive requests. The "LRU" aspect, in short, is used to remove items from the cache which have
not been accessed recently (least recently) - this only happens once one (or both?) of the limits
have been reached. By default the cache is permitted to grow without limit. If you want to limit
the cache by total storage space, then set a Size Limit. If you want to limit by total number of
items, then set a Count Limit. But unless you trust the Universe to not run you out of memory or
disk space, then you're going to need a limit in place. Once a limit is reached, attempting to add
a new item to the cache, we will find the least useful cache entries and remove them (note that we
might need to remove more than one item to make space if we a relimiting based on size in order for
the new thing to fit.)

As far as implementation, we're going to use a double-linked list and keep the oldest things at the
back of the list and push everything to the front each time it's accessed. This way, when we want
to remove the LRU items from the list, we can just pop them off the back.

ref: https://golang.org/pkg/container/list/

*/

import (
	"sync"
	"container/list"
)

type LruCacheIfc interface {
	Has(key string) bool
	Set(key, content string) bool
	Get(key string) *string
	Drop(key string) bool
	Size() int
	Count() int
	SetLimits(sizeLimit, countLimit int)
}

type lruCacheItem struct {
	Key, Content		string
	Size			int
}

type lruCache struct {
	count, countLimit	int
	size, sizeLimit		int
	ageList			*list.List
	elements		map[string]*list.Element
	mutex			*sync.Mutex
}

// Make a new one of these
func NewLRUCache() *lruCache {
	return &lruCache{
		count:		0,
		countLimit:	0,
		size:		0,
		sizeLimit:	0,
		ageList:	list.New(),
		elements:	make(map[string]*list.Element),
		mutex:		&sync.Mutex{},
	}
}

// -------------------------------------------------------------------------------------------------
// LruCacheIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Check whether we have an item in the cache with the supplied key
func (r lruCache) Has(key string) bool {
	return (nil != r.find(key, true))
}

// Add a content item to the cache with the supplied key
// return true if we set it, else false
func (r *lruCache) Set(key, content string) bool {
	r.lock()
	defer r.unlock()
	return r.set(key, content)
}

// Retrieve an item from the cache with the supplied key (or nil if there isn't one)
func (r lruCache) Get(key string) *string {
	r.lock()
	defer r.unlock()
	if element := r.find(key, true); nil != element {
		content := element.Value.(lruCacheItem).Content
		return &content
	}
	return nil
}

// Drop an item from the cache with the supplied key
// return true if we drop it, else false
func (r *lruCache) Drop(key string) bool {
	r.lock()
	defer r.unlock()
	return r.drop(key)
}

// Return the data size currently being held in this cache
func (r lruCache) Size() int {
	return r.size
}

// Return the count of entries currently being held in this cache
func (r lruCache) Count() int {
	return r.count
}

// Set the limits for size/count on this cache; 0 means unlimited (default for both)
func (r *lruCache) SetLimits(sizeLimit, countLimit int) {
	r.lock()
	defer r.unlock()
	r.sizeLimit = sizeLimit
	r.countLimit = countLimit
}

// -------------------------------------------------------------------------------------------------
// LruCacheIfc Private Interface
// -------------------------------------------------------------------------------------------------

// How many existing cache entries must be pruned to fit this new one?
func (r lruCache) numToPrune(key string, size int) int {
	var pruneCount, replaceCount, replaceSize int
	if element := r.find(key, false); nil != element {
		replaceCount = 1
		replaceSize = element.Value.(lruCacheItem).Size
	}

	// If there is a count limit in effect...
	if r.countLimit > 0 {
		futureCount := r.count + 1 - replaceCount
		if futureCount > r.countLimit { pruneCount = futureCount - r.countLimit }
	}

	// If there is a size limit in effect...
	if r.sizeLimit > 0 {
		// If we add this to the cache without pruning, future size would be...
		futureSize := r.size + size - replaceSize
		// If we break the size limit by adding...
		if futureSize > r.sizeLimit {
			pruneSize := futureSize - r.sizeLimit
			num := 0
			element := r.ageList.Back()
			for ; (nil != element) && (pruneSize > 0); num++ {
				pruneSize -= element.Value.(lruCacheItem).Size
				element = element.Next()
			}
			if num > pruneCount { pruneCount = num }
		}
	}
	return pruneCount
}

// Prune the currently cached element collection to fit the new element within limits
// return boolean true if it will fit, else false (true doesn't indicate whether we did any pruning)
func (r *lruCache) pruneToFit(key string, size int) bool {

	// Will it fit at all?
	if (r.sizeLimit > 0) && (size > r.sizeLimit) { return false }
	pruneCount := r.numToPrune(key, size)
	if 0 == pruneCount { return true }

	// TODO: Make sure we're not pruning more than some percentage threshold
	// (to minimize performance hits due to statistical outliers)

	// Prune starting at the back of the age list for the count we need to prune
	element := r.ageList.Back()
	for ; (nil != element) && (pruneCount > 0); pruneCount-- {
		key := element.Value.(lruCacheItem).Key
		element = element.Next()
		r.drop(key)
	}
	return true
}

// Add content to front of age List and remember it by key in elements map
// return true if we set it, else false
func (r *lruCache) set(key, content string) bool {
	if ! r.pruneToFit(key, len(content)) { return false }
	r.drop(key)
	r.elements[key] = r.ageList.PushFront(lruCacheItem{
		Key: key,
		Content: content,
		Size: len(content),
	})
	r.size += len(content)
	r.count++
	return true
}

// Drop if exists (don't bump on the find since we're going to drop it!)
// return bool true is we drop it, else false
func (r *lruCache) drop(key string) bool {
	if element := r.find(key, false); nil != element {
		r.size -= len(element.Value.(lruCacheItem).Content)
		r.count--
		r.ageList.Remove(element)
		delete(r.elements, key)
		return true
	}
	return false
}

func (r lruCache) find(key string, bump bool) *list.Element {
	if element, ok := r.elements[key]; ok {
		// If we got bumped, move to front of age List
		if bump { r.ageList.MoveToFront(element) }
		return element
	}
	return nil
}

func (r *lruCache) lock() {
	r.mutex.Lock()
}

func (r *lruCache) unlock() {
	r.mutex.Unlock()
}

