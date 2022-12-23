package memcache

/*

Memcached data access

TODO:
 * Add example(s) and test coverage
*/

type MemcacheClientIfc interface {

	// Memcache Item Factory Functions
	NewCacheItem(name string, value *[]byte, flags uint32, expiresIn int32) *memcacheItem

	// Memcached Primitives
	FlushAll() error
	Get(key string) (MemcacheItemIfc, error)
	Touch(key string, seconds int32) error
	Set(item MemcacheItemIfc) error		// Always sets the key=value
	Add(item MemcacheItemIfc) error		// Only adds key=value if ! exists key already
	Delete(key string) error
	Inc(key string, delta uint64) (uint64, error)
	Dec(key string, delta uint64) (uint64, error)
	Replace(item MemcacheItemIfc) error
	Append(item MemcacheItemIfc) error
	Prepend(item MemcacheItemIfc) error
	CompareAndSwap(item MemcacheItemIfc) error

	// Memcached Helpers
	Ping() error

}

