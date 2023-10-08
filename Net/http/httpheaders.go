package http

// TODO: Refactor into separate files for each of the data structures/method collections

// TODO: Add support for builder pattern, or chaining, or "functional options":
// ref: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
// ref: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

// TODO: Can we refactor this to just use lib.HashMapIfc instead?
type HttpHeadersIfc interface {
	Has(name string) bool
	GetHeaderNames() *[]string
	IsEmpty() bool

	// Deprecated: Single name, single value
	Get(name string) string
	Set(name string, value string)
	ToMap() *map[string]string
	Merge(headers HttpHeadersIfc)

	// Single-name, multi-value support
	Add(name string, value string)
	All(name string) *[]string
	MapAll() *map[string][]string
	MergeAll(headers HttpHeadersIfc)
}

// Name/value pair header map for Request or Response
type httpHeaders map[string][]string

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpHeaders() HttpHeadersIfc {
	r := make(httpHeaders)
	return &r
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
// Deprecated; use All() instead
func (r *httpHeaders) Get(name string) string {
	if values, ok := (*r)[name]; ok {
		if len(values) > 0 { return values[0] }
	}
	return ""
}

// Set a single header
// Deprecated; use Add() instead
func (r *httpHeaders) Set(name string, value string) {
	(*r)[name] = make([]string, 1)
	(*r)[name][0] = value
}

// Merge an HttpHeaders set into our own
// Deprecated; use MergeAll() instead
func (r *httpHeaders) Merge(headers HttpHeadersIfc) {
	names := headers.GetHeaderNames()
	for _, name := range *names {
		r.Set(name, headers.Get(name))
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
// Deprecated; use MapAll() instead
func (r *httpHeaders) ToMap() *map[string]string {
	// Copy it, don't just point to our internal data, or Bad Things Will Happen (tm)
	h := make(map[string]string)
	for n, vs := range *r {
		if len(vs) == 0 { continue }
		h[n] = vs[0]
	}
	return &h
}

// Single-name, multi-value support
// TODO: Consider variadic support for value(s) here
func (r *httpHeaders) Add(name string, value string) {
	if _, ok := (*r)[name]; ! ok { (*r)[name] = make([]string, 0) }
	(*r)[name] = append((*r)[name], value)
}

func (r *httpHeaders) All(name string) *[]string {
	if values, ok := (*r)[name]; ok {
		return &values
	}
	emptySet := make([]string, 0)
	return &emptySet
}

func (r *httpHeaders) MapAll() *map[string][]string {
	// Copy it, don't just point to our internal data, or Bad Things Will Happen (tm)
	h := make(map[string][]string)
	for n, vs := range *r { h[n] = vs }
	return &h
}

func (r *httpHeaders) MergeAll(headers HttpHeadersIfc) {
	names := headers.GetHeaderNames()
	for _, name := range *names {
		values := headers.All(name)
		if nil == values { continue }
		for _, value := range *values {
			r.Add(name, value)
		}
	}
}

