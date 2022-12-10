package restapi

// TODO: Refactor into separate files for each of the data structures/method collections

// TODO: Add support for builder pattern, or chaining, or "functional options":
// ref: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
// ref: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

// TODO: Can we refactor this to just use lib.HashMapIfc instead?
type HttpHeadersIfc interface {
	Has(name string) bool
	Get(name string) string
	Set(name string, value string)
	Merge(headers HttpHeadersIfc)
	GetHeaderNames() *[]string
	ToMap() *map[string]string
	IsEmpty() bool
}

// Name/value pair header map for Request or Response
type httpHeaders map[string]string

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpHeaders() HttpHeadersIfc {
	hh := make(httpHeaders)
	return &hh
}

// -------------------------------------------------------------------------------------------------
// HttpHeadersIfc Implementation
// -------------------------------------------------------------------------------------------------

// DO we have the named header?
func (r *httpHeaders) Has(name string) bool {
	if _, ok := (*r)[name]; ok { return true }
	return false
}

// Get a single header
// TODO: Change this to return nil (string pointer instead of string) if the value is not set - the difference between unset and set-but-empty
func (r *httpHeaders) Get(name string) string {
	if value, ok := (*r)[name]; ok { return value }
	return ""
}

// Set a single header
func (r *httpHeaders) Set(name string, value string) {
	(*r)[name] = value
}

// Merge an HttpHeaders set into our own
func (r *httpHeaders) Merge(headers HttpHeadersIfc) {
	names := headers.GetHeaderNames()
	for _, name := range *names {
		(*r)[name] = headers.Get(name)
	}
}

// Get the complete set of names
func (r *httpHeaders) GetHeaderNames() *[]string {
	names := make([]string, 0)
	for name, _ := range *r {
		names = append(names, name)
	}
	return &names
}

// Are there NO headers set?
func (r *httpHeaders) IsEmpty() bool {
	return 0 == len(*r)
}

// Some consumers need headers in the form of a simple data structure
func (r *httpHeaders) ToMap() *map[string]string {
	// Copy it, don't just point to our internal data, or Bad Things Will Happen (tm)
	h := make(map[string]string)
	for n, v := range *r { h[n] = v }
	return &h
}

