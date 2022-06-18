package objects

/*
type ObjectIfc interface {
        // Import
        FromString(content *string, encodingScheme EncodingScheme) error
        FromBytes(bytes *[]byte, encodingScheme EncodingScheme) error
        FromFile(path string, encodingScheme EncodingScheme) error

        // Export
        ToString(encodingScheme EncodingScheme) (*string, error)
        ToBytes(encodingScheme EncodingScheme) (*[]byte, error)
        ToFile(path string, encodingScheme EncodingScheme) error
        ToJson() (*string, error)

        // Fields
        AddField(fieldName string, value *string, ofType OFType) error
        SetFieldValue(fieldName string, value *string) error
        HasField(fieldName string)
        GetFieldType(fieldName string) *ObjectFieldType
}

*/

import(
	//"fmt"
	"testing"

	of "github.com/DigiStratum/GoLib/Object/field"
	xc "github.com/DigiStratum/GoLib/Data/transcoder"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Object_NewObject_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewObject()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Object_AddField_AddsField_WithoutError_ForGoodFieldType(t *testing.T) {
	// Setup
	sut := NewObject()
	expectedFieldName := "bogus-object-field"
	expectedValue := "bogus field value"
	expectedFieldType := of.OFT_NUMERIC

	// Test
	err := sut.AddField(expectedFieldName, &expectedValue, expectedFieldType)

	// Verify
	ExpectNoError(err, t)
}

func TestThat_Object_AddField_AddsField_WithoutError_ForUnknownFieldType(t *testing.T) {
	// Setup
	sut := NewObject()
	expectedFieldName := "bogus-object-field"
	expectedValue := "bogus field value"
	expectedFieldType := of.OFT_UNKNOWN

	// Test
	err := sut.AddField(expectedFieldName, &expectedValue, expectedFieldType)

	// Verify
	ExpectError(err, t)
}

//SetFieldValue(fieldName string, value *string) error
//ToString(encodingScheme xcode.EncodingScheme) (*string, error)
func TestThat_Object_SetFieldValue_SetsFieldValue_WithoutError_ForGoodValue(t *testing.T) {
	// Setup
	sut := NewObject()
	expectedFieldName := "bogus-object-field"
	originalValue := "222"
	expectedValue := "333"
	expectedFieldType := of.OFT_NUMERIC

	// Test
	err1 := sut.AddField(expectedFieldName, &originalValue, expectedFieldType)
	err2 := sut.SetFieldValue(expectedFieldName, &expectedValue)
	actualValue, err3 := sut.ToString(xc.ES_NONE)

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectNoError(err3, t)
	ExpectString(expectedValue, *actualValue, t)
}

