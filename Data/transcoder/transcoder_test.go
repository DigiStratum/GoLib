package transcoder

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Transcoder_NewTranscoder_ReturnsNewTranscoder_WithUnknownEncodingScheme(t *testing.T) {
	// Setup
	sut := NewTranscoder()

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(sut.encodingScheme == ES_UNKNOWN, t)
}
