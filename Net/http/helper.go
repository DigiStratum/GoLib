package http

import (
	"mime"
	"strings"
)

type HelperIfc interface {
	// Payload Helpers
	SingularizePostData(bodyData *HttpBodyData) map[string]string
	GetMimetype(uri string) string
}

type helper struct{}

var instance *helper

func init() {
	// Instantiate our singleton
	instance = NewHelper()
}

// Get our singleton
func GetHelper() HelperIfc {
	return instance
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHelper() *helper {
	return &helper{}
}

// -------------------------------------------------------------------------------------------------
// HelperIfc Implementation
// -------------------------------------------------------------------------------------------------

// Scan over the body data and, for each unique name, scrub out any duplicates
// TODO: refactor or make some variant which creates a value SET instead of tracking the dupes
func (hlpr *helper) SingularizePostData(bodyData *HttpBodyData) map[string]string {
	var data = make(map[string]string)
	for name, values := range *bodyData {
		if len(values) > 0 {
			data[name] = values[0]
		}
	}
	return data
}

// ref: https://golang.org/pkg/mime/#TypeByExtension
func (hlpr *helper) GetMimetype(uri string) string {
	dotpos := strings.LastIndex(uri, ".")
	if -1 == dotpos {
		return "application/octet-stream"
	}
	extension := uri[dotpos:]
	mimetype := mime.TypeByExtension(extension)
	if "" == mimetype {
		return "application/octet-stream"
	}
	return mimetype
}
