package http

/*
An HTTP Request structure and programmatic interface to it.

TODO:
 * Add support for builder pattern, or chaining, or "functional options":
 * Split setters off into a builder so that the Request becomes immutable, getters-only

ref: https://www.calhoun.io/using-functional-options-instead-of-method-chaining-in-go/
ref: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

*/

import(
	"strings"
	"net/url"

	log "github.com/DigiStratum/GoLib/Logger"
	"github.com/DigiStratum/GoLib/Data/metadata"
)

// Http Request public interface
type HttpRequestIfc interface {
	GetProtocol() string
	SetProtocol(protocol string)

	GetHost() string
	SetHost(host string)

	GetRemoteAddr() string
	SetRemoteAddr(remoteAddr string)

	GetScheme() string
	IsIdempotentMethod() bool

	GetURL() string
	SetURL(urlStr string)

	//GetMethod() string
	//SetMethod(method string)
	GetMethod() HttpRequestMethod
	SetMethod(method HttpRequestMethod)

	GetURI() string
	SetURI(uri string)

	GetQueryString() string
	SetQueryString(queryString string)

	SetQueryParameters(params metadata.MetadataIfc)
	GetQueryParameters() metadata.MetadataIfc

	GetBody() *string
	SetBody(body *string)

	GetBodyData() *HttpBodyData
	SetBodyData(bodyData *HttpBodyData)

	GetContext() HttpRequestContextIfc
	SetContext(context HttpRequestContextIfc)

	GetHeaders() HttpHeadersIfc
	SetHeaders(headers HttpHeadersIfc)

	GetAcceptableLanguages() *[]string

	getWeightedHeaderList(headerName string) *[]string

	SetPathParameters(params metadata.MetadataIfc)
	GetPathParameters() metadata.MetadataIfc
}

type httpRequest struct {
	protocol	string
	host		string
	remoteAddr	string
	scheme		string
	url		*url.URL
	//method		string
	method		HttpRequestMethod
	uri		string
	queryString	string
	queryParams	metadata.MetadataIfc
	headers		HttpHeadersIfc
	body		*string
	bodyData	*HttpBodyData
	context		HttpRequestContextIfc
	pathParams	metadata.MetadataIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpRequest() *httpRequest {
	bodyData := make(HttpBodyData)
	return &httpRequest{
		headers: NewHttpHeaders(),
		bodyData: &bodyData,
		context: NewHttpRequestContext(),
	}
}

// -------------------------------------------------------------------------------------------------
// HttpRequestIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *httpRequest) GetProtocol() string {
	return r.protocol
}

func (r *httpRequest) SetProtocol(protocol string) {
	r.protocol = protocol
}

func (r *httpRequest) GetHost() string {
	return r.host
}

func (r *httpRequest) SetHost(host string) {
	r.host = host
}

func (r *httpRequest) GetRemoteAddr() string {
	return r.remoteAddr
}

func (r *httpRequest) SetRemoteAddr(remoteAddr string) {
	r.remoteAddr = remoteAddr
}

func (r *httpRequest) GetScheme() string {
	return r.scheme
}

func (r *httpRequest) SetScheme(scheme string) {
	r.scheme = scheme
}

func (r *httpRequest) GetURL() string {
	if nil == r.url { return "" }
	return r.url.String()
}

func (r *httpRequest) SetURL(urlStr string) {
	u, err := url.Parse(urlStr)
	if nil != err {
		log.GetLogger().Warn("HttpRequest.SetUrl() - failed to parse as a URL: '%s'", u)
		return
	}
	r.url = u
}

func (r *httpRequest) GetMethod() HttpRequestMethod {
	// Force result to be one of the enumerated set
	hlpr := GetHelper()
	return hlpr.GetHttpRequestMethod(hlpr.GetHttpRequestMethodText(r.method))
}

func (r *httpRequest) SetMethod(method HttpRequestMethod) {
	// Force value to be one of the enumerated set
	hlpr := GetHelper()
	r.method = hlpr.GetHttpRequestMethod(hlpr.GetHttpRequestMethodText(method))
}

func (r *httpRequest) GetURI() string {
	return r.uri
}

func (r *httpRequest) SetURI(uri string) {
	r.uri = uri
}

func (r *httpRequest) GetQueryString() string {
	return r.queryString
}

func (r *httpRequest) SetQueryString(queryString string) {
	r.queryString = queryString
}

func (r *httpRequest) GetBody() *string {
	return r.body
}

func (r *httpRequest) SetBody(body *string) {
	r.body = body
}

func (r *httpRequest) GetBodyData() *HttpBodyData {
	return r.bodyData
}

func (r *httpRequest) SetBodyData(bodyData *HttpBodyData) {
	r.bodyData = bodyData
}

// TODO: Add helpers for Body Data to build it up one name/value at a time

// Get the Request Context
func (r *httpRequest) GetContext() HttpRequestContextIfc {
	return r.context
}

// Set the Request Context
func (r *httpRequest) SetContext(context HttpRequestContextIfc) {
	r.context = context
}

// Get the Request Headers
func (r *httpRequest) GetHeaders() HttpHeadersIfc {
	return r.headers
}

// Set the Request Headers
func (r *httpRequest) SetHeaders(headers HttpHeadersIfc) {
	r.headers = headers
}

// Extract a list of languages from the Accept-Language header (if any)
// ref: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
func (r *httpRequest) GetAcceptableLanguages() *[]string {
	languages := r.getWeightedHeaderList("Accept-Language")
	// TODO: filter results according to what we support
	// (i.e. code, code-locale, code-locale-orthography); remove orthography/anything after locale
	// TODO: convert "*" into "default"
	return languages
}

// Quick check if the current request is expected to be idempotent in implementation
func (r *httpRequest) IsIdempotentMethod() bool {
	if METHOD_POST == r.method { return false }
	return true
}

// Set the path parameters (should be endpointwrapper)
func (r *httpRequest) SetPathParameters(params metadata.MetadataIfc) {
	r.pathParams = params
}

// Get the path parameters (should be endpoint implementation)
func (r *httpRequest) GetPathParameters() metadata.MetadataIfc {
	return r.pathParams
}

// Set the query parameters
func (r *httpRequest) SetQueryParameters(params metadata.MetadataIfc) {
	r.queryParams = params
}

// Get the query parameters
func (r *httpRequest) GetQueryParameters() metadata.MetadataIfc {
	return r.queryParams
}

// -------------------------------------------------------------------------------------------------
// httpRequest Implementation
// -------------------------------------------------------------------------------------------------

// Extract a list of values from headers, ordered by preference expressed as quality value
// ref: https://developer.mozilla.org/en-US/docs/Glossary/Quality_values
func (r *httpRequest) getWeightedHeaderList(headerName string) *[]string {

	// Get the value of the header we're after
	headerValue := r.GetHeaders().Get(headerName)
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

