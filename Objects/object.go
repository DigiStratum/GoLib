package objects

/*

An Object represents a binary block of data, typically what one might consider a "file" on a disk,
which can be managed via ObjectStore. By abstracting Objects as the data set instead of as a named
file on disk, we can capture them in any number of places: files on disk, records in a database,
representations in an API, even codified chunks of data within our own executable.

TODO: Isolate the encode/decode so that other tools can build against it and have a function that
properly interacts with the same encoding scheme as us using ouo *Encoded* accessor methods.

*/

import (
	"encoding/base64"
	lib "github.com/DigiStratum/GoLib"
)

// A static Object that we're going to codify
type Object struct {
	isEncoded	bool	// Is the content encoded?
	content		*string
}

func NewObject() *Object {
	return &Object{}
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

