package objects

/*

An Object represents a binary block of data, typically what one might consider a "file" on a disk,
which can be managed via ObjectStore. By abstracting Objects as the data set instead of as a named
file on disk, we can capture them in any number of places: files on disk, records in a database,
representations in an API, even codified chunks of data within our own executable.

Object optionally support fields; if the fields map is empty, then they are not being used.

TODO:
 * Refactor Object.fields to use a new of.ObjectFieldMap which extends lib.HashMap with support for of.ObjectFieldTypes
 * Cache the ToString() result so that successive calls return the same value until it is cleared by
   some change made by another method
*/

import (
	"fmt"
	"regexp"
	"encoding/json"

	"github.com/DigiStratum/GoLib/FileIO"
	of "github.com/DigiStratum/GoLib/Object/field"
	xc "github.com/DigiStratum/GoLib/Data/transcoder"
	//enc "github.com/DigiStratum/GoLib/Data/transcoder/encodingscheme"
)

type ObjectIfc interface {
	// Import
	FromString(serialized *string) error
	FromBytes(serializedBytes *[]byte) error
	FromFile(path string) error

	// Export
	ToString() (*string, error)
	ToBytes() (*[]byte, error)
	ToFile(path string) error

	// Fields
	AddField(fieldName string, value *string, ofType of.OFType) error
	SetFieldValue(fieldName string, value *string) error
	HasField(fieldName string) bool
	GetFieldType(fieldName string) *of.ObjectFieldType
}

// An Object that we're going to codify
type Object struct {
	transcoder		xc.TranscoderIfc
	fields			map[string]*of.ObjectField		// Field name to value map
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObject(transcoder xc.TranscoderIfc) *Object {
	return &Object{
		transcoder:		transcoder,
		fields:			make(map[string]*of.ObjectField),
	}
}

// -------------------------------------------------------------------------------------------------
// ObjectIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) FromString(serialized *string) error {
	if nil == serialized { return fmt.Errorf("Cannot deserialize nil string value") }

	// Overall format: "ser[{Method}:{Type}:{Data}]"
	re := regexp.MustCompile(`^ser\[(?P<method>\w+):(?P<etype>\w+):(?P<edata>\w+\]$`)
	if matched := re.MatchString(*serialized); !matched {
		return fmt.Errorf("String does not match serialization requirements")
	}

	matches := re.FindStringSubmatch(*serialized)
	if (nil == matches) || (len(matches) < 3) {
		return fmt.Errorf("Unexpected mismatch on parameters for serialized value")
	}

	method := matches[re.SubexpIndex("method")]
	etype := matches[re.SubexpIndex("etype")]
	edata := matches[re.SubexpIndex("edata")]

	// We support multiple methods of object deserialization
	switch method {
		case "j64":
			// Encoding base64, JSON data
			if (nil == r.transcoder) || ("base64" != r.transcoder.GetDecoderSchemeName()) {
				return fmt.Errorf("Deserialization requires Transcoder with EncodingSchemeBase64")
			}

			utype, err := r.transcoder.Decode(&etype)
			if (nil != err) || (nil == utype) || ("Object" != *utype) {
				return fmt.Errorf("Error decoding serialized data type")
			}

			udata, err := r.transcoder.Decode(&edata)
			if (nil != err) || (nil == udata) { return fmt.Errorf("Error decoding serialized data") }

			return r.FromJson(udata)
	}

	return fmt.Errorf("Unsupported serialization method '%s'", method)
}

func (r *Object) FromBytes(serializedBytes *[]byte) error {
	str := string(*serializedBytes)
	return r.FromString(&str)
}

func (r *Object) FromFile(path string) error {
	serialized, err := fileio.ReadFileString(path)
	if nil != err { return err }
	return r.FromString(serialized)
}

func (r Object) ToString() (*string, error) {
	// Encoding base64, JSON data
	if (nil == r.transcoder) || ("base64" != r.transcoder.GetDecoderSchemeName()) {
		return nil, fmt.Errorf("Serialization requires Transcoder with EncodingSchemeBase64")
	}

	// Overall format: "ser[{Method}:{Type}:{Data}]"
	method := "j64"
	udata, err := r.ToJson()
	if nil != err { return  nil, err }

	edata, err := r.transcoder.Encode(udata)
	if (nil != err) || (nil == edata) {
		return nil, fmt.Errorf("Error encoding serialized data type")
	}

	utype := "Object"
	etype, err := r.transcoder.Encode(&utype)
	if (nil != err) || (nil == etype) {
		return nil, fmt.Errorf("Error encoding serialized data type")
	}

	serialized := fmt.Sprintf("ser[%s:%s:%s]", method, etype, edata)
	return &serialized, nil
}

func (r Object) ToBytes() (*[]byte, error) {
	serialized, err := r.ToString()
	if nil != err { return nil, err }
	serializedBytes := []byte(*serialized)
	return &serializedBytes, nil
}

func (r Object) ToFile(path string) error {
	serializedBytes, err := r.ToBytes()
	if nil != err { return err }
	return fileio.WriteFileBytes(path, serializedBytes)
}

func (r *Object) AddField(fieldName string, value *string, ofType of.OFType) error {

	// of.ObjectField Map needs a field with this name in place; create it if it's missing
	if ! r.HasField(fieldName) {
		objectField := of.NewObjectField()
		objectField.Type = of.NewObjectFieldTypeFromOFType(ofType)
		r.fields[fieldName] = objectField
	}

	return r.SetFieldValue(fieldName, value)
}

// Set the named field to the specified value (including nil!)
func (r *Object) SetFieldValue(fieldName string, value *string) error {

	// of.ObjectField Map needs a field with this name in place
	if objectField, ok := r.fields[fieldName]; ok {

		// Validate the new value against the field's type
		if ! objectField.Type.IsValid(value) {
			return fmt.Errorf(
				"Object Field [name: '%s', type: '%s'] does not match supplied value",
				fieldName, objectField.Type.ToString(),
			)
		}

		// Set the value, yey!
		objectField.Value = value
		r.fields[fieldName] = objectField
	} else {
		return fmt.Errorf("Object has no field named '%s'", fieldName)
	}
	return nil
}

func (r Object) HasField(fieldName string) bool {
	_, ok := r.fields[fieldName]
	return ok
}

// Return the type of the named field
func (r Object) GetFieldType(fieldName string) *of.ObjectFieldType {
	if ! r.HasField(fieldName) { return nil }
	return r.fields[fieldName].Type
}

// -------------------------------------------------------------------------------------------------
// JsonSerializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(r.fields)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes)
	return &jsonString, nil
}

/*
// FIXME: If fields are in use, the json should be of the field map, otherwise, it should be of a non-fielded content string
// We can probably just do a simple, hacky inspection of the JSON string to figure out if we are using fields by checking for '{'
// as the first character which indicates that an object (JSON) follows rather than a plain string which would start with '"'

func (r Object) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(r.fields)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}
*/

// -------------------------------------------------------------------------------------------------
// JsonDeserializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) FromJson(jsonString *string) error {
	r.fields = make(map[string]*of.ObjectField)
	return json.Unmarshal([]byte(*jsonString), &r.fields)
}

