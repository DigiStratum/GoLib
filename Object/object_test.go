package object

import(
	"testing"

	objf "github.com/DigiStratum/GoLib/Object/field"
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
	expectedFieldName := "valid-object-field"
	expectedValue := "bogus field value"
	expectedFieldType := objf.OFT_NUMERIC
	oldcontent := "testcontent"
	sut.SetContent(&oldcontent)

	newOF := objf.NewObjectField(expectedFieldName)
	newOF.SetType( objf.NewObjectFieldTypeFromOFType(expectedFieldType))

	// Test
	sut.AddField(newOF, &expectedValue)

	// Verify
	ExpectNil(sut.GetContent(), t)
}

func TestThat_Object_AddField_AddsField_WithoutError_ForUnknownFieldType(t *testing.T) {
	// Setup
	sut := NewObject()
	expectedFieldName := "valid-object-field"
	expectedValue := "bogus field value"
	expectedFieldType := objf.OFT_UNKNOWN

	newOF := objf.NewObjectField(expectedFieldName)
	newOF.SetType( objf.NewObjectFieldTypeFromOFType(expectedFieldType))

	// Test
	sut.AddField(newOF, &expectedValue)

	// Verify
	ExpectTrue(sut.HasField(expectedFieldName), t)
}

func TestThat_Object_HasField_ReturnsFalse(t *testing.T) {
	// Setup
	sut := NewObject()

	// Test
	actual := sut.HasField("bogusfield")

	// Verify
	ExpectFalse(actual, t)
}

func TestThat_Object_SetFieldValue_SetsFieldValue_WithoutError_ForGoodValue(t *testing.T) {
	// Setup
	sut := NewObject()
	transcoder := xc.NewTranscoder()
	transcoder.SetEncoderScheme(enc.NewEncodingSchemeBase64())
	sut.SetTranscoder(transcoder)
	expectedFieldName := "valid-object-field"
	originalValue := "222"
	expectedFieldType := objf.OFT_NUMERIC
	newValue := "333"
	// Serialization of: {"valid-object-field":{"Type":"numeric","Value":"333"}}
	expected := "ser[j64:T2JqZWN0:eyJ2YWxpZC1vYmplY3QtZmllbGQiOnsiVHlwZSI6Im51bWVyaWMiLCJWYWx1ZSI6IjMzMyJ9fQ==]"

	newOF := objf.NewObjectField(expectedFieldName)
	newOF.SetType( objf.NewObjectFieldTypeFromOFType(expectedFieldType))

	// Test
	sut.AddField(newOF, &originalValue)
	err2 := sut.SetFieldValue(expectedFieldName, &newValue)
	actual, err3 := sut.Serialize()
	hasIt := sut.HasField(expectedFieldName)

	// Verify
	ExpectNoError(err2, t)
	ExpectNoError(err3, t)
	ExpectTrue(hasIt, t)
	ExpectString(expected, *actual, t)
}

func TestThat_Object_SetFieldValue_ReturnsError_ForBadField(t *testing.T) {
	// Setup
	sut := NewObject()

	// Test
	err := sut.SetFieldValue("bogusfield", nil)

	// Verify
	ExpectError(err, t)
}

func TestThat_Object_GetField_ReturnsError_ForBadField(t *testing.T) {
	// Setup
	sut := NewObject()

	// Test
	actual, err := sut.GetField("bogusfield")

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

func TestThat_Object_GetField_ReturnsField_ForGoodField(t *testing.T) {
	// Setup
	sut := NewObject()
	expectedFieldName := "valid-object-field"
	originalValue := "222"
	expectedFieldType := objf.OFT_NUMERIC

	newOF := objf.NewObjectField(expectedFieldName)
	newOF.SetType( objf.NewObjectFieldTypeFromOFType(expectedFieldType))

	// Test
	sut.AddField(newOF, &originalValue)
	actual, err2 := sut.GetField(expectedFieldName)
	actualValue := *actual.GetValue()

	// Verify
	ExpectNoError(err2, t)
	ExpectNonNil(actual, t)
	ExpectString(originalValue, actualValue, t)
}

func TestThat_Object_GetFieldType_ReturnsNil_ForBadField(t *testing.T) {
	// Setup
	sut := NewObject()

	// Test
	actual := sut.GetFieldType("bogusfield")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Object_GetFieldType_ReturnsFieldType_ForGoodField(t *testing.T) {
	// Setup
	sut := NewObject()
	expectedFieldName := "valid-object-field"
	originalValue := "222"
	expectedFieldType := objf.OFT_NUMERIC

	newOF := objf.NewObjectField(expectedFieldName)
	newOF.SetType( objf.NewObjectFieldTypeFromOFType(expectedFieldType))

	// Test
	sut.AddField(newOF, &originalValue)
	actual := sut.GetFieldType(expectedFieldName)
	actualFieldType := actual.GetType()

	// Verify
	ExpectNonNil(actual, t)
	ExpectTrue(expectedFieldType == actualFieldType, t)
}

