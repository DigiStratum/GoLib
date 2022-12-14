package memcache

/*

Memcached data access

We're going to use someone else's library behind the scenes here for now, but abstract it so that we
can replace the library with something else without breaking consumers.

*/

import (
	mc "github.com/DigiStratum/go-bradfitz-gomemcache/memcache"
)

const MAX_KEY_LEN = 250

type MemcacheClientIfc interface {
	Ping() error
	FlushAll() error
	Get(key string) (*mc.Item, error) // FIXME: define our own item, don't use the one from the external library
	Touch(key string, seconds int32) error
	Set(item *Item) error		// Always sets the key=value
	Add(item *Item) error		// Only adds key=value if ! exists key already
	Delete(key string) error
	Inc(key string, delta uint64) (uint64, error)
	Dec(key string, delta uint64) (uint64, error)
	Replace(item *Item) error
	Append(item *Item) error
	Prepend(item *Item) error
	CompareAndSwap(item *Item) error
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
	for _, host in range hosts {
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

    func main() {
         mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
         mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

         it, err := mc.Get("foo")
         ...
    }

