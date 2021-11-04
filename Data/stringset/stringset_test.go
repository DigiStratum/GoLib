package stringset

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewStringSet_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewStringSet()

	// Verify
	ExpectNonNil(sut, t)
}
