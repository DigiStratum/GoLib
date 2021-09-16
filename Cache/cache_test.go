package cache

/*

Unit Tests for Cache

*/

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

	//cfg "github.com/DigiStratum/GoLib/Config"
	"github.com/DigiStratum/GoLib/Data/sizeable"
)

func TestThat_Cache_Size_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := NewCache()

	// Verify
	ExpectInt64(0, sut.Size(), t)
}

func TestThat_Cache_Size_IsCorrect_WithSomeEntries(t *testing.T) {
	// Setup
	sut := NewCache()

	// Test
	var expectedSize int64 = 0;
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		expectedSize += sizeable.Size(content)
		sut.Set(key, content)
		ExpectTrue(sut.Has(key), t)
	}

	// Verify
	ExpectInt64(expectedSize, sut.Size(), t)
}

func TestThat_Cache_Count_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := NewCache()

	// Verify
	ExpectInt(0, sut.Count(), t)
}

func TestThat_Cache_Get_ReturnsNil_ForMissingKeys(t *testing.T) {
	// Setup
	sut := NewCache()

	// Verify
	ExpectNil(sut.Get("boguskey"), t)
}

func TestThat_Cache_Has_ReturnsFalse_ForMissingKeys(t *testing.T) {
	// Setup
	sut := NewCache()

	// Verify
	ExpectFalse(sut.Has("boguskey"), t)
}

func TestThat_Cache_Set_AddsNew_UnlimitedWithFixedContent(t *testing.T) {
	// Setup
	sut := NewCache()

	// Test
	key := "fixedsizekey"
	content := "1234567890"
	sut.Set(key, content)	// Set our content
	res := sut.Get(key)	// And retrieve the same to check it out
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
	sut := NewCache()

	// Test
	key := "fixedsizekey"
	content := "0123456789"
	sut.Set(key, "oldgarbage")	// First set some old garbage
	sut.Set(key, content)		// Then replace it with our contet
	res := sut.Get(key)		// And retrieve the same to check it out
	ExpectNonNil(res, t)
	val := res.(string)

	// Verify
	ExpectInt(1, sut.Count(), t)
	ExpectInt64(sizeable.Size(content), sut.Size(), t)
	ExpectTrue(sut.Has(key), t)
	ExpectString(content, val, t)
}

/*
func TestThat_Cache_Drop_ReturnsFalse_ForMissingKeys(t *testing.T) {
	// Setup
	sut := NewCache()

	// Verify
	ExpectFalse(sut.Drop("boguskey"), t)
}

func TestThat_Cache_Drop_ReturnsTrue_WhenExistingKeyDropped(t *testing.T) {
	// Setup
	sut := NewCache()

	// Test
	key := "fixedsizekey"
	sut.Set(key, "0123456789")

	// Verify
	ExpectTrue(sut.Drop(key), t)
	ExpectInt64(0, sut.Size(), t)
}

func TestThat_Cache_SetLimits_LimitsCount_WhenSetAddsEntries(t *testing.T) {
	// Setup
	sut := NewCache()

	// Test
	countLimit := 5
	config := cfg.NewConfig()
	config.Set("totalCountLimit", fmt.Sprintf("%d", countLimit - 1)) // count limit is one less than we want
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	for i := 0; i < countLimit; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content%d", i)
		sut.Set(key, content)
	}

	// Verify
	ExpectInt(countLimit - 1, sut.Count(), t) // should only have count limit less one
	ExpectFalse(sut.Has("key0"), t)
	for i := 1; i < countLimit; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectNonNil(res, t)
		val := res.(string)
		ExpectString(content, val, t)
	}
}

func TestThat_Cache_SetLimits_LimitsSize_WhenSetAddsEntries(t *testing.T) {
	// Setup
	sut := NewCache()
	count := 5
	sizeLimit := (count - 1) * 10	// Limit size at 10 chars * our count, less one

	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		sut.Set(key, content)
	}

	// Verify
	ExpectInt(count - 1, sut.Count(), t) // should only have count limit less one
	ExpectFalse(sut.Has("key0"), t)
	for i := 1; i < count; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectNonNil(res, t)
		val := res.(string)
		ExpectString(content, val, t)
	}
}

func TestThat_Cache_SetLimits_LimitsCountAndSize_WhenSetAddsEntries(t *testing.T) {
	// Setup
	sut := NewCache()
	countLimit := 5
	sizeLimit := (countLimit - 1) * 10	// Limit size at 10 chars * our count, less one
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	config.Set("totalCountLimit", fmt.Sprintf("%d", countLimit - 1)) // count limit is one less than we want
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	for i := 0; i < countLimit; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		sut.Set(key, content)
	}
	// Drop in a double-sized item which should displace two regular ones
	expectedKey := fmt.Sprintf("key%d", countLimit)
	expectedContent := "12345678901234567890"
	sut.Set(expectedKey, expectedContent)
	ExpectTrue(sut.Has(expectedKey), t)
	res := sut.Get(expectedKey)
	ExpectNonNil(res, t)
	val := res.(string)
	ExpectString(expectedContent, val, t)

	// Verify
	ExpectInt(countLimit - 1, sut.Count(), t) // should only have count limit less TWO
	ExpectFalse(sut.Has("key0"), t)
	ExpectFalse(sut.Has("key1"), t)
	for i := 2; i < countLimit; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectNonNil(res, t)
		val := res.(string)
		ExpectString(content, val, t)
	}
}

func TestThat_Cache_SetLimits_PreventsSet_WhenEntryAddedExceedsLimit(t *testing.T) {
	// Setup
	sut := NewCache()

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

func TestThat_Cache_SetLimits_AllowsSet_WhenEntryIsExactlyLimit(t *testing.T) {
	// Setup
	sut := NewCache()

	// Test
	sizeLimit := 5
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Verify
	sut.Set("anykey", "12345")
	ExpectTrue(sut.Has("anykey"),t)
}

func TestThat_Cache_SetLimits_AllowsSet_WhenFullButEntryReplacesExisting(t *testing.T) {
	// Setup
	sut := NewCache()

	// Test
	key := "anykey"
	sizeLimit := 5
	config := cfg.NewConfig()
	config.Set("totalSizeLimit", fmt.Sprintf("%d", sizeLimit))
	err := sut.Configure(config)
	ExpectTrue((nil == err), t)

	// Verify
	sut.Set(key, "12345")
	ExpectTrue(sut.Has(key), t)
	ExpectInt(1, sut.Count(), t)

	content := "54321"
	sut.Set(key, content)
	ExpectTrue(sut.Has(key), t)
	ExpectInt(1, sut.Count(), t)
	res := sut.Get(key)
	ExpectNonNil(res, t)
//fmt.Printf("res:[%s]", *res)
	val := res.(string)
	ExpectString(content, val, t)
}
*/
