package resources

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

func (r *Resource) SetContentFromString(content *string) {
	encodedContent := base64.StdEncoding.EncodeToString([]byte(*content))
	r.content = &encodedContent
	r.isEncoded = true
}

func (r *Resource) SetContentFromFile(path string) error {
	s, err := lib.ReadFileString(path)
	if nil != err { return err }
	r.SetContentFromString(s)
	return nil
}

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

func (r *Resource) GetRawContent() *string {
	return r.content
}

