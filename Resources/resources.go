package resources

import (
	lib "github.com/DigiStratum/GoLib"
)

// A static Resource that we're going to codify
type Resource struct {
	content		*string	// Encoded content
}

// Map of resource path to the Resource and its properties
type ResourceMap map[string]*Resource

func NewResource() *Resource {
	return &Resource{}
}

func NewResourceFromString(content string) *Resource {
	r := NewResource()
	r.SetContentFromString(content)
}

func NewResourceFromFile(path string) *Resource {
	r := NewResource()
	r.SetContentFromFile(path)
}

func (r *Resource) SetContentFromString(content *string) {
	// FIXME: Encode the content!
	r.content = content
}

func (r *Resource) SetContentFromFile(path string) error {
	s, err := lib.ReadFileString(path)
	if nil != err { return err }
	r.SetContentFromString(s)
}

