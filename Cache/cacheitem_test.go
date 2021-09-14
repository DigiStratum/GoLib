package cache

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

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


func TestThat_CacheItem_GetSize_Returns_Sizeable_Value(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		sizeable_target{},
		chrono.NewTimeSource().Now(),
	)

	// Verify
	ExpectInt64(FIXED_SIZE, sut.GetSize(), t)
}