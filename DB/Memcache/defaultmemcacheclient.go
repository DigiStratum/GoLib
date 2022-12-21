package memcache

/*

MemcacheClient - Default Implementation

We're going to use someone else's library behind the scenes here for now, but abstract it so that we
can replace the library with something else without breaking consumers.


func main() {
	mc := memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	mc.Set(&memcache.MemcacheItemIfc{Key: "foo", Value: []byte("my value")})

	it, err := mc.Get("foo")
	...
}

TODO:
 * Port Brad Fitz library to our own code
*/

import (
	"fmt"

	mc "github.com/DigiStratum/go-bradfitz-gomemcache/memcache"

	chrono "github.com/DigiStratum/GoLib/Chrono"
)

const MAX_KEY_LEN = 250

type defaultMemcacheClient struct {
	hosts		[]string
	client		*mc.Client
	timeSource	chrono.TimeSourceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDefaultMemcacheClient(timeSource chrono.TimeSourceIfc, hosts ...string) (*defaultMemcacheClient, error) {
	if len(hosts) == 0 { return nil, fmt.Errorf("At least one memcached host must be specified") }
	if nil == timeSource { return nil, fmt.Errorf("TimeSource was nil!") }

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
	if nil != err { return nil, err }

	return &defaultMemcacheClient{
		hosts:		verifiedHosts,
		client:		client,
	}, nil
}

// -------------------------------------------------------------------------------------------------
// MemcacheClientIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *defaultMemcacheClient) NewCacheItem(key string, value *[]byte, flags uint32, expiresIn int32) *memcacheItem {
	var e chrono.TimeStampIfc = nil
	if 0 != expiresIn { e = r.timeSource.Now().Add(int64(expiresIn)) }
	return newMemcacheItem().
		SetKey(key).
		SetValue(value).
		SetFlags(flags).
		SetExpiresAt(e)
}

func (r *defaultMemcacheClient) Ping() error {
	return r.client.Ping()
}

func (r *defaultMemcacheClient) FlushAll() error {
	return r.client.FlushAll()
}

func (r *defaultMemcacheClient) Get(key string) (MemcacheItemIfc, error) {
	i, err := r.client.Get(key)
	if nil != err { return nil, err }
	return r.toItem(i), nil
}

func (r *defaultMemcacheClient) Touch(key string, seconds int32) error {
	return r.client.Touch(key, seconds)
}

// Always sets the key=value
func (r *defaultMemcacheClient) Set(item MemcacheItemIfc) error {
	return r.client.Set(r.fromItem(item))
}

// Only adds key=value if ! exists key already
func (r *defaultMemcacheClient) Add(item MemcacheItemIfc) error {
	return r.client.Add(r.fromItem(item))
}

func (r *defaultMemcacheClient) Delete(key string) error {
	return r.client.Delete(key)
}

func (r *defaultMemcacheClient) Inc(key string, delta uint64) (uint64, error) {
	return r.client.Increment(key, delta)
}

func (r *defaultMemcacheClient) Dec(key string, delta uint64) (uint64, error) {
	return r.client.Decrement(key, delta)
}

func (r *defaultMemcacheClient) Replace(item MemcacheItemIfc) error {
	return r.client.Replace(r.fromItem(item))
}

func (r *defaultMemcacheClient) Append(item MemcacheItemIfc) error {
	return r.client.Append(r.fromItem(item))
}

func (r *defaultMemcacheClient) Prepend(item MemcacheItemIfc) error {
	return r.client.Prepend(r.fromItem(item))
}

func (r *defaultMemcacheClient) CompareAndSwap(item MemcacheItemIfc) error {
	return r.client.CompareAndSwap(r.fromItem(item))
}

// -------------------------------------------------------------------------------------------------
// defaultMemcacheClient Implementation
// -------------------------------------------------------------------------------------------------

func (r *defaultMemcacheClient) toItem(i *mc.Item) MemcacheItemIfc {
	var e chrono.TimeStampIfc = nil
	if 0 != i.Expiration { e = r.timeSource.Now().Add(int64(i.Expiration)) }
	return newMemcacheItem().SetKey(i.Key).SetValue(&i.Value).SetFlags(i.Flags).SetExpiresAt(e)
}

func (r *defaultMemcacheClient) fromItem(memcacheItem MemcacheItemIfc) *mc.Item {
	var emptyValue []byte
	v := memcacheItem.GetValue()
	if nil == v { v = &emptyValue }
	return &mc.Item{
		Key:		memcacheItem.GetKey(),
		Value:		*v,
		Flags:		memcacheItem.GetFlags(),
		Expiration:	int32(memcacheItem.GetExpiresAt().DiffNow()),
	}
}

