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
	ToString(encodingScheme EncodingScheme) (*string, error)
	ToBytes(encodingScheme EncodingScheme) (*[]byte, error)
	ToFile(path string, encodingScheme EncodingScheme) error
}

type Transcoder struct {
	content		map[EncodingScheme][]byte
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTranscoder() *Transcoder {
	return &Transcoder{
		content:	make(map[EncodingScheme][]byte),
	}
}

// -------------------------------------------------------------------------------------------------
// TranscoderIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Transcoder) FromString(content *string, encodingScheme EncodingScheme) error {
	// TODO: Validate the encodingScheme and return an error if it doesn't check out
	// Reset the content encodings - we don't want some old encoding of some other content hanging around
	r.content = make(map[EncodingScheme][]byte)
	r.content[encodingScheme] = []byte(*content)
	return nil
}

func (r *Transcoder) FromBytes(bytes *[]byte, encodingScheme EncodingScheme) error {
	// TODO: Validate the encodingScheme and return an error if it doesn't check out
	// Reset the content encodings - we don't want some old encoding of some other content hanging around
	r.content = make(map[EncodingScheme][]byte)
	r.content[encodingScheme] = *bytes
	return nil
}

func (r *Transcoder) FromFile(path string, encodingScheme EncodingScheme) error {
	bytes, err	 := fileio.ReadFileBytes(path)
	if nil != err { return err }
	return r.FromBytes(bytes, encodingScheme)
}

func (r *Transcoder) ToString(requestedEncodingScheme EncodingScheme) (*string, error) {
	if 0 == len(r.content) { return nil, fmt.Errorf("Content not initialized") }
	contentBytes, ok := r.content[requestedEncodingScheme]
	if ! ok {
		// This encodingScheme veriant of the content is not in the map yet - let's get it (or fail)!
		return nil, nil
	}
	/*
	switch r.encodingScheme {
		case ES_AUTO:			// Automagically detect Encoding
		case ES_NONE:			// No Encoding
		case ES_BASE64:			// Base 64 Encoding
		case ES_UUENCODE:		// UU-Encoding (EMAIL)
		case ES_HTTPESCAPE:		// HTTP Escaped Encoding (HTTP/URL/form-post)
		case ES_UNKNOWN:		// Hmm!
			// TODO: Allow this to pass through if requestedEncodingScheme matches
			return nil, fmt.Errorf("Unknown encoding")
	}
	return string(), nil
	*/
	content := string(contentBytes)
	return &content, nil
}

func (r *Transcoder) ToBytes(encodingScheme EncodingScheme) (*[]byte, error) {
	// TODO: Implement!
	return nil, nil
}

func (r *Transcoder) ToFile(path string, encodingScheme EncodingScheme) error {
	// TODO: Implement!
	return nil
}

// -------------------------------------------------------------------------------------------------
// Transcoder Private Implementation
// -------------------------------------------------------------------------------------------------

func (r *Transcoder) convertEncodingScheme(targetEncodingScheme EncodingScheme)  (*[]byte, error) {
	// TODO: Implement!
	return nil, nil
}