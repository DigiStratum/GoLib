package restapi

// TODO: Add support for builder pattern, or chaining, or "functional options":
// ref: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
// ref: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

import(
	"fmt"
	"strings"
	"net/url"

	lib "github.com/DigiStratum/GoLib"
)

// Http Request public interface
type HttpRequestIfc struct {
	GetProtocol() string
	SetProtocol(protocol string)
	GetHost() string
	SetHost(host string)
	GetRemoteAddr() string
	SetRemoteAddr(remoteAddr string)
	GetScheme() string
	GetURL() string
	SetURL(urlStr string)
	GetMethod() string
	SetMethod(method string)
	GetURI() string
	SetURI(uri string)
	GetQueryString() string
	SetQueryString(queryString string)
	GetBody() *string
	SetBody(body *string)
	GetBodyData() *HttpBodyData
	SetBodyData(bodyData *HttpBodyData)
	GetContext() *HttpRequestContext
	SetContext(context *HttpRequestContext)
	GetHeaders() *HttpHeaders
	GetAcceptableLanguages() *[]string
	getWeightedHeaderList(headerName string) *[]string
	IsIdempotentMethod() bool
	SetPathParameters(params *lib.HashMap)
	GetPathParameters() *lib.HashMap
}

type httpRequest struct {
	protocol	string
	host		string
	remoteAddr	string
	scheme		string
	url		*url.URL
	method		string
	uri		string
	queryString	string
	headers		*HttpHeaders
	body		*string
	bodyData	*HttpBodyData
	context		*HttpRequestContext
	pathParams	*lib.HashMap
}

// Make a new one of these
func NewRequest() HttpRequestIfc {
	hdrs := make(HttpHeaders)
	bodyData := make(HttpBodyData)
	return &httpRequest{
		headers: &hdrs,
		bodyData: &bodyData,
		context: NewHttpRequestContext(),
	}
}

// -------------------------------------------------------------------------------------------------
// HttpRequestIfc

func (request *httpRequest) GetProtocol() string {
	return request.protocol
}

func (request *httpRequest) SetProtocol(protocol string) {
	request.protocol = protocol
}

func (request *httpRequest) GetHost() string {
	return request.host
}

func (request *httpRequest) SetHost(host string) {
	request.host = host
}

func (request *httpRequest) GetRemoteAddr() string {
	return request.remoteAddr
}

func (request *httpRequest) SetRemoteAddr(remoteAddr string) {
	request.remoteAddr = remoteAddr
}

func (request *httpRequest) GetScheme() string {
	return request.scheme
}

func (request *httpRequest) SetScheme(scheme string) {
	request.scheme = scheme
}

func (request *httpRequest) GetURL() string {
	if nil == request.url { return "" }
	return request.url.String()
}

func (request *httpRequest) SetURL(urlStr string) {
	u, err := url.Parse(urlStr)
	if nil != err {
		l := lib.GetLogger()
		l.Warn(fmt.Sprintf("HttpRequest.SetUrl() - failed to parse as a URL: '%s'", u))
	}
	request.url = u
}

func (request *httpRequest) GetMethod() string {
	return request.method
}

func (request *httpRequest) SetMethod(method string) {
	request.method = method
}

func (request *httpRequest) GetURI() string {
	return request.uri
}

func (request *httpRequest) SetURI(uri string) {
	request.uri = uri
}

func (request *httpRequest) GetQueryString() string {
	return request.queryString
}

func (request *httpRequest) SetQueryString(queryString string) {
	request.queryString = queryString
}

func (request *httpRequest) GetBody() *string {
	return request.body
}

func (request *httpRequest) SetBody(body *string) {
	request.body = body
}

func (request *httpRequest) GetBodyData() *HttpBodyData {
	return request.bodyData
}

func (request *httpRequest) SetBodyData(bodyData *HttpBodyData) {
	request.bodyData = bodyData
}

// TODO: Add helpers for Body Data to build it up one name/value at a time

// Get the Request Context
func (request *httpRequest) GetContext() *HttpRequestContext {
	return request.context
}

// Set the Request Context
func (request *httpRequest) SetContext(context *HttpRequestContext) {
	request.context = context
}

// Get the Request Headers
func (request *httpRequest) GetHeaders() *HttpHeaders {
	return request.headers
}

// Extract a list of languages from the Accept-Language header (if any)
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
func (request *httpRequest) GetAcceptableLanguages() *[]string {
	languages := request.getWeightedHeaderList("Accept-Language")
	// TODO: filter results according to what we support
	// (i.e. code, code-locale, code-locale-orthography); remove orthography/anything after locale
	// TODO: convert "*" into "default"
	return languages
}

// Quick check if the current request is expected to be idempotent in implementation
func (request *httpRequest) IsIdempotentMethod() bool {
	if request.method == "post" { return false }
	return true
}

// Set the path parameters (should be endpointwrapper)
func (request *httpRequest) SetPathParameters(params *lib.HashMap) {
	request.pathParams = params
}

// Get the path parameters (should be endpoint implementation)
func (request *httpRequest) GetPathParameters() *lib.HashMap {
	return request.pathParams
}

// -------------------------------------------------------------------------------------------------
// Private implementation

// Extract a list of values from headers, ordered by preference expressed as quality value
// ref: https://developer.mozilla.org/en-US/docs/Glossary/Quality_values
func (request *httpRequest) getWeightedHeaderList(headerName string) *[]string {

	// Get the value of the header we're after
	headerValue := request.GetHeaders().Get(headerName)
	if len(headerValue)  == 0 {
		// no header, no list!
		values := make([]string, 0)
		return &values
	}

	// Split the header on "," in case it has multiple values
	headerValues := strings.Split(headerValue, ",")
	values := make([]string, len(headerValues))

	// For each value we found...
	kept := 0
	for _, value := range headerValues {
		// If it has a ";", then there are extra details attached to split off
		if strings.Index(value, ";") > -1 {
			valueParts := strings.Split(value, ";")
			// Keep the first part, the rest is metadata
			values[kept] = valueParts[0]
			kept++
			// TODO: Check out the other parts; they should be formatted as "name=value"
			// TODO: If the name = "q" and the value is a float from 0.0 to 1.0, use it to sort
			// TODO: If ANY value has a "q" specified, then sort must be performed, assume 1.0 for unspecified q
		} else {
			// Keep this value
			values[kept] = value
			kept++
		}
	}
	return &values
}

