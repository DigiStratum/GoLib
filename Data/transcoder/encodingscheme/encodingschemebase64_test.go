package encodingscheme

import(
        "testing"

        . "github.com/DigiStratum/GoLib/Testing"
)

/*
type EncodingSchemeIfc interface {
        SetEncodedValue(source *string) error
        GetEncodedValue() (*string, error)
        SetDecodedValue(source *string) error
        GetDecodedValue() (*string, error)
}
*/

func TestThat_EncodingSchemeBase64_NewEncodingSchemeBase64_ReturnsEncodingScheme(t *testing.T) {
        // Test
        sut := NewEncodingSchemeBase64()
	var suti interface{} = sut
	_, ok := suti.(EncodingSchemeIfc)

        // Verify
        ExpectNonNil(sut, t)
        ExpectTrue(ok, t)
}

func TestThat_EncodingSchemeBase64_SetEncodedValue_ReturnsError_WhenSourceIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()

	// Test
	err := sut.SetEncodedValue(nil)

        // Verify
        ExpectError(err, t)
}

func TestThat_EncodingSchemeBase64_GetEncodedValue_ReturnsError_WhenSourceIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()

	// Test
	actual, err := sut.GetEncodedValue()

        // Verify
        ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_EncodingSchemeBase64_GetEncodedValue_ReturnsOriginalValue_WhenSourceIsGood(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()
	expected := "YWJjMTIzIT8kKiYoKSctPUB+"
	err1 := sut.SetEncodedValue(&expected)

	// Test
	actual, err2 := sut.GetEncodedValue()

        // Verify
        ExpectNoError(err1, t)
        ExpectNoError(err2, t)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_EncodingSchemeBase64_GetDecodedValue_ReturnsError_WhenSourceIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()

	// Test
	actual, err := sut.GetDecodedValue()

        // Verify
        ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_EncodingSchemeBase64_GetDecodedValue_ReturnsOriginalValue_WhenSourceIsGood(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()
	expected := "abc123!?$*&()'-=@~"
	sut.SetDecodedValue(&expected)

	// Test
	actual, err := sut.GetDecodedValue()

        // Verify
        ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

func TestThat_EncodingSchemeBase64_GetDecodedValue_ReturnsDecodedValue_WhenSourceIsEncoded(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()
	source:= "YWJjMTIzIT8kKiYoKSctPUB+"
	sut.SetEncodedValue(&source)
	expected := "abc123!?$*&()'-=@~"

	// Test
	actual, err := sut.GetDecodedValue()

        // Verify
        ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_EncodingSchemeBase64_GetEncodedValue_ReturnsEncodedValue_WhenSourceIsDecoded(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()
	source := "abc123!?$*&()'-=@~"
	sut.SetDecodedValue(&source)
	expected := "YWJjMTIzIT8kKiYoKSctPUB+"

	// Test
	actual, err := sut.GetEncodedValue()

        // Verify
        ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

