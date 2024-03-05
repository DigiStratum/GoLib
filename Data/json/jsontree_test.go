package json

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_JsonTree_NewJson_ReturnsInstance(t *testing.T) {
	// Setup
	jsonString := "{}"
	sut := NewJsonTree(&jsonString)

	// Verify
	ExpectNonNil(sut, t)
}


