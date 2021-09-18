package cache

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

	"github.com/DigiStratum/GoLib/Chrono"
)

const FIXED_SIZE = 333

type regularTarget struct {
	Buffer	[50]int
}

type sizeableTarget struct {
	Tag	string
}

func (r sizeableTarget) Size() int64 {
	return FIXED_SIZE
}


func TestThat_CacheItem_Size_Returns_Sizeable_Value(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		sizeableTarget{},
		chrono.NewTimeSource().Now(),
	)

	// Verify
	ExpectInt64(FIXED_SIZE, sut.Size(), t)
}

func TestThat_CacheItem_Size_Returns_Regular_Value(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		regularTarget{},
		chrono.NewTimeSource().Now(),
	)

	// Verify
	ExpectInt64(120, sut.Size(), t)
}

func TestThat_CacheItem_IsExpired_Returns_False(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		regularTarget{},
		chrono.NewTimeSource().Now().Add(1), // Expires 1 second from now
	)

	// Verify
	ExpectFalse(sut.IsExpired(), t)
}

func TestThat_CacheItem_IsExpired_Returns_True(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		regularTarget{},
		chrono.NewTimeSource().Now().Add(-1), // Expired 1 second ago
	)

	// Verify
	ExpectTrue(sut.IsExpired(), t)
}

func TestThat_CacheItem_SetExpires_Pushes_Causes_Expiration(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		regularTarget{},
		chrono.NewTimeSource().Now().Add(1), // Expires 1 second from now
	)

	// Test
	sut.SetExpires(chrono.NewTimeSource().Now().Add(-1)) // Expired 1 second ago

	// Verify
	ExpectTrue(sut.IsExpired(), t)
}

func TestThat_CacheItem_SetExpires_Pushes_Prevents_Expiration(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		regularTarget{},
		chrono.NewTimeSource().Now().Add(-1), // Expired 1 second ago
	)

	// Test
	sut.SetExpires(chrono.NewTimeSource().Now().Add(1)) // Expires 1 from now

	// Verify
	ExpectFalse(sut.IsExpired(), t)
}

func TestThat_CacheItem_GetValue_Returns_Original(t *testing.T) {
	// Setup
	inputItem := sizeableTarget{ Tag: "verificationtag!" }
	sut := NewCacheItem(
		"boguscacheitemkey",
		inputItem,
		chrono.NewTimeSource().Now(),
	)

	// Test
	outputItem := sut.GetValue()
	sizeableTargetItem, ok := outputItem.(sizeableTarget)

	// Verify
	ExpectTrue(ok && (sizeableTargetItem.Tag == "verificationtag!"), t)
}

func TestThat_CacheItem_GetKey_Returns_Key(t *testing.T) {
	// Setup
	sut := NewCacheItem(
		"boguscacheitemkey",
		regularTarget{},
		chrono.NewTimeSource().Now(),
	)

	// Verify
	ExpectString("boguscacheitemkey", sut.GetKey(), t)
}
