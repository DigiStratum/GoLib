package http

/*

TODO
 * Add "BodyBuilder helper for Body Data to build it up one name/value at a time
*/

import (
	"net/url"

	"github.com/DigiStratum/GoLib/Data/metadata"
)

// Http Request public interface
type HttpRequestBuilderIfc interface {
	// Our own interface
	SetProtocol(protocol string) *httpRequestBuilder
	SetHost(host string) *httpRequestBuilder
	SetRemoteAddr(remoteAddr string) *httpRequestBuilder
	SetURL(urlStr string) *httpRequestBuilder
	SetURI(uri string) *httpRequestBuilder
	SetMethod(method HttpRequestMethod) *httpRequestBuilder
	SetQueryString(queryString string) *httpRequestBuilder
	SetQueryParameters(params metadata.MetadataIfc) *httpRequestBuilder
	SetBody(body *string) *httpRequestBuilder
	SetBodyData(bodyData *HttpBodyData) *httpRequestBuilder
	SetHeaders(headers HttpHeadersIfc) *httpRequestBuilder
	SetPathParameters(params metadata.MetadataIfc) *httpRequestBuilder
	GetHttpRequest() *httpRequest
}

type httpRequestBuilder struct {
	request *httpRequest
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpRequestBuilder() *httpRequestBuilder {
	bodyData := make(HttpBodyData)
	return &httpRequestBuilder{
		request: &httpRequest{
			headers:  NewHttpHeaders(),
			bodyData: &bodyData,
		},
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

func (r *httpRequestBuilder) SetRemoteAddr(remoteAddr string) *httpRequestBuilder {
	r.request.remoteAddr = remoteAddr
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
	r.request.queryString = queryString
	return r
}

func (r *httpRequestBuilder) SetBody(body *string) *httpRequestBuilder {
	r.request.body = body
	return r
}

func (r *httpRequestBuilder) SetBodyData(bodyData *HttpBodyData) *httpRequestBuilder {
	r.request.bodyData = bodyData
	return r
}

// Set the Request Headers
func (r *httpRequestBuilder) SetHeaders(headers HttpHeadersIfc) *httpRequestBuilder {
	r.request.headers = headers
	return r
}

// Set the path parameters (should be endpointwrapper)
func (r *httpRequestBuilder) SetPathParameters(params metadata.MetadataIfc) *httpRequestBuilder {
	r.request.pathParams = params
	return r
}

// Set the query parameters
func (r *httpRequestBuilder) SetQueryParameters(params metadata.MetadataIfc) *httpRequestBuilder {
	r.request.queryParams = params
	return r
}

func (r *httpRequestBuilder) GetHttpRequest() *httpRequest {
	return r.request
}
