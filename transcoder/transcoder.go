// DigiStratum GoLib - Transcoder for plain text content
package transcoder

import (
	"fmt"
	"encoding/base64"
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
	content		[]byte
	encodingScheme	EncodingScheme
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTranscoder() *Transcoder {
	return &Transcoder{
		encodingScheme:		ES_UNKNOWN,
	}
}

// -------------------------------------------------------------------------------------------------
// TranscoderIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Transcoder) FromString(content *string, encodingScheme EncodingScheme) error {
	// TODO: Validate the encodingScheme and return an error if it doesn't check out
	r.content = []byte(*content)
	r.encodingScheme = encodingScheme
	return nil
}

func (r *Transcoder) FromBytes(bytes *[]byte, encodingScheme EncodingScheme) error {
	// TODO: Validate the encodingScheme and return an error if it doesn't check out
	r.content = *bytes
	r.encodingScheme = encodingScheme
	return nil
}

func (r *Transcoder) FromFile(path string, encodingScheme EncodingScheme) error {
	bytes, err := lib.ReadFileBytes(path)
	if nil != err { return err }
	return r.FromString(bytes, encodingScheme)
}

func (r *Transcoder) ToString(encodingScheme EncodingScheme) (*string, error) {
	// TODO: Implement!
	return nil, nil
}

func (r *Transcoder) ToBytes(encodingScheme EncodingScheme) (*[]byte, error) {
	// TODO: Implement!
	return nil, nil
}

func (r *Transcoder) ToFile(path string, encodingScheme EncodingScheme) error {
	// TODO: Implement!
	return nil
}
