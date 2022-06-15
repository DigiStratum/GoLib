package encodingscheme

import(
        "testing"

        . "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_EncodingSchemeHTMLEscape_NewEncodingSchemeHTMLEscape_ReturnsEncodingScheme(t *testing.T) {
        // Test
        sut := NewEncodingSchemeHTMLEscape()
	var suti interface{} = sut
	_, ok := suti.(EncodingSchemeIfc)

        // Verify
        ExpectNonNil(sut, t)
        ExpectTrue(ok, t)
}

func TestThat_EncodingSchemeHTMLEscape_SetEncodedValue_ReturnsError_WhenSourceIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()

	// Test
	err := sut.SetEncodedValue(nil)

        // Verify
        ExpectError(err, t)
}

func TestThat_EncodingSchemeHTMLEscape_SetDecodedValue_ReturnsError_WhenSourceIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()

	// Test
	err := sut.SetDecodedValue(nil)

        // Verify
        ExpectError(err, t)
}

func TestThat_EncodingSchemeHTMLEscape_GetEncodedValue_ReturnsError_WhenSourceIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()

	// Test
	actual, err := sut.GetEncodedValue()

        // Verify
        ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_EncodingSchemeHTMLEscape_GetEncodedValue_ReturnsOriginalValue_WhenSourceIsGood(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()
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

func TestThat_EncodingSchemeHTMLEscape_GetDecodedValue_ReturnsError_WhenSourceIsNil(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()

	// Test
	actual, err := sut.GetDecodedValue()

        // Verify
        ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_EncodingSchemeHTMLEscape_GetDecodedValue_ReturnsOriginalValue_WhenSourceIsGood(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()
	expected := "<&>"
	sut.SetDecodedValue(&expected)

	// Test
	actual, err := sut.GetDecodedValue()

        // Verify
        ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

func TestThat_EncodingSchemeHTMLEscape_GetDecodedValue_ReturnsDecodedValue_WhenSourceIsEncoded(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()
	source:= "&lt;&amp;&gt;"
	sut.SetEncodedValue(&source)
	expected := "<&>"

	// Test
	actual, err := sut.GetDecodedValue()

        // Verify
        ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_EncodingSchemeHTMLEscape_GetEncodedValue_ReturnsEncodedValue_WhenSourceIsDecoded(t *testing.T) {
	// Setup
        sut := NewEncodingSchemeHTMLEscape()
	source := "<&>"
	sut.SetDecodedValue(&source)
	expected := "&lt;&amp;&gt;"

	// Test
	actual, err := sut.GetEncodedValue()

        // Verify
        ExpectNoError(err, t)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

