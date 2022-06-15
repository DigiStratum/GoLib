package encodingscheme

import(
	"fmt"
	"html"
)

/*
Implement EncodingSchemeIfc for HTML Escaping

We capture two values, encoded and raw so that we can cache the encode/decode results for this value
to avoid re-running expensive en|decode operations on repeat requests. When we set either value,
encoded or deccoded, we wipe out the cached value of the other to prevent returning stale results
and cause the conversion to run again and re-cache.
*/

type EncodingSchemeHTMLEscape struct {
	contentEncoded		*string
	contentRaw		*string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------
func NewEncodingSchemeHTMLEscape() *EncodingSchemeHTMLEscape {
	return &EncodingSchemeHTMLEscape{}
}

// -------------------------------------------------------------------------------------------------
// EncodingSchemeIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *EncodingSchemeHTMLEscape) SetEncodedValue(source *string) error {
	if nil == source { return fmt.Errorf("Cannot set: source is nil") }
	r.contentEncoded = source
	r.contentRaw = nil
	return nil
}

func (r *EncodingSchemeHTMLEscape) GetEncodedValue() (*string, error) {
	if nil == r.contentEncoded {
		if nil == r.contentRaw { return nil, fmt.Errorf("Cannot encode: raw content is nil") }
		encoded := html.EscapeString(*r.contentRaw)
		r.contentEncoded = &encoded
	}
	return r.contentEncoded, nil
}

func (r *EncodingSchemeHTMLEscape) SetDecodedValue(source *string) error {
	if nil == source { return fmt.Errorf("Cannot set: source is nil") }
	r.contentRaw = source
	r.contentEncoded = nil
	return nil
}

func (r *EncodingSchemeHTMLEscape) GetDecodedValue() (*string, error) {
	if nil == r.contentRaw {
		if nil == r.contentEncoded { return nil, fmt.Errorf("Cannot decode: encoded content is nil") }
		decoded := html.UnescapeString(*r.contentEncoded)
		r.contentRaw = &decoded
	}
	return r.contentRaw, nil
}

