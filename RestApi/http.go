package restapi

// TODO: Refactor into separate files for each of the data structures/method collections

// TODO: Add support for builder pattern, or chaining, or "functional options":
// ref: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
// ref: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

import(
	"fmt"
	"strings"
	"net/url"

	lib "github.com/DigiStratum/GoLib"
)

// TODO: Add more interesting properties such as which User is logged
// in, which Account/Customer/Business/etc is being requested
type HttpRequestContext struct {
	serverPath	string	// The path that the Server matched on
	modulePath	string	// The path that the Module matched on
	prefixPath	string	// ServerPath/ModulePath
	requestId	string	// UUID for this request
}

// Name/value pair header map for Request or Response
type HttpHeaders map[string]string

// Name/value set for HTTP response data for Request (typically form-encoded data)
type HttpBodyData map[string][]string

type HttpRequest struct {
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
}

type HttpStatus int

type HttpResponse struct {
	status		HttpStatus
	headers		*HttpHeaders
	body		*string
}

// -------------------------------------------
// HTTP Request

func NewRequest() *HttpRequest {
	hdrs := make(HttpHeaders)
	bodyData := make(HttpBodyData)
	return &HttpRequest{
		headers: &hdrs,
		bodyData: &bodyData,
		context: NewHttpRequestContext(),
	}
}

func (request *HttpRequest) GetProtocol() string {
	return request.protocol
}

func (request *HttpRequest) SetProtocol(protocol string) {
	request.protocol = protocol
}

func (request *HttpRequest) GetHost() string {
	return request.host
}

func (request *HttpRequest) SetHost(host string) {
	request.host = host
}

func (request *HttpRequest) GetRemoteAddr() string {
	return request.remoteAddr
}

func (request *HttpRequest) SetRemoteAddr(remoteAddr string) {
	request.remoteAddr = remoteAddr
}

func (request *HttpRequest) GetScheme() string {
	return request.scheme
}

func (request *HttpRequest) SetScheme(scheme string) {
	request.scheme = scheme
}

func (request *HttpRequest) GetURL() string {
	if nil == request.url { return "" }
	return request.url.String()
}

func (request *HttpRequest) SetURL(urlStr string) {
	u, err := url.Parse(urlStr)
	if nil != err {
		l := lib.GetLogger()
		l.Warn(fmt.Sprintf("HttpRequest.SetUrl() - failed to parse as a URL: '%s'", u))
	}
	request.url = u
}

func (request *HttpRequest) GetMethod() string {
	return request.method
}

func (request *HttpRequest) SetMethod(method string) {
	request.method = method
}

func (request *HttpRequest) GetURI() string {
	return request.uri
}

func (request *HttpRequest) SetURI(uri string) {
	request.uri = uri
}

func (request *HttpRequest) GetQueryString() string {
	return request.queryString
}

func (request *HttpRequest) SetQueryString(queryString string) {
	request.queryString = queryString
}

func (request *HttpRequest) GetBody() *string {
	return request.body
}

func (request *HttpRequest) SetBody(body *string) {
	request.body = body
}

func (request *HttpRequest) GetBodyData() *HttpBodyData {
	return request.bodyData
}

func (request *HttpRequest) SetBodyData(bodyData *HttpBodyData) {
	request.bodyData = bodyData
}

// TODO: Add helpers for Body Data to build it up one name/value at a time

// Get the Request Context
func (request *HttpRequest) GetContext() *HttpRequestContext {
	return request.context
}

// Set the Request Context
func (request *HttpRequest) SetContext(context *HttpRequestContext) {
	request.context = context
}

// Get the Request Headers
func (request *HttpRequest) GetHeaders() *HttpHeaders {
	return request.headers
}

// Extract a list of languages from the Accept-Language header (if any)
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
func (request *HttpRequest) GetAcceptableLanguages() *[]string {
	languages := request.getWeightedHeaderList("Accept-Language")
	// TODO: filter results according to what we support
	// (i.e. code, code-locale, code-locale-orthography); remove orthography/anything after locale
	// TODO: convert "*" into "default"
	return languages
}

// Extract a list of values from headers, ordered by preference expressed as quality value
// ref: https://developer.mozilla.org/en-US/docs/Glossary/Quality_values
func (request *HttpRequest) getWeightedHeaderList(headerName string) *[]string {

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

func (request *HttpRequest) IsIdempotentMethod() bool {
	if request.method == "post" { return false }
	return true
}

// -------------------------------------------
// HTTP Headers

// Get a single header
func (hdrs *HttpHeaders) Get(name string) string {
	if value, ok := (*hdrs)[name]; ok { return value }
	return ""
}

// Set a single header
func (hdrs *HttpHeaders) Set(name string, value string) {
	(*hdrs)[name] = value
}

// Merge an HttpHeaders set into our own
func (hdrs *HttpHeaders) Merge(headers *HttpHeaders) {
	for name, value := range *headers {
		(*hdrs)[name] = value
	}
}

// -------------------------------------------
// HTTP Request Context

func NewHttpRequestContext() *HttpRequestContext {
	return &HttpRequestContext{}
}

func (ctx *HttpRequestContext) SetServerPath(serverPath string) {
	ctx.serverPath = serverPath
}

func (ctx *HttpRequestContext) GetServerPath() string {
	return ctx.serverPath
}

func (ctx *HttpRequestContext) SetModulePath(modulePath string) {
	ctx.modulePath = modulePath
}

func (ctx *HttpRequestContext) GetModulePath() string {
	return ctx.modulePath
}

func (ctx *HttpRequestContext) SetPrefixPath(prefixPath string) {
	ctx.prefixPath = prefixPath
}

func (ctx *HttpRequestContext) GetPrefixPath() string {
	return ctx.prefixPath
}

func (ctx *HttpRequestContext) SetRequestId(requestId string) {
	ctx.requestId = requestId
}

func (ctx *HttpRequestContext) GetRequestId() string {
	return ctx.requestId
}

// -------------------------------------------
// HTTP Response

func NewResponse() *HttpResponse {
	hdrs := make(HttpHeaders)
	res := HttpResponse{
		headers: &hdrs,
	}
	return &res
}

func (response *HttpResponse) GetBody() *string {
	return response.body
}

func (response *HttpResponse) SetBody(body *string) {
	response.body = body
}

func (response *HttpResponse) GetStatus() HttpStatus {
	return response.status
}

func (response *HttpResponse) SetStatus(status HttpStatus) {
	response.status = status
}

// Get the Response Headers
func (response *HttpResponse) GetHeaders() *HttpHeaders {
	return response.headers
}

