package objects

/*

An Object represents a binary block of data, typically what one might consider a "file" on a disk,
which can be managed via ObjectStore. By abstracting Objects as the data set instead of as a named
file on disk, we can capture them in any number of places: files on disk, records in a database,
representations in an API, even codified chunks of data within our own executable.

TODO: Isolate the encode/decode so that other tools can build against it and have a function that
      properly interacts with the same encoding scheme as us using ouo *Encoded* accessor methods.

TODO: Add support for populating ObjectField.Value according to ObjectField.Type validation rules,
      either individually or in sets (feed from JSON? CSV? MySQL? Other?)

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

type ObjectField struct {
	Type		ObjectFieldType
	Value		*string			// Significance varies with Type
}

type ObjectFieldMap	map[string]ObjectField  // Field name to value map - every value

// A static Object that we're going to codify
// We optionally support fields; if the fields map is empty, then they are not being used
type Object struct {
	isEncoded	bool			// Is the content encoded?
	content		*string			// Non-fielded Object "BLOB" representation
	fields		*ObjectFieldMap		// Field name to value map - every value
}

// Make a new one of these
func NewObject() *Object {
	return &Object{}
}

// Make a new one of these with mapped fields (yey!)
// Note that field map could be just names & types (spec), or could also include values (record)
func NewFIeldMappedObject(objectFieldMap *ObjectFieldMap) *Object {
	return &Object{
		fields:	objectFieldMap,
	}
}

func NewObjectFromString(content string) *Object {
	object := NewObject()
	object.SetContentFromString(&content)
	return object
}

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

