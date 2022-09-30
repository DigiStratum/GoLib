package field

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_ObjectField_NewObjectField_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewObjectField("bogusfield")

	// Verify
	ExpectNonNil(sut, t)
}

