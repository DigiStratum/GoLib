// DigiStratum GoLib - Transcoder for plain text content
package transcoder

import (
	"fmt"
	//"encoding/base64"

	"github.com/DigiStratum/GoLib/FileIO"
)


type TranscoderIfc interface {
	FromString(content *string, encodingScheme EncodingScheme) error
	FromBytes(bytes *[]byte, encodingScheme EncodingScheme) error
	FromFile(path string, encodingScheme EncodingScheme) error
	ToString(requestedEncodingScheme EncodingScheme) (*string, error)
	ToBytes(requestedEncodingScheme EncodingScheme) (*[]byte, error)
	ToFile(path string, encodingScheme EncodingScheme) error
}

type Transcoder struct {
	content		map[EncodingScheme]*[]byte
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTranscoder() *Transcoder {
	return &Transcoder{
		content:	make(map[EncodingScheme]*[]byte),
	}
}

// -------------------------------------------------------------------------------------------------
// TranscoderIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Transcoder) FromString(content *string, encodingScheme EncodingScheme) error {
	// TODO: Validate the encodingScheme and return an error if it doesn't check out
	// Reset the content encodings - we don't want some old encoding of some other content hanging around
	r.content = make(map[EncodingScheme]*[]byte)
	contentBytes := []byte(*content)
	// TODO: if encoding scheme is ES_AUTO, then examine the content and determine the actual
	r.content[encodingScheme] = &contentBytes
	return nil
}

func (r *Transcoder) FromBytes(content *[]byte, encodingScheme EncodingScheme) error {
	// TODO: Validate the encodingScheme and return an error if it doesn't check out
	// Reset the content encodings - we don't want some old encoding of some other content hanging around
	r.content = make(map[EncodingScheme]*[]byte)
	// TODO: if encoding scheme is ES_AUTO, then examine the content and determine the actual
	r.content[encodingScheme] = content
	return nil
}

func (r *Transcoder) FromFile(path string, encodingScheme EncodingScheme) error {
	bytes, err := fileio.ReadFileBytes(path)
	if nil != err { return err }
	return r.FromBytes(bytes, encodingScheme)
}

func (r *Transcoder) ToString(requestedEncodingScheme EncodingScheme) (*string, error) {
	contentBytes, err := r.ToBytes(requestedEncodingScheme)
	if nil != err { return nil, err }
	content := string(*contentBytes)
	return &content, nil
}

func (r *Transcoder) ToBytes(requestedEncodingScheme EncodingScheme) (*[]byte, error) {
	if 0 == len(r.content) { return nil, fmt.Errorf("Content not initialized") }
	contentBytes, ok := r.content[requestedEncodingScheme]
	if ! ok {
		// This encodingScheme veriant of the content is not in the map yet - let's get it (or fail)!
		var err error
		contentBytes, err = r.convertEncodingScheme(requestedEncodingScheme)
		if nil != err { return nil, err }
	}
	return contentBytes, nil
}

func (r *Transcoder) ToFile(path string, encodingScheme EncodingScheme) error {
	// TODO: Implement!
	return nil
}

// -------------------------------------------------------------------------------------------------
// Transcoder Private Implementation
// -------------------------------------------------------------------------------------------------

func (r *Transcoder) convertEncodingScheme(targetEncodingScheme EncodingScheme)  (*[]byte, error) {
	pContentBytes, err := r.decodeContent()
	if nil != err { return nil, err }
	// TODO: Now call another function to encode the pContentBytes in the targetEncodingScheme
	return nil, nil
}

func (r *Transcoder) decodeContent()  (*[]byte, error) {
	// Do we already have non-encoded content cached?
	if contentBytes, ok := r.content[ES_NONE]; ok {
		return &contentBytes, nil
	}

	// TODO: find some cache entry to decode, possibly in order of least to greatest computational cost
/*
	// Decode existing content in cache to contentBytes
	var contentBytes []byte
	switch (targetEncodingScheme) {
		// TODO: Dec
		case ES_BASE64:		// Base 64 Encoding
		case ES_UUENCODE:	// UU-Encoding (EMAIL)
		case ES_HTTPESCAPE:	// HTTP Escaped Encoding (HTTP/URL/form-post)
		case ES_UNKNOWN:	// error!
	}
*/
	return nil, nil
}

