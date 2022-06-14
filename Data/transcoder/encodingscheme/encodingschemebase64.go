package encodingscheme

import(
	"fmt"
	"encoding/base64"
)

/*
Implement EncodingSchemeIfc for Base64 (RFC 4648)

We capture two values, encoded and raw so that we can cache the encode/decode results for this value
to avoid re-running expensive en|decode operations on repeat requests. When we set either value,
encoded or deccoded, we wipe out the cached value of the other to prevent returning stale results
and cause the conversion to run again and re-cache.
*/

type EncodingSchemeBase64 struct {
	contentEncoded		*string
	contentRaw		*string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------
func NewEncodingSchemeBase64() *EncodingSchemeBase64 {
	return &EncodingSchemeBase64{}
}

// -------------------------------------------------------------------------------------------------
// EncodingSchemeIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *EncodingSchemeBase64) SetEncodedValue(source *string) error {
	if nil == source { return fmt.Errorf("Cannot set: source is nil") }
	r.contentEncoded = source
	r.contentRaw = nil
	return nil
}

func (r *EncodingSchemeBase64) GetEncodedValue() (*string, error) {
	if nil == r.contentEncoded {
		if nil == r.contentRaw { return nil, fmt.Errorf("Cannot encode: raw content is nil") }
		encoded := base64.StdEncoding.EncodeToString([]byte(*r.contentRaw))
		r.contentEncoded = &encoded
	}
	return r.contentEncoded, nil
}

func (r *EncodingSchemeBase64) SetDecodedValue(source *string) error {
	if nil == source { return fmt.Errorf("Cannot set: source is nil") }
	r.contentRaw = source
	r.contentEncoded = nil
	return nil
}

func (r *EncodingSchemeBase64) GetDecodedValue() (*string, error) {
	if nil == r.contentRaw {
		if nil == r.contentEncoded { return nil, fmt.Errorf("Cannot decode: encoded content is nil") }
		decodedBytes, err := base64.StdEncoding.DecodeString(*r.contentEncoded)
		if nil != err { return nil, err }
		decodedStr := string(decodedBytes)
		r.contentRaw = &decodedStr
	}
	return r.contentRaw, nil
}

