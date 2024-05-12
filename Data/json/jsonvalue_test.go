package json

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_JsonValue_NewJsonValue_ReturnsInstance(t *testing.T) {
	// Setup
	jsonString := "{}"
	sut := NewJsonValue(&jsonString)

	// Verify
	ExpectNonNil(sut, t)
}


