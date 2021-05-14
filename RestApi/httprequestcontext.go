package restapi

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

// Make a new one of these!
func NewhttpRequestContext() httpRequestContextIfc {
	return &httpRequestContext{}
}

func (ctx *httpRequestContext) SetServerPath(serverPath string) {
	ctx.serverPath = serverPath
}

func (ctx *httpRequestContext) GetServerPath() string {
	return ctx.serverPath
}

func (ctx *httpRequestContext) SetModulePath(modulePath string) {
	ctx.modulePath = modulePath
}

func (ctx *httpRequestContext) GetModulePath() string {
	return ctx.modulePath
}

func (ctx *httpRequestContext) SetPrefixPath(prefixPath string) {
	ctx.prefixPath = prefixPath
}

func (ctx *httpRequestContext) GetPrefixPath() string {
	return ctx.prefixPath
}

func (ctx *httpRequestContext) SetRequestId(requestId string) {
	ctx.requestId = requestId
}

func (ctx *httpRequestContext) GetRequestId() string {
	return ctx.requestId
}

