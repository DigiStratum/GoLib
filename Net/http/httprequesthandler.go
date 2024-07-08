package http

type RequestHandlerIfc interface {
	HandleRequest(request HttpRequestIfc) HttpResponseIfc
}

