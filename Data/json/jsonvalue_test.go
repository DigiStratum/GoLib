package json

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Interface

func TestThat_JsonValue_NewJsonValue_ReturnsInstance(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(
		sut.IsNull() || sut.IsString() || sut.IsObject() || sut.IsArray() || sut.IsInteger() || sut.IsFloat() || sut.IsValid(),
		t,
	) { return }
}

// Validity

func TestThat_JsonValue_IsValid_Returns_false_for_new_value(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(sut.IsValid(), t) { return }
}

// Nulls

func TestThat_JsonValue_IsNull_Returns_true_after_setting_null(t *testing.T) {
	// Test
	sut := NewJsonValue().SetNull()

	// Verify
	if ! ExpectTrue(sut.IsNull(), t) { return }
}

// Strings

func TestThat_JsonValue_IsString_Returns_true_after_setting_string(t *testing.T) {
	// Setup
	expected1 := "hiyee!"
	expected2 := "BYee!"

	// Test
	sut := NewJsonValue().SetString(expected1)

	// Verify
	if ! ExpectTrue(sut.IsString(), t) { return }
	if ! ExpectString(expected1, sut.GetString(), t) { return }
	if ! ExpectString(expected2, sut.SetString(expected2).GetString(), t) { return }
}

// Objects

func TestThat_JsonValue_IsObject_Returns_true_after_preparing_object(t *testing.T) {
	// Test
	sut := NewJsonValue().PrepareObject()

	// Verify
	if ! ExpectTrue(sut.IsObject(), t) { return }
}

func TestThat_JsonValue_SetObjectProperty_Returns_error_for_non_object(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()

	// Test
	err := sut.SetObjectProperty("name", NewJsonValue())

	// Verify
	if ! ExpectError(err, t) { return }
}

func TestThat_JsonValue_SetObjectProperty_Returns_successfully_sets_value(t *testing.T) {
	// Setup
	sut := NewJsonValue().PrepareObject()
	expectedName := "name"
	expectedValue := "value"
	value := NewJsonValue()
	value.SetString(expectedValue)

	// Test
	err := sut.SetObjectProperty(expectedName, value)
	actualValue := sut.GetObjectProperty(expectedName)

	// Verify
	if ! ExpectNoError(err, t) { return }
	if ! ExpectTrue(sut.HasObjectProperty(expectedName), t) { return }
	if ! ExpectNonNil(actualValue, t) { return }
	if ! ExpectString(expectedValue, actualValue.GetString(), t) { return }
}

func TestThat_JsonValue_HasObjectProperty_Returns_false_for_non_property(t *testing.T) {
	// Setup
	sut := NewJsonValue().PrepareObject()

	// Verify
	if ! ExpectFalse(sut.HasObjectProperty("Nope!"), t) { return }
}

func TestThat_JsonValue_GetObjectPropertyNames_Returns_Empty_set(t *testing.T) {
	// Test
	actual := NewJsonValue().PrepareObject().GetObjectPropertyNames()

	// Verify
	if ! ExpectInt(0, len(actual), t) { return }
}

func TestThat_JsonValue_GetObjectPropertyNames_Returns_Expected_Names(t *testing.T) {
	// Setup
	sut := NewJsonValue().PrepareObject()
	expectedName := "name"
	value := NewJsonValue().SetString("value")
	sut.SetObjectProperty(expectedName, value)

	// Test
	actual := sut.GetObjectPropertyNames()

	// Verify
	if ! ExpectInt(1, len(actual), t) { return }
	if ! ExpectString(expectedName, actual[0], t) { return }
}

func TestThat_JsonValue_GetObjectProperty_Returns_nil_for_missing_property(t *testing.T) {
	// Test
	actual := NewJsonValue().PrepareObject().GetObjectProperty("missing property")

	// Verify
	if ! ExpectNil(actual, t) { return }
}

// Booleans

func TestThat_JsonValue_IsBoolean_Returns_true(t *testing.T) {
	// Setup
	sut := NewJsonValue().SetBoolean(true)

	// Verify
	if ! ExpectTrue(sut.IsBoolean(), t) { return }
	if ! ExpectTrue(sut.GetBoolean(), t) { return }
}

func TestThat_JsonValue_GetBoolean_Returns_false(t *testing.T) {
	// Setup
	sut := NewJsonValue().SetBoolean(false)

	// Verify
	if ! ExpectFalse(sut.GetBoolean(), t) { return }
}

// Arrays

func TestThat_JsonValue_IsArray_Returns_true(t *testing.T) {
	// Setup
	sut := NewJsonValue().PrepareArray()

	// Verify
	if ! ExpectTrue(sut.IsArray(), t) { return }
}

func TestThat_JsonValue_GetArraySize_Returns_zero_by_default(t *testing.T) {
	// Setup
	sut := NewJsonValue().PrepareArray()

	// Verify
	if ! ExpectInt(0, sut.GetArraySize(), t) { return }
}

func TestThat_JsonValue_GetArraySize_Returns_non_zero(t *testing.T) {
	// Setup
	sut := NewJsonValue().PrepareArray()
	sut.AppendArrayValue(NewJsonValue())

	// Verify
	if ! ExpectInt(1, sut.GetArraySize(), t) { return }
	sut.AppendArrayValue(NewJsonValue())
	if ! ExpectInt(2, sut.GetArraySize(), t) { return }
}

func TestThat_JsonValue_GetArrayValue_Returns_nil_for_non_array_or_mising_index(t *testing.T) {
	// Setup
	sut := NewJsonValue()

	// Verify
	if ! ExpectNil(sut.GetArrayValue(0), t) { return }
	sut.PrepareArray()
	if ! ExpectNil(sut.GetArrayValue(0), t) { return }
}

func TestThat_JsonValue_GetArrayValue_Returns_value_for_good_index(t *testing.T) {
	// Setup
	sut := NewJsonValue().PrepareArray()
	expectedValue := "value"
	sut.AppendArrayValue(NewJsonValue().SetString(expectedValue))

	// Test
	actual := sut.GetArrayValue(0)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expectedValue, actual.GetString(), t) { return }
}

// Floats

func TestThat_JsonValue_IsFloat_Returns_true(t *testing.T) {
	// Setup
	sut := NewJsonValue().SetFloat(9.9)

	// Verify
	if ! ExpectTrue(sut.IsFloat(), t) { return }
}

func TestThat_JsonValue_GetFloat_Returns_expected_value(t *testing.T) {
	// Setup
	var expectedValue float64 = 9.9
	sut := NewJsonValue().SetFloat(expectedValue)

	// Verify
	if ! ExpectFloat64(expectedValue, sut.GetFloat(), t) { return }
}

// Integers

func TestThat_JsonValue_IsInteger_Returns_true(t *testing.T) {
	// Setup
	sut := NewJsonValue().SetInteger(777)

	// Verify
	if ! ExpectTrue(sut.IsInteger(), t) { return }
}

func TestThat_JsonValue_GetInteger_Returns_expected_value(t *testing.T) {
	// Setup
	var expectedValue int64 = 777
	sut := NewJsonValue().SetInteger(expectedValue)

	// Verify
	if ! ExpectInt64(expectedValue, sut.GetInteger(), t) { return }
}

