package json

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func newJsonLexer() *jsonLexer {
	return &jsonLexer{}
}

// Validity

func TestThat_JsonLexer_LexDataValue_Returns_error_for_broken_UTF8_encoding(t *testing.T) {
	// Setup
	// ref: https://stackoverflow.com/questions/36426327/why-utf8-validstring-function-not-detecting-invalid-unicode-characters
	json := string([]byte{237, 159, 193})
	sut := newJsonLexer()

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_invalid_value_for_blank_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()

	// Test
	actual, actualErr := sut.LexDataValue("")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectFalse(actual.IsValid(), t) { return }
}

// Strings

func TestThat_JsonLexer_LexDataValue_Returns_string_value_for_string_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	expected := "Hello, \\\"Json\\\"!"
	json := "\"" + expected + "\""

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsString(), t) { return }
	if ! ExpectString(expected, actual.GetString(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_unclosed_string_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	expected := "Bogus Json!"
	json := "\"" + expected

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

// Objects

func TestThat_JsonLexer_LexDataValue_Returns_object_value_for_empty_object_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := " { } "

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsObject(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_unclosed_object_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := " {  "

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_object_value_for_object_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	expectedName := "prop1"
	expectedValue := "value1"
	json := "{\n\t\"" + expectedName + "\":\n\"" + expectedValue + "\"\n}"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsObject(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_object_value_with_expected_property_names(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	expectedName := "prop2"
	expectedValue := "value2"
	json := "{\"" + expectedName + "\":\"" + expectedValue + "\"}"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.HasObjectProperty(expectedName), t) { return }
}

// Arrays

func TestThat_JsonLexer_LexDataValue_Returns_array_value_for_empty_array_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := " [ ] "

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_unclosed_array_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := " [  "

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_array_value_for_array_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := "[\n\t\"A\",\n\t\"B\",\n\t\"C\"\n]"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
	if ! ExpectInt(3, actual.GetArraySize(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_array_with_hanging_comma(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := " [ 1, ] "

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

// Nulls

func TestThat_JsonLexer_LexDataValue_Returns_null_value_for_null_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := "null"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsNull(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_broken_null_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := "nu"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

// Booleans

func TestThat_JsonLexer_LexDataValue_Returns_boolean_values_for_various_booleans_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := "[ true, TRUE, TrUe, false, FALSE, FaLsE ]"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
	if ! ExpectInt(6, actual.GetArraySize(), t) { return }
	for i := 0; i < 6; i++ {
		boolValue := actual.GetArrayValue(i)
		if ! ExpectNonNil(boolValue, t) { return }
		if ! ExpectTrue(boolValue.IsBoolean(), t) { return }
		if i < 3 {
			if ! ExpectTrue(boolValue.GetBoolean(), t) { return }
		} else {
			if ! ExpectFalse(boolValue.GetBoolean(), t) { return }
		}
	}
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_broken_boolean_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := "tr"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

// Integers

func TestThat_JsonLexer_LexDataValue_Returns_integer_values_for_various_integers_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := "[ -9223372036854775808, 0, 9223372036854775807 ]"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
	if ! ExpectInt(3, actual.GetArraySize(), t) { return }

	intValue := actual.GetArrayValue(0)
	if ! ExpectNonNil(intValue, t) { return }
	if ! ExpectTrue(intValue.IsInteger(), t) { return }
	if ! ExpectInt64(-9223372036854775808, intValue.GetInteger(), t) { return }

	intValue = actual.GetArrayValue(1)
	if ! ExpectNonNil(intValue, t) { return }
	if ! ExpectTrue(intValue.IsInteger(), t) { return }
	if ! ExpectInt64(0, intValue.GetInteger(), t) { return }

	intValue = actual.GetArrayValue(2)
	if ! ExpectNonNil(intValue, t) { return }
	if ! ExpectTrue(intValue.IsInteger(), t) { return }
	if ! ExpectInt64(9223372036854775807, intValue.GetInteger(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_integer_overflow(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := make([]string, 2)
	json[0] = "-92233720368547758080"
	json[1] = "92233720368547758070"

	for _, js := range json {
		// Test
		actual, actualErr := sut.LexDataValue(js)

		// Verify
		if ! ExpectNil(actual, t) { return }
		if ! ExpectError(actualErr, t) { return }
	}
}

// Floats

func TestThat_JsonLexer_LexDataValue_Returns_float_values_for_various_floats_json(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := "[ 0.0, -3.14159, 2.9979E8, 6.62607015e-34]"

	// Test
	actual, actualErr := sut.LexDataValue(json)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(actual.IsValid(), t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
	if ! ExpectInt(4, actual.GetArraySize(), t) { return }

	// Precisely zero (... zero!)
	floatValue := actual.GetArrayValue(0)
	if ! ExpectNonNil(floatValue, t) { return }
	if ! ExpectTrue(floatValue.IsFloat(), t) { return }
	if ! ExpectFloat64(float64(0.0), floatValue.GetFloat(), t) { return }

	// Negative PI (negative float)
	floatValue = actual.GetArrayValue(1)
	if ! ExpectNonNil(floatValue, t) { return }
	if ! ExpectTrue(floatValue.IsFloat(), t) { return }
	if ! ExpectFloat64(float64(-3.14159), floatValue.GetFloat(), t) { return }

	// Speed of light (m/s) (positive exponent)
	floatValue = actual.GetArrayValue(2)
	if ! ExpectNonNil(floatValue, t) { return }
	if ! ExpectTrue(floatValue.IsFloat(), t) { return }
	if ! ExpectFloat64(float64(2.9979E8), floatValue.GetFloat(), t) { return }

	// Planck's Constant (negative exponent)
	floatValue = actual.GetArrayValue(3)
	if ! ExpectNonNil(floatValue, t) { return }
	if ! ExpectTrue(floatValue.IsFloat(), t) { return }
	if ! ExpectFloat64(float64(6.62607015e-34), floatValue.GetFloat(), t) { return }
}

func TestThat_JsonLexer_LexDataValue_Returns_error_for_float_overflow(t *testing.T) {
	// Setup
	sut := newJsonLexer()
	json := make([]string, 1)
	json[0] = "1E400" // Borrowed from strconv's atof_test.go as an error case

	for _, js := range json {
		// Test
		actual, actualErr := sut.LexDataValue(js)

		// Verify
		if ! ExpectNil(actual, t) { return }
		if ! ExpectError(actualErr, t) { return }
	}
}

