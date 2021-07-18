// DigiStratum GoLib - Cache
package golib

/*

This in-memory cache class wraps a basic Go map with essential helper functions to make life easier
for dealing with simple key/value pair data where the values are objects (interface{}) instead of
strings. An important differentiator of a cache vs. a hashmap is that the cache entries have time-
based expirations to ensure that none of the entries outlive their prescribed freshness.

TODO: Put some multi-threaded protections around the accessors here

*/

import (
	"time"
)

type cacheitem struct {
	Value	interface{}
	Expires	int64
}

func (ci cacheitem) IsExpired() bool {
	return time.Now().Unix() < ci.Expires
}

type cache	map[string]cacheitem

type CacheIfc interface {
	IsEmpty() bool
	Size() int
	Set(key string, value interface{}, expires int64)
	Get(key string) interface{}
	Has(key string) bool
	HasAll(keys *[]string) bool
	Flush()
}

// Make a new one of these!
func NewCache() CacheIfc {
	return &cache{}
}

// -------------------------------------------------------------------------------------------------
// Cache Public Interface
// -------------------------------------------------------------------------------------------------

// Flush all the items out of the cache
func (c *cache) Flush() {
	c = &cache{}
}

// Check whether this Cache is empty (has no properties)
func (c *cache) IsEmpty() bool {
	return 0 == c.Size()
}

// Get the number of properties in this Cache
func (c *cache) Size() int {
	c.purgeExpired()
	return len(*c)
}

// Set a single cache element key to the specified value
func (c *cache) Set(key string, value interface{}, expires int64) {
	(*c)[key] = cacheitem{
		Value:		value,
		Expires:	expires,
	}
}

// Get a single cache element by key name
func (c *cache) Get(key string) interface{} {
	c.purgeExpired()
	if ci, ok := (*c)[key]; ok { return ci.Value }
	return nil
}

// Check whether we have a configuration element by key name
func (c *cache) Has(key string) bool {
	c.purgeExpired()
	_, ok := (*c)[key];
	return ok
}

// Check whether we have configuration elements for all the key names
func (c *cache) HasAll(keys *[]string) bool {
	for _, key := range *keys {
		_, ok := (*c)[key];
		if ! ok { return false }
	}
	return true
}

// Purge expired cache items
func (c *cache) purgeExpired() {
	// Find which keys we need to purge because their cacheitem is expired
	purgeKeys := []string{}
	for key, ci := range *c {
		if ci.IsExpired() { purgeKeys = append(purgeKeys, key) }
	}
	// Purge them!
	for _, key := range purgeKeys { delete(*c, key) }
}

