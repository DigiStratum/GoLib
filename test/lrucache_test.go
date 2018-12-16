package golib_test

/*

Unit Tests for LRUCache

*/

import(
        "testing"

        test "github.com/DigiStratum/GoTools/test"
        lib "github.com/DigiStratum/GoLib"
)

func TestThat_LRUCache_Size_Is0_WhenNew(t *testing.T) {
        // Setup
        sut := lib.NewLRUCache()

        // Verify
        test.ExpectInt(0, sut.Size(), t)
}

func TestThat_LRUCache_Count_Is0_WhenNew(t *testing.T) {
        // Setup
        sut := lib.NewLRUCache()

        // Verify
        test.ExpectInt(0, sut.Count(), t)
}

func TestThat_LRUCache_Set_Works_WithFixedContent(t *testing.T) {
        // Setup
        sut := lib.NewLRUCache()

	// Test
	key := "fixedsizekey"
	content := "1234567890"
	sut.Set(key, content)

        // Verify
        test.ExpectInt(1, sut.Count(), t)
        test.ExpectInt(10, sut.Size(), t)
        test.ExpectBool(true, sut.Has(key), t)
        //test.ExpectString(true, sut.Has(key), t)
}

