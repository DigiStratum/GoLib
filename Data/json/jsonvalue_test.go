package json

import(
	"fmt"
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
	sut := NewJsonValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(sut.IsValid(), t) { return }
}

func TestThat_JsonValue_GetType_returns_expected_type(t *testing.T) {
	// Setup
	sut := NewJsonValue()

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

// Conveniences

func TestThat_JsonValue_Select_Returns_Errors(t *testing.T) {
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

func TestThat_JsonValue_Select_Returns_Values(t *testing.T) {
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


func TestThat_JsonValue_GetIterator_Returns_nil(t *testing.T) {
	// Setup
	sut := NewJsonValue()

	// Test
	actual := sut.GetIterator()

	// Verify
	// Note: Can't use ExpectNil() here because there's something about the nil zero-value of
	// actual once passed as an argument to the Expect*() func that prevents it from being seen
	// as truly nil any longer.
	if ! ExpectTrue(nil == actual, t) { return }
}

func TestThat_JsonValue_GetIterator_Returns_func(t *testing.T) {
	// Setup
	sut := NewJsonValue()
	sut.PrepareObject()

	// Test
	actual := sut.GetIterator()

	// Verify
	// Note: Can't use ExpectNonNil() here because there's something about the nil zero-value of
	// actual once passed as an argument to the Expect*() func that prevents it from being seen
	// as truly nil any longer.
	if ! ExpectFalse(nil == actual, t) { return }
}

func TestThat_JsonValue_GetIterator_Returns_Good_Iterator(t *testing.T) {
	// Setup
	sut := NewJsonValue()
	sut.PrepareObject()
	numprops := 3
	for i := 1; i <= numprops; i++ {
		sut.SetObjectProperty(fmt.Sprintf("p%d", i), NewJsonValue().SetInteger(int64(i)))
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

// TODO: Test Array iteration as well!

