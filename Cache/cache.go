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
 * Support optional Logger dependency injection (pass configuration in through DI as well?) so that we
   can log errors, stats, and more
 * Capture stats for things like sets, drops, hits, misses, purge operations, etc.
 * Add iterator for cache entry keys
 * After correcting the set/drop/pruning and mutex/async timing issues we would need to make a second
   pass to support prevention of the cache from growing beyond the configured size limit. As it is,
   we allow the size to break the limit temporarily and then prune it back down to the limit. We had
   some bits of ligic to attempt to acocunt for this but it was wrong and causing issues, so we have
   simplified in favor of circling back around to this concern later.

*/

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	chrono "github.com/DigiStratum/GoLib/Chrono"
	cfg "github.com/DigiStratum/GoLib/Config"
	"github.com/DigiStratum/GoLib/Data/sizeable"
)

type expiringItems []*cacheItem

type CacheIfc interface {
	Configure(config cfg.ConfigIfc) error // cfg.ConfigurableIfc
	SetTimeSource(timeSource chrono.TimeSourceIfc)
	IsEmpty() bool
	Size() int64
	Count() int
	Set(key string, value interface{}) bool
	SetExpires(key string, expires chrono.TimeStampIfc) bool
	GetExpires(key string) chrono.TimeStampIfc
	Get(key string) interface{}
	GetKeys() []string
	Has(key string) bool
	HasAll(keys *[]string) bool
	Drop(key string) (bool, error)
	DropAll(keys *[]string) (int, error)
	Flush()
	Close() error
}

type Cache struct {
	cache           map[string]*cacheItem
	totalCountLimit int

	// A list of cache keys with least recently used at back
	usageList *list.List
	// A list of cache keys sorted by expiration
	expiresList expiringItems

	// Default TimeSource; can change to a different TimeSource, but cannot be nil
	timeSource     chrono.TimeSourceIfc
	newItemExpires int64

	totalSize      int64
	totalSizeLimit int64

	mutex  sync.Mutex
	closed bool
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
	if nil == config {
		return fmt.Errorf("Cache.Configure() - Configuration was nil")
	}

	// New items added to cache will expire in this count of seconds; 0 (default) = no expiration
	if config.Has("newItemExpires") {
		newItemExpires := config.GetInt64("newItemExpires")
		if nil != newItemExpires {
			r.newItemExpires = *newItemExpires
		}
		//fmt.Printf("Cache::Configure() - newItemExpires = %d\n", r.newItemExpires)
	}

	// New items added to cache won't drive total count above this; 0 (default) = unlimited
	// When a limit is in place, the Least Recently Used (LRU) item will be evicted to make room for the new one
	if config.Has("totalCountLimit") {
		totalCountLimit := config.GetInt64("totalCountLimit")
		if nil != totalCountLimit {
			r.totalCountLimit = int(*totalCountLimit)
		}
		//fmt.Printf("Cache::Configure() - totalCountLimit = %d\n", r.totalCountLimit)
	}

	// New items  added to cache we won't drive total size of all items above this; 0 = unlimited
	// When a limit is in place, the Least Recently Used (LRU) item(s) will be evicted to make room for the new one
	if config.Has("totalSizeLimit") {
		totalSizeLimit := config.GetInt64("totalSizeLimit")
		if nil != totalSizeLimit {
			r.totalSizeLimit = *totalSizeLimit
		}
		//fmt.Printf("Cache::Configure() - totalSizeLimit = %d\n", r.totalSizeLimit)
	}

	return nil
}

// -------------------------------------------------------------------------------------------------
// CacheIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) SetTimeSource(timeSource chrono.TimeSourceIfc) {
	if nil != timeSource {
		r.timeSource = timeSource
	}
}

// Check whether this Cache is empty (has no properties)
func (r *Cache) IsEmpty() bool {
	return 0 == r.Count()
}

// Get the number of properties in this Cache
func (r *Cache) Size() int64 {
	return r.totalSize
}

// Return the count of entries currently being held in this cache
func (r *Cache) Count() int {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return len(r.cache)
}

// Set a single cache element key to the specified value
func (r *Cache) Set(key string, value interface{}) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.set(key, value)
}

// Set the expiration timestamp for a given Cache item; returns true if set, else false
func (r *Cache) SetExpires(key string, expires chrono.TimeStampIfc) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if !r.has(key) {
		return false
	}
	if nil == expires {
		return false
	}
	r.cache[key].SetExpires(expires)
	return true
}

// Get the expiration timestamp for a given Cache item; returns nil if not set
func (r *Cache) GetExpires(key string) chrono.TimeStampIfc {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if !r.has(key) {
		return nil
	}
	return r.cache[key].GetExpires()
}

// Get a single cache element by key name
func (r *Cache) Get(key string) interface{} {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if ci, ok := r.cache[key]; ok {
		if !ci.IsExpired() {
			return ci.GetValue()
		}
	}
	return nil
}

func (r *Cache) GetKeys() []string {
	keys := make([]string, len(r.cache))
	i := 0
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for key, _ := range r.cache {
		//fmt.Printf("Key: '%s'\n", key)
		keys[i] = key
		i++
	}
	return keys
}

// Check whether we have a configuration element by key name
func (r *Cache) Has(key string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.has(key)
}

// Check whether we have configuration elements for all the key names
func (r *Cache) HasAll(keys *[]string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, key := range *keys {
		if !r.has(key) {
			return false
		}
	}
	return true
}

// Drop an item from the cache with the supplied key
// return true if we drop it, else false
func (r *Cache) Drop(key string) (bool, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.drop(key)
}

// Check whether we have configuration elements for all the key names
// return count of items actually dropped
func (r *Cache) DropAll(keys *[]string) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	numDropped := 0
	for _, key := range *keys {
		dropped, err := r.drop(key)
		if nil != err {
			return 0, err
		}
		if dropped {
			numDropped++
		}
	}
	return numDropped, nil
}

// Flush all the items out of the cache
func (r *Cache) Flush() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.flush()
}

// -------------------------------------------------------------------------------------------------
// io.Closer Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) Close() error {
	r.closed = true
	r.flush()
	return nil
}

// -------------------------------------------------------------------------------------------------
// GoLib/Process/runnable/RunnableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) Run() {
	if r.IsRunning() {
		return
	}
	r.init()
	go r.runLoop()
}

func (r *Cache) IsRunning() bool {
	return !r.closed
}

func (r *Cache) Stop() {
	r.Close()
}

// -------------------------------------------------------------------------------------------------
// Cache Private Interface
// -------------------------------------------------------------------------------------------------

func (r *Cache) init() {
	r.flush()
	r.timeSource = chrono.NewTimeSource()
	r.closed = false
}

func (r *Cache) flush() {
	r.cache = make(map[string]*cacheItem)
	r.usageList = list.New()
	r.expiresList = make(expiringItems, 0)
	r.totalSize = 0
}

func (r *Cache) runLoop() {
	// While the Cache has not been closed...
	for r.IsRunning() {
		r.pruneExpired()
		time.Sleep(60)
	}
}
func (r *Cache) has(key string) bool {
	_, ok := r.cache[key]
	return ok
}

// Purge expired cache items
func (r *Cache) pruneExpired() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	//fmt.Printf("Cache::pruneExpired() ... \n")

	// Find which keys we need to purge because their cacheItem is expired
	purgeKeys := []string{}
	for key, ci := range r.cache {
		//fmt.Printf("Key: '%s' ...", key)
		if ci.IsExpired() {
			//fmt.Printf("Drop!\n")
			// Expired items should be removed
			purgeKeys = append(purgeKeys, key)
		} else {
			//fmt.Printf("Keep!\n")
			// The first non-expired one we find means all others after it are non-expired!
			// FIXME: ^^^ this seems like a lie! This is a hashmap, there is no sequencing of the keys - why would we think the ones that follow based on key iteration would have any newer/older timestamp - the collection is unsorted by the nature of the type of data structure!
			//break
		}
	}
	// Purge them!
	for _, key := range purgeKeys {
		if _, err := r.drop(key); nil != err {
			return err
		}
	}
	//fmt.Printf("Cache::pruneExpired() purged %d keys \n", len(purgeKeys))
	return nil
}

// Determine whether an item of this size fits our cache if it were the ONLY item
func (r *Cache) itemCanFit(size int64) bool {
	if (r.totalSizeLimit > 0) && (size > r.totalSizeLimit) {
		return false
	}
	return true
}

// How many existing cache entries must be pruned
func (r *Cache) numToPrune() int {
	var pruneCount int

	// If there is a count limit in effect...
	if r.totalCountLimit > 0 {
		if len(r.cache) > r.totalCountLimit {
			pruneCount = len(r.cache) - r.totalCountLimit
		}
	}

	// If there is a size limit in effect...
	if r.totalSizeLimit > 0 {
		// If we break the size limit by adding...
		if r.totalSize > r.totalSizeLimit {
			pruneSize := r.totalSize - r.totalSizeLimit
			num := 0
			element := r.usageList.Back()
			for ; (nil != element) && (pruneSize > 0); num++ {
				pruneKey := element.Value.(string)
				ci := r.cache[pruneKey]
				pruneSize -= ci.Size()
				element = element.Next()
			}
			if num > pruneCount {
				pruneCount = num
			}
		}
	}

	return pruneCount
}

// Prune the currently cached element collection to established limits
func (r *Cache) pruneToLimits() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	var err error
	//fmt.Printf("Cache::pruneToLimits() ... \n")

	// Are we within limits right now without any prune/purge?
	if (r.totalSizeLimit == 0) || (r.totalSize <= r.totalSizeLimit) {
		if (r.totalCountLimit == 0) || (len(r.cache) <= r.totalCountLimit) {
			//fmt.Printf("Cache::pruneToLimits() Nothing to do!\n")
			return nil
		}
	}

	pruneCount := r.numToPrune()
	if 0 < pruneCount {
		//fmt.Printf("Cache::pruneToLimits() pruning %d keys \n", pruneCount)

		// TODO: Make sure we're not pruning more than some percentage threshold
		// (to minimize performance hits due to statistical outliers)

		// Prune starting at the back of the age list for the count we need to prune
		element := r.usageList.Back()
		prunedCount := 0
		for ; (nil != element) && (prunedCount < pruneCount); prunedCount++ {
			dropKey := element.Value.(string)
			//fmt.Printf("Cache::pruneToLimits() dropping key %s \n", dropKey)
			if _, err = r.drop(dropKey); nil != err {
				break
			}
			element = element.Next()
		}
		//fmt.Printf("Cache::pruneToLimits() pruned %d keys \n", prunedCount)
	}

	return err
}

// Add content to front of age List and remember it by key in elements map
// return true if we set it, else false
func (r *Cache) set(key string, value interface{}) bool {

	// Get the size of the value
	newSize := sizeable.Size(value)
	//fmt.Printf("Cache::set() - key: '%s', size: %d\n", key, newSize)

	// If size limit is in play and this value is bigger than that, then it won't fit
	if !r.itemCanFit(newSize) {
		//fmt.Printf("Cache::set() - item too big to fit in cache\n")
		return false
	}

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
	if r.has(key) {
		// subtract out the old size so that we're adjusting totalSize to be the difference between the two
		oldSize = r.cache[key].Size()
		r.rejuvenateListElementByKey(key)
	} else {
		// Add the new item
		r.usageList.PushFront(key)
	}
	r.cache[key] = ci
	r.totalSize += (newSize - oldSize)

	// Go do some pruning, async so that we can get back to the caller now
	// TODO: Capture error from pruneToLimits() in stats... or log... or something!
	go r.pruneToLimits()

	return true
}

// Drop if exists
// return bool true if we drop it, else false
// If something strange causes a desync between the map and the list, then we will
// return true if either one has it and we drop it.
func (r *Cache) drop(key string) (bool, error) {
	var ret bool
	var err error
	//fmt.Printf("Cache::drop() - key: '%s'\n", key)
	foundMap := r.has(key)
	if foundMap {
		//fmt.Printf("Cache::drop() - foundMap\n")
		// Drop from the cache map
		size := r.cache[key].Size()
		r.totalSize -= size
		delete(r.cache, key)
		ret = true
	}

	// Don't rejuvenate on the find since we're going to drop it!
	element := r.findUsageListElementByKey(key, false)
	if nil != element {
		//fmt.Printf("Cache::drop() - found element\n")
		// Drop from the ordered usage list
		r.usageList.Remove(element)
		ret = true
	}

	// If foundMap XOR element (i.e. they don't match), then there was a desync
	if foundMap != (nil != element) {
		err = fmt.Errorf(
			"Cache.drop() - WARN: cache desync; found '%s' in either map or list but not both!",
			key,
		)
	}

	return ret, err
}

// Find the usageList element whose e.Value == key
func (r *Cache) findUsageListElementByKey(key string, rejuvenate bool) *list.Element {
	for e := r.usageList.Front(); e != nil; e = e.Next() {
		if ek, ok := e.Value.(string); ok {
			if ek != key {
				continue
			}
			return e
		}
	}
	return nil
}

func (r *Cache) rejuvenateListElementByKey(key string) {
	if element := r.findUsageListElementByKey(key, true); nil != element {
		// Pull the element forward in the ageList
		r.usageList.MoveToFront(element)
		// Also touch the expiration time
		expires := r.timeSource.Now().Add(r.newItemExpires)
		r.cache[key].SetExpires(expires)
	}
}
