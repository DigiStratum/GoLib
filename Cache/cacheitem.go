package cache

import(
	"fmt"

	"github.com/DigiStratum/GoLib/Chrono"
	"github.com/DigiStratum/GoLib/Data/sizeable"
)
/*

An item in the cache which may hold any arbitrary value with an expiration timestamp. The value is
stored as an interface{}, so the consumer must be able to assert into the form it needs upon retrieval.
*/

type cacheItemIfc interface {
	IsExpired() bool
	SetExpires(expires chrono.TimeStampIfc)
	GetValue() interface{}
	GetKey() string
	Size() int64
}

type cacheItem struct {
	key		string
	value		interface{}
	expires		chrono.TimeStampIfc
	size		int64
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------
func NewCacheItem(key string, value interface{}, expires chrono.TimeStampIfc) *cacheItem {
	return &cacheItem{
		key:		key,
		value: 		value,
		expires:	expires,
		size:		sizeable.Size(value),
	}
}

// -------------------------------------------------------------------------------------------------
// cacheItem Public Interface
// -------------------------------------------------------------------------------------------------

func (r cacheItem) IsExpired() bool {
	if nil != r.expires { return r.expires.IsPast() }
fmt.Printf("ci.IsExpired() -> false")
	return false
}

func (r *cacheItem) SetExpires(expires chrono.TimeStampIfc) {
	r.expires = expires
}

func (r cacheItem) GetKey() string {
	return r.key
}

func (r cacheItem) GetValue() interface{} {
	return r.value
}

// -------------------------------------------------------------------------------------------------
// SizeableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r cacheItem) Size() int64 {
	return r.size
}
