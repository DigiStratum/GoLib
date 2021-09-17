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
 * Replace expiresList with a Binary Tree implementation that will allow us to quickly insert new items
   and find the expired ones for purging without having to scan the entire collection
 * Add support and/or change interface to make expires relative offset from NOW. This would allow a
   Touch() to use the same value if we store the offset with the cacheItem.

*/

import (
	"fmt"
	"sync"
	"time"
	"container/list"

	"github.com/DigiStratum/GoLib/Chrono"
	"github.com/DigiStratum/GoLib/Data/sizeable"
	cfg "github.com/DigiStratum/GoLib/Config"
	//"github.com/DigiStratum/GoLib/Process/runnable"
)

type expiringItems []*cacheItem

type CacheIfc interface {
	Configure(config cfg.ConfigIfc) error			// cfg.ConfigurableIfc
	SetTimeSource(timeSource chrono.TimeSourceIfc)
	IsEmpty() bool
	Size() int64
	Count() int
	Set(key string, value interface{}) bool
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
	totalCountLimit		int

	// A list of cache keys with least recently used at back
	usageList			*list.List
	// A list of cache keys sorted by expiration
	expiresList			expiringItems

	// Default TimeSource; can change to a different TimeSource, but cannot be nil
	timeSource		chrono.TimeSourceIfc
	newItemExpires		int64

	totalSize		int64
	totalSizeLimit		int64

	mutex			sync.Mutex
	pruneMutex		sync.Mutex
	closed			bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewCache() *Cache {
	cache := Cache{}
	cache.init()
	return &cache
}

// -------------------------------------------------------------------------------------------------
// cfg.ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) Configure(config cfg.ConfigIfc) error {
	if nil == config { return fmt.Errorf("Cache.Configure() - Configuration was nil") }

	// New items added to cache will expire in this count of seconds; 0 (default) = no expiration
	if config.Has("newItemExpires") {
		newItemExpires := config.GetInt64("newItemExpires")
		if nil != newItemExpires { r.newItemExpires = *newItemExpires }
	}

	// New items added to cache won't drive total count above this; 0 (default) = unlimited
	// When a limit is in place, the Least Recently Used (LRU) item will be evicted to make room for the new one
	if config.Has("totalCountLimit") {
		totalCountLimit := config.GetInt64("totalCountLimit")
		if nil != totalCountLimit { r.totalCountLimit = int(*totalCountLimit) }
	}

	// New items  added to cache we won't drive total size of all items above this; 0 = unlimited
	// When a limit is in place, the Least Recently Used (LRU) item(s) will be evicted to make room for the new one
	if config.Has("totalSizeLimit") {
		totalSizeLimit := config.GetInt64("totalSizeLimit")
		if nil != totalSizeLimit {
			r.totalSizeLimit = *totalSizeLimit
//fmt.Printf("Setting totalSizeLimit=%d\n", r.totalSizeLimit)
		}
	}

	return nil
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
func (r *Cache) Set(key string, value interface{}) bool {
	r.mutex.Lock(); defer r.mutex.Unlock()
	return r.set(key, value)
}

func (r *Cache) SetExpires(key string, expires chrono.TimeStampIfc) {
	if r.Has(key) {
		ci := (*r).cache[key]
		ci.SetExpires(expires)
	}
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
	r.mutex.Lock(); defer r.mutex.Unlock()
	return r.drop(key)
}

// Check whether we have configuration elements for all the key names
// return count of items actually dropped
func (r *Cache) DropAll(keys *[]string) int {
	r.mutex.Lock(); defer r.mutex.Unlock()
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
	return nil
}

// -------------------------------------------------------------------------------------------------
// GoLib/Process/runnable/RunnableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) Run() {
	if r.IsRunning() { return }
	r.init()
	go r.runLoop()
}

func (r Cache) IsRunning() bool {
	return ! r.closed
}

func (r *Cache) Stop() {
	r.Close()
}

// -------------------------------------------------------------------------------------------------
// Cache Private Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) init() {
	r.cache = make(map[string]cacheItem)
	r.usageList = list.New()
	r.expiresList = make(expiringItems, 0)
	r.timeSource = chrono.NewTimeSource()
	r.closed = false
}

func (r *Cache) runLoop() {
	// While the Cache has not been closed...
	for r.IsRunning() {
		r.pruneExpired()
		time.Sleep(60)
	}
}

// Purge expired cache items
func (r *Cache) pruneExpired() {
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
}

func (r *Cache) itemCanFit(key string, size int64) bool {
	// If it's bigger than the size limit, then it's impossible
	if (r.totalSizeLimit > 0) && (size > r.totalSizeLimit) { return false }
	return true
}

// How many existing cache entries must be pruned to fit one of this size?
func (r Cache) numToPrune(key string, size int64) int {
fmt.Printf("numToPrune() - START!\n")
	var pruneCount, replaceCount int
	var replaceSize int64
	if element := r.findUsageListElementByKey(key, false); nil != element {
		replaceCount = 1
		ci := r.cache[key]
		replaceSize = ci.Size()
	}

	// If there is a count limit in effect...
	if r.totalCountLimit > 0 {
		futureCount := len(r.cache) + 1 - replaceCount
		if futureCount > r.totalCountLimit { pruneCount = futureCount - r.totalCountLimit }
fmt.Printf("numToPrune() - (count limit) prune count=%d\n", pruneCount)
	}

	// If there is a size limit in effect...
	if r.totalSizeLimit > 0 {
		// If we add this to the cache without pruning, future size would be...
		futureSize := r.totalSize + size - replaceSize
		// If we break the size limit by adding...
		if futureSize > r.totalSizeLimit {
			pruneSize := futureSize - r.totalSizeLimit
			num := 0
			element := r.usageList.Back()
			for ; (nil != element) && (pruneSize > 0); num++ {
				pruneKey := element.Value.(string)
				ci := r.cache[pruneKey]
				pruneSize -= ci.Size()
				element = element.Next()
			}
			if num > pruneCount { pruneCount = num }
fmt.Printf(
	"numToPrune() - (size limit=%d, current size=%d, new size=%d) prune count=%d\n",
	r.totalSizeLimit,
	r.totalSize,
	r.totalSize+pruneSize,
	pruneCount,
)
		}
	}

fmt.Printf("numToPrune() - DONE!\n")
	return pruneCount
}

// Prune the currently cached element collection to established limits
func (r *Cache) pruneToLimits(key string, size int64) {
	r.pruneMutex.Lock(); defer r.pruneMutex.Unlock()

	// Does this item fit right now without any prune/purge?
	if (r.totalSizeLimit == 0) || (r.totalSize + size < r.totalSizeLimit) {
		if (r.totalCountLimit == 0) || (len(r.cache) + 1 < r.totalCountLimit) { return }
	}

	pruneCount := r.numToPrune(key, size)
	if 0 == pruneCount { return }

	// TODO: Make sure we're not pruning more than some percentage threshold
	// (to minimize performance hits due to statistical outliers)

	// Prune starting at the back of the age list for the count we need to prune
	element := r.usageList.Back()
	for ; (nil != element) && (pruneCount > 0); pruneCount-- {
		dropKey := element.Value.(string)
		r.drop(dropKey)
		element = element.Next()
	}
}

// Add content to front of age List and remember it by key in elements map
// return true if we set it, else false
func (r *Cache) set(key string, value interface{}) bool {

	// Get the size of the value
	newSize := sizeable.Size(value)

	// If size limit is in play and this value is bigger than that, then it won't fit
	if (0 < r.totalSizeLimit) && (newSize > r.totalSizeLimit) { return false }
	if ! r.itemCanFit(key, newSize) { return false }

	// Make a new cacheItem...

	// Set up an expiration time for this new item
	var expires chrono.TimeStampIfc
	if 0 == r.newItemExpires {
		expires = chrono.NewTimeStampForever()
	} else {
		expires = r.timeSource.Now().Add(r.newItemExpires)
	}

	ci := NewCacheItem(key, value, expires)

	// If this key already exists...
	var oldSize int64 = 0
	if r.Has(key) {
		// Replace the existing item
		// subtract out the old size so that we're adjusting totalSize to be the difference between the two
		oldSize = r.cache[key].Size()
		r.findUsageListElementByKey(key, true)
	} else {
		// Add the new item
		r.usageList.PushFront(key)
	}
	r.cache[key] = *ci
	r.totalSize += (newSize - oldSize)

	// Go do some pruning, async so that we can get back to the caller now
	go r.pruneToLimits(key, newSize)

	return true
}

// Drop if exists
// return bool true if we drop it, else false
func (r *Cache) drop(key string) bool {
	if ! r.Has(key) { return false }
	// Don't rejuvenate on the find since we're going to drop it!
	element := r.findUsageListElementByKey(key, false)
	if nil == element {
		// ERROR: Somehow we have desynched r.cache[] with r.usageList; they should have the same keys!
//fmt.Printf("Failed to drop key '%s'\n", key)
		return false
	}

	// Drop from the ordered usage list
	r.usageList.Remove(element)
	// Drop from the cache map
	size := r.cache[key].Size()
	r.totalSize -= size
//fmt.Printf("Size dropped=[%d]\n", size)
	delete(r.cache, key)
	return true
}

func (r *Cache) findUsageListElementByKey(key string, rejuvenate bool) *list.Element {
	// If the key is in the cache at all...
	if _, ok := r.cache[key]; ok {
		// Find the usageList element whose e.Value == key
		for e := r.usageList.Front(); e != nil; e = e.Next() {
//fmt.Printf("@Here\n")
			if ek, ok := e.Value.(string); ok {
				if ek != key { continue }
				// Found it!
				if rejuvenate {
					// Pull the element forward in the ageList
					r.usageList.MoveToFront(e)
					// Also touch the expiration time
					expires := r.timeSource.Now().Add(r.newItemExpires)
					ci := r.cache[key]
					ci.SetExpires(expires)
				}
				return e
			} else {
				// e.Value is not a string? Strange problem to have... ignore!
			}
		}
	}
	return nil
}
