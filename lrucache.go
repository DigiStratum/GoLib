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

In addition to the explicit imports below, we use the following classes from this same package:
 * Logger

*/

import (
	"container/list"
)

type lruCache struct {
	countLimit	int
	sizeLimit	int
	ageList		List
	elements	map[string]*Element
}

// Make a new one of these
func NewLRUCache() {
	return &lruCache{
		countLimit:	0,
		sizeLimit:	0,
		ageList:	list.New(),
		elements:	make(map[string]*Element),
	}
}

// THREAD SAFE:

func (lru *lruCache) set(key string, content string) {
	// Drop if exists already
	lru.drop(key)
	// Add content to front of age List and remember it by key in elements map
	*lru.elements[key] = lru.ageList.PushFront(content)
}

func (lru *lruCache) drop(key string) {
	// Drop if exists
	if element := lru.find(key); nil != element {
		lru.ageList.Remove(element)
		delete(lru.elements, key)
	}
}

func (lru *lruCache) find(key string) *Element {
	if element, ok := *lru.elements[key]; ok {
		// Just got bumped; move to front of age List
		lru.ageList.MoveToFront(element)
		return element
	}
	return nil
}

// UNSAFE!
func (lru *lruCache) Has(key string) bool {
	element := lru.find(key)
	return (nil != element)
}

func (lru *lruCache) Set(key string, content string) interface{} {
	lru.set(key, content)
}

func (lru *lruCache) Get(key string) interface{} {
	if element := lru.find(key); nil != element {
		return element.Value()
	}
	return nil
}

func (lru *lruCache) Drop(key string) {
	lru.drop(key)
}

