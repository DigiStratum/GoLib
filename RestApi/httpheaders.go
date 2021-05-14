package restapi

// TODO: Refactor into separate files for each of the data structures/method collections

// TODO: Add support for builder pattern, or chaining, or "functional options":
// ref: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
// ref: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

// TODO: Can we refactor this to just use lib.HashMapIfc instead?
type HttpHeadersIfc interface {
	Get(name string) string
	Set(name string, value string)
	Merge(headers *HttpHeaders)
	GetHeaderNames() *[]string
}

// Name/value pair header map for Request or Response
type httpHeaders [string]string

// Make a new one of these
func NewHttpHeaders() HttpHeadersIfc {
	hh := make(httpHeaders)
	return &hh
}

// Get a single header
// TODO: Change this to return nil (string pointer instead of string) if the value is not set - the difference between unset and set-but-empty
func (hdrs *httpHeaders) Get(name string) string {
	if value, ok := (*hdrs)[name]; ok { return value }
	return ""
}

// Set a single header
func (hdrs *httpHeaders) Set(name string, value string) {
	(*hdrs)[name] = value
}

// Merge an HttpHeaders set into our own
func (hdrs *httpHeaders) Merge(headers HttpHeadersIfc) {
	names := headers.GetHeaderNames()
	for name := range *names {
		(*hdrs)[name] = headers.Get(name)
	}
}

// Get the complete set of names
func (hdrs *httpHeaders) GetHeaderNames() *[]string {
	names := make([]string, 0)
	for name, _ := range *headers {
		names = append(names, name)
	}
	return &names
}

