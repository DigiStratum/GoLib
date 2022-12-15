package memcache

/*

Memcached data access

We're going to use someone else's library behind the scenes here for now, but abstract it so that we
can replace the library with something else without breaking consumers.


from bradfitz/memcache/memcache.go:

// MemcacheItemIfc is an item to be got or stored in a memcached server.
type MemcacheItemIfc struct {
        // Key is the MemcacheItemIfc's key (250 bytes maximum).
        Key string

        // Value is the MemcacheItemIfc's value.
        Value []byte

        // Flags are server-opaque flags whose semantics are entirely
        // up to the app.
        Flags uint32

        // Expiration is the cache expiration time, in seconds: either a relative
        // time from now (up to 1 month), or an absolute Unix epoch time.
        // Zero means the MemcacheItemIfc has no expiration time.
        Expiration int32

        // Compare and swap ID.
        casid uint64
}

*/

import (
	mc "github.com/DigiStratum/go-bradfitz-gomemcache/memcache"
)

const MAX_KEY_LEN = 250

type MemcacheClientIfc interface {
	Ping() error
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
}

type memcacheClient struct {
	hosts		[]string
	client		*mc.Client
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMemcacheClient(hosts ...string) *memcacheClient {
	var verifiedHosts []string
	for _, host := range hosts {
		// TODO: Validate host network specifier as ip|hostname:port (check out net.Addr)
		// TODO: check name resolution and convert to IP (check out mc.Selector.SetServers() which has this already...
		// TODO: check name host reachability
		// TODO: fail if host is unreachable... (log an error)
		verifiedHosts = append(verifiedHosts, host)
	}

	// Connect and check!
	client := mc.New(hosts...)
	err := client.Ping()
	if nil != err {
		// TODO: Log error
		return nil
	}

	return &memcacheClient{
		hosts:		verifiedHosts,
		client:		client,
	}
}

/*
    func main() {
         mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
         mc.Set(&memcache.MemcacheItemIfc{Key: "foo", Value: []byte("my value")})

         it, err := mc.Get("foo")
         ...
    }
*/
