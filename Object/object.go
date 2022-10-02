package object

/*

An Object represents a binary block of data, typically what one might consider a "file" on a disk,
which can be managed via ObjectStore. By abstracting Objects as the data set instead of as a named
file on disk, we can capture them in any number of places: files on disk, records in a database,
representations in an API, even codified chunks of data within our own executable.


Object optionally support fields; if the fields map is empty, then they are not being used. In lieu
of a populated field map, the content string is used. If it is nil then the entire Object's value
is nil. Setting a content value purges the field map. Setting a field map nils the content.

Objects are defined with fields by using AddField() to add each one. The defined fields may have a
new value assigned at any time. The fields/values may also be defined by way of deserialization.
But once the fields are defined for an object, the field definitions themselves are immutable: only
the values ay be updated. In the case that the consumer wants to modify the field definitions of
stored data for purposes of some sort of data transformation and/or data structure migration, it
must first create a new Object and AddField() for all the new fields, and read the fields out of
the old Object into the new one. We could potentially add a helper method with some sort of
transform function to facilitate this should the need arise.

TODO:
 * Refactor Object.fields to use a new objf.ObjectFieldMap which extends lib.HashMap with support for
   objf.ObjectFieldTypes
 * Cache Serialize() result so that successive calls return the same value until it is cleared by
   some change made by another method
 * Add support for arbitrary field values, not just string (see SetFieldValue(); see also go 1.18
   which supports generics)
 * Add support for Clone() method (does some new "clonable" interface make sense? Can we use 1.18
   generics for this?) 
*/

import (
	"fmt"
	"encoding/json"

	"github.com/DigiStratum/GoLib/Data/serializable"
	objf "github.com/DigiStratum/GoLib/Object/field"
	xc "github.com/DigiStratum/GoLib/Data/transcoder"
)

type ObjectIfc interface {
	SetTranscoder(transcoder xc.TranscoderIfc)
	SetContent(content *string)
	GetContent() *string
	AddField(objectField objf.ObjectFieldIfc, value *string)
	HasField(fieldName string) bool
	GetFieldType(fieldName string) *objf.ObjectFieldType
	SetFieldValue(fieldName string, value *string) error
	GetField(fieldName string) (*objf.ObjectField, error)
}

// An Object that we're going to codify
type Object struct {
	transcoder		xc.TranscoderIfc		// Transcoder for De|Serialization
	fields			map[string]*objf.ObjectField	// Field name to value map
	content			*string				// Non-Field-Mapped Object Content
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObject() *Object {
	return &Object{
		fields:			make(map[string]*objf.ObjectField),
	}
}

// -------------------------------------------------------------------------------------------------
// ObjectIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) SetTranscoder(transcoder xc.TranscoderIfc) {
	r.transcoder = transcoder
}

func (r *Object) SetContent(content *string) {
	// Purge field map on setting content
	r.fields = make(map[string]*objf.ObjectField)
	r.content = content
}

func (r Object) GetContent() *string {
	return r.content
}

func (r *Object) DefineField(objectField objf.ObjectFieldIfc) {
	newOF := objf.NewObjectField(objectField.GetName())
	newOF.SetType(objectField.GetType())
	r.fields[objectField.GetName()] = newOF
	// Invalidate string content when we touch the fieldmap
	r.content = nil
}

func (r *Object) AddField(objectField objf.ObjectFieldIfc, value *string) {
	r.DefineField(objectField)
	r.SetFieldValue(objectField.GetName(), value)
}

func (r Object) HasField(fieldName string) bool {
	_, ok := r.fields[fieldName]
	return ok
}

func (r Object) GetFieldType(fieldName string) *objf.ObjectFieldType {
	if ! r.HasField(fieldName) { return nil }
	return r.fields[fieldName].GetType()
}

// Set the named field to the specified value (including nil!)
func (r *Object) SetFieldValue(fieldName string, value *string) error {
	if objectField, ok := r.fields[fieldName]; ok {
		if ! objectField.GetType().IsValid(value) {
			return fmt.Errorf(
				"Object Field [name: '%s', type: '%s'] does not match supplied value",
				fieldName, objectField.GetType().ToString(),
			)
		}
		objectField.SetValue(value)
		return nil
	}
	return fmt.Errorf("Object has no field named '%s'", fieldName)
}

func (r *Object) GetField(fieldName string) (*objf.ObjectField, error) {
	if objectField, ok := r.fields[fieldName]; ok { return objectField, nil }
	return nil, fmt.Errorf("Object has no field named '%s'", fieldName)
}

// -------------------------------------------------------------------------------------------------
// JsonSerializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// JSON String result will be nil for error, start with doublequote
// for string value, and open curly brace for fieldmapped value
func (r *Object) ToJson() (*string, error) {
	var jsonString string
	if len(r.fields) > 0 {
		jsonBytes, err := json.Marshal(r.fields)
		if nil != err { return nil, err }
		jsonString = string(jsonBytes)
	} else if nil != r.content {
		jsonBytes, err := json.Marshal(r.content)
		if nil != err { return nil, err }
		jsonString = string(jsonBytes)
	} else {
		jsonString = ""
	}
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// JsonDeserializableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Object) FromJson(jsonString *string) error {
	// Purge field map and content
	r.fields = make(map[string]*objf.ObjectField)
	r.content = nil
	// nil
	if nil == jsonString { return nil }
	// content string
	if "\"" == string([]rune(*jsonString)[0]) { return json.Unmarshal([]byte(*jsonString), &r.content) }
	// field mapped object
	if "{" == string([]rune(*jsonString)[0]) { return json.Unmarshal([]byte(*jsonString), &r.fields) }
	return fmt.Errorf("Unsupported JSON value")
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

