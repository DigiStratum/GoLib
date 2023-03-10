package mysql

import(
	"testing"
	"encoding/json"

	nulls "github.com/DigiStratum/GoLib/DB/MySQL/nullables"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewResultRow_ReturnsEmptyResultRow(t *testing.T) {
	// Test
	var sut *ResultRow = NewResultRow()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_ResultRow_Fields_ReturnsEmptySet_WhenNoFieldsSet(t *testing.T) {
	// Setup
	sut := NewResultRow()

	// Test
	actualFields := sut.Fields()

	// Verify
	ExpectInt(0, len(actualFields), t)
}

func TestThat_ResultRow_Fields_ReturnsFieldNames_WhenFieldsSet(t *testing.T) {
	// Setup
	sut := NewResultRow()
	expectedFields := []string{"one", "two", "three"}
	for _, field := range expectedFields {
		nullableValue := nulls.NewNullable(field)
		sut.Set(field, *nullableValue)
	}

	// Test
	actualFields := sut.Fields()

	// Verify
	ExpectNonNil(sut, t)
	ExpectInt(len(expectedFields), len(actualFields), t)
}

func TestThat_ResultRow_Get_ReturnsNothing_ForUnsetFields(t *testing.T) {
	// Setup
	sut := NewResultRow()

	// Test
	actualNullableValue := sut.Get("bogusfield")

	// Verify
	ExpectNil(actualNullableValue, t)
}

func TestThat_ResultRow_Fields_ReturnsValue_ForSetFields(t *testing.T) {
	// Setup
	sut := NewResultRow()
	expectedFields := []string{"one", "two", "three"}
	for _, field := range expectedFields {
		nullableValue := nulls.NewNullable(field)
		sut.Set(field, *nullableValue)
	}

	// Test / Verify
	for _, field := range expectedFields {
		actualNullableValue := sut.Get(field)
		ExpectNonNil(actualNullableValue, t)
		//ExpectTrue(actualNullableValue.IsString(), t)
		ExpectTrue(nulls.NULLABLE_STRING == actualNullableValue.GetType(), t)
		actualStringValue := actualNullableValue.GetString()
		ExpectNonNil(actualStringValue, t)
		ExpectString(field, *actualStringValue, t)
	}
}

func TestThat_ResultRow_ToJson_ReturnsEmptyObject_WhenNoFieldsSet(t *testing.T) {
	// Setup
	sut := NewResultRow()

	// Test
	actualJson, err := sut.ToJson()

	// Verify
	ExpectNonNil(actualJson, t)
	ExpectNoError(err, t)
	ExpectString("{}", *actualJson, t)
}

func TestThat_ResultRow_ToJson_ReturnsPopulatedObject_WhenFieldsSet(t *testing.T) {
	// Setup
	sut := NewResultRow()
	expectedFields := []string{"one", "two", "three"}
	for _, field := range expectedFields {
		nullableValue := nulls.NewNullable(field)
		sut.Set(field, *nullableValue)
	}

	// Test
	actualJson, err := sut.ToJson()

	// Verify
	ExpectString("", *actualJson, t)
	ExpectNonNil(actualJson, t)
	ExpectNoError(err, t)
	// We cannot expect JSON-serialized fields in a specific order, so we deserialize and check for expected results
	actualFields := make(map[string]string)
	err = json.Unmarshal([]byte(*actualJson), &actualFields)
	ExpectNoError(err, t)
	for _, field := range expectedFields {
		actualValue, ok := actualFields[field]
		ExpectTrue(ok, t)
		ExpectString(field, actualValue, t)
	}
}
