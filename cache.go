// DigiStratum GoLib - Cache
package golib

/*

This in-memory cache class wraps a basic Go map with essential helper functions to make life easier
for dealing with simple key/value pair data where the values are objects (interface{}) instead of
strings. An important differentiator of a cache vs. a hashmap is that the cache entries have time-
based expirations to ensure that none of the entries outlive their prescribed freshness.

TODO: Set up go routine thread to regularly purgeExpired; track the next nearest expire time so that
we can land right on it (certainly no sooner than necessary?) How do we stop the go routine when the
underlying struct is dead (like the caller creates/uses cache, then throws it away)? the go routine
should stop if the struct is "dead"
ref: https://stackoverflow.com/questions/6807590/how-to-stop-a-goroutine
ref: https://yourbasic.org/golang/wait-for-goroutines-waitgroup/

*/

import (
	"time"
	"sync"
)

type cacheItem struct {
	Value	interface{}
	Expires	int64
}

func (ci cacheItem) IsExpired() bool {
	return time.Now().Unix() < ci.Expires
}

type Cache struct {
	cache			map[string]cacheItem
	mutex			sync.Mutex
}

type CacheIfc interface {
	IsEmpty() bool
	Size() int
	Set(key string, value interface{}, expires int64)
	Get(key string) interface{}
	Has(key string) bool
	HasAll(keys *[]string) bool
	Flush()
}

// Factory Functions
func NewCache() Cache {
	return Cache{
		cache:	make(map[string]cacheItem),
	}
}

// -------------------------------------------------------------------------------------------------
// Cache Public Interface
// -------------------------------------------------------------------------------------------------

// Flush all the items out of the cache
func (r *Cache) Flush() {
	r.mutex.Lock(); defer r.mutex.Unlock()
	r.cache = make(map[string]cacheItem)
}

// Check whether this Cache is empty (has no properties)
func (r Cache) IsEmpty() bool {
	return 0 == r.Size()
}

// Get the number of properties in this Cache
func (r Cache) Size() int {
	r.purgeExpired()
	return len(r.cache)
}

// Set a single cache element key to the specified value
func (r *Cache) Set(key string, value interface{}, expires int64) {
	r.mutex.Lock(); defer r.mutex.Unlock()
	r.cache[key] = cacheItem{
		Value:		value,
		Expires:	expires,
	}
}

// Get a single cache element by key name
func (r Cache) Get(key string) interface{} {
	if ci, ok := r.cache[key]; ok {
		if ! ci.IsExpired() { return ci.Value }
	}
	return nil
}

// Check whether we have a configuration element by key name
func (r Cache) Has(key string) bool {
	_, ok := r.cache[key];
	return ok
}

// Check whether we have configuration elements for all the key names
func (r Cache) HasAll(keys *[]string) bool {
	for _, key := range *keys {
		if ! r.Has(key) { return false }
	}
	return true
}

// -------------------------------------------------------------------------------------------------
// Cache Private Interface
// -------------------------------------------------------------------------------------------------

// Purge expired cache items
func (r *Cache) purgeExpired() {
	// Find which keys we need to purge because their cacheItem is expired
	purgeKeys := []string{}
	for key, ci := range r.cache {
		if ci.IsExpired() { purgeKeys = append(purgeKeys, key) }
	}
	// Purge them!
	for _, key := range purgeKeys { delete(r.cache, key) }
}
