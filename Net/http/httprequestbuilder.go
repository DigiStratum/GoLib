package http

/*

TODO
  * Add "BodyBuilder helper for Body Data to build it up one name/value at a time
  * Fix up the SetQueryString and SetQueryParameters funcs so that they match values, cannot diverge
  * Get rid of all the URL component operations, just use SetURL(), but keep QueryString conveniences
  * Update the URL when querystring/params are set
  * Add support for page fragment URL component (not something that we can see for incoming requests,
    but useful for outgoing requests)
  * Add Factory Function(s) to convert other request types to this type
    - PathParameters are useful to parse requests from other sources, but not for building a new request

*/

import (
	"net/url"

	"github.com/DigiStratum/GoLib/Data/metadata"
)

// Http Request public interface
type HttpRequestBuilderIfc interface {
	// URL bits
	SetProtocol(protocol string) *httpRequestBuilder
	SetHost(host string) *httpRequestBuilder
	SetURL(urlStr string) *httpRequestBuilder
	SetURI(uri string) *httpRequestBuilder
	SetMethod(method HttpRequestMethod) *httpRequestBuilder
	SetQueryString(queryString string) *httpRequestBuilder
	SetQueryParameters(params metadata.MetadataIfc) *httpRequestBuilder

	// Body Bits
	SetBody(body *string) *httpRequestBuilder
	SetBodyData(bodyData *httpRequestBody) *httpRequestBuilder

	// Headers
	SetHeaders(headers HttpHeadersIfc) *httpRequestBuilder

	// Build
	GetHttpRequest() *httpRequest
}

type httpRequestBuilder struct {
	request *httpRequest
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpRequestBuilder() *httpRequestBuilder {
	return &httpRequestBuilder{
		request: &httpRequest{},
	}
}

// -------------------------------------------------------------------------------------------------
// HttpRequestBuilderIfc
// -------------------------------------------------------------------------------------------------

func (r *httpRequestBuilder) SetProtocol(protocol string) *httpRequestBuilder {
	r.request.protocol = protocol
	return r
}

func (r *httpRequestBuilder) SetHost(host string) *httpRequestBuilder {
	r.request.host = host
	return r
}

func (r *httpRequestBuilder) SetScheme(scheme string) *httpRequestBuilder {
	r.request.scheme = scheme
	return r
}

func (r *httpRequestBuilder) SetURL(urlStr string) *httpRequestBuilder {
	u, err := url.Parse(urlStr)
	if nil != err {
		//log.GetLogger().Warn("HttpRequest.SetUrl() - failed to parse as a URL: '%s'", u)
		return nil
	}
	r.request.url = u
	return r
}

func (r *httpRequestBuilder) SetMethod(method HttpRequestMethod) *httpRequestBuilder {
	r.request.method = method
	return r
}

func (r *httpRequestBuilder) SetURI(uri string) *httpRequestBuilder {
	r.request.uri = uri
	return r
}

func (r *httpRequestBuilder) SetQueryString(queryString string) *httpRequestBuilder {
	// TODO: SetQueryParameters() should be called to match this
	r.request.queryString = queryString
	return r
}

func (r *httpRequestBuilder) SetBody(body *string) *httpRequestBuilder {
	r.request.body = body
	return r
}

func (r *httpRequestBuilder) SetBodyData(bodyData *httpRequestBody) *httpRequestBuilder {
	r.request.bodyData = bodyData
	return r
}

// Set the Request Headers
func (r *httpRequestBuilder) SetHeaders(headers *httpHeaders) *httpRequestBuilder {
	r.request.headers = headers
	return r
}

// Set the query parameters
func (r *httpRequestBuilder) SetQueryParameters(params metadata.MetadataIfc) *httpRequestBuilder {
	// TODO: SetQueryString() should be called to match this
	r.request.queryParams = params
	return r
}

func (r *httpRequestBuilder) GetHttpRequest() *httpRequest {
	return r.request
}
