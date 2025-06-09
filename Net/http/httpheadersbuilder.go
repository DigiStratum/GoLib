package http

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
// HttpHeadersBuilderIfc Implementation
// -------------------------------------------------------------------------------------------------

// Single-name, multi-value support
// TODO: Consider variadic support for value(s) here
func (r *httpHeadersBuilder) Set(name string, values ...string) *httpHeadersBuilder {
	// If the header is not set, then create it
	if _, ok := (*r.headers)[name]; !ok {
		(*r.headers)[name] = make([]string, len(values))
	}

	if _, ok := (*r.headers)[name]; !ok {
		(*r.headers)[name] = make([]string, 0)
	}
	for i, value := range values {
		(*r.headers)[name][i] = value
	}
	return r
}

func (r *httpHeadersBuilder) Merge(headers HttpHeadersIfc) *httpHeadersBuilder {
	names := headers.GetHeaderNames()
	for _, name := range *names {
		values := headers.All(name)
		if nil == values {
			continue
		}
		r.Set(name, *values...)
	}
	return r
}

func (r *httpHeadersBuilder) GetHttpHeaders() *httpHeaders {
	return r.headers
}
