package golib_test

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
	sizeLimit := 0
	countLimit := 5
	sut.SetLimits(sizeLimit, countLimit)
	for i := 0; i < countLimit; i++ {
		key := fmt.Sprintf("key%d", i)
		content := fmt.Sprintf("content%d", i)
		ExpectTrue(sut.Set(key, content), t)
	}

	// Verify
	ExpectFalse(sut.Set("extrakey", "this attempt to add an entry should be blocked by the limit"), t)
}

