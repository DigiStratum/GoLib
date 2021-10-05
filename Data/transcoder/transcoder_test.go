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

func TestThat_Transcoder_FromString_SetterRetainsProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"

	// Testing
	sut.FromString(&expected, ES_NONE)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(sut.encodingScheme == ES_NONE, t)
	actual := string(sut.content)
	ExpectString(expected, actual, t)
}

func TestThat_Transcoder_FromBytes_SetterRetainsProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"
	expectedBytes := []byte(expected)

	// Testing
	sut.FromBytes(&expectedBytes, ES_NONE)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(sut.encodingScheme == ES_NONE, t)
	actual := string(sut.content)
	ExpectString(expected, actual, t)
}

func TestThat_Transcoder_FromFile_SetterRetainsProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"

	// Testing
	sut.FromFile("transcoder_test.vegetable_soup.txt", ES_NONE)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(sut.encodingScheme == ES_NONE, t)
	actual := string(sut.content)
	ExpectString(expected, actual, t)
}
