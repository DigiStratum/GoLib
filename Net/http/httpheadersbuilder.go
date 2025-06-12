package http

/*

Build a set of HTTP Headers for Request or Response

TODO:
 * Consider direct support with validation for specific, known header names
   ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers
 * Also validate custom header names against RFC 7230 for acceptable characters < 40 length
   ref: https://datatracker.ietf.org/doc/html/rfc7230
 * Deduplicate provided values for Set() against existing ones
*/

type HttpHeadersBuilderIfc interface {
	Set(name string, values ...string)
	Merge(headers HttpHeadersIfc)
	GetHttpHeaders() HttpHeadersIfc
}

// Name/value pair header map for Request or Response
type httpHeadersBuilder struct {
	headers *httpHeaders
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpHeadersBuilder() *httpHeadersBuilder {
	r := httpHeadersBuilder{
		headers: NewHttpHeaders(),
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// HttpHeadersBuilderIfc
// -------------------------------------------------------------------------------------------------

// Single-name, multi-value support
func (r *httpHeadersBuilder) Set(name string, values ...string) *httpHeadersBuilder {
	// If the named header is not set, then create it

	if _, ok := (*r.headers)[name]; !ok {
		(*r.headers)[name] = make([]string, 0)
	}

	for _, value := range values {
		(*r.headers)[name] = append((*r.headers)[name], value)
	}
	return r
}

func (r *httpHeadersBuilder) Merge(headers *httpHeaders) *httpHeadersBuilder {
	if headers != nil {
		for name, values := range *headers {
			// Use Set() to merge provided values with existing, instead of overwriting
			r.Set(name, values...)
		}
	}
	return r
}

func (r *httpHeadersBuilder) GetHttpHeaders() *httpHeaders {
	return r.headers
}
