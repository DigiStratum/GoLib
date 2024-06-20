package cache

/*

Unit Tests for Cache

*/

import(
	"fmt"
	"time"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

	cfg "github.com/DigiStratum/GoLib/Config"
	"github.com/DigiStratum/GoLib/Chrono"
	"github.com/DigiStratum/GoLib/Data/sizeable"
	"github.com/DigiStratum/GoLib/Process/runnable"
)

const GOROUTINE_WAIT_MSEC	= 25

func TestThat_Cache_SetTimeSource_SetsTimeSource_WhenNonNil(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	ts := chrono.NewTimeSource()
	// Reset the timeSource to nil so that we can observe it change
	sut.timeSource = nil

	// Test
	sut.SetTimeSource(ts)

	// Verify
	ExpectNonNil(sut.timeSource, t)
}

func TestThat_Cache_SetTimeSource_DoesNotSetTimeSource_WhenNil(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	sut.SetTimeSource(nil)

	// Verify
	ExpectNonNil(sut.timeSource, t)
}

func TestThat_Cache_IsEmpty_ReturnsTrue_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectTrue(sut.IsEmpty(), t)
}

func TestThat_Cache_IsEmpty_ReturnsFalse_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Set("boguskey", "bogus value")

	// Verify
	ExpectFalse(sut.IsEmpty(), t)
}

func TestThat_Cache_Size_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectInt64(0, sut.Size(), t)
}

func TestThat_Cache_Size_IsCorrect_WithSomeEntries(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	var expectedSize int64 = 0;
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		expectedSize += sizeable.Size(content)
		ExpectTrue(sut.Set(key, content), t)
		ExpectTrue(sut.Has(key), t)
	}

	// Verify
	ExpectInt64(expectedSize, sut.Size(), t)
}

func TestThat_Cache_Count_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectInt(0, sut.Count(), t)
}

func TestThat_Cache_Get_ReturnsNil_ForMissingKey(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectNil(sut.Get("boguskey"), t)
}

func TestThat_Cache_Has_ReturnsFalse_ForMissingKey(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectFalse(sut.Has("boguskey"), t)
}

func TestThat_Cache_Set_AddsNew_UnlimitedWithFixedContent(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	key := "fixedsizekey"
	content := "1234567890"
	ExpectTrue(sut.Set(key, content), t)	// Set our content
	res := sut.Get(key)			// And retrieve the same to check it out
	ExpectNonNil(res, t)
	val := res.(string)

	// Verify
	ExpectInt(1, sut.Count(), t)
	ExpectInt64(sizeable.Size(content), sut.Size(), t)
	ExpectTrue(sut.Has(key), t)
	ExpectString(content, val, t)
}

func TestThat_Cache_Set_ReplacesExisting_WithFixedContent(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	key := "fixedsizekey"
	content := "0123456789"
	ExpectTrue(sut.Set(key, "oldgarbage"), t)	// First set some old garbage
	ExpectTrue(sut.Set(key, content), t)		// Then replace it with our contet
	res := sut.Get(key)				// And retrieve the same to check it out
	ExpectNonNil(res, t)
	val := res.(string)

	// Verify
	ExpectInt(1, sut.Count(), t)
	ExpectInt64(sizeable.Size(content), sut.Size(), t)
	ExpectTrue(sut.Has(key), t)
	ExpectString(content, val, t)
}

func TestThat_Cache_SetExpires_ReturnsFalse_ForMissingKey(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectFalse(sut.SetExpires("boguskey", nil), t)
}

func TestThat_Cache_SetExpires_ReturnsFalse_ForGoodKeyBadTimeStamp(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Set("boguskey", "bogus value")

	// Verify
	ExpectFalse(sut.SetExpires("boguskey", nil), t)
}

func TestThat_Cache_SetExpires_ReturnsFalse_ForGoodKeyGoodTimeStamp(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Set("boguskey", "bogus value")
	ts := chrono.NewTimeSource()

	// Verify
	ExpectTrue(sut.SetExpires("boguskey", ts.Now()), t)
}

func TestThat_Cache_Drop_ReturnsFalse_ForMissingKeys(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	res, err := sut.Drop("boguskey")

	// Verify
	ExpectNil(err, t)
	ExpectFalse(res, t)
}

func TestThat_Cache_Drop_ReturnsTrue_WhenExistingKeyDropped(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	key := "fixedsizekey"
	sut.Set(key, "0123456789")
	res, err := sut.Drop(key)

	// Verify
	ExpectNil(err, t)
	ExpectTrue(res, t)
	ExpectInt64(0, sut.Size(), t)
}

func TestThat_Cache_Set_CausesPruning_WhenCountOverLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	countLimit := 5
	config := cfg.NewConfig()
	config.Set("totalCountLimit", fmt.Sprintf("%d", countLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	for i := 0; i <= countLimit; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content%d", i)
		sut.Set(key, content)
		time.Sleep(GOROUTINE_WAIT_MSEC * time.Millisecond)	// pruning is an asynchronous operation - it needs time to run!
	}

	// Verify
	ExpectInt(countLimit, sut.Count(), t)	// should only have count limit
	ExpectFalse(sut.Has("key0"), t)		// The oldest on should be gone (pruned out due to limit)
	for i := 1; i <= countLimit; i++ {	// We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectNonNil(res, t)
		val := res.(string)
		ExpectString(content, val, t)
	}
}

func TestThat_Cache_Set_CausesPruning_WhenSizeOverLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	count := 5
	contentFormat := "content--##"
	sizeLimit := (count - 1) * int(sizeable.Size(contentFormat) + 1)	// Limit size at 10 chars * our count, less one
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%2d", i)
		sut.Set(key, content)
		time.Sleep(GOROUTINE_WAIT_MSEC * time.Millisecond)	// pruning is an asynchronous operation - it needs time to run!
	}

	// Verify
	ExpectInt(count - 1, sut.Count(), t)	// should cap out at the count-1 because of the size limit
	ExpectFalse(sut.Has("key0"), t)
	for i := 1; i < count; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%2d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectNonNil(res, t)
		val := res.(string)
		ExpectString(content, val, t)
	}
}

/*
FIXME: This test fails occasionally when system is busy/multi-tasking
=== RUN   TestThat_Cache_SetCausesPruning_WhenBothOverLimit
    expect.go:49:

        @/home/digistratum/Documents/Development/GoProjects/src/github.com/DigiStratum/GoLib/Cache/cache_test.go:298
        Expect: '4', Actual: '3'
    expect.go:64:

        @/home/digistratum/Documents/Development/GoProjects/src/github.com/DigiStratum/GoLib/Cache/cache_test.go:303:
        Expect: 'true', Actual: 'false'
    expect.go:88:

        @/home/digistratum/Documents/Development/GoProjects/src/github.com/DigiStratum/GoLib/Cache/cache_test.go:305:
        Expect: non-nil, Actual: nil
--- FAIL: TestThat_Cache_SetCausesPruning_WhenBothOverLimit (0.15s)
panic: interface conversion: interface {} is nil, not string [recovered]
	panic: interface conversion: interface {} is nil, not string

goroutine 82 [running]:
testing.tRunner.func1.2({0x520c00, 0xc00019a6f0})
	/usr/local/go/src/testing/testing.go:1209 +0x24e
testing.tRunner.func1()
	/usr/local/go/src/testing/testing.go:1212 +0x218
panic({0x520c00, 0xc00019a6f0})
	/usr/local/go/src/runtime/panic.go:1038 +0x215
command-line-arguments.TestThat_Cache_SetCausesPruning_WhenBothOverLimit(0xc000182340)
	/home/digistratum/Documents/Development/GoProjects/src/github.com/DigiStratum/GoLib/Cache/cache_test.go:306 +0x913
testing.tRunner(0xc000182340, 0x549738)
	/usr/local/go/src/testing/testing.go:1259 +0x102
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:1306 +0x35a
FAIL	command-line-arguments	0.434s
FAIL
ERROR!
*/
func TestThat_Cache_Set_CausesPruning_WhenBothOverLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	countLimit := 5
	//var sizeLimit int64 = int64(countLimit) * 10	// Limit size at 10 * our count
	contentFormat := "content--##"
	sizeLimit := (countLimit - 1) * int(sizeable.Size(contentFormat) + 1)	// Limit size at 10 chars * our count, less one

	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	config.Set("totalCountLimit", fmt.Sprintf("%d", countLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	for i := 0; i <= countLimit; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%2d", i)
		sut.Set(key, content)
		time.Sleep(GOROUTINE_WAIT_MSEC * time.Millisecond)	// pruning is an asynchronous operation - it needs time to run!
	}
	// Drop in a double-sized item which should displace two regular ones
	expectedKey := fmt.Sprintf("key%d", countLimit)
	expectedContent := "12345678901234567890"			// Items will be size 20+
	sut.Set(expectedKey, expectedContent)
	ExpectTrue(sut.Has(expectedKey), t)
	res := sut.Get(expectedKey)
	ExpectNonNil(res, t)
	val := res.(string)
	ExpectString(expectedContent, val, t)

	// Verify
	ExpectInt(countLimit - 1, sut.Count(), t) // should only have count limit less ONE
	ExpectFalse(sut.Has("key0"), t)
	ExpectFalse(sut.Has("key1"), t)
	for i := 2; i < countLimit; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectNonNil(res, t)
		val := res.(string)
		content := fmt.Sprintf("content--%2d", i)
		ExpectString(content, val, t)
	}
}

func TestThat_Cache_Configure_PreventsSet_WhenEntryAddedExceedsLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	sizeLimit := 5
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)
	val := "1234567890"
	sut.Set("anykey", val)

	// Verify
	ExpectFalse(sut.Has("anykey"), t)
}

func TestThat_Cache_Configure_AllowsSet_WhenEntryIsExactlyLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	content := "12345"
	sizeLimit := sizeable.Size(content)
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Verify
	sut.Set("anykey", content)
	ExpectTrue(sut.Has("anykey"),t)
}

func TestThat_Cache_Configure_AllowsSet_WhenFullButEntryReplacesExisting(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Test
	key := "anykey"
	content := "12345"
	sizeLimit := sizeable.Size(content)
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Verify
	sut.Set(key, "12345")
	ExpectTrue(sut.Has(key), t)
	ExpectInt(1, sut.Count(), t)

	content = "54321"
	sut.Set(key, content)
	ExpectTrue(sut.Has(key), t)
	ExpectInt(1, sut.Count(), t)
	res := sut.Get(key)
	ExpectNonNil(res, t)
	val := res.(string)
	ExpectString(content, val, t)
}

func TestThat_Cache_HasAll_ReturnsTrue_IfCacheHasAllKeys(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	keys := []string{ "key1", "key2" }
	sut.Set("key1", "value1")
	sut.Set("key2", "value2")

	// Verify
	ExpectTrue(sut.HasAll(&keys), t)
}

func TestThat_Cache_HasAll_ReturnsFalse_IfCacheDoesNotHaveAllKeys(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	keys := []string{ "key1", "key2" }
	sut.Set("key1", "value1")

	// Verify
	ExpectFalse(sut.HasAll(&keys), t)
}

func TestThat_Cache_DropAll_DropsSomeKeys_ForPartiallyMatchingKeySet(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	keys := []string{ "key1", "key2" }
	sut.Set("key1", "value1")

	// Test
	res, err := sut.DropAll(&keys)

	// Verify
	ExpectNil(err, t)
	ExpectInt(len(keys) - 1, res, t)
}
func TestThat_Cache_DropAll_DropsAllKeys_ForFullyMatchingKeySet(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	keys := []string{ "key1", "key2" }
	sut.Set("key1", "value1")
	sut.Set("key2", "value2")

	// Test
	res, err := sut.DropAll(&keys)

	// Verify
	ExpectNil(err, t)
	ExpectInt(len(keys), res, t)
}

func TestThat_Cache_Flush_DropsAllKeys(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Set("key1", "value1")
	sut.Set("key2", "value2")

	// Test
	sut.Flush()

	// Verify
	ExpectTrue(sut.IsEmpty(), t)
	ExpectInt64(0, sut.Size(), t)
	ExpectInt(0, sut.Count(), t)
}

func TestThat_Cache_Implements_RunnableIfc(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	var ifc interface{} = sut
	_, ok := ifc.(runnable.RunnableIfc)

	// Verify
	ExpectTrue(ok, t)
}

func TestThat_Cache_IsRunning_ReturnsTrue_AfterInitialization(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectTrue(sut.IsRunning(), t)
}

func TestThat_Cache_IsRunning_ReturnsFalse_AfterClose(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Close()

	// Verify
	ExpectFalse(sut.IsRunning(), t)
}

func TestThat_Cache_IsRunning_ReturnsFalse_AfterStop(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Stop()

	// Verify
	ExpectFalse(sut.IsRunning(), t)
}

func TestThat_Cache_IsRunning_ReturnsTrue_AfterRun(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Stop()
	sut.Run()

	// Verify
	ExpectTrue(sut.IsRunning(), t)
}

func TestThat_Cache_Set_SetsForeverTimeStamp_WhenCacheHasDefaultExpiresSetting(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	sut.Set("testkey", "value1")

	// Test
	timeStamp := sut.GetExpires("testkey")

	// Verify
	ExpectNonNil(timeStamp, t)
	ExpectTrue(timeStamp.IsForever(), t)
}

func TestThat_Cache_pruneExpired_LeavesItemsThatNeverExpire(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	sut.Set("foreverkey", "value1")

	// Test
	sut.pruneExpired()

	// Verify
	ExpectInt(1, sut.Count(), t)
}

/*
FIXME: This test fails occasionally when system is busy/multi-tasking
=== RUN   TestThat_Cache_pruneExpired_PurgesExpiredItems
    expect.go:49:

        @/home/digistratum/Documents/Development/GoProjects/src/github.com/DigiStratum/GoLib/Cache/cache_test.go:519
        Expect: '1', Actual: '2'
--- FAIL: TestThat_Cache_pruneExpired_PurgesExpiredItems (0.00s)

FIXME: Got another one because of bad luck when the system was not at all busy:
=== RUN   TestThat_Cache_pruneExpired_PurgesExpiredItems
    expect.go:49:

        @/home/digistratum/Documents/Development/GoProjects/src/github.com/DigiStratum/GoLib/Cache/cache_test.go:578
        Expect: '1', Actual: '2'
--- FAIL: TestThat_Cache_pruneExpired_PurgesExpiredItems (0.00s)

FIXME: Got another one, not busy at all. Seem like just timing fragility
=== RUN   TestThat_Cache_pruneExpired_PurgesExpiredItems
    expect.go:71:

        @/home/digistratum/Documents/Development/GoProjects/src/github.com/DigiStratum/GoLib/Cache/cache_test.go:587
        Expect: '1', Actual: '2'

*/
func TestThat_Cache_pruneExpired_PurgesExpiredItems(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	ts := chrono.NewTimeSource()
	config := cfg.NewConfig()
	config.Set("newItemExpires", "100")
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)
	sut.Set("expiredkey", "value1")
	expiration := ts.Now().Add(-10000)
	sut.SetExpires("expiredkey", expiration)
	sut.Set("futurekey", "value2")

	// Test
	sut.pruneExpired()

	// Verify
	ExpectInt(1, sut.Count(), t)
}

func TestThat_Cache_itemCanFit_ReturnsTrue_WhenUnlimited(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	var size int64 = 1000

	// Verify
	ExpectTrue(sut.itemCanFit(size), t)
}

func TestThat_Cache_itemCanFit_ReturnsTrue_WhenItemSizeUnderOrAtLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	var size int64 = 1000
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", size))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Verify
	ExpectTrue(sut.itemCanFit(size), t)
	ExpectTrue(sut.itemCanFit(size / 2), t)
}

func TestThat_Cache_itemCanFit_ReturnsFalse_WhenItemSizeOverLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	var size int64 = 1000
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", size))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Verify
	ExpectFalse(sut.itemCanFit(size + 1), t)
	ExpectFalse(sut.itemCanFit(size * 2), t)
}

func TestThat_Cache_numToPrune_ReturnsZero_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	key := "newkey"
	var size int64 = 1000

	// Verify
	ExpectInt(0, sut.numToPrune(key, size), t)
}

func TestThat_Cache_numToPrune_ReturnsZero_ForExistingItemKeyUnderSizeLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	key := "key"
	content := "12345"
	size := sizeable.Size(content)
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", size + (size / 2)))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Test
	sut.Set(key, content)

	// Verify
	ExpectInt(0, sut.numToPrune(key, size + 1), t)
}

func TestThat_Cache_numToPrune_ReturnsOne_ForNewItemKeyUnderSizeLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	key := "existingkey"
	content := "12345"
	size := sizeable.Size(content)
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", size + (size / 2)))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Test
	sut.Set(key, content)

	// Verify
	ExpectInt(1, sut.numToPrune("newkey", size + 1), t)
}

func TestThat_Cache_pruneToLimits_DropsOneItem_ForNewItemKeyUnderSizeLimit(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	key := "existingkey"
	content := "12345"
	size := sizeable.Size(content)
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", size + (size / 2)))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)
	newKey := "newkey"

	// Test
	sut.Set(key, content)
	sut.Set(newKey, content)
	sut.pruneToLimits(newKey, size)

	// Verify
	ExpectInt(1, sut.Count(), t)
}

func TestThat_Cache_findUsageListElementByKey_ReturnsNil_ForMissingKey(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()

	// Verify
	ExpectNil(sut.findUsageListElementByKey("boguskey", false), t)
}

func TestThat_Cache_findUsageListElementByKey_ReturnsNonNil_ForExistingKey(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	key := "key"

	// Test
	sut.Set(key, "content")

	// Verify
	ExpectNonNil(sut.findUsageListElementByKey(key, false), t)
}

func TestThat_Cache_findUsageListElementByKey_BumpsFirstItem_CausingSecondOneToDrop(t *testing.T) {
	// Setup
	sut := NewCache(); defer sut.Close()
	content := "12345"
	size := sizeable.Size(content)
	sizeLimit := (size * 2)+ (size / 2)	// Big enough to hold two of these, but not three!
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Test
	sut.Set("firstkey", content)
	sut.Set("secondkey", content)
	sut.findUsageListElementByKey("firstkey", true)
	sut.Set("thirdkey", content)
	sut.pruneToLimits("thirdkey", size)

	// Verify
	ExpectInt(2, sut.Count(), t)
	// LRU: with firstkey bumped, we expect secondkey pruned when thirdkey set
	ExpectFalse(sut.Has("secondkey"), t)
}
