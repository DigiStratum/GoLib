package objectfield

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_ObjectFieldType_NewObjectFieldType_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewObjectFieldType()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_ObjectFieldType_NewObjectFieldTypeFromString_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewObjectFieldTypeFromString("numeric")

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_ObjectFieldType_NewObjectFieldTypeFromOFType_RetrunsSomething(t *testing.T) {
	// Setup
	expected := OFT_NUMERIC

	// Test
	sut := NewObjectFieldTypeFromOFType(expected)
	actualt := sut.GetType()
	actuals := sut.ToString()

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(actualt == expected, t)
	ExpectString("numeric", actuals, t)
}

func TestThat_ObjectFieldType_IsValue_ReturnsFalse_ForUnknownType(t *testing.T) {
	// Setup
	sut := NewObjectFieldTypeFromString("unknown")
	value := "bogus"

	// Test
	actual := sut.IsValid(&value)

	// Verify
	ExpectFalse(actual, t)
}

// TODO: Fix this up once we fix up the IsValid function itself to properly check the value by type
func TestThat_ObjectFieldType_IsValue_ReturnsTrue_ForGoodType(t *testing.T) {
	// Setup
	sut := NewObjectFieldTypeFromString("string")
	value := "bogus"

	// Test
	actual := sut.IsValid(&value)

	// Verify
	ExpectTrue(actual, t)
}

