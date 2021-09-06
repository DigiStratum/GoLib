// DigiStratum GoLib - Cache
package cache

/*

This in-memory cache class wraps a basic Go map with essential helper functions to make life easier
for dealing with simple key/value pair data where the values are objects (interface{}) instead of
strings. An important differentiator of a cache vs. a hashmap is that the cache entries have time-
based expirations to ensure that none of the entries outlive their prescribed freshness.

TODO:
 * Set up go routine thread to regularly purgeExpired; track the next nearest expire time so that
we can land right on it (certainly no sooner than necessary?) How do we stop the go routine when the
underlying struct is dead (like the caller creates/uses cache, then throws it away)? the go routine
should stop if the struct is "dead"
ref: https://stackoverflow.com/questions/6807590/how-to-stop-a-goroutine
ref: https://yourbasic.org/golang/wait-for-goroutines-waitgroup/


TODO:
 * Add SetLogger() to set a logger for output; don't just assume default logger in a library. Consumer
   gets to control. purgeExpired() should be logging when it does work, and maybe Trace() log output
   from every operation.
*/

import (
	"fmt"
	"time"
	"sync"
	"container/list"

	"github.com/DigiStratum/GoLib/Chrono"
	"github.com/DigiStratum/GoLib/Data/sizeable"
)

type CacheIfc interface {
	Configure(config cfg.ConfigIfc) error			// cfg.ConfigurableIfc
	SetTimeSource(timeSource chrono.TimeSourceIfc)
	IsEmpty() bool
	Size() int64
	Count() int
	Set(key string, value interface{})
	SetExpires(key string, expires chrono.TimeStampIfc)
	Get(key string) interface{}
	Has(key string) bool
	HasAll(keys *[]string) bool
	Drop(key string) bool
	DropAll(keys *[]string) int
	Flush()
	Close() error
}

type Cache struct {
	cache			map[string]cacheItem
	ageList			*list.List

	totalCountLimit		int
	totalSizeLimit		int
	newItemExpires		int64
	timeSource		chrono.TimeSourceIfc

	totalSize		int64

	mutex			sync.Mutex
	closed			bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewCache() *Cache {
	cache := Cache{
		cache:		make(map[string]cacheItem),
		ageList:	list.New(),
		timeSource:	chrono.NewTimeSource(),
	}
	// Set up a go routine that will run continuously until we get Close()ed
	go cache.purgeExpired()
	return &cache
}

// -------------------------------------------------------------------------------------------------
// cfg.ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) Configure(config cfg.ConfigIfc) error {
	if nil == config { return fmt.Errorf("Cache.Configure() - Configuration was nil") }

	// New items added to cache will expire in this count of seconds; 0 (default) = no expiration
	if config.Has["newItemExpires"] {
		newItemExpires := config.GetInt64("newItemExpires")
		if nil != newItemExpires { r.newItemExpires = *newItemExpires }
	}

	// New items added to cache won't drive total count above this; 0 (default) = unlimited
	// When a limit is in place, the Least Recently Used (LRU) item will be evicted to make room for the new one
	if config.Has["totalCountLimit"] {
		totalCountLimit := config.GetInt64("totalCountLimit")
		if nil != totalCountLimit { r.totalCountLimit = int(*totalCountLimit) }
	}

	// New items  added to cache we won't drive total size of all items above this; 0 = unlimited
	// When a limit is in place, the Least Recently Used (LRU) item(s) will be evicted to make room for the new one
	if config.Has["totalSizeLimit"] {
		totalSizeLimit := config.GetInt64("totalSizeLimit")
		if nil != totalSizeLimit { r.totalSizeLimit = int(*totalSizeLimit) }
	}
}

// -------------------------------------------------------------------------------------------------
// CacheIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) SetTimeSource(timeSource chrono.TimeSourceIfc) {
	if nil != timeSource { r.timeSource = timeSource }
}

// Check whether this Cache is empty (has no properties)
func (r Cache) IsEmpty() bool {
	return 0 == r.Size()
}

// Get the number of properties in this Cache
func (r Cache) Size() int64 {
	return r.totalSize
}
// Return the count of entries currently being held in this cache
func (r Cache) Count() int {
	return len(r.cache)
}

// Set a single cache element key to the specified value
func (r *Cache) Set(key string, value interface{}) {
	r.mutex.Lock(); defer r.mutex.Unlock()

	// Get the size of the value
	size := sizeable.Size(value)

	// If size limit is in play and this value is bigger than that, then it won't fit
	if (0 < r.sizeLimit) && (size > r.sizeLimit) { return false }

	_ = r.drop(key)

	var expires chrono.TimeStampIfc
	if 0 == r.newItemExpires {
		expires = chrono.NewTimeStampForever()
	} else {
		expires = r.TimeSource.NewTimeStamp().Add(r.newItemExpires)
	}
	ci := NewCacheItem(value, expires)
	r.cache[key] = *item
}

func (r *Cache) SetExpires(key string, expires chrono.TimeStampIfc) {
	if r.Has(key) { r.cache[key].SetExpires(chrono.TimeStampIfc) }
}

// Get a single cache element by key name
func (r Cache) Get(key string) interface{} {
	if ci, ok := r.cache[key]; ok {
		if ! ci.IsExpired() { return ci.GetValue() }
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

// Drop an item from the cache with the supplied key
// return true if we drop it, else false
func (r *Cache) Drop(key string) bool {
	r.mutex.lock(); defer r.mutex.unlock()
	return r.drop(key)
}

// Check whether we have configuration elements for all the key names
// return count of items actually dropped
func (r *Cache) DropAll(keys *[]string) int {
	r.mutex.lock(); defer r.mutex.unlock()
	numDropped := 0
	for _, key := range *keys { if r.drop(key) { numDropped++ } }
	return numDropped
}

// Flush all the items out of the cache
func (r *Cache) Flush() {
	r.mutex.Lock(); defer r.mutex.Unlock()
	r.cache = make(map[string]cacheItem)
}

// -------------------------------------------------------------------------------------------------
// io.Closer Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) Close() error {
	r.closed = true
	r.Flush()
}

// -------------------------------------------------------------------------------------------------
// Cache Private Interface
// -------------------------------------------------------------------------------------------------

func (r Cache) isClosed() bool {
	return r.closed
}

// Purge expired cache items
func (r *Cache) pruneExpired() {
	// While the Cache has not been closed...
	for (! r.isClosed() {
		r.mutex.Lock(); defer r.mutex.Unlock()

		// Find which keys we need to purge because their cacheItem is expired
		purgeKeys := []string{}
		for key, ci := range r.cache {
			if ci.IsExpired() {
				// Expired items should be removed
				purgeKeys = append(purgeKeys, key)
			} else {
				// The first non-expired one we find means all others after it are non-expired!
				break
			}
		}
		// Purge them!
		for _, key := range purgeKeys { _ = r.drop(key) }

		sleep(60)
	}
}

// Prune the currently cached element collection to fit the new element within limits
// return boolean true if it will fit, else false (true doesn't indicate whether we did any pruning)
func (r *Cache) pruneToFit(key string, size int) bool {

	// Will it fit at all?
	if (r.sizeLimit > 0) && (size > r.sizeLimit) { return false }
	pruneCount := r.numToPrune(key, size)
	if 0 == pruneCount { return true }

	// TODO: Make sure we're not pruning more than some percentage threshold
	// (to minimize performance hits due to statistical outliers)

	// Prune starting at the back of the age list for the count we need to prune
	element := r.ageList.Back()
	for ; (nil != element) && (pruneCount > 0); pruneCount-- {
		dropKey := element.Value.(cacheItem).Key
		element = element.Next()
		_ = r.drop(dropKey)
	}
	return true
}

// Add content to front of age List and remember it by key in elements map
// return true if we set it, else false
func (r *Cache) set(key string, ci cacheItem) bool {
	if ! r.pruneToFit(key, ci.GetSize() { return false }
	_ = r.drop(key)
	r.ageListElements[key] = r.ageList.PushFront(ci)
	r.size += ci.GetSize()
	r.count++
	return true
}

// Drop if exists (don't bump on the find since we're going to drop it!)
// return bool true if we drop it, else false
func (r *lruCache) drop(key string) bool {
	if element := r.find(key, false); nil != element {
		r.size -= sizeable.Size(element.Value.(lruCacheItem))
		r.count--
		r.ageList.Remove(element)
		delete(r.ageListElements, key)
		return true
	}
	return false
}
