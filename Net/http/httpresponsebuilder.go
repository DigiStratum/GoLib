package http

import (
	ver "github.com/DigiStratum/GoLib/Version"
)

type HttpResponseBuilderIfc interface {
	SetBinBody(body *[]byte) *httpResponseBuilder
	SetBody(body *string) *httpResponseBuilder
	SetStatus(status HttpStatus) *httpResponseBuilder
	SetProtocolVersion(version string) *httpResponseBuilder
	SetHeaders(headers HttpHeadersIfc) *httpResponseBuilder
	GetHttpResponse() *httpResponse
}

type httpResponseBuilder struct {
	response *httpResponse
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpResponseBuilder() *httpResponseBuilder {
	r := httpResponseBuilder{
		response: &httpResponse{},
	}
	return &r
}

func (r *httpResponseBuilder) SetBinBody(body *[]byte) *httpResponseBuilder {
	r.response.bodybytes = body
	r.response.bodystring = nil
	return r
}

func (r *httpResponseBuilder) SetBody(body *string) *httpResponseBuilder {
	r.response.bodystring = body
	r.response.bodybytes = nil
	return r
}

func (r *httpResponseBuilder) SetStatus(status HttpStatus) *httpResponseBuilder {
	r.response.status = status
	return r
}

func (r *httpResponseBuilder) SetProtocolVersion(version string) *httpResponseBuilder {
	if v := ver.NewMajorMinor(version); nil != v {
		r.response.protocolVersion = v
	}
	return r
}

func (r *httpResponseBuilder) SetHeaders(headers HttpHeadersIfc) *httpResponseBuilder {
	r.response.headers = headers
	return r
}

func (r *httpResponseBuilder) GetHttpResponse() *httpResponse {
	return r.response
}
