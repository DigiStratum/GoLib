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
	"encoding/json"

	"github.com/DigiStratum/GoLib/Data/serializable"
	of "github.com/DigiStratum/GoLib/Object/field"
	xc "github.com/DigiStratum/GoLib/Data/transcoder"
)

type ObjectIfc interface {
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

// -------------------------------------------------------------------------------------------------
// SerializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) Serialize() (*string, error) {
	data, err := r.ToJson()
	if nil != err { return nil, err }
	serializer := serializable.NewSerializer(r.transcoder)
	return serializer.Serialize(data, "Object")
}

// -------------------------------------------------------------------------------------------------
// DeserializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) Deserialize(data *string) error {
	serializer := serializable.NewSerializer(r.transcoder)
	udata, err := serializer.Deserialize(data)
	if nil != err { return err }
	return r.FromJson(udata)
}

