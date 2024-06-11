package json

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_JsonLexer_NewJsonLexer_ReturnsInstance(t *testing.T) {
	// Setup
	var sut JsonLexerIfc = NewJsonLexer()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_error_for_broken_UTF8_encoding(t *testing.T) {
	// Setup
	// ref: https://stackoverflow.com/questions/36426327/why-utf8-validstring-function-not-detecting-invalid-unicode-characters
	json := string([]byte{237, 159, 193})
	sut := NewJsonLexer()

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNil(actual, t)
	ExpectError(actualErr, t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_invalid_value_for_blank_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()

	// Test
	actual, actualErr := sut.LexJsonValue("")

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectFalse(actual.IsValid(), t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_string_value_for_string_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	expected := "Hello, \\\"Json\\\"!"
	json := "\"" + expected + "\""

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsString(), t)
	ExpectString(expected, actual.GetString(), t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_error_for_unclosed_string_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	expected := "Bogus Json!"
	json := "\"" + expected

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNil(actual, t)
	ExpectError(actualErr, t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_object_value_for_empty_object_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := " { } "

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsObject(), t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_error_for_unclosed_object_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := " {  "

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNil(actual, t)
	ExpectError(actualErr, t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_object_value_for_object_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	expectedName := "prop1"
	expectedValue := "value1"
	json := "{\n\t\"" + expectedName + "\":\n\"" + expectedValue + "\"\n}"

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsObject(), t)
	ExpectTrue(actual.HasObjectProperty(expectedName), t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_array_value_for_empty_array_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := " [ ] "

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsArray(), t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_error_for_unclosed_array_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := " [  "

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNil(actual, t)
	ExpectError(actualErr, t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_arrar_array_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := "[\n\t\"A\",\n\t\"B\",\n\t\"C\",\n]"

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsArray(), t)
	ExpectInt(3, actual.GetArraySize(), t)
}

