package data

import(
	"fmt"
	"testing"

	. "GoLib/Testing"
)

// Interface

func TestThat_DataValue_NewDataValue_ReturnsInstance(t *testing.T) {
	// Setup
	var sut DataValueIfc = NewDataValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(
		sut.IsNull() || sut.IsString() || sut.IsObject() || sut.IsArray() || sut.IsInteger() || sut.IsFloat() || sut.IsValid(),
		t,
	) { return }
}

// Validity

func TestThat_DataValue_IsValid_Returns_false_for_new_value(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(sut.IsValid(), t) { return }
}

func TestThat_DataValue_GetType_returns_expected_type(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Test
	if ! ExpectTrue(VALUE_TYPE_INVALID == sut.GetType(), t) { return }
	if ! ExpectTrue(VALUE_TYPE_INTEGER == sut.SetInteger(0).GetType(), t) { return }
	if ! ExpectTrue(VALUE_TYPE_STRING == sut.SetString("howdy!").GetType(), t) { return }
	if ! ExpectTrue(VALUE_TYPE_NULL == sut.SetNull().GetType(), t) { return }
	if ! ExpectTrue(VALUE_TYPE_FLOAT == sut.SetFloat(3.14159).GetType(), t) { return }
	if ! ExpectTrue(VALUE_TYPE_ARRAY == sut.PrepareArray().GetType(), t) { return }
	if ! ExpectTrue(VALUE_TYPE_OBJECT == sut.PrepareObject().GetType(), t) { return }

}

// Nulls

func TestThat_DataValue_IsNull_Returns_true_after_setting_null(t *testing.T) {
	// Test
	sut := NewDataValue().SetNull()

	// Verify
	if ! ExpectTrue(sut.IsNull(), t) { return }
}

// Strings

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

// Objects

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
	err := sut.SetObjectProperty("name", NewDataValue())

	// Verify
	if ! ExpectError(err, t) { return }
}

func TestThat_DataValue_SetObjectProperty_Returns_successfully_sets_value(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()
	expectedName := "name"
	expectedValue := "value"
	value := NewDataValue()
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

func TestThat_DataValue_HasObjectProperty_Returns_false_for_non_property(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()

	// Verify
	if ! ExpectFalse(sut.HasObjectProperty("Nope!"), t) { return }
}

func TestThat_DataValue_GetObjectPropertyNames_Returns_Empty_set(t *testing.T) {
	// Test
	actual := NewDataValue().PrepareObject().GetObjectPropertyNames()

	// Verify
	if ! ExpectInt(0, len(actual), t) { return }
}

func TestThat_DataValue_GetObjectPropertyNames_Returns_Expected_Names(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareObject()
	expectedName := "name"
	value := NewDataValue().SetString("value")
	sut.SetObjectProperty(expectedName, value)

	// Test
	actual := sut.GetObjectPropertyNames()

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

// Booleans

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

func TestThat_DataValue_IsArray_Returns_true(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray()

	// Verify
	if ! ExpectTrue(sut.IsArray(), t) { return }
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

func TestThat_DataValue_GetArrayValue_Returns_nil_for_non_array_or_mising_index(t *testing.T) {
	// Setup
	sut := NewDataValue()

	// Verify
	if ! ExpectNil(sut.GetArrayValue(0), t) { return }
	sut.PrepareArray()
	if ! ExpectNil(sut.GetArrayValue(0), t) { return }
}

func TestThat_DataValue_GetArrayValue_Returns_value_for_good_index(t *testing.T) {
	// Setup
	sut := NewDataValue().PrepareArray()
	expectedValue := "value"
	sut.AppendArrayValue(NewDataValue().SetString(expectedValue))

	// Test
	actual := sut.GetArrayValue(0)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expectedValue, actual.GetString(), t) { return }
}

// Floats

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

// Integers

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

// Conveniences

// FIXME: Don't use Json Lexer - we are a base class, can't depend on a subclass!
/*
func TestThat_DataValue_Select_Returns_Errors(t *testing.T) {
	// Setup
	json := `
{
	"shape": "arc",
	"vectors": [
		{
			"radians": 3.14159,
			"radius": 2,
			"color": "red",
			"hidden": false
		},
		{
			"radians": 6.28318,
			"radius": 7,
			"color": "blue",
			"hidden": false
		}
	]
}
`
	lex := NewJsonLexer()
	sut, err := lex.LexJsonValue(json)
	if ! ExpectNoError(err, t) { return }
	if ! ExpectNonNil(sut, t) { return }

	// Test
	selectors := []string{ "[]", "[0]", "bogusproperty", "vectors[2]", "vectors[0]radians" }
	for _, selector := range selectors {
		actual1, err1 := sut.Select(selector)
		// Verify
		if ! ExpectNil(actual1, t) { return }
		if ! ExpectError(err1, t) { return }
	}
}

func TestThat_DataValue_Select_Returns_Values(t *testing.T) {
	// Setup
	json := `
{
	"shape": "arc",
	"vectors": [
		{
			"radians": 3.14159,
			"radius": 2,
			"color": "red",
			"hidden": false
		},
		{
			"radians": 6.28318,
			"radius": 7,
			"color": "blue",
			"hidden": false
		}
	]
}
`
	lex := NewJsonLexer()
	sut, err := lex.LexJsonValue(json)
	if ! ExpectNoError(err, t) { return }
	if ! ExpectNonNil(sut, t) { return }

	// Test
	selectors := []string{ ".shape", ".vectors", ".vectors[1]", ".vectors[1].radians" }
	for _, selector := range selectors {
		//t.Logf("Testing good selector '%s'\n", selector)
		// Verify
		actual1, err1 := sut.Select(selector)
		if ! ExpectNonNil(actual1, t) { return }
		if ! ExpectNoError(err1, t) { return }

		switch actual1.GetType() {
			case VALUE_TYPE_STRING: if ! ExpectString("arc", actual1.GetString(), t) { return }
			case VALUE_TYPE_ARRAY:
				if ! ExpectTrue(actual1.IsArray(), t) { return }
				if ! ExpectInt(2, actual1.GetArraySize(), t) { return }
			case VALUE_TYPE_OBJECT:
				if ! ExpectTrue(actual1.IsObject(), t) { return }
				if ! ExpectTrue(actual1.HasObjectProperty("radians"), t) { return }
			case VALUE_TYPE_FLOAT: if ! ExpectFloat64(float64(6.28318), actual1.GetFloat(), t) { return }
		}
	}
}
*/

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
		vi := actual()
		v, ok := vi.(*DataValue)
		// Expect a KeyValuePair to be the result of calling the Iterator func
		if ! ExpectTrue(ok, t) { return }
		if ! ExpectNonNil(v, t) { return }
		if ! ExpectInt64(int64(i), v.GetInteger(),t) { return }
	}
}

