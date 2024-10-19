package data

import(
	"fmt"
	"strings"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Interface

func TestThat_DataValue_NewDataValue_ReturnsInstance(t *testing.T) {
	// Setup
	var sut DataValueIfc = NewDataValue() // Verifies that result satisfies IFC

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(
		sut.IsNull() || sut.IsString() || sut.IsObject() || sut.IsArray() || sut.IsInteger() || sut.IsFloat() || sut.IsValid(),
		t,
	) { return }
}

// Factory Functions

func TestThat_NewNull_clears_error_sets_type_and_returns_value(t *testing.T) {

	// Test
	actual := NewNull()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue((DATA_TYPE_NULL == actual.GetType()), t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
}

func TestThat_NewString_clears_error_sets_type_and_returns_value(t *testing.T) {
	// Setup
	expectedString := "howdy!"

	// Test
	actual := NewString(expectedString)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue((DATA_TYPE_STRING == actual.GetType()), t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
	if ! ExpectString(expectedString, actual.GetString(), t) { return }
}

func TestThat_NewObject_clears_error_sets_type_and_returns_value(t *testing.T) {
	// Test
	actual := NewObject()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue((DATA_TYPE_OBJECT == actual.GetType()), t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
}

func TestThat_NewBoolean_clears_error_sets_type_and_returns_value(t *testing.T) {
	// Test
	actual := NewBoolean(true)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue((DATA_TYPE_BOOLEAN == actual.GetType()), t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
	if ! ExpectTrue(actual.GetBoolean(), t) { return }
}

func TestThat_NewArray_clears_error_sets_type_and_returns_value(t *testing.T) {
	// Test
	actual := NewArray()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue((DATA_TYPE_ARRAY == actual.GetType()), t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
	if ! ExpectInt(0, actual.GetArraySize(), t) { return }
}

func TestThat_NewFloat_clears_error_sets_type_and_returns_value(t *testing.T) {
	// Setup
	var expectedFloat float64 = 3.14159

	// Test
	actual := NewFloat(expectedFloat)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue((DATA_TYPE_FLOAT == actual.GetType()), t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
	if ! ExpectFloat64(expectedFloat, actual.GetFloat(), t) { return }
}

func TestThat_NewInteger_clears_error_sets_type_and_returns_value(t *testing.T) {
	// Setup
	var expectedInt int64 = 333

	// Test
	actual := NewInteger(expectedInt)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue((DATA_TYPE_INTEGER == actual.GetType()), t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
	if ! ExpectInt64(expectedInt, actual.GetInteger(), t) { return }
}

// State

func TestThat_DataValue_IsValid_Returns_false_for_new_value(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(sut.IsValid(), t) { return }
}

func TestThat_DataValue_IsValid_Returns_true(t *testing.T) {
	// Setup
	sut := NewDataValue().SetString("whatever")

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectTrue(sut.IsValid(), t) { return }
}

func TestThat_DataValue_GetType_returns_expected_type(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Test
	if ! ExpectTrue(DATA_TYPE_INVALID == sut.GetType(), t) { return }
	if ! ExpectTrue(DATA_TYPE_INTEGER == sut.SetInteger(0).GetType(), t) { return }
	if ! ExpectTrue(DATA_TYPE_STRING == sut.SetString("howdy!").GetType(), t) { return }
	if ! ExpectTrue(DATA_TYPE_NULL == sut.SetNull().GetType(), t) { return }
	if ! ExpectTrue(DATA_TYPE_FLOAT == sut.SetFloat(3.14159).GetType(), t) { return }
	if ! ExpectTrue(DATA_TYPE_ARRAY == sut.PrepareArray().GetType(), t) { return }
	if ! ExpectTrue(DATA_TYPE_OBJECT == sut.PrepareObject().GetType(), t) { return }
}

func TestThat_DataValue_IsImmutable_returns_false_by_default(t *testing.T) {
	// Test
	sut := NewDataValue()

	// Verify
	if ! ExpectFalse(sut.IsImmutable(), t) { return }
}

func TestThat_DataValue_SetImmutable_changes_immutable_state(t *testing.T) {
	// Test
	sut := NewDataValue()

	// Verify
	if ! ExpectTrue(sut.SetImmutable().IsImmutable(), t) { return }
}

// Nulls

func TestThat_DataValue_IsNull_Returns_true_after_setting_null(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectTrue(sut.SetNull().IsNull(), t) { return }
}

func TestThat_DataValue_IsNull_Returns_false_for_non_null_value(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectFalse(sut.SetString("whatever").IsNull(), t) { return }
}

func TestThat_DataValue_SetNull_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewDataValue().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.SetNull()
	if ! ExpectError(sut.GetError(), t) { return }
}

// Strings

func TestThat_DataValue_GetString_Returns_empty_for_non_string(t *testing.T) {
	// Test
	sut := NewDataValue().SetNull()

	// Verify
	if ! ExpectString("", sut.GetString(), t) { return }
}

func TestThat_DataValue_IsString_Returns_false_for_non_string(t *testing.T) {
	// Test
	sut := NewDataValue().SetNull()

	// Verify
	if ! ExpectFalse(sut.IsString(), t) { return }
}

func TestThat_DataValue_IsString_Returns_true_after_setting_string(t *testing.T) {
	// Setup
	expected1 := "hiyee!"
	expected2 := "BYee!"

	// Test
	sut := NewDataValue().SetString(expected1)

	// Verify
	if ! ExpectTrue(sut.IsString(), t) { return }
	if ! ExpectString(expected1, sut.GetString(), t) { return }
	if ! ExpectString(expected2, sut.SetString(expected2).GetString(), t) { return }
}

func TestThat_DataValue_SetString_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewDataValue().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.SetString("whatever")
	if ! ExpectError(sut.GetError(), t) { return }
}

// Objects

func TestThat_DataValue_PrepareObject_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewDataValue().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.PrepareObject()
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_IsObject_Returns_false_for_non_object(t *testing.T) {
	// Test
	sut := NewDataValue().SetString("whatever")

	// Verify
	if ! ExpectFalse(sut.IsObject(), t) { return }
}

func TestThat_DataValue_IsObject_Returns_true_after_preparing_object(t *testing.T) {
	// Test
	sut := NewDataValue().PrepareObject()

	// Verify
	if ! ExpectTrue(sut.IsObject(), t) { return }
}

func TestThat_DataValue_SetObjectProperty_Returns_error_for_non_object(t *testing.T) {
	// Setup
	var sut DataValueIfc = NewDataValue()

	// Test
	actual := sut.SetObjectProperty("name", NewDataValue())
	actualErr := sut.GetError()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectError(actualErr, t) { return }
}

func TestThat_DataValue_SetObjectProperty_Returns_successfully_sets_value(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()
	expectedName := "name"
	expectedValue := "value"
	value := NewDataValue()
	value.SetString(expectedValue)

	// Test
	actual := sut.SetObjectProperty(expectedName, value)
	actualErr := sut.GetError()

	// Verify
	if ! ExpectNoError(actualErr, t) { return }
	if ! ExpectTrue(sut.HasObjectProperty(expectedName), t) { return }
	if ! ExpectNonNil(actual, t) { return }
	actualValue := sut.GetObjectProperty(expectedName)
	if ! ExpectNonNil(actualValue, t) { return }
	if ! ExpectString(expectedValue, actualValue.GetString(), t) { return }
}

func TestThat_DataValue_SetObjectProperty_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewObject().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.SetObjectProperty("prop", NewNull())
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_HasObjectProperty_Returns_false_for_non_property(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()

	// Verify
	if ! ExpectFalse(sut.HasObjectProperty("Nope!"), t) { return }
}

func TestThat_DataValue_HasObjectProperty_Sets_Error_For_Non_Objects(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(true)

	// Test
	errBefore := sut.GetError()
	actual := sut.HasObjectProperty("invalid")
	errAfter := sut.GetError()

	// Verify
	if ! ExpectNoError(errBefore, t) { return }
	if ! ExpectFalse(actual, t) { return }
	if ! ExpectError(errAfter, t) { return }
}

func TestThat_DataValue_HasObjectProperty_Returns_true(t *testing.T) {
	// Setup
	expectedName := "name"
	sut := NewDataValue().PrepareObject().SetObjectProperty(expectedName, NewDataValue())

	// Verify
	if ! ExpectTrue(sut.HasObjectProperty(expectedName), t) { return }
}

func TestThat_DataValue_DropObjectProperty_Returns_DataValue_And_Sets_Error_For_Non_Objects(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(true)

	// Test
	errBefore := sut.GetError()
	actual := sut.DropObjectProperty("invalid")
	errAfter := sut.GetError()

	// Verify
	if ! ExpectNoError(errBefore, t) { return }
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actual.GetBoolean(), t) { return }
	if ! ExpectError(errAfter, t) { return }
}

func TestThat_DataValue_DropObjectProperty_Returns_DataValue_And_Clears_Error_And_Drops_Property(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().SetObjectProperty("prop", NewDataValue().SetBoolean(true))

	// Test
	sut.AppendArrayValue(nil)			// Cause an error
	errBefore := sut.GetError()			// Verify that's the case
	actual := sut.DropObjectProperty("prop")	// Our unit to test
	errAfter := sut.GetError()			// Check the error state

	// Verify
	if ! ExpectError(errBefore, t) { return }
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actual.IsObject(), t) { return }
	if ! ExpectNoError(errAfter, t) { return }
}

func TestThat_DataValue_DropObjectProperty_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewObject().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.DropObjectProperty("prop")
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_GetObjectProperties_Returns_Empty_set(t *testing.T) {
	// Test
	actual := NewDataValue().PrepareObject().GetObjectProperties()

	// Verify
	if ! ExpectInt(0, len(actual), t) { return }
}

func TestThat_DataValue_GetObjectProperties_Returns_empty_set_and_Sets_Error_For_Non_Objects(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(true)

	// Test
	errBefore := sut.GetError()
	actual := sut.GetObjectProperties()
	errAfter := sut.GetError()

	// Verify
	if ! ExpectNoError(errBefore, t) { return }
	if ! ExpectInt(0, len(actual), t) { return }
	if ! ExpectError(errAfter, t) { return }
}

func TestThat_DataValue_GetObjectProperties_Returns_Expected_Names(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()
	expectedName := "name"
	value := NewDataValue().SetString("value")
	sut.SetObjectProperty(expectedName, value)

	// Test
	actual := sut.GetObjectProperties()

	// Verify
	if ! ExpectInt(1, len(actual), t) { return }
	if ! ExpectString(expectedName, actual[0], t) { return }
}

func TestThat_DataValue_GetObjectProperty_Returns_nil_for_missing_property(t *testing.T) {
	// Test
	actual := NewDataValue().PrepareObject().GetObjectProperty("missing property")

	// Verify
	if ! ExpectNil(actual, t) { return }
}

func TestThat_DataValue_GetObjectProperty_Returns_datavalue(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()
	expectedName := "name"
	expectedValue := "value"
	sut.SetObjectProperty(expectedName, NewDataValue().SetString(expectedValue))

	// Test
	actual := sut.GetObjectProperty(expectedName)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expectedValue, actual.GetString(), t) { return }
}

func TestThat_DataValue_HasAllObjectProperties_returns_false_for_non_objects(t *testing.T) {
	// Setup
	sut := NewDataValue().SetString("whatever")

	// Verify
	if ! ExpectFalse(sut.HasAllObjectProperties("bogus"), t) { return }
}

func TestThat_DataValue_HasAllObjectProperties_returns_false_for_missing_items(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()

	// Verify
	if ! ExpectFalse(sut.HasAllObjectProperties("bogus"), t) { return }
}

func TestThat_DataValue_HasAllObjectProperties_returns_true(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("name1", NewDataValue()).
		SetObjectProperty("name2", NewDataValue())

	// Verify
	if ! ExpectTrue(sut.HasAllObjectProperties("name1", "name2"), t) { return }
}

func TestThat_DataValue_GetMissingObjectProperties_returns_nil_for_non_objects(t *testing.T) {
	// Setup
	sut := NewDataValue().SetString("whatever")

	// Verify
	if ! ExpectNil(sut.GetMissingObjectProperties("bogus"), t) { return }
}

func TestThat_DataValue_GetMissingObjectProperties_returns_empty_set(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("name1", NewDataValue()).
		SetObjectProperty("name2", NewDataValue())

	actual := sut.GetMissingObjectProperties("name1", "name2")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(*actual), t) { return }
}

func TestThat_DataValue_GetMissingObjectProperties_returns_properties(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()

	actual := sut.GetMissingObjectProperties("name1", "name2")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(2, len(*actual), t) { return }
}

func TestThat_DataValue_DropObjectProperties_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewObject().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.DropObjectProperties("prop")
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_DropObjectProperties_returns_datavalue(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("name1", NewDataValue())

	// Verify
	if ! ExpectTrue(sut.DropObjectProperties().HasObjectProperty("name1"), t) { return }
}

func TestThat_DataValue_DropObjectProperties_returns_datavalue_and_sets_error_for_non_objects(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(true)

	// Test
	errBefore := sut.GetError()
	actual := sut.DropObjectProperties()
	errAfter := sut.GetError()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNil(errBefore, t) { return }
	if ! ExpectTrue(actual.GetBoolean(), t) { return }
	if ! ExpectNonNil(errAfter, t) { return }
}

func TestThat_DataValue_DropObjectProperties_drops_named_properties(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("name1", NewDataValue()).
		SetObjectProperty("name2", NewDataValue()).
		SetObjectProperty("name3", NewDataValue())

	// Test
	actual := sut.DropObjectProperties("name1", "name3")

	// Verify
	if ! ExpectTrue(actual.HasObjectProperty("name2"), t) { return }
	if ! ExpectFalse(actual.HasObjectProperty("name1"), t) { return }
	if ! ExpectFalse(actual.HasObjectProperty("name3"), t) { return }
}

func TestThat_DataValue_GetObjectProperty_Returns_nil_and_Sets_Error_For_Non_Objects(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(true)

	// Test
	errBefore := sut.GetError()
	actual := sut.GetObjectProperty("prop")
	errAfter := sut.GetError()

	// Verify
	if ! ExpectNoError(errBefore, t) { return }
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(errAfter, t) { return }
}

// Booleans

func TestThat_DataValue_SetBoolean_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewDataValue().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.SetBoolean(true)
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_IsBoolean_Returns_true(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(true)

	// Verify
	if ! ExpectTrue(sut.IsBoolean(), t) { return }
	if ! ExpectTrue(sut.GetBoolean(), t) { return }
}

func TestThat_DataValue_GetBoolean_Returns_false(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(false)

	// Verify
	if ! ExpectFalse(sut.GetBoolean(), t) { return }
}

// Arrays

func TestThat_DataValue_PrepareArray_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewDataValue().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.PrepareArray()
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_PrepareArray_clears_error_and_Returns_array(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Test
	actualBefore := sut.GetType()
	actualBeforeErr := sut.GetError()
	actualAfter := sut.PrepareArray().GetType()
	actualAfterErr := sut.GetError()

	// Verify
	if ! ExpectTrue((DATA_TYPE_INVALID == actualBefore), t) { return }
	if ! ExpectNoError(actualBeforeErr, t) { return }
	if ! ExpectTrue((DATA_TYPE_ARRAY == actualAfter), t) { return }
	if ! ExpectNoError(actualAfterErr, t) { return }
}

func TestThat_DataValue_IsArray_clears_error_and_Returns_true(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray()

	// Verify
	if ! ExpectTrue(sut.IsArray(), t) { return }
	if ! ExpectNoError(sut.GetError(), t) { return }
}

func TestThat_DataValue_GetArraySize_sets_error_and_Returns_zero_for_non_arrays(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectInt(0, sut.GetArraySize(), t) { return }
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_GetArraySize_Returns_zero_by_default(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray()

	// Verify
	if ! ExpectInt(0, sut.GetArraySize(), t) { return }
}

func TestThat_DataValue_GetArraySize_Returns_non_zero(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray()
	sut.AppendArrayValue(NewDataValue())

	// Verify
	if ! ExpectInt(1, sut.GetArraySize(), t) { return }
	sut.AppendArrayValue(NewDataValue())
	if ! ExpectInt(2, sut.GetArraySize(), t) { return }
}

func TestThat_DataValue_GetArrayValue_sets_error_and_Returns_nil_for_non_array_or_bad_index(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	if ! ExpectNil(sut.GetArrayValue(0), t) { return }
	if ! ExpectError(sut.GetError(), t) { return }

	sut.PrepareArray()
	if ! ExpectNoError(sut.GetError(), t) { return }
	if ! ExpectNil(sut.GetArrayValue(0), t) { return }
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_GetArrayValue_clears_error_and_Returns_value_for_good_index(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray()
	expectedValue := "value"
	sut.AppendArrayValue(NewDataValue().SetString(expectedValue))

	// Test
	actual := sut.GetArrayValue(0)

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expectedValue, actual.GetString(), t) { return }
}

func TestThat_DataValue_AppendArrayValue_sets_error_and_returns_original_for_non_array_or_nil_value(t *testing.T) {
	// Setup
	var expectedInt int64 = 333
	sut := NewDataValue().SetInteger(expectedInt)
	expectedValue := "value"

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	actual := sut.AppendArrayValue(NewDataValue().SetString(expectedValue))
	if ! ExpectError(sut.GetError(), t) { return }
	if ! ExpectInt64(expectedInt, actual.GetInteger(), t) { return }

	// Test
	sut.PrepareArray()
	if ! ExpectNoError(sut.GetError(), t) { return }
	actual = sut.AppendArrayValue(nil)
	if ! ExpectError(sut.GetError(), t) { return }
	if ! ExpectTrue((DATA_TYPE_ARRAY == actual.GetType()), t) { return }
}

func TestThat_DataValue_AppendArrayValue_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewArray().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.AppendArrayValue(NewNull())
	if ! ExpectError(sut.GetError(), t) { return }
}

// ReplaceArrayValue(index int, dataValue *DataValue) *DataValue

func TestThat_DataValue_ReplaceArrayValue_Replaces_existing_value(t *testing.T) {
	// Setup
	sut := NewArray().AppendArrayValue(NewString("original"))
	expected := "replaced"

	// Test
	sut.ReplaceArrayValue(0, NewString(expected))

	// Verify
	if ! ExpectString(expected, sut.GetArrayValue(0).GetString(), t) { return }
}

func TestThat_DataValue_ReplaceArrayValue_sets_error_for_index_too_high(t *testing.T) {
	// Setup
	sut := NewArray()

	// Test
	sut.ReplaceArrayValue(0, NewString("boing!"))

	// Verify
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_ReplaceArrayValue_sets_error_for_negative_index(t *testing.T) {
	// Setup
	sut := NewArray()

	// Test
	sut.ReplaceArrayValue(-1, NewString("boing!"))

	// Verify
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_ReplaceArrayValue_replaces_existing_with_null_value(t *testing.T) {
	// Setup
	sut := NewArray().AppendArrayValue(NewString("original"))

	// Test
	sut.ReplaceArrayValue(0, nil)

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	if ! ExpectTrue(sut.GetArrayValue(0).IsNull(), t) { return }
}

func TestThat_DataValue_ReplaceArrayValue_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewArray().AppendArrayValue(NewString("original")).SetImmutable()

	// Test
	sut.ReplaceArrayValue(0, NewString("boing!"))

	// Verify
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_ReplaceArrayValue_sets_error_for_non_array(t *testing.T) {
	// Setup
	sut := NewObject()

	// Test
	sut.ReplaceArrayValue(0, NewString("boing!"))

	// Verify
	if ! ExpectError(sut.GetError(), t) { return }
}

// Floats

func TestThat_DataValue_SetFloat_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewArray().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.SetFloat(3.14159)
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_IsFloat_Returns_true(t *testing.T) {
	// Setup
	sut := NewDataValue().SetFloat(9.9)

	// Verify
	if ! ExpectTrue(sut.IsFloat(), t) { return }
}

func TestThat_DataValue_GetFloat_Returns_expected_value(t *testing.T) {
	// Setup
	var expectedValue float64 = 9.9
	sut := NewDataValue().SetFloat(expectedValue)

	// Verify
	if ! ExpectFloat64(expectedValue, sut.GetFloat(), t) { return }
}

func TestThat_DataValue_GetFloat_returns_0value_for_non_integers(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectFloat64(0.0, sut.GetFloat(), t) { return }
}

// Integers

func TestThat_DataValue_SetInteger_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewArray().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.SetInteger(333)
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_IsInteger_Returns_true(t *testing.T) {
	// Setup
	sut := NewDataValue().SetInteger(777)

	// Verify
	if ! ExpectTrue(sut.IsInteger(), t) { return }
}

func TestThat_DataValue_GetInteger_Returns_expected_value(t *testing.T) {
	// Setup
	var expectedValue int64 = 777
	sut := NewDataValue().SetInteger(expectedValue)

	// Verify
	if ! ExpectInt64(expectedValue, sut.GetInteger(), t) { return }
}

func TestThat_DataValue_GetInteger_returns_0value_for_non_integers(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectInt64(0, sut.GetInteger(), t) { return }
}

// Conveniences

func makeBigDataValue() *DataValue {
	return NewObject().
		SetObjectProperty("shape", NewDataValue().SetString("arc")).
		SetObjectProperty( "vectors", NewArray().
			AppendArrayValue(
				NewObject().
					SetObjectProperty("radians", NewDataValue().SetFloat(float64(3.14159))).
					SetObjectProperty("radius", NewDataValue().SetInteger(int64(2))).
					SetObjectProperty("color", NewDataValue().SetString("red")).
					SetObjectProperty("hidden", NewDataValue().SetBoolean(false)),
			).
			AppendArrayValue(
				NewObject().
					SetObjectProperty("radians", NewDataValue().SetFloat(float64(6.28318))).
					SetObjectProperty("radius", NewDataValue().SetInteger(int64(7))).
					SetObjectProperty("color", NewDataValue().SetString("blue")).
					SetObjectProperty("hidden", NewDataValue().SetBoolean(true)),
			),
		)
}

// Select

func TestThat_DataValue_Select_Returns_Values_For_Good_Selectors(t *testing.T) {
	// Setup
	sut := makeBigDataValue()

	// Test
	selectors := []string{ "shape", "vectors", "vectors[1]", "vectors[1].radians" }
	for _, selector := range selectors {
		//t.Logf("Testing good selector '%s'\n", selector)
		// Verify
		actual1 := sut.Select(selector)
		if ! ExpectNonNil(actual1, t) { return }
		if ! ExpectNoError(sut.GetError(), t) { return }

		switch actual1.GetType() {
			case DATA_TYPE_STRING: if ! ExpectString("arc", actual1.GetString(), t) { return }
			case DATA_TYPE_ARRAY:
				if ! ExpectTrue(actual1.IsArray(), t) { return }
				if ! ExpectInt(2, actual1.GetArraySize(), t) { return }
			case DATA_TYPE_OBJECT:
				if ! ExpectTrue(actual1.IsObject(), t) { return }
				if ! ExpectTrue(actual1.HasObjectProperty("radians"), t) { return }
			case DATA_TYPE_FLOAT: if ! ExpectFloat64(float64(6.28318), actual1.GetFloat(), t) { return }
		}
	}
}

func TestThat_DataValue_Select_Returns_Errors_For_Bad_Selectors(t *testing.T) {
	// Setup
	sut := makeBigDataValue()

	// Test
	selectors := []string{ "[]", "[0]", "bogusproperty", "vectors[2]", "vectors[0]radians" }
	for _, selector := range selectors {
		actual := sut.Select(selector)
		// Verify
		if ! ExpectNil(actual, t) { return }
	}
}

func TestThat_DataValue_Select_Returns_nil_and_error_for_Missing_Object_Property(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()

	// Test
	actual := sut.Select(".missingprop")

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_Select_Returns_nil_and_error_for_bad_array_indexes(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray().
		AppendArrayValue(NewDataValue().SetString("apples")).
		AppendArrayValue(NewDataValue().SetString("oranges"))

	// Verify
	actual1 := sut.Select("[0]")
	if ! ExpectNonNil(actual1, t) { return }
	if ! ExpectNoError(sut.GetError(), t) { return }
	actual2 := sut.Select("[]")
	if ! ExpectNil(actual2, t) { return }
	if ! ExpectError(sut.GetError(), t) { return }
	actual3 := sut.Select("[A]")
	if ! ExpectNil(actual3, t) { return }
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_Select_Returns_nil_and_error_for_bad_object_property_names(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("prop", NewDataValue().SetString("value"))

	// Verify
	actual1 := sut.Select("prop")		// <- ok
	if ! ExpectNonNil(actual1, t) { return }
	if ! ExpectNoError(sut.GetError(), t) { return }
	actual2 := sut.Select(".\nprop")		// <- white space = fail
	if ! ExpectNil(actual2, t) { return }
	if ! ExpectError(sut.GetError(), t) { return }
	actual3 := sut.Select(".")		// <- no name = fail
	if ! ExpectNil(actual3, t) { return }
	if ! ExpectError(sut.GetError(), t) { return }
}

// Merge

func TestThat_DataValue_Merge_sets_error_when_immutable(t *testing.T) {
	// Setup
	sut := NewArray().SetImmutable()

	// Verify
	if ! ExpectNoError(sut.GetError(), t) { return }
	sut.Merge(NewArray())
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_Merge_sets_error_when_given_nil(t *testing.T) {
	// Setup
	sut := NewArray()

	// Test
	sut.Merge(nil)

	// Verify
	if ! ExpectError(sut.GetError(), t) { return }
}

func TestThat_DataValue_Merge_Returns_DataValue_for_unstructured_values(t *testing.T) {
	// Setup
	sut := NewDataValue().SetBoolean(true)
	m := NewDataValue().SetInteger(0)

	// Test
	actual := sut.Merge(m)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actual.GetBoolean(), t) { return }
}

func TestThat_DataValue_Merge_Merges_DataValue_for_Objects(t *testing.T) {
	// Setup
	sut := NewObject()
	mergeData := NewObject().SetObjectProperty("prop", NewDataValue().SetString("value"))

	// Test
	actual := sut.Merge(mergeData)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actual.HasObjectProperty("prop"), t) { return }
}

func TestThat_DataValue_Merge_Merges_DataValue_for_Arrays(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray()
	m := NewDataValue().PrepareArray().AppendArrayValue(NewDataValue().SetInteger(33))

	// Test
	actual := sut.Merge(m)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(1, actual.GetArraySize(), t) { return }
	dv := actual.GetArrayValue(0)
	if ! ExpectNonNil(dv, t) { return }
	if ! ExpectInt64(33, dv.GetInteger(), t) { return }
}

// HasAll

func TestThat_DataValue_HasAll_returns_true_for_empty_list(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectTrue(sut.HasAll(), t) { return }
}

func TestThat_DataValue_HasAll_returns_false_for_missing_selector(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectFalse(sut.HasAll("bogus"), t) { return }
}

// GetMissing

func TestThat_DataValue_GetMissing_returns_empty_set(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Test
	actual := sut.GetMissing()

	// Verify
	if ! ExpectInt(0, len(actual), t) { return }
}

func TestThat_DataValue_GetMissing_returns_missing_selectors(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Test
	actual := sut.GetMissing("shape1", "shape2")

	// Verify
	if ! ExpectInt(2, len(actual), t) { return }
}

func TestThat_DataValue_GetMissing_notices_matching_selectors(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("shape1", NewDataValue().SetString("arc")).
		SetObjectProperty("shape2", NewDataValue().SetString("arc"))

	// Test
	actual := sut.GetMissing("shape1", "shape2")

	// Verify
	if ! ExpectInt(0, len(actual), t) { return }
}

// GetIterator

func TestThat_DataValue_GetIterator_Returns_nil(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Test
	actual := sut.GetIterator()

	// Verify
	// Note: Can't use ExpectNil() here because there's something about the nil zero-value of
	// actual once passed as an argument to the Expect*() func that prevents it from being seen
	// as truly nil any longer.
	if ! ExpectTrue(nil == actual, t) { return }
}

func TestThat_DataValue_GetIterator_Returns_func(t *testing.T) {
	// Setup
	sut := NewDataValue()
	sut.PrepareObject()

	// Test
	actual := sut.GetIterator()

	// Verify
	// Note: Can't use ExpectNonNil() here because there's something about the nil zero-value of
	// actual once passed as an argument to the Expect*() func that prevents it from being seen
	// as truly nil any longer.
	if ! ExpectFalse(nil == actual, t) { return }
}

func TestThat_DataValue_GetIterator_Returns_Object_Iterator(t *testing.T) {
	// Setup
	sut := NewDataValue()
	sut.PrepareObject()
	numprops := 3
	for i := 1; i <= numprops; i++ {
		sut.SetObjectProperty(fmt.Sprintf("p%d", i), NewDataValue().SetInteger(int64(i)))
	}

	// Test
	actual := sut.GetIterator()

	// Verify
	// Note: Can't use ExpectNonNil() here because there's something about the nil zero-value of
	// actual once passed as an argument to the Expect*() func that prevents it from being seen
	// as truly nil any longer.
	if ! ExpectFalse(nil == actual, t) { return }

	actualProps := make([]bool, numprops)
	for i := 1; i <= numprops; i++ {
		kvpi := actual()
		kvp, ok := kvpi.(KeyValuePair)
		// Expect a KeyValuePair to be the result of calling the Iterator func
		if ! ExpectTrue(ok, t) { return }
		actualValue := kvp.Value.GetInteger()
		// Expect the KVP value to be an integer betwee 0 and numprops, inclusive
		if ! ExpectTrue((actualValue >= 0) && (actualValue <= int64(numprops)), t) { return }
		// Expect the KVP Key to be a string starting with "p" then the value
		if ! ExpectString(fmt.Sprintf("p%d", actualValue), kvp.Key, t) { return }
		actualProps[actualValue - 1] = true
	}
	final := actual()
	// Expect the Iterator to return nothing after consuming all the expected KVP values
	if ! ExpectNil(final, t) { return }
	// Expect all the actual properties to have been represented (they come in unordered as object props)
	for i := 1; i <= numprops; i++ {
		if ! ExpectTrue(actualProps[i - 1], t) { return }
	}
}

func TestThat_DataValue_GetIterator_Returns_Array_Iterator(t *testing.T) {
	// Setup
	sut := NewDataValue()
	sut.PrepareArray()
	numvals := 3
	for i := 1; i <= numvals; i++ {
		sut.AppendArrayValue(NewDataValue().SetInteger(int64(i)))
	}

	// Test
	actual := sut.GetIterator()

	// Verify
	// Note: Can't use ExpectNonNil() here because there's something about the nil zero-value of
	// actual once passed as an argument to the Expect*() func that prevents it from being seen
	// as truly nil any longer.
	if ! ExpectFalse(nil == actual, t) { return }

	for i := 1; i <= numvals; i++ {
		ivpi := actual()
		ivp, ok := ivpi.(IndexValuePair)
		// Expect a KeyValuePair to be the result of calling the Iterator func
		if ! ExpectTrue(ok, t) { return }
		if ! ExpectNonNil(ivp, t) { return }
		if ! ExpectInt64(int64(i), ivp.Value.GetInteger(),t) { return }
	}
}

// ToString

func TestThat_DataValue_ToString_Returns_empty_string_for_invalid(t *testing.T) {
	// Verify
	if ! ExpectString("", NewDataValue().ToString(), t) { return }
}

func TestThat_DataValue_ToString_Returns_string_for_primitives(t *testing.T) {
	// Verify
	if ! ExpectString("null", NewDataValue().SetNull().ToString(), t) { return }
	if ! ExpectString("apple", NewDataValue().SetString("apple").ToString(), t) { return }
	if ! ExpectString("50", NewDataValue().SetInteger(50).ToString(), t) { return }
	if ! ExpectString("6.28318", NewDataValue().SetFloat(6.28318).ToString(), t) { return }
	if ! ExpectString("false", NewDataValue().SetBoolean(false).ToString(), t) { return }
}

func TestThat_DataValue_ToString_Returns_object_properties(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("shape", NewDataValue().SetString("circle")).
		SetObjectProperty("size", NewDataValue().SetInteger(50)).
		SetObjectProperty("outline", NewDataValue().SetBoolean(true)).
		SetObjectProperty("pi", NewDataValue().SetFloat(3.14159))

	// Test
	actual := sut.ToString()

	// Verify
	// Note: because the underlying structure is a map, we cannot guarantee the order of output.
	if ! ExpectTrue(0 <= strings.Index(actual, "\"shape\":\"circle\""), t) { return }
	if ! ExpectTrue(0 <= strings.Index(actual, "\"size\":50"), t) { return }
	if ! ExpectTrue(0 <= strings.Index(actual, "\"outline\":true"), t) { return }
	if ! ExpectTrue(0 <= strings.Index(actual, "\"pi\":3.14159"), t) { return }
}

func TestThat_DataValue_ToString_Returns_array_elements(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray().
		AppendArrayValue(NewDataValue().SetInteger(13)).
		AppendArrayValue(NewDataValue().SetString("bananas"))

	// Test
	actual := sut.ToString()

	// Verify
	if ! ExpectString("[13,\"bananas\"]", actual, t) { return }
}

func TestThat_DataValue_ToString_Returns_array_in_object(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject().
		SetObjectProperty("prop",
			NewDataValue().PrepareArray().
				AppendArrayValue(NewDataValue().SetString("elem")),
		)

	// Test
	actual := sut.ToString()

	// Verify
	if ! ExpectString("{\"prop\":[\"elem\"]}", actual, t) { return }
}

func TestThat_DataValue_ToString_Returns_object_in_array(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray().
		AppendArrayValue(
			NewDataValue().PrepareObject().
				SetObjectProperty("prop", NewDataValue().SetString("val")),
		)

	// Test
	actual := sut.ToString()

	// Verify
	if ! ExpectString("[{\"prop\":\"val\"}]", actual, t) { return }
}

// ToJson

// The only difference in string encodings for Json is that we containstrings with quotes
func TestThat_DataValue_ToJson_Quotes_Strings(t *testing.T) {
	// Verify
	if ! ExpectString("\"apple\"", NewDataValue().SetString("apple").ToJson(), t) { return }
}

// Clone

func TestThat_DataValue_Clone_Returns_deep_copy(t *testing.T) {
	// Setup
	var expectedInt int64 = 13
	var expectedFloat float64 = 3.14159
	expectedString := "bananas"
	sut := NewDataValue().PrepareArray().
		AppendArrayValue(NewInteger(expectedInt)).
		AppendArrayValue(NewString(expectedString)).
		AppendArrayValue(NewObject().
			SetObjectProperty("float", NewFloat(expectedFloat)).
			SetObjectProperty("boolean", NewBoolean(true)).
			SetObjectProperty("null", NewNull()),
	)

	// Test
	actual := sut.Clone()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
	if ! ExpectInt(3, actual.GetArraySize(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(0).IsInteger(), t) { return }
	if ! ExpectInt64(expectedInt, actual.GetArrayValue(0).GetInteger(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(1).IsString(), t) { return }
	if ! ExpectString(expectedString, actual.GetArrayValue(1).GetString(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).IsObject(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).HasObjectProperty("float"), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).GetObjectProperty("float").IsFloat(), t) { return }
	if ! ExpectFloat64(expectedFloat, actual.GetArrayValue(2).GetObjectProperty("float").GetFloat(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).HasObjectProperty("boolean"), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).GetObjectProperty("boolean").GetBoolean(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).GetObjectProperty("boolean").IsBoolean(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).HasObjectProperty("null"), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(2).GetObjectProperty("null").IsNull(), t) { return }
}

