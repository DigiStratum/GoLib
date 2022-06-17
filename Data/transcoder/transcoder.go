// DigiStratum GoLib - Transcoder for plain text content
package transcoder

import (
	"fmt"

	enc "github.com/DigiStratum/GoLib/Data/transcoder/encodingscheme"
)

type TranscoderIfc interface {
	SetEncoderScheme(encoder enc.EncodingSchemeIfc)
	SetDecoderScheme(decoder enc.EncodingSchemeIfc)
	Encode(source * string) (*string, error)
	Decode(source * string) (*string, error)
	Transcode(source *string) (*string, error)
}

type Transcoder struct {
	encoder, decoder	enc.EncodingSchemeIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewTranscoder() *Transcoder {
	return &Transcoder{}
}

// -------------------------------------------------------------------------------------------------
// TranscoderIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Transcoder) SetEncoderScheme(encoder enc.EncodingSchemeIfc) {
	r.encoder = encoder
}

func (r *Transcoder) SetDecoderScheme(decoder enc.EncodingSchemeIfc) {
	r.decoder = decoder
}

func (r *Transcoder) Encode(source * string) (*string, error) {
	if nil == r.encoder { return nil, fmt.Errorf("Encoder not set") }
	err := r.encoder.SetDecodedValue(source)
	if nil != err { return nil, err }
	return r.encoder.GetEncodedValue()
}

func (r *Transcoder) Decode(source * string) (*string, error) {
	if nil == r.decoder { return nil, fmt.Errorf("Decoder not set") }
	err := r.decoder.SetEncodedValue(source)
	if nil != err { return nil, err }
	return r.decoder.GetDecodedValue()
}

func (r *Transcoder) Transcode(source *string) (*string, error) {
	if (nil == r.encoder) || (nil == r.decoder) {
		return nil, fmt.Errorf("Transcoder not initialized")
	}
	var err error
	var decoded, encoded *string
	if err = r.decoder.SetEncodedValue(source); nil == err {
		if decoded, err = r.decoder.GetDecodedValue(); nil == err {
			if err = r.encoder.SetDecodedValue(decoded); nil == err {
				if encoded, err = r.encoder.GetEncodedValue(); nil == err {
					return encoded, nil
				}
			}
		}
	}
	return nil, err
}

