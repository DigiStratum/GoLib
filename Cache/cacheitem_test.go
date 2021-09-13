package cache

import(
	"testing"

	. "github.com/DigiStratum/GoTools/test"

	"github.com/DigiStratum/GoLib/Data/sizeable"
	"github.com/DigiStratum/GoLib/Chrono"
)

const FIXED_SIZE = 333

type regular_target struct {
	Buffer	[50]int
}

type sizeable_target struct {
}

func (r sizeable_target) Size() int64 {
	return FIXED_SIZE
}


func TestThat_CacheItem_Uses_Sizeable(t *testing.T) {
	// Setup
	timestamp := NewTimeStamp(timeSource TimeSourceIfc) *TimeStamp
	value := sizeable_target{}
	sut := NewCacheItem("bogus", value, expires chrono.TimeStampIfc) *cacheItem
	// Verify
}