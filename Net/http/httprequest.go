package http

/*

An Immutable HTTP Request structure and programmatic interface to it.

Use NewHttpRequestBuilder() to create one of these.

*/

import (
	"net/url"

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
	GetQueryString() string
	GetQueryParameters() metadata.MetadataIfc
	GetPathParameters() metadata.MetadataIfc
	GetBody() *string
	GetBodyData() *httpRequestBody
	GetHeaders() *httpHeaders
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
	headers     *httpHeaders
	body        *string
	bodyData    *httpRequestBody
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
