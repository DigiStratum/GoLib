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

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Object_NewObject_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewObject()

	// Verify
	ExpectNonNil(sut, t)
}

