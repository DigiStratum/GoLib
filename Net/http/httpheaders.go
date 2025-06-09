package http

type HttpHeadersIfc interface {
	Has(name string) bool
	GetHeaderNames() *[]string
	IsEmpty() bool

	// Deprecated: Single name, single value
	Get(name string) string
	ToMap() *map[string]string

	// Single-name, multi-value support
	All(name string) *[]string
	MapAll() *map[string][]string
}

// Name/value pair header map for Request or Response
type httpHeaders map[string][]string

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpHeaders() *httpHeaders {
	r := make(httpHeaders)
	return &r
}

// -------------------------------------------------------------------------------------------------
// HttpHeadersIfc Implementation
// -------------------------------------------------------------------------------------------------

// DO we have the named header?
func (r *httpHeaders) Has(name string) bool {
	if _, ok := (*r)[name]; ok {
		return true
	}
	return false
}

// Get a single header
// TODO: Change this to return nil (string pointer instead of string) if the value is not set - the difference between unset and set-but-empty
// Deprecated; use All() instead
func (r *httpHeaders) Get(name string) string {
	if values, ok := (*r)[name]; ok {
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
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
	return len(*r) == 0
}

// Some consumers need headers in the form of a simple data structure
// Deprecated; use MapAll() instead
func (r *httpHeaders) ToMap() *map[string]string {
	// Copy it, don't just point to our internal data, or Bad Things Will Happen (tm)
	h := make(map[string]string)
	for n, vs := range *r {
		if len(vs) == 0 {
			continue
		}
		h[n] = vs[0]
	}
	return &h
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
	for n, vs := range *r {
		h[n] = vs
	}
	return &h
}
