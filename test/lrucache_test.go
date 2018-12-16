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


