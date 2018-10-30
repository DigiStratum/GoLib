package resources

import (
	"encoding/base64"
	lib "github.com/DigiStratum/GoLib"
)

// A static Resource that we're going to codify
type Resource struct {
	content		*string	// Encoded content
}

// Map of resource path to the Resource and its properties
// TODO: Add some supporting funcs to Resource to get a list of Resources below a given path (i.e. everything in a dir)
type ResourceMap map[string]*Resource

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

func (r *Resource) SetContentFromString(content *string) {
	// ref: https://golang.org/pkg/encoding/base64/#pkg-examples
	encodedContent := base64.StdEncoding.EncodeToString([]byte(*content))
	r.content = &encodedContent
}

func (r *Resource) SetContentFromFile(path string) error {
	s, err := lib.ReadFileString(path)
	if nil != err { return err }
	r.SetContentFromString(s)
	return nil
}

func (r *Resource) GetContent() *string {
	return r.content
}

func (r *Resource) GetDecodedContent() *string {
	decodedContentBytes, err := base64.StdEncoding.DecodeString(*r.content)
	if nil != err {
		// TODO: Handle errors
		s := ""
		return &s
	}
	decodedContent := string(decodedContentBytes)
	return &decodedContent
}

