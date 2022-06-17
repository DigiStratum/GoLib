package transcoder

import(
	"strings"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Transcoder_NewTranscoder_ReturnsNewTranscoder_WithUnknownEncodingScheme(t *testing.T) {
	// Setup
	sut := NewTranscoder()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Transcoder_Encode_ReturnsError_WithNoEncoderSet(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "bogus"

	// Test
	actual, err := sut.Encode(&expected)

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_Transcoder_Encode_ReturnsEncodedString_WithEncoderSet(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	sut.SetEncoderScheme(&MockEncodingScheme{GoUpper: true})
	expected := "bogus"

	// Test
	actual, err := sut.Encode(&expected)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString("BOGUS", *actual, t)
}

func TestThat_Transcoder_Decode_ReturnsError_WithNoDecoderSet(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "BOGUS"

	// Test
	actual, err := sut.Decode(&expected)

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_Transcoder_Decode_ReturnsDecodedString_WithDecoderSet(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	sut.SetDecoderScheme(&MockEncodingScheme{GoUpper: true})
	expected := "BOGUS"

	// Test
	actual, err := sut.Decode(&expected)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString("bogus", *actual, t)
}

func TestThat_Transcoder_Transcode_ReturnsError_WithNoDecoderSet(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	sut.SetEncoderScheme(&MockEncodingScheme{GoUpper: true})
	expected := "BOGUS"

	// Test
	actual, err := sut.Transcode(&expected)

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_Transcoder_Transcode_ReturnsError_WithNoEncoderSet(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	sut.SetDecoderScheme(&MockEncodingScheme{GoUpper: true})
	expected := "BOGUS"

	// Test
	actual, err := sut.Transcode(&expected)

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_Transcoder_Transcode_ReturnsTranscodedString_WithGoodEncodingSchemes(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	// A little counter intuitice, but the decoder goes from A -> a and encoder goes from a -> A
	// therefore a full transcode goes from A -> a -> A
	sut.SetDecoderScheme(&MockEncodingScheme{GoUpper: true})
	sut.SetEncoderScheme(&MockEncodingScheme{GoUpper: true})
	expected := "BOGUS"

	// Test
	actual, err := sut.Transcode(&expected)

	// Verify
	ExpectNoError(err, t)
	ExpectString(expected, *actual, t)
}

type MockEncodingScheme struct {
	GoUpper			bool
	encoded, decoded	*string
}

func (r *MockEncodingScheme) SetEncodedValue(source *string) error {
	r.encoded = source
	return nil
}

func (r *MockEncodingScheme) GetEncodedValue() (*string, error) {
	if nil == r.encoded {
		var encoded string
		if r.GoUpper {
			encoded = strings.ToUpper(*r.decoded)
		} else {
			encoded = strings.ToLower(*r.decoded)
		}
		r.encoded = &encoded
	}
	return r.encoded, nil
}

func (r *MockEncodingScheme) SetDecodedValue(source *string) error {
	r.decoded = source
	return nil
}

func (r *MockEncodingScheme) GetDecodedValue() (*string, error) {
	if nil == r.decoded {
		var decoded string
		if r.GoUpper {
			decoded = strings.ToLower(*r.encoded)
		} else {
			decoded = strings.ToUpper(*r.encoded)
		}
		r.decoded = &decoded
	}
	return r.decoded, nil
}

