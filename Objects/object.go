package objects

/*

An Object represents a binary block of data, typically what one might consider a "file" on a disk,
which can be managed via ObjectStore. By abstracting Objects as the data set instead of as a named
file on disk, we can capture them in any number of places: files on disk, records in a database,
representations in an API, even codified chunks of data within our own executable.

Object optionally support fields; if the fields map is empty, then they are not being used.

TODO:
 * Refactor Object.fields to use a new ObjectFieldMap which extends lib.HashMap with support for ObjectFieldTypes

*/

import (
	"fmt"
	"encoding/json"

	xcode "github.com/DigiStratum/GoLib/transcoder"
)

type ObjectIfc interface {
	SetContentFromString(content *string)
	SetEncodedContentFromString(encodedContent *string)
	SetContentFromFile(path string) error
	GetContent() *string
	GetEncodedContent() *string
	GetFieldType(fieldName string) string
}

// A static Object that we're going to codify
type Object struct {
	contentTranscoder	xcode.Transcoder
	fields			map[string]ObjectField		// Field name to value map
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObject() *Object {
	return &Object{
		contentTranscoder:	xcode.NewTranscoder(),
		fields:			make(map[string]ObjectField),
	}
}

// -------------------------------------------------------------------------------------------------
// ObjectIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) FromString(content *string, encodingScheme EncodingScheme) error {
	return r.contentTranscoder.FromString(content, encodingScheme)
}

func (r *Object) FromBytes(bytes *[]byte, encodingScheme EncodingScheme) error {
	return r.contentTranscoder.FromBytes(bytes, encodingScheme)
}

func (r *Object) FromFile(path string, encodingScheme EncodingScheme) error {
	return r.contentTranscoder.FromFile(path, encodingScheme)
}

func (r Object) ToString(encodingScheme EncodingScheme) (*string, error) {
	return r.contentTranscoder.ToString(encodingScheme)
}

func (r Object) ToBytes(encodingScheme EncodingScheme) (*[]byte, error) {
	return r.contentTranscoder.ToBytes(encodingScheme)
}

func (r Object) ToFile(path string, encodingScheme EncodingScheme) error {
	return r.contentTranscoder.ToFile(encodingScheme)
}

// If fields are in use, we can pop out JSON
// TODO: Should this be part of the Transcoder? It has no sense of fielded properties...
func (r Object) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(r.fields)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

func (r *Object) AddField(fieldName string, value *string, ofType OFType) error {

	// ObjectField Map needs a field with this name in place; create it if it's missing
	if ! r.HasField(fieldName) {
		objectField := NewObjectField()
		objectField.Type = NewObjectFieldType()
		objectField.Type.SetType(ofType)
		r.fields[fieldName] = objectField
	}

	return r.SetFieldValue(fieldName, value)
}

// Set the named field to the specified value (including nil!)
func (r *Object) SetFieldValue(fieldName string, value *string) error {

	// ObjectField Map needs a field with this name in place
	if objectField, ok := r.fields[fieldName]; ok {

		// Validate the new value against the field's type
		if ! objectField.Type.IsValid(value) {
			return fmt.Errorf(
				"Object Field [name: '%s', type: '%s'] does not match supplied value",
				name, objectField.Type.ToString(),
			)
		}

		// Set the value, yey!
		objectField.Value = value
		r.fields[fieldName] = objectField
	} else {
		return fmt.Errorf("Object has no field named '%s'", name)
	}
	return nil
}

func (r Object) HasField(fieldName string) {
	_, ok := r.fields[fieldName]
	return ok
}

// Return the type of the named field
func (r Object) GetFieldType(fieldName string) *ObjectFieldType {
	if ! r.HasField(fieldName) { return nil }
	return r.fields[fieldName].Type
}
