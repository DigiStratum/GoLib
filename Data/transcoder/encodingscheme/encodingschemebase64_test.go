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

func TestThat_EncodingSchemeBase64_SetEncodedValue_ReturnsNoError_WhenSourceIsString(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()
	testBase64 := "YWJjMTIzIT8kKiYoKSctPUB+"

	// Test
	err := sut.SetEncodedValue(&testBase64)

        // Verify
        ExpectNoError(err, t)
}

func TestThat_EncodingSchemeBase64_GetEncodedValue_ReturnsError_WhenRawIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeBase64()

	// Test
	actual, err := sut.GetEncodedValue()

        // Verify
        ExpectError(err, t)
	ExpectNil(actual, t)
}

