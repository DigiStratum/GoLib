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
	ExpectInt(0, len(sut.content), t)
}

func TestThat_Transcoder_FromString_SetterRetainsProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"

	// Testing
	res := sut.FromString(&expected, ES_NONE)

	// Verify
	ExpectNil(res, t)
	ExpectInt(1, len(sut.content), t)
	actual, ok := sut.content[ES_NONE]
	ExpectTrue(ok, t)
	ExpectString(expected, string(actual), t)
}

func TestThat_Transcoder_FromString_SetterReplacesProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"

	// Testing
	sut.FromString(&expected, ES_NONE)
	res := sut.FromString(&expected, ES_UNKNOWN)

	// Verify
	ExpectNil(res, t)
	ExpectInt(1, len(sut.content), t)
	actual, ok := sut.content[ES_UNKNOWN]
	ExpectTrue(ok, t)
	ExpectString(expected, string(actual), t)
}

func TestThat_Transcoder_FromBytes_SetterRetainsProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"
	expectedBytes := []byte(expected)

	// Testing
	res := sut.FromBytes(&expectedBytes, ES_NONE)

	// Verify
	ExpectNil(res, t)
	ExpectInt(1, len(sut.content), t)
	actual, ok := sut.content[ES_NONE]
	ExpectTrue(ok, t)
	ExpectString(expected, string(actual), t)
}

func TestThat_Transcoder_FromBytes_SetterReplacesProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"
	expectedBytes := []byte(expected)

	// Testing
	sut.FromBytes(&expectedBytes, ES_NONE)
	res := sut.FromBytes(&expectedBytes, ES_UNKNOWN)

	// Verify
	ExpectNil(res, t)
	ExpectInt(1, len(sut.content), t)
	actual, ok := sut.content[ES_UNKNOWN]
	ExpectTrue(ok, t)
	ExpectString(expected, string(actual), t)
}

func TestThat_Transcoder_FromFile_SetterRetainsProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"

	// Testing
	res := sut.FromFile("transcoder_test.vegetable_soup.txt", ES_NONE)

	// Verify
	ExpectNil(res, t)
	ExpectInt(1, len(sut.content), t)
	actual, ok := sut.content[ES_NONE]
	ExpectTrue(ok, t)
	ExpectString(expected, string(actual), t)
}

func TestThat_Transcoder_FromFile_SetterReplacesProperties_WithoutErrorResult(t *testing.T) {
	// Setup
	sut := NewTranscoder()
	expected := "Vegetable Soup!"

	// Testing
	sut.FromFile("transcoder_test.vegetable_soup.txt", ES_NONE)
	res := sut.FromFile("transcoder_test.vegetable_soup.txt", ES_UNKNOWN)

	// Verify
	ExpectNil(res, t)
	ExpectInt(1, len(sut.content), t)
	actual, ok := sut.content[ES_UNKNOWN]
	ExpectTrue(ok, t)
	ExpectString(expected, string(actual), t)
}

func TestThat_Transcoder_FromFile_SetterChangesNothing_WithErrorForMissingFile(t *testing.T) {
	// Setup
	sut := NewTranscoder()

	// Testing
	res := sut.FromFile("bogus_filename.txt", ES_NONE)

	// Verify
	ExpectNonNil(res, t)
}
