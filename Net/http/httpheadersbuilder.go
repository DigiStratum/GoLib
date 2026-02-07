package http

/*

Build a set of HTTP Headers for Request or Response

TODO:
  * Consider direct support with validation for specific, known header names
    ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers
  * Also validate custom header names against RFC 7230 for acceptable characters < 40 length
    ref: https://datatracker.ietf.org/doc/html/rfc7230
  * Deduplicate provided values for Set() against existing ones
  * Merge() may need to optionally overwrite existing values instead of append
*/

// Name/value pair header map for Request or Response
type httpHeadersBuilder struct {
	headers *httpHeaders
}

type HttpHeadersBuilderIfc interface {
	Set(name string, values ...string) *httpHeadersBuilder
	Append(name string, values ...string) *httpHeadersBuilder
	Merge(headers HttpHeadersIfc) *httpHeadersBuilder
	Override(headers HttpHeadersIfc) *httpHeadersBuilder
	GetHttpHeaders() *httpHeaders
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpHeadersBuilder() *httpHeadersBuilder {
	r := httpHeadersBuilder{
		headers: &httpHeaders{
			headers: make(httpHeadersData),
		},
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// HttpHeadersBuilderIfc
// -------------------------------------------------------------------------------------------------

// Single-name, multi-value support
func (r *httpHeadersBuilder) Set(name string, values ...string) *httpHeadersBuilder {
	// Create the named header (and clear the value!)
	r.headers.headers[name] = make([]string, 0)
	return r.Append(name, values...)
}

func (r *httpHeadersBuilder) Append(name string, values ...string) *httpHeadersBuilder {
	// If the named header is not set, then create it
	if _, ok := r.headers.headers[name]; !ok {
		r.headers.headers[name] = make([]string, 0)
	}

	r.headers.headers[name] = append(r.headers.headers[name], values...)
	return r
}

func (r *httpHeadersBuilder) Merge(headers HttpHeadersIfc) *httpHeadersBuilder {
	if headers != nil {
		names := headers.GetNames()
		if names != nil {
			for _, name := range *names {
				// Use Set() to merge provided values with existing, instead of overwriting
				values := headers.Get(name)
				if values == nil {
					continue
				}
				r.Set(name, *values...)
			}
		}
	}
	return r
}

func (r *httpHeadersBuilder) Override(headers HttpHeadersIfc) *httpHeadersBuilder {
	if headers != nil {
		names := headers.GetNames()
		if names != nil {
			for index := range *names {
				name := (*names)[index]
				// Directly set the header, overwriting any existing values
				values := headers.Get(name)
				if values == nil { // If no values, then skip
					continue
				}
				r.Set(name, *values...)
			}
		}
	}
	return r
}

func (r *httpHeadersBuilder) GetHttpHeaders() *httpHeaders {
	return r.headers
}
