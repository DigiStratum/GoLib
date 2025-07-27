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
  * Reject attempts to set Body/Data for methods that do not support a body
*/

import (
	"net/url"

	"github.com/DigiStratum/GoLib/Data/metadata"
)

// Http Request public interface
type HttpRequestBuilderIfc interface {
	// URL bits
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

func NewHttpRequestBuilder(method HttpRequestMethod, url string) *httpRequestBuilder {
	builder := httpRequestBuilder{
		request: &httpRequest{},
	}
	builder.setURL(url)
	builder.setMethod(method)
	return &builder
}

func NewHttpRequestBuilderFromRequest(request *httpRequest) *httpRequestBuilder {
	requestCopy := *request
	builder := httpRequestBuilder{
		request: &requestCopy,
	}
	return &builder
}

// -------------------------------------------------------------------------------------------------
// HttpRequestBuilderIfc
// -------------------------------------------------------------------------------------------------

func (r *httpRequestBuilder) SetQueryString(queryString string) *httpRequestBuilder {
	// TODO: SetQueryParameters() should be called to match this
	r.request.url.RawQuery = queryString
	return r
}

// Set the query parameters
func (r *httpRequestBuilder) SetQueryParameters(params metadata.MetadataIfc) *httpRequestBuilder {
	// TODO: SetQueryString() should be called to match this
	r.request.queryParams = params
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

func (r *httpRequestBuilder) GetHttpRequest() *httpRequest {
	// A built request should always have headers; pick sane defaults if unset
	builtRequest := r.request
	if r.request.headers == nil {
		builtRequest.headers = NewHttpHeadersBuilder().GetHttpHeaders()
	}
	return builtRequest
}

// -------------------------------------------------------------------------------------------------
// httpRequestBuilder
// -------------------------------------------------------------------------------------------------

func (r *httpRequestBuilder) setURL(urlStr string) {
	if url, err := url.Parse(urlStr); err == nil {
		r.request.url = url
	}
}

func (r *httpRequestBuilder) setMethod(method HttpRequestMethod) {
	r.request.method = method
}
