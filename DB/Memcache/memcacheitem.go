package memcache

/*
A cache item representation for memcache data store.

*/

import (
	chrono "github.com/DigiStratum/GoLib/Chrono"
)

type MemcacheItemIfc  interface {
	GetKey() string
	GetValue() *[]byte
	GetFlags() uint32
	GetExpiration() chrono.TimeStampIfc

	// Chainable setters
	SetKey(key string) *memcacheItem
	SetValue(value *[]byte) *memcacheItem
	SetFlags(flags uint32) *memcacheItem
	SetExpiration(expiration chrono.TimeStampIfc) *memcacheItem
}

type memcacheItem struct {
	key		string			// Item key, unique id, max length 250
	value		[]byte			// Item value, max size based on memcached config
	flags		uint32			// 32 bit binary flag; app-defined, optional
	expiration	chrono.TimeStampIfc	// Expiration time in seconds
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMemcacheItem() *memcacheItem {
	return &memcacheItem{}
}

func (r *memcacheItem) GetKey() string {
	return r.key
}

func (r *memcacheItem) GetValue() *[]byte {
	return &r.value
}

func (r *memcacheItem) GetFlags() uint32 {
	return r.flags
}

func (r *memcacheItem) GetExpiration() chrono.TimeStampIfc {
	return r.expiration
}

func (r *memcacheItem) SetKey(key string) *memcacheItem {
	r.key = key
	return r
}

func (r *memcacheItem) SetValue(value *[]byte) *memcacheItem {
	r.value = *value
	return r
}

func (r *memcacheItem) SetFlags(flags uint32) *memcacheItem {
	r.flags = flags
	return r
}

func (r *memcacheItem) SetExpiration(expiration chrono.TimeStampIfc) *memcacheItem {
	r.expiration = expiration
	return r
}

