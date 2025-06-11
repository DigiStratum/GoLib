package http

/*

Build a set of HTTP Headers for Request or Response

FIXME:
 * Set function got entirely broken by refactoring

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

	// FIXME: This implementation is wrong. originally it created 0-length then appended values.
	// Now it creates exact length and then sets, but doesn't account for new values with append
	if _, ok := (*r.headers)[name]; !ok {
		(*r.headers)[name] = make([]string, len(values))
	}

	// FIXME: This is unreachable code, it will never be !ok here:
	if _, ok := (*r.headers)[name]; !ok {
		(*r.headers)[name] = make([]string, 0)
	}

	// FIXME: How do we append new values to existing ones? This is wrong!
	for i, value := range values {
		(*r.headers)[name][i] = value
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
