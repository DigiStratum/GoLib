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


// Public Interface

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

// Check whether we have an item in the cache with the supplied key
func (lru *lruCache) Has(key string) bool {
	lru.lock()
	defer lru.unlock()
	return (nil != lru.find(key, true))
}

// Add a content item to the cache with the supplied key
// return true if we set it, else false
func (lru *lruCache) Set(key, content string) bool {
	lru.lock()
	defer lru.unlock()
	return lru.set(key, content)
}

// Retrieve an item from the cache with the supplied key (or nil if there isn't one)
func (lru *lruCache) Get(key string) *string {
	lru.lock()
	defer lru.unlock()
	if element := lru.find(key, true); nil != element {
		content := element.Value.(lruCacheItem).Content
		return &content
	}
	return nil
}

// Drop an item from the cache with the supplied key
// return true if we drop it, else false
func (lru *lruCache) Drop(key string) bool {
	lru.lock()
	defer lru.unlock()
	return lru.drop(key)
}

// Return the data size currently being held in this cache
func (lru *lruCache) Size() int {
	lru.lock()
	defer lru.unlock()
	return lru.size
}

// Return the count of entries currently being held in this cache
func (lru *lruCache) Count() int {
	lru.lock()
	defer lru.unlock()
	return lru.count
}

// Set the limits for size/count on this cache; 0 means unlimited (default for both)
func (lru *lruCache) SetLimits(sizeLimit, countLimit int) {
	lru.lock()
	defer lru.unlock()
	lru.sizeLimit = sizeLimit
	lru.countLimit = countLimit
}

// Private Implementation

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

// How many existing cache entries must be pruned to fit this new one?
func (lru *lruCache) numToPrune(key string, size int) int {
	var pruneCount, replaceCount, replaceSize int
	if element := lru.find(key, false); nil != element {
		replaceCount = 1
		replaceSize = element.Value.(lruCacheItem).Size
	}

	// If there is a count limit in effect...
	if lru.countLimit > 0 {
		futureCount := lru.count + 1 - replaceCount
		if futureCount > lru.countLimit { pruneCount = futureCount - lru.countLimit }
	}

	// If there is a size limit in effect...
	if lru.sizeLimit > 0 {
		// If we add this to the cache without pruning, future size would be...
		futureSize := lru.size + size - replaceSize
		// If we break the size limit by adding...
		if futureSize > lru.sizeLimit {
			pruneSize := futureSize - lru.sizeLimit
			num := 0
			element := lru.ageList.Back()
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
func (lru *lruCache) pruneToFit(key string, size int) bool {

	// Will it fit at all?
	if (lru.sizeLimit > 0) && (size > lru.sizeLimit) { return false }
	pruneCount := lru.numToPrune(key, size)
	if 0 == pruneCount { return true }

	// TODO: Make sure we're not pruning more than some percentage threshold
	// (to minimize performance hits due to statistical outliers)

	// Prune starting at the back of the age list for the count we need to prune
	element := lru.ageList.Back()
	for ; (nil != element) && (pruneCount > 0); pruneCount-- {
		key := element.Value.(lruCacheItem).Key
		element = element.Next()
		lru.drop(key)
	}
	return true
}

// Add content to front of age List and remember it by key in elements map
// return true if we set it, else false
func (lru *lruCache) set(key, content string) bool {
	if ! lru.pruneToFit(key, len(content)) { return false }
	lru.drop(key)
	(*lru).elements[key] = lru.ageList.PushFront(lruCacheItem{
		Key: key,
		Content: content,
		Size: len(content),
	})
	(*lru).size += len(content)
	(*lru).count++
	return true
}

// Drop if exists (don't bump on the find since we're going to drop it!)
// return bool true is we drop it, else false
func (lru *lruCache) drop(key string) bool {
	if element := lru.find(key, false); nil != element {
		(*lru).size -= len(element.Value.(lruCacheItem).Content)
		(*lru).count--
		lru.ageList.Remove(element)
		delete(lru.elements, key)
		return true
	}
	return false
}

func (lru *lruCache) find(key string, bump bool) *list.Element {
	if element, ok := (*lru).elements[key]; ok {
		// If we got bumped, move to front of age List
		if bump { lru.ageList.MoveToFront(element) }
		return element
	}
	return nil
}

func (lru *lruCache) lock() {
	lru.mutex.Lock()
}

func (lru *lruCache) unlock() {
	lru.mutex.Unlock()
}

