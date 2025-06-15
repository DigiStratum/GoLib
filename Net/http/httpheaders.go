package http

/*

A set of HTTP Headers for Request or Response

Note that a resource server will typically support a maximum size for the request payload. This
varies by server (for example, Apache defaults to 8K, IIS 16K, and Nginx 1M), and is configurable.

Because this implementation is server-agnostic, we do not enforce a maximum size here, but we do
provide a Size() function to allow the caller to determine the size of the headers in bytes. If a
resource server implementation decides to reject a request due to the size breaking the configured
limit, then it should return an HTTP 413 Payload Too Large response status to the client

TODO:
  * Make this iterable so that we can iterate over the name-value pairs
*/

import (
	"strings"
)

// Name/value pair header map for Request or Response
type httpHeadersData map[string][]string

type httpHeaders struct {
	headers httpHeadersData
}

type HttpHeadersIfc interface {
	Has(name string) bool
	GetNames() *[]string
	IsEmpty() bool
	Get(name string) *[]string
	ToMap() *httpHeadersData
	Size() int
}

// -------------------------------------------------------------------------------------------------
// HttpHeadersIfc Implementation
// -------------------------------------------------------------------------------------------------

// DO we have the named header?
func (r *httpHeaders) Has(name string) bool {
	_, ok := r.headers[name]
	return ok
}

// Get the complete set of names
func (r *httpHeaders) GetNames() *[]string {
	names := make([]string, 0)
	for name, _ := range r.headers {
		names = append(names, name)
	}
	return &names
}

// Are there NO headers set?
func (r *httpHeaders) IsEmpty() bool {
	return len(r.headers) == 0
}

func (r *httpHeaders) Get(name string) *[]string {
	if values, ok := r.headers[name]; ok {
		return &values
	}
	return nil
}

func (r *httpHeaders) ToMap() *httpHeadersData {
	// Copy it, don't just point to our internal data, or caller gets control of our content
	h := make(httpHeadersData)
	for n, vs := range r.headers {
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
	for _, values := range r.headers {
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

// TODO: Move header-specific support to header-specific source files
// Extract a list of languages from the Accept-Language header (if any)
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
func (r *httpHeaders) GetAcceptableLanguages() *[]string {
	if !r.Has("Accept-Language") {
		return nil
	}
	languages := r.getWeightedHeaderList("Accept-Language")
	// TODO: filter results according to what we support
	// (i.e. code, code-locale, code-locale-orthography); remove orthography/anything after locale
	// TODO: convert "*" into "default"
	return languages
}

// -------------------------------------------------------------------------------------------------
// httpHeaders
// -------------------------------------------------------------------------------------------------

// Extract a list of values from headers, ordered by preference expressed as quality value
// ref: https://developer.mozilla.org/en-US/docs/Glossary/Quality_values
func (r *httpHeaders) getWeightedHeaderList(headerName string) *[]string {

	// Get the value of the header we're after
	headerValue := r.Get(headerName)

	if (headerValue == nil) || (len(*headerValue) == 0) {
		// no header, no list!
		values := make([]string, 0)
		return &values
	}

	// Split the header on "," in case it has multiple values
	headerValues := strings.Split((*headerValue)[0], ",")
	values := make([]string, len(headerValues))

	// For each value we found...
	kept := 0
	for _, value := range headerValues {
		// If it has a ";", then there are extra details attached to split off
		if strings.Contains(value, ";") {
			valueParts := strings.Split(value, ";")
			// Keep the first part, the rest is metadata
			values[kept] = valueParts[0]
			kept++
			// TODO: Check out the other parts; they should be formatted as "name=value"
			// TODO: If the name = "q" and the value is a float from 0.0 to 1.0, use it to sort
			// TODO: If ANY value has a "q" specified, then sort must be performed, assume 1.0 for unspecified q
		} else {
			// Keep this value
			values[kept] = value
			kept++
		}
	}
	return &values
}
