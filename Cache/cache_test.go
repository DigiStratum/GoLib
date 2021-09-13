package cache

/*

Unit Tests for LRUCache

*/

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoTools/test"
	lib "github.com/DigiStratum/GoLib"
)

func TestThat_LRUCache_Size_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Verify
	ExpectInt(0, sut.Size(), t)
}

func TestThat_LRUCache_Size_IsCorrect_WithSomeEntries(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	expectedSize := 0;
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		expectedSize += len(content)
		ExpectTrue(sut.Set(key, content), t)
	}

	// Verify
	ExpectInt(expectedSize, sut.Size(), t)
}

func TestThat_LRUCache_Count_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Verify
	ExpectInt(0, sut.Count(), t)
}

func TestThat_LRUCache_Get_ReturnsNil_ForMissingKeys(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Verify
	ExpectNil(sut.Get("boguskey"), t)
}

func TestThat_LRUCache_Has_ReturnsFalse_ForMissingKeys(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Verify
	ExpectFalse(sut.Has("boguskey"), t)
}

func TestThat_LRUCache_Set_AddsNew_UnlimitedWithFixedContent(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	key := "fixedsizekey"
	content := "1234567890"
	sut.Set(key, content)	// Set our content
	res := sut.Get(key)	// And retrieve the same to check it out

	// Verify
	ExpectInt(1, sut.Count(), t)
	ExpectInt(10, sut.Size(), t)
	ExpectTrue(sut.Has(key), t)
	ExpectString(content, *res, t)
}

func TestThat_LRUCache_Set_ReplacesExisting_WithFixedContent(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	key := "fixedsizekey"
	content := "0123456789"
	sut.Set(key, "oldgarbage")	// First set some old garbage
	sut.Set(key, content)		// Then replace it with our contet
	res := sut.Get(key)		// And retrieve the same to check it out

	// Verify
	ExpectInt(1, sut.Count(), t)
	ExpectInt(10, sut.Size(), t)
	ExpectTrue(sut.Has(key), t)
	ExpectString(content, *res, t)
}

func TestThat_LRUCache_Drop_ReturnsFalse_ForMissingKeys(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Verify
	ExpectFalse(sut.Drop("boguskey"), t)
}

func TestThat_LRUCache_Drop_ReturnsTrue_WhenExistingKeyDropped(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	key := "fixedsizekey"
	sut.Set(key, "0123456789")

	// Verify
	ExpectTrue(sut.Drop(key), t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_LRUCache_SetLimits_LimitsCount_WhenSetAddsEntries(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	countLimit := 5
	sizeLimit := 0
	sut.SetLimits(sizeLimit, countLimit - 1) // count limit is one less than we want
	for i := 0; i < countLimit; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content%d", i)
		ExpectTrue(sut.Set(key, content), t)
	}

	// Verify
	ExpectInt(countLimit - 1, sut.Count(), t) // should only have count limit less one
	ExpectFalse(sut.Has("key0"), t)
	for i := 1; i < countLimit; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectString(content, *res, t)
	}
}

func TestThat_LRUCache_SetLimits_LimitsSize_WhenSetAddsEntries(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	count := 5
	sizeLimit := (count - 1) * 10	// Limit size at 10 chars * our count, less one
	sut.SetLimits(sizeLimit, 0)
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		ExpectTrue(sut.Set(key, content), t)
	}

	// Verify
	ExpectInt(count - 1, sut.Count(), t) // should only have count limit less one
	ExpectFalse(sut.Has("key0"), t)
	for i := 1; i < count; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectString(content, *res, t)
	}
}

func TestThat_LRUCache_SetLimits_LimitsCountAndSize_WhenSetAddsEntries(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	countLimit := 5
	sizeLimit := (countLimit - 1) * 10	// Limit size at 10 chars * our count, less one
	sut.SetLimits(sizeLimit, countLimit - 1) // count limit is one less than we want
	for i := 0; i < countLimit; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		ExpectTrue(sut.Set(key, content), t)
	}
	// Drop in a double-sized item which should displace two regular ones
	expectedKey := fmt.Sprintf("key%d", countLimit)
	expectedContent := "12345678901234567890"
	ExpectTrue(sut.Set(expectedKey, expectedContent), t)
	ExpectTrue(sut.Has(expectedKey), t)
	res := sut.Get(expectedKey)
	ExpectString(expectedContent, *res, t)

	// Verify
	ExpectInt(countLimit - 1, sut.Count(), t) // should only have count limit less TWO
	ExpectFalse(sut.Has("key0"), t)
	ExpectFalse(sut.Has("key1"), t)
	for i := 2; i < countLimit; i++ { // We expect the lowest (oldest one) is replaced with the newest
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content--%d", i)
		ExpectTrue(sut.Has(key), t)
		res := sut.Get(key)
		ExpectString(content, *res, t)
	}
}

func TestThat_LRUCache_SetLimits_PreventsSet_WhenEntryAddsExceedsLimit(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	sizeLimit := 5
	sut.SetLimits(sizeLimit, 0)

	// Verify
	ExpectFalse(sut.Set("anykey", "1234567890"),t)
}

func TestThat_LRUCache_SetLimits_AllowsSet_WhenEntryIsExactlyLimit(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	sizeLimit := 5
	sut.SetLimits(sizeLimit, 0)

	// Verify
	ExpectTrue(sut.Set("anykey", "12345"),t)
}

func TestThat_LRUCache_SetLimits_AllowsSet_WhenFullButEntryReplacesExisting(t *testing.T) {
	// Setup
	sut := lib.NewLRUCache()

	// Test
	key := "anykey"
	sizeLimit := 5
	sut.SetLimits(sizeLimit, 0)

	// Verify
	ExpectTrue(sut.Set(key, "12345"), t)
	content := "54321"
	ExpectTrue(sut.Set(key, content), t)
	ExpectInt(1, sut.Count(), t)
	res := sut.Get(key)
	ExpectString(content, *res, t)
}

