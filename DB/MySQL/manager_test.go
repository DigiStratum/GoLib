package mysql

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewManager_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewManager()

	// Verify
	ExpectNonNil(sut, t)
}
