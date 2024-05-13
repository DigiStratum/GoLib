package json

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_JsonValue_NewJsonValue_ReturnsInstance(t *testing.T) {
	// Setup
	json := []rune("{}")
	sut := NewJsonValue(&json)

	// Verify
	ExpectNonNil(sut, t)
}


