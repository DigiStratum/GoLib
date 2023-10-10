package http

import (
	ver "github.com/DigiStratum/GoLib/Version"
)

// HTTP Response public interface
type HttpResponseIfc interface {
	// Body supports any media type, text or binary, so []byte is the common storage structure
	GetBinBody() *[]byte
	SetBinBody(body *[]byte)

	// Some Body reponses might be text, so we support conversion between []byte and string
	GetBody() *string
	SetBody(body *string)

	// Status is a code, but only specific, standards based statuses are supported
	GetStatus() HttpStatus
	SetStatus(status HttpStatus)

	GetHeaders() HttpHeadersIfc
	SetHeaders(headers HttpHeadersIfc)
}

type httpResponse struct {
	httpRequest		HttpRequestIfc
	status			HttpStatus
	protocolVersion		ver.VersionIfc
	headers			HttpHeadersIfc

	// Retain both forms of body representation to avoid wasteful repeat conversions
	bodybytes		*[]byte
	bodystring		*string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------


func NewHttpResponse() *httpResponse {
	// Initialize with properties that will generate valid result for all Getter methods
	bodystring := ""
	return &httpResponse{
		status:		STATUS_UNKNOWN,
		headers:	NewHttpHeaders(),
		bodystring:	&bodystring,
	}
}

// -------------------------------------------------------------------------------------------------
// HttpResponseIfc Implementation
// -------------------------------------------------------------------------------------------------

// Get the original request that resulted in this HttpResponse
func (r *httpResponse) GetRequest() HttpRequestIfc {
	return r.httpRequest
}

// Supply the original request that resulted in this HttpResponse
func (r *httpResponse) SetRequest(httpRequest HttpRequestIfc) {
	r.httpRequest = httpRequest
}

func (r *httpResponse) GetBinBody() *[]byte {
	if nil == r.bodybytes {
		if nil == r.bodystring { return nil }
		bb := []byte(*r.bodystring)
		r.bodybytes = &bb
	}
	return r.bodybytes
}

func (r *httpResponse) SetBinBody(body *[]byte) {
	r.bodybytes = body
	r.bodystring = nil
}

func (r *httpResponse) GetBody() *string {
	if nil == r.bodystring {
		if nil == r.bodybytes { return nil }
		bs := string((*r.bodybytes)[:])
		r.bodystring = &bs
	}
	return r.bodystring
}

func (r *httpResponse) SetBody(body *string) {
	r.bodystring = body
	r.bodybytes = nil
}

func (r *httpResponse) GetStatus() HttpStatus {
	return r.status
}

func (r *httpResponse) SetStatus(status HttpStatus) {
	r.status = status
}

func (r *httpResponse) GetProtocolVersion() ver.VersionIfc{
	return r.protocolVersion
}

func (r *httpResponse) SetProtocolVersion(version string) {
	if v := ver.NewMajorMinor(version); nil != v { r.protocolVersion = v }
}

func (r *httpResponse) GetHeaders() HttpHeadersIfc {
	return r.headers
}

func (r *httpResponse) SetHeaders(headers HttpHeadersIfc) {
	r.headers = headers
}


