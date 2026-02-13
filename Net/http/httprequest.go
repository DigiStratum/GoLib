package http

/*

An Immutable HTTP Request structure and programmatic interface to it.

Use NewHttpRequestBuilder() to create one of these.

FIXME:
 * PathParameters support is incomplete; we need to specify some sort of pattern that cna then be
   used to match the path parameters to the request URI.

TODO:
  * Get rid of all the URL component operations, just use GetURL(), but keep QueryString conveniences
*/

import (
	"net/url"

	"github.com/DigiStratum/GoLib/Data/metadata"
)

// RequestContextIfc provides request-scoped metadata (request ID, paths, identity, etc.)
// This is distinct from Go's context.Context - it's for HTTP request metadata.
type RequestContextIfc interface {
	GetRequestId() string
	SetRequestId(requestId string)
	GetServerPath() string
	SetServerPath(serverPath string)
	GetModulePath() string
	SetModulePath(modulePath string)
	GetPrefixPath() string
	SetPrefixPath(prefixPath string)
	GetModuleId() string
	SetModuleId(moduleId string)
	// Identity stores the authenticated identity for the request (if any)
	// The value should implement an identity interface from the auth layer
	GetIdentity() interface{}
	SetIdentity(identity interface{})
}

// requestContext is the concrete implementation of RequestContextIfc
type requestContext struct {
	requestId  string
	serverPath string
	modulePath string
	prefixPath string
	moduleId   string
	identity   interface{}
}

// NewRequestContext creates a new request context
func NewRequestContext() RequestContextIfc {
	return &requestContext{}
}

func (r *requestContext) GetRequestId() string       { return r.requestId }
func (r *requestContext) SetRequestId(id string)     { r.requestId = id }
func (r *requestContext) GetServerPath() string      { return r.serverPath }
func (r *requestContext) SetServerPath(p string)     { r.serverPath = p }
func (r *requestContext) GetModulePath() string      { return r.modulePath }
func (r *requestContext) SetModulePath(p string)     { r.modulePath = p }
func (r *requestContext) GetPrefixPath() string      { return r.prefixPath }
func (r *requestContext) SetPrefixPath(p string)     { r.prefixPath = p }
func (r *requestContext) GetModuleId() string        { return r.moduleId }
func (r *requestContext) SetModuleId(id string)      { r.moduleId = id }
func (r *requestContext) GetIdentity() interface{}   { return r.identity }
func (r *requestContext) SetIdentity(id interface{}) { r.identity = id }

type HttpRequestIfc interface {
	GetHost() string
	GetScheme() string
	GetURL() string
	GetURI() string
	GetMethod() HttpRequestMethod
	GetQueryString() string
	GetQueryParameters() metadata.MetadataIfc
	GetPathParameters() metadata.MetadataIfc
	SetPathParameters(params metadata.MetadataIfc)
	GetBody() *string
	GetBodyData() *httpRequestBody
	GetHeaders() *httpHeaders
	GetBuilder() *httpRequestBuilder
	// Request context for metadata (request ID, paths, etc.)
	GetContext() RequestContextIfc
	SetContext(ctx RequestContextIfc)
}

type httpRequest struct {
	url         *url.URL
	method      HttpRequestMethod
	queryParams metadata.MetadataIfc
	headers     *httpHeaders
	body        *string
	bodyData    *httpRequestBody
	pathParams  metadata.MetadataIfc
	context     RequestContextIfc
}

// -------------------------------------------------------------------------------------------------
// HttpRequestIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *httpRequest) GetHost() string {
	return r.url.Host
}

func (r *httpRequest) GetScheme() string {
	return r.url.Scheme
}

func (r *httpRequest) GetURL() string {
	if nil == r.url {
		return ""
	}
	return r.url.String()
}

func (r *httpRequest) GetMethod() HttpRequestMethod {
	return r.method
}

func (r *httpRequest) GetURI() string {
	return r.url.Path
}

func (r *httpRequest) GetQueryString() string {
	return r.url.RawQuery
}

func (r *httpRequest) GetBody() *string {
	return r.body
}

func (r *httpRequest) GetBodyData() *httpRequestBody {
	return r.bodyData
}

// Get the Request Headers
func (r *httpRequest) GetHeaders() *httpHeaders {
	return r.headers
}

// Get the path parameters (should be endpoint implementation)
func (r *httpRequest) GetPathParameters() metadata.MetadataIfc {
	return r.pathParams
}

// Set the path parameters (used by endpoint wrappers to populate extracted values)
func (r *httpRequest) SetPathParameters(params metadata.MetadataIfc) {
	r.pathParams = params
}

// Get the query parameters
func (r *httpRequest) GetQueryParameters() metadata.MetadataIfc {
	return r.queryParams
}

func (r *httpRequest) GetBuilder() *httpRequestBuilder {
	builder := NewHttpRequestBuilder(
		r.method,
		r.url.String(),
	).
		SetQueryParameters(r.queryParams).
		SetHeaders(r.headers).
		SetBody(r.body).
		SetBodyData(r.bodyData).
		SetContext(r.context)

	return builder
}

// Get the request context (request ID, paths, etc.)
func (r *httpRequest) GetContext() RequestContextIfc {
	// Lazily initialize context if nil
	if r.context == nil {
		r.context = NewRequestContext()
	}
	return r.context
}

// Set the request context
func (r *httpRequest) SetContext(ctx RequestContextIfc) {
	r.context = ctx
}
