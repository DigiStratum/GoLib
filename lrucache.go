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
func (lru *lruCache) Set(key, content string) {
	lru.lock()
	defer lru.unlock()
	lru.set(key, content)
}

// Retrieve an item from the cache with the supplied key (or nil if there isn't one)
func (lru *lruCache) Get(key string) interface{} {
	lru.lock()
	defer lru.unlock()
	if element := lru.find(key, true); nil != element {
		return element.Value
	}
	return nil
}

// Drop an item from the cache with the supplied key
func (lru *lruCache) Drop(key string) {
	lru.lock()
	defer lru.unlock()
	lru.drop(key)
}


// Private Implementation

type cacheItem struct {
	Key, Content		string
}

type lruCache struct {
	count, countLimit	int
	size, sizeLimit		int
	ageList			*list.List
	elements		map[string]*list.Element
	mutex			*sync.Mutex
}

// Check if this item will even fit within our cache limits
// If should also not displace more than some threshold of current cache elements
// TODO: Merge with prune function below?
func (lru *lruCache) isCacheable(key string, size int) bool {
	return true
}

// Prune the currently cached element collection to fit the new element within limits
func (lru *lruCache) pruneToFit(key string, size int) bool {
	if ! lru.isCacheable(key, size) { return false }

	// Based on limits, how many elements would we need to prune to fit this new thing?
	elementsToPrune := 0

	// TODO: Account for there already being a cacheItem with this key:
	// subract the existing item's size from consideration and deduct
	// one from the total count for limits checks

	// If there is a size limit in effect...
	if (lru.sizeLimit > 0) && ((lru.size + size) >= lru.sizeLimit) {
		pruneSize := (lru.size + size) - lru.sizeLimit
		if pruneSize > elementsToPrune { elementsToPrune = pruneSize }
	}

	// If there is a count limit in effect...
	if (lru.countLimit > 0) && ((lru.count + 1) >= lru.countLimit) {
		pruneCount := (lru.count + 1) - lru.countLimit
		if pruneCount > elementsToPrune { elementsToPrune = pruneCount }
	}

	// What percentage of our cache needs to be pruned?
	if 0 == elementsToPrune { return true } // No pruning needed to fit!

	// TODO: Make sure we're not pruning more than some percentage threshold
	// (to minimize performance hits due to statistical outliers)

	// Prune starting at the back of the age list for the count we need to prune
	element := lru.ageList.Back()
	for i := 0; i < elementsToPrune; i++ {
		keyIfc := element.Value		// get the key for this element
		keyCI := keyIfc.(cacheItem)
		key := keyCI.Key
		element = element.Next()	// get the next element
		lru.drop(key)			// drop this one from the cache
	}
	return true
}

func (lru *lruCache) set(key, content string) {
	if ! lru.pruneToFit(key, len(content)) { return }
	// Drop if exists already
	lru.drop(key)
	// Add content to front of age List and remember it by key in elements map
	(*lru).elements[key] = lru.ageList.PushFront(cacheItem{ Key: key, Content: content })
	(*lru).size += len(content)
	(*lru).count++
}

func (lru *lruCache) drop(key string) {
	// Drop if exists (don't bump on the find since we're going to drop it!)
	if element := lru.find(key, false); nil != element {
		contentIfc := element.Value
		contentCI := contentIfc.(cacheItem)
		(*lru).size -= len(contentCI.Content)
		(*lru).count--
		lru.ageList.Remove(element)
		delete(lru.elements, key)
	}
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

