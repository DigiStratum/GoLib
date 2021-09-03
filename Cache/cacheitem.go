package cache

/*

An item in the cache which may hold any arbitrary value with an expiration timestamp. The value is
stored as an interface{}, so the consumer must be able to assert into the form it needs upon retrieval.

The expires value is specified as a UTC() timestamp. If 0, then the item never expires.
*/

type cacheItemIfc interface {
	IsExpired() bool
	SetExpires(expires int64)
	GetValue() interface{}
	GetSize() int64
}

type cacheItem struct {
	value	interface{}
	expires	int64
	size	int64
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------
func NewCacheItem(value interface{}, expires int64) *cacheItem {
	return &cacheItem{
		value: 		value,
		expires:	expires,
		size:		len(value),
	}
}

// -------------------------------------------------------------------------------------------------
// cacheItem Public Interface
// -------------------------------------------------------------------------------------------------

func (r cacheItem) IsExpired() bool {
	return (expires > 0) && (time.Now().Unix() < r.expires)
}

func (r *cacheItem) SetExpires(expires int64) {
	r.expires = expires
}

func (r cacheItem) GetValue() interface{} {
	return r.value
}

func (r cacheItem) GetSize() int64 {
	return r.size
}
