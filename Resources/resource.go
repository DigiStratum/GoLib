package resources

/*

TODO: Isolate the encode/decode so that other tools can build against it and have a function that
properly interacts with the same encoding scheme as us using our *Encoded* accessor methods.

*/

import (
	"encoding/base64"
	lib "github.com/DigiStratum/GoLib"
)

// A static Resource that we're going to codify
type Resource struct {
	isEncoded	bool	// Is the content encoded?
	content		*string
}

func NewResource() *Resource {
	return &Resource{}
}

func NewResourceFromString(content string) *Resource {
	r := NewResource()
	r.SetContentFromString(&content)
	return r
}

func NewResourceFromFile(path string) *Resource {
	r := NewResource()
	r.SetContentFromFile(path)
	return r
}

// Set the Resource Content from a plain text string (it will be encoded!)
func (r *Resource) SetContentFromString(content *string) {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(*content))
	r.content = &encodedContent
	r.isEncoded = true
}

// Set the Resource Content from a text string which is already endcoded
// (This is used by callers such as res2go that know how to pre-encode)
func (r *Resource) SetEncodedContentFromString(encodedContent *string) {
	r.content = encodedContent
	r.isEncoded = true
}

// Set the Resource Content from a source file path (it will be encoded!)
// (This is used to anything froma simple text file to full binary assets)
func (r *Resource) SetContentFromFile(path string) error {
	s, err := lib.ReadFileString(path)
	if nil != err { return err }
	r.SetContentFromString(s)
	return nil
}

// Get the Resource Content as a string (could be anything!)
func (r *Resource) GetContent() *string {
	// For non-encoded, raw content (probably loaded from disk, DB, service, etc)
	if ! r.isEncoded { return r.content }

	// For encoded content (probably compiled)
	decodedContentBytes, err := base64.StdEncoding.DecodeString(*r.content)
	if nil != err {
		// TODO: Handle errors
		s := ""
		return &s
	}
	decodedContent := string(decodedContentBytes)
	return &decodedContent
}

// Get the Resource Content as an Encoded string (you better know what to do with it!)
func (r *Resource) GetEncodedContent() *string {
	return r.content
}

