package json

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Inerface

func TestThat_JsonLexer_NewJsonLexer_ReturnsInstance(t *testing.T) {
	// Setup
	var sut JsonLexerIfc = NewJsonLexer()

	// Verify
	ExpectNonNil(sut, t)
}

// Validity

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

// Strings

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

// Objects

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

// Arrays

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

func TestThat_JsonLexer_LexJsonValue_Returns_array_value_for_array_json(t *testing.T) {
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

func TestThat_JsonLexer_LexJsonValue_Returns_error_for__array_with_trailing_comma(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := " [ 1, ] "

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNil(actual, t)
	ExpectError(actualErr, t)
}

// Nulls

func TestThat_JsonLexer_LexJsonValue_Returns_null_value_for_null_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := "null"

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsNull(), t)
}

func TestThat_JsonLexer_LexJsonValue_Returns_error_for_broken_null_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := "nu"

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNil(actual, t)
	ExpectError(actualErr, t)
}

// Booleans
func TestThat_JsonLexer_LexJsonValue_Returns_boolean_values_for_various_booleans_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := "[ true, TRUE, TrUe, false, FALSE, FaLsE ]"

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsArray(), t)
	ExpectInt(6, actual.GetArraySize(), t)
	for i := 0; i < 6; i++ {
		boolValue := actual.GetArrayElement(i)
		ExpectNonNil(boolValue, t)
		ExpectTrue(boolValue.IsBoolean(), t)
		if i < 3 {
			ExpectTrue(boolValue.GetBoolean(), t)
		} else {
			ExpectFalse(boolValue.GetBoolean(), t)
		}
	}
}

func TestThat_JsonLexer_LexJsonValue_Returns_error_for_broken_boolean_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := "tr"

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNil(actual, t)
	ExpectError(actualErr, t)
}

// Integers

func TestThat_JsonLexer_LexJsonValue_Returns_integer_values_for_various_integers_json(t *testing.T) {
	// Setup
	sut := NewJsonLexer()
	json := "[ -9223372036854775808, 0, 9223372036854775807 ]"

	// Test
	actual, actualErr := sut.LexJsonValue(json)

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(actualErr, t)
	ExpectTrue(actual.IsValid(), t)
	ExpectTrue(actual.IsArray(), t)
	ExpectInt(3, actual.GetArraySize(), t)

	intValue := actual.GetArrayElement(0)
	ExpectNonNil(intValue, t)
	ExpectTrue(intValue.IsInteger(), t)
	ExpectInt64(-9223372036854775808, intValue.GetInteger(), t)

	intValue = actual.GetArrayElement(1)
	ExpectNonNil(intValue, t)
	ExpectTrue(intValue.IsInteger(), t)
	ExpectInt64(0, intValue.GetInteger(), t)

	intValue = actual.GetArrayElement(2)
	ExpectNonNil(intValue, t)
	ExpectTrue(intValue.IsInteger(), t)
	ExpectInt64(9223372036854775807, intValue.GetInteger(), t)
}
