package restapi

// TODO: Add support for builder pattern, or chaining, or "functional options":
// ref: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
// ref: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

import(
	"fmt"
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
type httpBodyData map[string][]string

type HttpRequest struct {
	protocol	string
	host		string
	remoteAddr	string
	scheme		string
	url		*url.URL
	method		string
	uri		string
	queryString	string
	headers		HttpHeaders
	body		string
	bodyData	httpBodyData
	context		HttpRequestContext
}

type HttpStatus int

type HttpResponse struct {
	status		HttpStatus
	headers		HttpHeaders
	body		string
}

// -------------------------------------------
// HTTP Request

func NewRequest() HttpRequest {
	return HttpRequest{
		headers: make(HttpHeaders),
		bodyData: make(httpBodyData),
		context: HttpRequestContext{},
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

func (request *HttpRequest) GetBody() string {
	return request.body
}

func (request *HttpRequest) SetBody(body string) {
	request.body = body
}

func (request *HttpRequest) GetBodyData() httpBodyData {
	return request.bodyData
}

func (request *HttpRequest) SetBodyData(bodyData httpBodyData) {
	request.bodyData = bodyData
}

// TODO: Add helpers for Body Data to build it up one name/value at a time

// Get the Request Context
func (request *HttpRequest) GetContext() *HttpRequestContext {
	return &request.context
}

// Set the Request Context
func (request *HttpRequest) SetContext(context HttpRequestContext) {
	request.context = context
}

// Get the Request Headers
func (request *HttpRequest) GetHeaders() *HttpHeaders {
	return &request.headers
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
func (hdrs *HttpHeaders) Merge(headers HttpHeaders) {
	for name, value := range headers {
		(*hdrs)[name] = value
	}
}

// -------------------------------------------
// HTTP Request Context

func (context *HttpRequestContext) SetServerPath(serverPath string) {
	context.serverPath = serverPath
}

func (context *HttpRequestContext) GetServerPath() string {
	return context.serverPath
}

func (context *HttpRequestContext) SetModulePath(modulePath string) {
	context.modulePath = modulePath
}

func (context *HttpRequestContext) GetModulePath() string {
	return context.modulePath
}

func (context *HttpRequestContext) SetPrefixPath(prefixPath string) {
	context.prefixPath = prefixPath
}

func (context *HttpRequestContext) GetPrefixPath() string {
	return context.prefixPath
}

func (context *HttpRequestContext) SetRequestId(requestId string) {
	context.requestId = requestId
}

func (context *HttpRequestContext) GetRequestId() string {
	return context.requestId
}

// -------------------------------------------
// HTTP Response

func NewResponse() HttpResponse {
	res := HttpResponse{}
	res.headers = make(HttpHeaders)
	return res
}

func (response *HttpResponse) GetBody() string {
	return response.body
}

func (response *HttpResponse) SetBody(body string) {
	response.body = body
}

func (response *HttpResponse) GetStatus() HttpStatus {
	return response.status
}

func (response *HttpResponse) SetStatus(status HttpStatus) {
	response.status = status
}

// Get the Response Headers
func (response *HttpResponse) GetHeaders() HttpHeaders {
	return response.headers
}

