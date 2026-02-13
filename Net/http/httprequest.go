package http

/*

An Immutable HTTP Request structure and programmatic interface to it.

Use NewHttpRequestBuilder() to create one of these.

FIXME:
 * PathParameters support is incomplete; we need to specify some sort of pattern that cna then be
   used to match the path parameters to the request URI.

TODO:
  * Get rid of all the URL component operations, just use GetURL(), but keep QueryString conveniences
  * Add general support for request metadata, such as request ID, etc.
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
	SetPathParameters(params metadata.MetadataIfc)
	GetBody() *string
	GetBodyData() *httpRequestBody
	GetHeaders() *httpHeaders
	GetBuilder() *httpRequestBuilder
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

// Set the path parameters (used by endpoint wrappers to populate extracted values)
func (r *httpRequest) SetPathParameters(params metadata.MetadataIfc) {
	r.pathParams = params
}

// Get the query parameters
func (r *httpRequest) GetQueryParameters() metadata.MetadataIfc {
	return r.queryParams
}

func (r *httpRequest) GetBuilder() *httpRequestBuilder {
	builder := NewHttpRequestBuilder(
		r.method,
		r.url.String(),
	).
		SetQueryParameters(r.queryParams).
		SetHeaders(r.headers).
		SetBody(r.body).
		SetBodyData(r.bodyData)

	return builder
}
