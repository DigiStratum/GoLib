package http

import (
	"fmt"

	obj "github.com/DigiStratum/GoLib/Object"
	ver "github.com/DigiStratum/GoLib/Version"
)

// HTTP Response public interface
type HttpResponseIfc interface {
	// Body supports any media type, text or binary, so []byte is the common storage structure
	GetBinBody() *[]byte
	// Some Body reponses might be text, so we support conversion between []byte and string
	GetBody() *string
	// Status is a code, but only specific, standards based statuses are supported
	GetStatus() HttpStatus
	GetProtocolVersion() ver.VersionIfc
	GetHeaders() HttpHeadersIfc
}

type httpResponse struct {
	status          HttpStatus
	protocolVersion ver.VersionIfc
	headers         HttpHeadersIfc

	// Retain both forms of body representation to avoid wasteful repeat conversions
	bodybytes  *[]byte
	bodystring *string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Produce an HTTP response with standard headers
func NewHttpResponseStandard(status HttpStatus, body *string, contentType string) *httpResponse {
	hdrs := NewHttpHeadersBuilder().
		Set("content-type", contentType).
		GetHttpHeaders()
	return NewHttpResponseWithHeaders(status, body, hdrs)
}

// Produce an HTTP response, code only, no headers/body
func NewHttpReponseCode(status HttpStatus) *httpResponse {
	body := ""
	return NewHttpResponseWithHeaders(status, &body, NewHttpHeadersBuilder().GetHttpHeaders())
}

// Produce an HTTP response, code and default status text, JSON format
func NewHttpResponseSimpleJson(status HttpStatus) *httpResponse {
	var staticResponse string
	if status.IsStatus2xx() || status.IsStatus3xx() || status.IsStatus1xx() {
		staticResponse = fmt.Sprintf("[ { \"msg\": \"%s\" } ]", status.ToString())
	} else {
		staticResponse = fmt.Sprintf("[ { \"error\": { \"msg\": \"%s\" } } ]", status.ToString())
	}
	return NewHttpResponseStandard(status, &staticResponse, "application/json")
}

// Produce an HTTP error response by HTTP status code only
func NewHttpResponseError(status HttpStatus) *httpResponse {
	hdrs := NewHttpHeadersBuilder().
		Set("content-type", "text/plain").
		GetHttpHeaders()
	body := status.ToString()
	return NewHttpResponseWithHeaders(status, &body, hdrs)
}

// Produce an OK HTTP response with standard headers
func NewHttpResponseOk(body *string, contentType string) *httpResponse {
	return NewHttpResponseStandard(STATUS_OK, body, contentType)
}

// Produce an ERROR HTTP response with JSON message body and standard headers
func NewHttpResponseErrorJson(status HttpStatus, message string) *httpResponse {
	staticResponse := fmt.Sprintf("[ { \"error\": { \"msg\": \"%s\" } } ]", message)
	return NewHttpResponseStandard(status, &staticResponse, "application/json")
}

// Produce an HTTP response from an Object (200 OK)
func NewHttpResponseObject(object *obj.Object, uri string) *httpResponse {
	hlpr := GetHelper()
	return NewHttpResponseStandard(STATUS_OK, object.GetContent(), hlpr.GetMimetype(uri))
}

// Produce an HTTP response from an Object (200 OK)
func NewHttpResponseObjectCacheable(object *obj.Object, uri string, maxAgeSeconds int) *httpResponse {
	hlpr := GetHelper()
	hdrs := NewHttpHeadersBuilder().
		Set("content-type", hlpr.GetMimetype(uri)).
		// ref: https://varvy.com/pagespeed/cache-control.html
		Set("cache-control", fmt.Sprintf("max-age=%d,public", maxAgeSeconds)).
		GetHttpHeaders()
	return NewHttpResponseWithHeaders(STATUS_OK, object.GetContent(), hdrs)
}

// Produce an HTTP response with custom headers
func NewHttpResponseWithHeaders(status HttpStatus, body *string, headers HttpHeadersIfc) *httpResponse {
	rb := NewHttpResponseBuilder().
		SetStatus(status).
		SetBody(body)
	if !headers.IsEmpty() {
		rb.SetHeaders(headers)
	}
	return rb.GetHttpResponse()
}

// Produce an HTTP redirect (TEMPORARY) response to the supplied URL
func NewHttpResponseRedirect(URL string) *httpResponse {
	hdrs := NewHttpHeadersBuilder().
		Set("location", URL).
		GetHttpHeaders()
	body := ""
	return NewHttpResponseWithHeaders(STATUS_TEMPORARY_REDIRECT, &body, hdrs)
}

// Produce an HTTP redirect (PERMANENT) response to the supplied URL
func NewHttpResponseRedirectPermanent(URL string) HttpResponseIfc {
	hdrs := NewHttpHeadersBuilder().
		Set("location", URL).
		GetHttpHeaders()
	body := ""
	return NewHttpResponseWithHeaders(STATUS_MOVED_PERMANENTLY, &body, hdrs)
}

// -------------------------------------------------------------------------------------------------
// HttpResponseIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *httpResponse) GetBinBody() *[]byte {
	if nil == r.bodybytes {
		if nil == r.bodystring {
			return nil
		}
		bb := []byte(*r.bodystring)
		r.bodybytes = &bb
	}
	return r.bodybytes
}

func (r *httpResponse) GetBody() *string {
	if nil == r.bodystring {
		if nil == r.bodybytes {
			return nil
		}
		bs := string((*r.bodybytes)[:])
		r.bodystring = &bs
	}
	return r.bodystring
}

func (r *httpResponse) GetStatus() HttpStatus {
	return r.status
}

func (r *httpResponse) GetProtocolVersion() ver.VersionIfc {
	return r.protocolVersion
}

func (r *httpResponse) GetHeaders() HttpHeadersIfc {
	return r.headers
}
