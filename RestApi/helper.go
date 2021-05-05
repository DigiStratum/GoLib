package restapi

import(
	"fmt"
	"sync"
	"strings"
	"mime"

	obj "github.com/DigiStratum/GoLib/Objects"
)

type helper struct { }

var instance *helper

// TODO: Make this the same as we do other "singletons"; this sync.Once thing may
// only be needed for multithreaded confusion that we don't need (or do we?)
// Get the singleton instance of our helper
func GetHelper() *helper {
	var once sync.Once
	once.Do(func() {
		instance = &helper{ }
	})
	return instance
}

// Produce an HTTP response with standard headers
func (hlpr *helper) Response(status HttpStatus, body *string, contentType string) *HttpResponse {
	hdrs := HttpHeaders{}
	hdrs.Set("content-type", contentType)
	return hlpr.ResponseWithHeaders(status, body, &hdrs)
}

// Produce an HTTP response, code only, no headers/body
func (hlpr *helper) ReponseCode(status HttpStatus) *HttpResponse {
	hdrs := HttpHeaders{}
	body := ""
	return hlpr.ResponseWithHeaders(status, &body, &hdrs)
}

// Produce an HTTP response, code and default status text, JSON format
func (hlpr *helper) ResponseSimpleJson(status HttpStatus) *HttpResponse {
	message := hlpr.GetHttpStatusText(status)
	var staticResponse string
	if hlpr.IsStatus2xx(status) {
		staticResponse = fmt.Sprintf("[ { \"msg\": \"%s\" } ]", message)
	} else {
		staticResponse = fmt.Sprintf("[ { \"error\": { \"msg\": \"%s\" } } ]", message)
	}
	return hlpr.Response(status, &staticResponse, "application/json")
}

// Produce an HTTP error response by HTTP status code only
func (hlpr *helper) ResponseError(status HttpStatus) *HttpResponse {
	hdrs := HttpHeaders{}
	hdrs.Set("content-type", "text/plain")
	body := hlpr.GetHttpStatusText(status)
	return hlpr.ResponseWithHeaders(status, &body, &hdrs)
}

// Produce an OK HTTP response with standard headers
func (hlpr *helper) ResponseOk(body *string, contentType string) *HttpResponse {
	return hlpr.Response(STATUS_OK, body, contentType)
}

// Produce an ERROR HTTP response with JSON message body and standard headers
func (hlpr *helper) ResponseErrorJson(status HttpStatus, message string) *HttpResponse {
	staticResponse := fmt.Sprintf("[ { \"error\": { \"msg\": \"%s\" } } ]", message)
	return hlpr.Response(status, &staticResponse, "application/json")
}

// Produce an HTTP response from an Object (200 OK)
func (hlpr *helper) ResponseObject(object *obj.Object, uri string) *HttpResponse {
	return hlpr.Response(STATUS_OK, object.GetContent(), hlpr.GetMimetype(uri))
}

// Produce an HTTP response from an Object (200 OK)
func (hlpr *helper) ResponseObjectCacheable(object *obj.Object, uri string, maxAgeSeconds int) *HttpResponse {
	hdrs := HttpHeaders{}
	hdrs.Set("content-type", hlpr.GetMimetype(uri))
	// ref: https://varvy.com/pagespeed/cache-control.html
	hdrs.Set("cache-control", fmt.Sprintf("max-age=%d,public", maxAgeSeconds))
	return hlpr.ResponseWithHeaders(STATUS_OK, object.GetContent(), &hdrs)
}

// Produce an HTTP response with custom headers
func (hlpr *helper) ResponseWithHeaders(status HttpStatus, body *string, headers *HttpHeaders) *HttpResponse {
	response := NewResponse()
	response.SetStatus(status)
	response.SetBody(body)
	if len(*headers) > 0 {
		hdrs := response.GetHeaders()
		hdrs.Merge(headers)
	}
	return response
}

// Produce an HTTP redirect (TEMPORARY) response to the supplied URL
func (hlpr *helper) ResponseRedirect(URL string) *HttpResponse {
	hdrs := HttpHeaders{}
	hdrs.Set("location", URL)
	body := ""
	return hlpr.ResponseWithHeaders(STATUS_TEMPORARY_REDIRECT, &body, &hdrs)
}

// Produce an HTTP redirect (PERMANENT) response to the supplied URL
func (hlpr *helper) ResponseRedirectPermanent(URL string) *HttpResponse {
	hdrs := HttpHeaders{}
	hdrs.Set("location", URL)
	body := ""
	return hlpr.ResponseWithHeaders(STATUS_MOVED_PERMANENTLY, &body, &hdrs)
}

// Scan over the body data and, for each unique name, scrub out any duplicates
// TODO: refactor or make some variant which creates a value SET instead of tracking the dupes
func (hlpr *helper) SingularizePostData(bodyData *HttpBodyData) map[string]string {
	var data = make(map[string]string)
	for name, values := range *bodyData {
		if len(values) > 0 { data[name] = values[0] }
	}
	return data
}

// ref: https://golang.org/pkg/mime/#TypeByExtension
func (hlpr *helper) GetMimetype(uri string) string {
	dotpos := strings.LastIndex(uri, ".")
	if -1 == dotpos {
		return "application/octet-stream"
	}
	extension := uri[dotpos:]
	mimetype := mime.TypeByExtension(extension)
	if "" == mimetype {
		return "application/octet-stream"
	}
	return mimetype
}

// Is the supplied HttpStatus a 2xx?
func (hlpr *helper) IsStatus2xx(httpStatus HttpStatus) bool {
	statusCode := hlpr.GetHttpStatusCode(httpStatus)
	return (hlpr.GetHttpStatusCode(STATUS_OK) <= statusCode) && (hlpr.GetHttpStatusCode(STATUS_MULTIPLE_CHOICES) > statusCode)
}

// Is the supplied HttpStatus a 3xx?
func (hlpr *helper) IsStatus3xx(httpStatus HttpStatus) bool {
	statusCode := hlpr.GetHttpStatusCode(httpStatus)
	return (hlpr.GetHttpStatusCode(STATUS_MULTIPLE_CHOICES) <= statusCode) && (hlpr.GetHttpStatusCode(STATUS_BAD_REQUEST) > statusCode)
}

// Is the supplied HttpStatus a 4xx?
func (hlpr *helper) IsStatus4xx(httpStatus HttpStatus) bool {
	statusCode := hlpr.GetHttpStatusCode(httpStatus)
	return (hlpr.GetHttpStatusCode(STATUS_BAD_REQUEST) <= statusCode) && (hlpr.GetHttpStatusCode(STATUS_INTERNAL_SERVER_ERROR) > statusCode)
}

// Is the supplied HttpStatus a 5xx?
func (hlpr *helper) IsStatus5xx(httpStatus HttpStatus) bool {
	statusCode := hlpr.GetHttpStatusCode(httpStatus)
	return hlpr.GetHttpStatusCode(STATUS_INTERNAL_SERVER_ERROR) <= statusCode
}

const (
	STATUS_CONTINUE	HttpStatus = iota
	STATUS_SWITCHING_PROTOCOLS
	STATUS_OK
	STATUS_CREATED
	STATUS_ACCEPTED
	STATUS_NON_AUTHORITATIVE_INFORMATION
	STATUS_NO_CONTENT
	STATUS_RESET_CONTENT
	STATUS_PARTIAL_CONTENT
	STATUS_MULTIPLE_CHOICES
	STATUS_MOVED_PERMANENTLY
	STATUS_FOUND
	STATUS_SEE_OTHER
	STATUS_NOT_MODIFIED
	STATUS_USE_PROXY
	STATUS_TEMPORARY_REDIRECT
	STATUS_BAD_REQUEST
	STATUS_UNAUTHORIZED
	STATUS_FORBIDDEN
	STATUS_NOT_FOUND
	STATUS_METHOD_NOT_ALLOWED
	STATUS_NOT_ACCEPTABLE
	STATUS_PROXY_AUTHENTICATION_REQUIRED
	STATUS_REQUEST_TIMEOUT
	STATUS_CONFLICT
	STATUS_GONE
	STATUS_LENGTH_REQUIRED
	STATUS_PRECONDITION_FAILED
	STATUS_REQUEST_ENTITY_TOO_LARGE
	STATUS_REQUEST_URI_TOO_LONG
	STATUS_UNSUPPORTED_MEDIA_TYPE
	STATUS_REQUESTED_RANGE_NOT_SATISFIABLE
	STATUS_EXPECTATION_FAILED
	STATUS_INTERNAL_SERVER_ERROR
	STATUS_NOT_IMPLEMENTED
	STATUS_BAD_GATEWAY
	STATUS_SERVICE_UNAVAILABLE
	STATUS_GATEWAY_TIMEOUT
	STATUS_HTTP_VERSION_NOT_SUPPORTED
)

func (hlpr *helper) GetHttpStatusCode(httpStatus HttpStatus) int {
	switch (httpStatus) {
		case STATUS_CONTINUE:				return 100
		case STATUS_SWITCHING_PROTOCOLS:		return 101
		case STATUS_OK:					return 200
		case STATUS_CREATED:				return 201
		case STATUS_ACCEPTED:				return 202
		case STATUS_NON_AUTHORITATIVE_INFORMATION:	return 203
		case STATUS_NO_CONTENT:				return 204
		case STATUS_RESET_CONTENT:			return 205
		case STATUS_PARTIAL_CONTENT:			return 206
		case STATUS_MULTIPLE_CHOICES:			return 300
		case STATUS_MOVED_PERMANENTLY:			return 301
		case STATUS_FOUND:				return 302
		case STATUS_SEE_OTHER:				return 303
		case STATUS_NOT_MODIFIED:			return 304
		case STATUS_USE_PROXY:				return 305
		case STATUS_TEMPORARY_REDIRECT:			return 307
		case STATUS_BAD_REQUEST:			return 400
		case STATUS_UNAUTHORIZED:			return 401
		case STATUS_FORBIDDEN:				return 403
		case STATUS_NOT_FOUND:				return 404
		case STATUS_METHOD_NOT_ALLOWED:			return 405
		case STATUS_NOT_ACCEPTABLE:			return 406
		case STATUS_PROXY_AUTHENTICATION_REQUIRED:	return 407
		case STATUS_REQUEST_TIMEOUT:			return 408
		case STATUS_CONFLICT:				return 409
		case STATUS_GONE:				return 410
		case STATUS_LENGTH_REQUIRED:			return 411
		case STATUS_PRECONDITION_FAILED:		return 412
		case STATUS_REQUEST_ENTITY_TOO_LARGE:		return 413
		case STATUS_REQUEST_URI_TOO_LONG:		return 414
		case STATUS_UNSUPPORTED_MEDIA_TYPE:		return 415
		case STATUS_REQUESTED_RANGE_NOT_SATISFIABLE:	return 416
		case STATUS_EXPECTATION_FAILED:			return 417
		case STATUS_INTERNAL_SERVER_ERROR:		return 500
		case STATUS_NOT_IMPLEMENTED:			return 501
		case STATUS_BAD_GATEWAY:			return 502
		case STATUS_SERVICE_UNAVAILABLE:		return 503
		case STATUS_GATEWAY_TIMEOUT:			return 504
		case STATUS_HTTP_VERSION_NOT_SUPPORTED:		return 505
	}
	return 0
}

func (hlpr *helper) GetHttpStatusText(httpStatus HttpStatus) string {
	switch (httpStatus) {
		case STATUS_CONTINUE:				return "CONTINUE"
		case STATUS_SWITCHING_PROTOCOLS:		return "SWITCHING PROTOCOLS"
		case STATUS_OK:					return "OK"
		case STATUS_CREATED:				return "CREATED"
		case STATUS_ACCEPTED:				return "ACCEPTED"
		case STATUS_NON_AUTHORITATIVE_INFORMATION:	return "NON_AUTHORITATIVE INFORMATION"
		case STATUS_NO_CONTENT:				return "NO CONTENT"
		case STATUS_RESET_CONTENT:			return "RESET CONTENT"
		case STATUS_PARTIAL_CONTENT:			return "PARTIAL CONTENT"
		case STATUS_MULTIPLE_CHOICES:			return "MULTIPLE CHOICES"
		case STATUS_MOVED_PERMANENTLY:			return "MOVED PERMANENTLY"
		case STATUS_FOUND:				return "FOUND"
		case STATUS_SEE_OTHER:				return "SEE OTHER"
		case STATUS_NOT_MODIFIED:			return "NOT MODIFIED"
		case STATUS_USE_PROXY:				return "USE PROXY"
		case STATUS_TEMPORARY_REDIRECT:			return "TEMPORARY REDIRECT"
		case STATUS_BAD_REQUEST:			return "BAD REQUEST"
		case STATUS_UNAUTHORIZED:			return "UNAUTHORIZED"
		case STATUS_FORBIDDEN:				return "FORBIDDEN"
		case STATUS_NOT_FOUND:				return "NOT FOUND"
		case STATUS_METHOD_NOT_ALLOWED:			return "METHOD NOT ALLOWED"
		case STATUS_NOT_ACCEPTABLE:			return "NOT ACCEPTABLE"
		case STATUS_PROXY_AUTHENTICATION_REQUIRED:	return "PROXY AUTHENTICATION REQUIRED"
		case STATUS_REQUEST_TIMEOUT:			return "REQUEST TIMEOUT"
		case STATUS_CONFLICT:				return "CONFLICT"
		case STATUS_GONE:				return "GONE"
		case STATUS_LENGTH_REQUIRED:			return "LENGTH REQUIRED"
		case STATUS_PRECONDITION_FAILED:		return "PRECONDITION FAILED"
		case STATUS_REQUEST_ENTITY_TOO_LARGE:		return "REQUEST ENTITY TOO LARGE"
		case STATUS_REQUEST_URI_TOO_LONG:		return "REQUEST URI TOO LONG"
		case STATUS_UNSUPPORTED_MEDIA_TYPE:		return "UNSUPPORTED MEDIA TYPE"
		case STATUS_REQUESTED_RANGE_NOT_SATISFIABLE:	return "REQUESTED RANGE NOT SATISFIABLE"
		case STATUS_EXPECTATION_FAILED:			return "EXPECTATION FAILED"
		case STATUS_INTERNAL_SERVER_ERROR:		return "INTERNAL SERVER ERROR"
		case STATUS_NOT_IMPLEMENTED:			return "NOT IMPLEMENTED"
		case STATUS_BAD_GATEWAY:			return "BAD GATEWAY"
		case STATUS_SERVICE_UNAVAILABLE:		return "SERVICE UNAVAILABLE"
		case STATUS_GATEWAY_TIMEOUT:			return "GATEWAY TIMEOUT"
		case STATUS_HTTP_VERSION_NOT_SUPPORTED:		return "HTTP VERSION NOT SUPPORTED"
	}
	return "UNKNOWN STATUS CODE"
}

