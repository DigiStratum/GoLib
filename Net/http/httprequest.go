package http

/*

An Immutable HTTP Request structure and programmatic interface to it.

Use NewHttpRequestBuilder() to create one of these.

TODO:
  * Get rid of all the URL component operations, just use GetURL(), but keep QueryString conveniences
*/

import (
	"net/url"

	"github.com/DigiStratum/GoLib/Data/metadata"
)

type HttpRequestIfc interface {
	GetHost() string
	GetScheme() string
	GetURL() string
	GetURI() string
	GetMethod() HttpRequestMethod
	GetQueryString() string
	GetQueryParameters() metadata.MetadataIfc
	GetPathParameters() metadata.MetadataIfc
	GetBody() *string
	GetBodyData() *httpRequestBody
	GetHeaders() *httpHeaders
}

type httpRequest struct {
	url         *url.URL
	method      HttpRequestMethod
	queryParams metadata.MetadataIfc
	headers     *httpHeaders
	body        *string
	bodyData    *httpRequestBody
	pathParams  metadata.MetadataIfc
}

// -------------------------------------------------------------------------------------------------
// HttpRequestIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *httpRequest) GetHost() string {
	return r.url.Host
}

func (r *httpRequest) GetScheme() string {
	return r.url.Scheme
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
	return r.url.Path
}

func (r *httpRequest) GetQueryString() string {
	return r.url.RawQuery
}

func (r *httpRequest) GetBody() *string {
	return r.body
}

func (r *httpRequest) GetBodyData() *httpRequestBody {
	return r.bodyData
}

// Get the Request Headers
func (r *httpRequest) GetHeaders() *httpHeaders {
	return r.headers
}

// Get the path parameters (should be endpoint implementation)
func (r *httpRequest) GetPathParameters() metadata.MetadataIfc {
	return r.pathParams
}

// Get the query parameters
func (r *httpRequest) GetQueryParameters() metadata.MetadataIfc {
	return r.queryParams
}
