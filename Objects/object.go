package objects

/*

An Object represents a binary block of data, typically what one might consider a "file" on a disk,
which can be managed via ObjectStore. By abstracting Objects as the data set instead of as a named
file on disk, we can capture them in any number of places: files on disk, records in a database,
representations in an API, even codified chunks of data within our own executable.

Object optionally support fields; if the field map is nil, then they are not being used.

TODO: Isolate the encode/decode so that other tools can build against it and have a function that
      properly interacts with the same encoding scheme as us using ouo *Encoded* accessor methods.

*/

import (
	"fmt"
	"errors"
	"encoding/base64"
	lib "github.com/DigiStratum/GoLib"
)

type ObjectFieldType int

const (
	OFT_UNKNOWN ObjectFieldType = iota
	OFT_NUMERIC	// Any base 10 numeric form
	OFT_TEXTUAL	// Any string/text form
	OFT_DATETIME	// Any valid date and/or time form
	OFT_BOOLEAN	// Any boolean form
	OFT_BYTE	// any 8 bit form
	OFT_SHORT	// any 16 bit form
	OFT_INT		// any 32 bit form
	OFT_LONG	// any 64 but form
	OFT_FLOAT	// any floating point "real" value
	OFT_DOUBLE	// any double-precision "real" value
	OFT_FIXED	// any fixed point "real" value
	OFT_STRING	// any ASCII string
	OFT_CHAR	// any ASCII single character
	OFT_MBSTRING	// any multibyte string
	OFT_MBCHAR	// any multibyte single character
)

type ObjectTemplate map[string]ObjectFieldType

// Association of Type and value for a single Object Field
type ObjectField struct {
	Type		ObjectFieldType
	Value		*string			// Significance varies with Type
}

// Map name to ObjectField
type ObjectFieldMap	map[string]ObjectField

// Map name to value which may be nil for a given field
type ObjectFieldValues	map[string]*string

// A static Object that we're going to codify
type Object struct {
	isEncoded	bool			// Is the content encoded?
	content		*string			// Non-fielded Object "BLOB" representation
	fields		*ObjectFieldMap		// Field name to value map
}

// Make a new one of these
func NewObject() *Object {
	return &Object{}
}

// Make a new one of these with mapped fields (yey!)
// Note that field map could be just names & types (spec), or could also include values (record)
func NewObjectFromTemplate(objectTemplate *ObjectTemplate) *Object {

	// No template makes this the same as a plain Object
	if nil == objectTemplate {
		return NewObject()
	}

	// Transfer the names/types of our Template to a new ObjectFieldMap
	objectFieldMap := make(ObjectFieldMap)
	for name := range *objectTemplate {
		objectFieldMap[name] = ObjectField{
			Type: (*objectTemplate)[name],
		}
	}

	return &Object{
		fields:	&objectFieldMap,
	}
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
func (o *Object) GetEncodedContent() *string {
	return o.content
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

// Set all field values as supplied (they much match field map for this object!)
func (o *Object) SetFieldValuesFromMap(fieldValues *ObjectFieldValues) {
	// TODO: Implement Me; iterate over ObjectFieldValues; fail if any supplied field is not part of object fieldmap, use SetFieldValue() above for each
}

// Set all field values as supplied (they much match field map for this object!)
func (o *Object) SetFieldValuesFromJson(fieldJson string) {
	// TODO: IMplement Me; convert the json into ObjectFieldValues struct and call SetFieldValuesFromMap() above
}

// Determine whethere the value passes all the rules for the specified field type
func (o *Object) IsValueType(value *string, fieldType ObjectFieldType) bool {
	// TODO switch on type and run the value through the wringer here
	return true
}

// Return a readable string for each one
func (o *Object) GetObjectFieldTypeReadable(fieldType ObjectFieldType) string {
	switch (fieldType) {
		case OFT_UNKNOWN:
			return "unknown"
		case OFT_NUMERIC:
			return "numeric"
		case OFT_TEXTUAL:
			return "textual"
		case OFT_DATETIME:
			return "datetime"
		case OFT_BOOLEAN:
			return "boolean"
		case OFT_BYTE:
			return "byte"
		case OFT_SHORT:
			return "short"
		case OFT_INT:
			return "int"
		case OFT_LONG:
			return "long"
		case OFT_FLOAT:
			return "float"
		case OFT_DOUBLE:
			return "double"
		case OFT_FIXED:
			return "fixed"
		case OFT_STRING:
			return "string"
		case OFT_CHAR:
			return "char"
		case OFT_MBSTRING:
			return "mbstring"
		case OFT_MBCHAR:
			return "mbchar"
		default:
			return "unknown"
	}
}

