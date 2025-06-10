package http

/*

An Immutable HTTP Request structure and programmatic interface to it.

Use NewHttpRequestBuilder() to create one of these.

FIXME:
 * Getters for non-primitive data types are returning mutable structures. We need Builders for ALL
   of these to make them immutable as well.

*/

import (
	"net/url"
	"strings"

	"github.com/DigiStratum/GoLib/Data/metadata"
)

type HttpRequestIfc interface {
	GetProtocol() string
	GetHost() string
	GetRemoteAddr() string
	GetScheme() string
	GetURL() string
	GetURI() string
	GetMethod() HttpRequestMethod
	IsIdempotentMethod() bool
	GetQueryString() string
	GetQueryParameters() metadata.MetadataIfc
	GetPathParameters() metadata.MetadataIfc
	GetBody() *string
	GetBodyData() *HttpBodyData
	GetHeaders() HttpHeadersIfc
	GetAcceptableLanguages() *[]string
}

type httpRequest struct {
	protocol    string
	host        string
	remoteAddr  string
	scheme      string
	url         *url.URL
	method      HttpRequestMethod
	uri         string
	queryString string
	queryParams metadata.MetadataIfc
	headers     HttpHeadersIfc
	body        *string
	bodyData    *HttpBodyData
	pathParams  metadata.MetadataIfc
}

// -------------------------------------------------------------------------------------------------
// HttpRequestIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *httpRequest) GetProtocol() string {
	return r.protocol
}

func (r *httpRequest) GetHost() string {
	return r.host
}

func (r *httpRequest) GetRemoteAddr() string {
	return r.remoteAddr
}

func (r *httpRequest) GetScheme() string {
	return r.scheme
}

func (r *httpRequest) GetURL() string {
	if nil == r.url {
		return ""
	}
	return r.url.String()
}

func (r *httpRequest) GetMethod() HttpRequestMethod {
	return r.method
}

func (r *httpRequest) GetURI() string {
	return r.uri
}

func (r *httpRequest) GetQueryString() string {
	return r.queryString
}

func (r *httpRequest) GetBody() *string {
	return r.body
}

func (r *httpRequest) GetBodyData() *HttpBodyData {
	return r.bodyData
}

// Get the Request Headers
func (r *httpRequest) GetHeaders() HttpHeadersIfc {
	return r.headers
}

// Extract a list of languages from the Accept-Language header (if any)
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
// TODO: Move this to a convenience method of HttpHeadersIfc - it doesn't belong here!
func (r *httpRequest) GetAcceptableLanguages() *[]string {
	languages := r.getWeightedHeaderList("Accept-Language")
	// TODO: filter results according to what we support
	// (i.e. code, code-locale, code-locale-orthography); remove orthography/anything after locale
	// TODO: convert "*" into "default"
	return languages
}

// Quick check if the current request is expected to be idempotent in implementation
// TODO: Move this to a convenience method of HttpRequestMethod - it doesn't belong here!
func (r *httpRequest) IsIdempotentMethod() bool {
	return (METHOD_PATCH != r.method) && (METHOD_POST != r.method)
}

// Get the path parameters (should be endpoint implementation)
func (r *httpRequest) GetPathParameters() metadata.MetadataIfc {
	return r.pathParams
}

// Get the query parameters
func (r *httpRequest) GetQueryParameters() metadata.MetadataIfc {
	return r.queryParams
}

// -------------------------------------------------------------------------------------------------
// httpRequest Implementation
// -------------------------------------------------------------------------------------------------

// Extract a list of values from headers, ordered by preference expressed as quality value
// ref: https://developer.mozilla.org/en-US/docs/Glossary/Quality_values
func (r *httpRequest) getWeightedHeaderList(headerName string) *[]string {

	// Get the value of the header we're after
	headerValue := r.GetHeaders().Get(headerName)
	if len(headerValue) == 0 {
		// no header, no list!
		values := make([]string, 0)
		return &values
	}

	// Split the header on "," in case it has multiple values
	headerValues := strings.Split(headerValue, ",")
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
