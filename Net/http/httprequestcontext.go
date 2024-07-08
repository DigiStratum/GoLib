package http

type HttpRequestContextIfc interface {
	SetServerPath(serverPath string)
	GetServerPath() string
	SetModulePath(modulePath string)
	GetModulePath() string
	SetPrefixPath(prefixPath string)
	GetPrefixPath() string
	SetRequestId(requestId string)
	GetRequestId() string
}

// TODO: Add more interesting properties such as which User is logged
// in, which Account/Customer/Business/etc is being requested
type httpRequestContext struct {
	serverPath	string	// The path that the Server matched on
	modulePath	string	// The path that the Module matched on
	prefixPath	string	// ServerPath/ModulePath
	requestId	string	// UUID for this request
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHttpRequestContext() HttpRequestContextIfc {
	return &httpRequestContext{}
}

// -------------------------------------------------------------------------------------------------
// HttpRequestContextIfc Implementation
// -------------------------------------------------------------------------------------------------

func (r *httpRequestContext) SetServerPath(serverPath string) {
	r.serverPath = serverPath
}

func (r *httpRequestContext) GetServerPath() string {
	return r.serverPath
}

func (r *httpRequestContext) SetModulePath(modulePath string) {
	r.modulePath = modulePath
}

func (r *httpRequestContext) GetModulePath() string {
	return r.modulePath
}

func (r *httpRequestContext) SetPrefixPath(prefixPath string) {
	r.prefixPath = prefixPath
}

func (r *httpRequestContext) GetPrefixPath() string {
	return r.prefixPath
}

func (r *httpRequestContext) SetRequestId(requestId string) {
	r.requestId = requestId
}

func (r *httpRequestContext) GetRequestId() string {
	return r.requestId
}

