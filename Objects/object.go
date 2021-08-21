package objects

/*

An Object represents a binary block of data, typically what one might consider a "file" on a disk,
which can be managed via ObjectStore. By abstracting Objects as the data set instead of as a named
file on disk, we can capture them in any number of places: files on disk, records in a database,
representations in an API, even codified chunks of data within our own executable.

Object optionally support fields; if the field map is nil, then they are not being used.

TODO: Isolate the encode/decode so that other tools can build against it and have a function that
      properly interacts with the same encoding scheme as us using our *Encoded* accessor methods.

*/

import (
	"fmt"
	"errors"
	"encoding/base64"
	lib "github.com/DigiStratum/GoLib"
)

type ObjectIfc interface {
	SetContentFromString(content *string)
	SetEncodedContentFromString(encodedContent *string)
	SetContentFromFile(path string) error
	GetContent() *string
	GetEncodedContent() *string
	GetFieldType(fieldName string) string
}

type objectEncodingScheme int

const (
	OES_UNKNOWN objectEncodingScheme = iota
	OES_BASE64					// Base64 Encoding
)

// A static Object that we're going to codify
type Object struct {
	isEncoded	bool				// Is the content encoded?
	encodingScheme	objectEncodingScheme		// What method of encoding is used?
	content		*string				// Non-fielded Object "BLOB" representation
	fields		map[string]ObjectField		// Field name to value map
}

// Factory Functions
func NewObject() *Object {
	return &Object{}
}

// Make a new one of these from a simple string
func NewObjectFromString(content string) *Object {
	object := NewObject()
	object.SetContentFromString(&content)
	return object
}

// Make a new one of these from a file on disk
func NewObjectFromFile(path string) *Object {
	object := NewObject()
	object.SetContentFromFile(path)
	return object
}

// -------------------------------------------------------------------------------------------------
// ObjectIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Set the Object Content from a plain text string (it will be encoded!)
func (o *Object) SetContentFromString(content *string) {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(*content))
	o.content = &encodedContent
	o.isEncoded = true
}

// Set the Object Content from a text string which is already endcoded
// (This is used by callers such as res2go that know how to pre-encode)
func (o *Object) SetEncodedContentFromString(encodedContent *string) {
	o.content = encodedContent
	o.isEncoded = true
}

// Set the Object Content from a source file path (it will be encoded!)
// (This is used to anything froma simple text file to full binary assets)
func (o *Object) SetContentFromFile(path string) error {
	s, err := lib.ReadFileString(path)
	if nil != err { return err }
	o.SetContentFromString(s)
	return nil
}

// Get the Object Content as a string (could be anything!)
func (o *Object) GetContent() *string {
	// For non-encoded, raw content (probably loaded from disk, DB, service, etc)
	if ! o.isEncoded { return o.content }

	// For encoded content (probably compiled)
	decodedContentBytes, err := base64.StdEncoding.DecodeString(*o.content)
	if nil != err {
		// TODO: Handle errors
		s := ""
		return &s
	}
	decodedContent := string(decodedContentBytes)
	return &decodedContent
}

// Get the Object Content as an Encoded string (you better know what to do with it!)
func (r Object) GetEncodedContent() *string {
	copy := *(r.content) // Make an immutable copy of the content
	return &copy
}

// Reset the content of this object; preserves field map for fielded objects 
func (o *Object) ResetObject() {
	// For fielded objects...
	if nil != o.fields {
		// Reset all the field values to nil
		for name := range (*o.fields) {
			// Screwy golang workaround: can't index fields of structs in a map because
			// they freak out about memory management, where things live and meaning
			// ref: https://github.com/golang/go/issues/3117
			objectField := (*o.fields)[name]	// So: grab the objectField...
			objectField.Value = nil			// ... nil out the value...
			(*o.fields)[name] = objectField		// ... and jam it back into place
		}
	} else {
		// reset content for non-fielded objects
		o.content = nil
		o.isEncoded = false
	}
}

func (r *Object) SetField(fieldName string, objectField ObjectFieldIfc) {
	r.fields[fieldName] = ojectField
}

// Set the named field to the specified value (including nil!)
func (o *Object) SetFieldValue(name string, value *string) error {

	// Object needs an ObjectFieldMap in place
	if nil == o.fields {
		return errors.New("Object has no field map")
	}

	// ObjectFieldMap needs a field with this name in place
	if of, ok := (*o.fields)[name]; ok {

		// Validate the new value against the field's type
		if ! o.IsValueType(value, of.Type) {
			return errors.New(fmt.Sprintf(
				"Value does not match object field '%s' with type '%s'",
				name, o.GetObjectFieldTypeReadable(of.Type),
			))
		}

		// Set the value, yey!
		of.Value = value
		(*o.fields)[name] = of
	} else {
		return errors.New(fmt.Sprintf("Object has no field named '%s'", name))
	}
	return nil
}

// Determine whethere the value passes all the rules for the specified field type
func (o *Object) IsValueType(value *string, fieldType ObjectFieldType) bool {
	// TODO switch on type and run the value through the wringer here
	return true
}

func (r Object) HasField(fieldName string) {
	_, ok := r.fields[name]
	return ok
}

// Return a readable string for the named field
func (o *Object) GetFieldType(fieldName string) string {
	if ! r.HasField(fieldName) { return "unknown field" }
	return r.fields[name].Type.ToString()
}

