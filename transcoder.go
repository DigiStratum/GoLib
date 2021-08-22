// DigiStratum GoLib - Transcoder
package golib

type EncodingScheme int

const (
	ES_UNKNOWN EncodingScheme = iota
	ES_NONE				// No Encoding
	ES_BASE64			// Base64 Encoding
)

type TranscoderIfc interface {
	LoadFromString(content string) error
	LoadFromFile(path string) error
}

type Transcoder struct {
	content		[]byte
}

func NewTranscoder() *Transcoder {
	return &Transcoder{}
}

// -------------------------------------------------------------------------------------------------
// TranscoderIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Transcoder) LoadFromString(content string) error {
	r.content = []byte(content)
	return nil
}

func (r *Transcoder) LoadFromFile(path string) error {
	b, err := lib.ReadFileBytes(path)
	if nil != err { return err }
	r.content = *b
	return nil
}