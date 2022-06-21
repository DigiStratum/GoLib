package objects

/*
type ObjectIfc interface {
	GetFieldType(fieldName string) *of.ObjectFieldType
	SetFieldValue(fieldName string, value *string) error
	GetField(fieldName string) (*of.ObjectField, error)
}
*/

import(
	//"fmt"
	"testing"

	of "github.com/DigiStratum/GoLib/Object/field"
	xc "github.com/DigiStratum/GoLib/Data/transcoder"
	enc "github.com/DigiStratum/GoLib/Data/transcoder/encodingscheme"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Object_NewObject_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewObject()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Object_SetContent_SetsTheContent(t *testing.T) {
	// Setup
	sut := NewObject()
	expected := "testcontent"

	// Test
	sut.SetContent(&expected)
	actual := sut.GetContent()

	// Verify
	ExpectString(expected, *actual, t)
}

func TestThat_Object_GetContent_ReturnsNil_WhenContentIsUnset(t *testing.T) {
	// Setup
	sut := NewObject()

	// Test
	actual := sut.GetContent()

	// Verify
	ExpectNil(actual, t)
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

//HasField(fieldName string) bool
func TestThat_Object_HasField_ReturnsFalse(t *testing.T) {
	// Setup
	sut := NewObject()

	// Test
	actual := sut.HasField("bogusfield")

	// Verify
	ExpectFalse(actual, t)
}


//SetFieldValue(fieldName string, value *string) error
func TestThat_Object_SetFieldValue_SetsFieldValue_WithoutError_ForGoodValue(t *testing.T) {
	// Setup
	sut := NewObject()
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	sut.SetTranscoder(transcoder)
	expectedFieldName := "bogus-object-field"
	originalValue := "222"
	newValue := "333"
	expected := "ser[j64:T2JqZWN0:eyJib2d1cy1vYmplY3QtZmllbGQiOnsiVHlwZSI6e30sIlZhbHVlIjoiMzMzIn19]"
	expectedFieldType := of.OFT_NUMERIC

	// Test
	err1 := sut.AddField(expectedFieldName, &originalValue, expectedFieldType)
	err2 := sut.SetFieldValue(expectedFieldName, &newValue)
	actual, err3 := sut.Serialize()
	hasIt := sut.HasField(expectedFieldName)

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectNoError(err3, t)
	ExpectTrue(hasIt, t)
	ExpectString(expected, *actual, t)
}
