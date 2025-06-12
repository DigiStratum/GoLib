package http

/*

A set of HTTP Headers for Request or Response

Note that a resource server will typically support a maximum size for the request payload. This
varies by server (for example, Apache defaults to 8K, IIS 16K, and Nginx 1M), and is configurable.

Because this implementation is server-agnostic, we do not enforce a maximum size here, but we do
provide a Size() function to allow the caller to determine the size of the headers in bytes. If a
resource server implementation decides to reject a request due to the size breaking the configured
limit, then it should return an HTTP 413 Payload Too Large response status to the client

FIXME:
 * Make httpHeaders a struct that HAS a map, not IS a map! being a map makes it possible for any
   layer of code to modify the map contents directly. by hiding it in a private struct member, we
   can prevent consumers from touching it.

*/

type HttpHeadersIfc interface {
	Has(name string) bool
	GetNames() *[]string
	IsEmpty() bool
	Get(name string) *[]string
	ToMap() *map[string][]string
	Size() int
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
	_, ok := (*r)[name]
	return ok
}

// Get the complete set of names
func (r *httpHeaders) GetNames() *[]string {
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

func (r *httpHeaders) Get(name string) *[]string {
	if values, ok := (*r)[name]; ok {
		return &values
	}
	return nil
}

func (r *httpHeaders) ToMap() *map[string][]string {
	// Copy it, don't just point to our internal data, or caller gets control of our content
	h := make(map[string][]string)
	for n, vs := range *r {
		h[n] = vs
	}
	return &h
}

/*
A set of headers will render out as a text block like:
---
Header-Name: value1; value2; value3\n
Header-Name2: value1; value2\n
\n\n
---
Our Size() function will return the length of this text block, including the separators
*/
func (r *httpHeaders) Size() int {
	l := 1
	for _, values := range *r {
		for _, value := range values {
			// Length of the value plus colon-space and space+semicolon separators
			l += (len(value) + 4)
		}
		l++ // Length of a newline bewtween headers
	}
	// If there are no headers, we return 2 for final, double newline
	if l == 1 {
		l = 2
	}
	return l
}
