package objectcollection

import(
        "testing"

        . "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_ObjectCollection_NewObjectCollection_ReturnsSomething(t *testing.T) {
        // Test
        sut := NewObjectCollection()

        // Verify
        ExpectNonNil(sut, t)
}

