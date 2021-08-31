package cache

type cacheItemIfc interface {
	IsExpired() bool
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
	return time.Now().Unix() < r.Expires
}

func (r cacheItem) GetValue() interface{} {
	return r.value
}

func (r cacheItem) GetSize() int64 {
	return r.size
}
