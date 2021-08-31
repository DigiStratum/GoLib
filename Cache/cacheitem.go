package cache

type cacheItemIfc interface {
	IsExpired() bool
	GetValue() interface{}
}

type cacheItem struct {
	value	interface{}
	expires	int64
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------
func NewCacheItem(value interface{}, expires int64) *cacheItem {
	return &cacheItem{
		value: 		value,
		expires:	expires,
	}
}

// -------------------------------------------------------------------------------------------------
// cacheItem Public Interface
// -------------------------------------------------------------------------------------------------

func (r cacheItem) IsExpired() bool {
	return time.Now().Unix() < r.Expires
}

func (r cacheItem) GetValue() interface{} {
	return r.value
}
