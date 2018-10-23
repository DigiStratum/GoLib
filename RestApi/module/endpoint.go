package module

/*

Our definition of an Endpoint is the implementation which responds to all supported HTTP requests
that map to the Endpoint's URL path. Which HTTP requests are supported are determined by which
HTTP request method verbs are implemented (i.e. get, post, etc). The Endpoint implementation
declares which methods it supports - this cannot be overridden. If the Endpoint does not support a
method and the server receives a request for that method against that Endpoint's mapped path, then
the server will automatically response with a 405 METHOD NOT SUPPORTED; the Endpoint implementation
does not need to support responses to unsupported methods.

Here's how the Endpoint Path maps out:

http://Host/HttpServerPath/ModulePath/EndpointPath

If the HttpServer is running in standalone mode, then the path map simplifies:

http://localhost:{port}/ModulePath/EndpointPath

Configuration data exists for every element in the path; each element provides its own defaults, but
is subject to override from the level above. For example, the Endpoint may map itself by default to
"/test", but the Module may override this and remap it to "/modulewins", and the server may further
override and remap it to "/serverwins".

Configuration data is simplified to generic name-value pairs as a string map. There is no support
for structured data as it would be unnecessarily complex to support this in a generalized way (for
now, anyway!)

The name+version of a Module implementation is a form of readable, globally unique ID (GUID). By
design, this allows us to load multiple Module implementations with the same name, but different
versions. Because only one Module may map to a given Path, the Server must reconcile which versioned
Module implementation is used to satisfy HTTP requests to that Path. The normal behavior should be
such that, by default, the newest version of the Module implementation is used. If the HTTP request
includes an X-Version header to specify a version other than the default, then that version will be
used instead if possible. If the requested Module version is not loaded into the Server, then it
will respond with a 501 NOT IMPLEMENTED.

Modules are built as Go "plugins". Go supports loading plugins which are unique by the name of the
file from which the plugin is loaded. Thus, to differentiate one odule version from another, we
build the plugin into a file which includes the name+version to facilitate loading more than one
plugin version for a given name. Modules are thus packaged collections of one or more Endpoints.

Endpoints themselves are compiled into a given Module and are not separately loaded as plugins. Thus
any version information supplied with an Endpoint is for the client to have visibility on the state
of the server, not to give the client a choice - that is done at the Module layer.

ref: https://www.w3.org/Protocols/rfc2616/rfc2616-sec9.html

*/

import (
	"fmt"
	"strings"
	"strconv"

	lib "../../golib"
	rest "../restpai"
)

var supportedMethods	[]string

func init() {
	supportedMethods = []string{ "get", "post", "put", "options", "head", "patch", "delete" }
}

// TODO: Apply these types to other properties below...
type EndpointMethod		string
type EndpointPattern		string
type EndpointVersion		string
type EndpointName		string

type AbstractEndpointIfc interface {
	Configure(serverConfig lib.Config, moduleConfig lib.Config)
	GetSecurityPolicy() *SecurityPolicy
	GetName() string
	GetVersion() string
	GetMethods() []string
	GetPattern() string
	SetPattern(pattern string)
	HandleRequest(request rest.HttpRequest, endpoint *AbstractEndpointIfc) *rest.HttpResponse
}

type AbstractEndpoint struct {
        serverConfig    lib.Config		// Server configuration copy
        moduleConfig    lib.Config		// Module configuration copy
        name            string			// Unique name of this Endpoint
        version         string			// Version of this Endpoint
        pattern         string			// Pattern which matches URI's to us (relative to Module)
        methods         []string		// List of HTTP request methods that we respond to
	securityPolicy	SecurityPolicy	// Security Policy for this Endpoint
}

// Initialize
func (ep *AbstractEndpoint) Init(endpoint AbstractEndpointIfc, name string, version string, pattern string) {

	// Capture basic properties
	ep.SetName(name)
	ep.SetVersion(version)
	ep.SetPattern(pattern)

	ep.securityPolicy = NewSecurityPolicy()
	ep.methods = []string{}

	// Find which methods this Endpoint actually implements
	l := lib.GetLogger()
	implementedMethods := make(map[string]bool)
	for _, method := range supportedMethods {
		implemented := false
		if implementsMethod(method, endpoint) {
			ep.methods = append(ep.methods, method)
			implemented = true
		}
		implementedMethods[method] = implemented
	}

	// If GET is implemented, but not HEAD, enable HEAD so we receive the default behavior
	if implementedMethods["get"] && !implementedMethods["head"] {
		ep.methods = append(ep.methods, "head")
	}

	l.Trace(fmt.Sprintf("Endpoint: Methods Implemented: [%s]", strings.Join(ep.methods, ",")))
}

// Does the supplied Endpoint implement the interface for the specified Method?
func implementsMethod(method string, endpoint interface{}) bool {
	switch (method) {
		case "get": if _, ok := endpoint.(AbstractGetEndpointIfc); ok { return true }
		case "post": if _, ok := endpoint.(AbstractPostEndpointIfc); ok { return true }
		case "put": if _, ok := endpoint.(AbstractPutEndpointIfc); ok { return true }
		case "options": if _, ok := endpoint.(AbstractOptionsEndpointIfc); ok { return true }
		case "head": if _, ok := endpoint.(AbstractHeadEndpointIfc); ok { return true }
		case "delete": if _, ok := endpoint.(AbstractPatchEndpointIfc); ok { return true }
		case "patch": if _, ok := endpoint.(AbstractDeleteEndpointIfc); ok { return true }
	}
	return false
}

// Capture the configuration data for this endpoint
func (ep *AbstractEndpoint) Configure(serverConfig lib.Config, moduleConfig lib.Config) {
	ep.serverConfig = serverConfig
	ep.moduleConfig = moduleConfig
}

// Endpoint needs to be able to access its own Security Policy
func (ep *AbstractEndpoint) GetSecurityPolicy() *SecurityPolicy {
	return &ep.securityPolicy
}

// Override the default path matching pattern
// This is used by Module to override our path in the case that we don't have a default defined;
// This is useful for endpoints which are generally useful in many Modules, and may need to be
// mapped differently, depending on the application
func (ep *AbstractEndpoint) SetPattern(pattern string) {
	// We only allow the Module to set our pattern if one is not already set
	if "" == ep.pattern {
		// TODO: Validate this somehow? Module is responsible for capturing this change for itself and passing on to us
		ep.pattern = pattern
		return
	}
	ident := fmt.Sprintf("Module: '%s', Endpoint: '%s'", ep.moduleConfig.Get("module.name"), ep.name)
	message := fmt.Sprintf("Cannot set pattern for Endpoint (%s) to (%s) as it is already set to (%s)", ident, pattern, ep.pattern)
	l := lib.GetLogger()
	l.Warn(message)
}

// Return our pattern
func (ep *AbstractEndpoint) GetPattern() string {
	return ep.pattern
}

// Return our version
func (ep *AbstractEndpoint) GetVersion() string {
	return ep.version
}

// Set the version
func (ep *AbstractEndpoint) SetVersion(version string) {
	ep.version = version
}

// Return our name
func (ep *AbstractEndpoint) GetName() string {
	return ep.name
}

// Set the name
func (ep *AbstractEndpoint) SetName(name string) {
	ep.name = name
}

// Return our list of methods
// This is used by the Controller to add us to the map to send us requests.
func (ep *AbstractEndpoint) GetMethods() []string {
	return ep.methods
}

// Request handler
// TODO: Pass around request as a pointer to minimize the memory copying for a potentially large data structure
func (ep *AbstractEndpoint) HandleRequest(request rest.HttpRequest, endpoint *AbstractEndpointIfc) *rest.HttpResponse {

	// Will our SecurityPolicy reject this Request?
	epsp := ep.GetSecurityPolicy()
	if rej := epsp.HandleRejection(request); nil != rej { return rej } // REJECT!

	method := request.GetMethod()
	l := lib.GetLogger()
	ctx := request.GetContext()
	l.Trace(fmt.Sprintf(
		"[%s] Endpoint (%s): Dispatching %s Request",
		ctx.GetRequestId(),
		ep.name,
		method,
	))
	// Note that checking requestMethod against ep.methods would be redundant
	// because Controller should already be doing this for us via ep.GetMethods()
	switch (method) {
		case "get": return handleGet(request, *endpoint)
		case "head": return handleHead(request, *endpoint)
		case "post": return handlePost(request, *endpoint)
		case "put": return handlePut(request, *endpoint)
		case "options": return handleOptions(request, *endpoint)
		case "delete": return handleDelete(request, *endpoint)
		case "patch": return handlePatch(request, *endpoint)
	}
	l.Error(fmt.Sprintf(
		"[%s] Endpoint (%s): Controller passed us a non-implemented Request Method '%s'",
		ctx.GetRequestId(),
		ep.name,
		method,
	))

	// Default response handling for a request we are not prepared to receive
	hlpr := rest.GetHelper()
	return hlpr.ResponseError(rest.STATUS_NOT_IMPLEMENTED)
}

// Default Options handler for endpoints
func (endpoint *AbstractEndpoint) HandleOptions(request rest.HttpRequest) *rest.HttpResponse {
	hdrs := rest.HttpHeaders{}
	hdrs.Set("allow", strings.Join(endpoint.methods, ","))
	hlpr := rest.GetHelper()
	return hlpr.ResponseWithHeaders(rest.STATUS_OK, "", hdrs)
}

type AbstractGetEndpointIfc interface {
	HandleGet(request rest.HttpRequest) *rest.HttpResponse
}

type AbstractPostEndpointIfc interface {
	HandlePost(request rest.HttpRequest) *rest.HttpResponse
}

type AbstractPutEndpointIfc interface {
	HandlePut(request rest.HttpRequest) *rest.HttpResponse
}

type AbstractOptionsEndpointIfc interface {
	HandleOptions(request rest.HttpRequest) *rest.HttpResponse
}

type AbstractHeadEndpointIfc interface {
	HandleHead(request rest.HttpRequest) *rest.HttpResponse
}

type AbstractDeleteEndpointIfc interface {
	HandleDelete(request rest.HttpRequest) *rest.HttpResponse
}

type AbstractPatchEndpointIfc interface {
	HandlePatch(request rest.HttpRequest) *rest.HttpResponse
}

func handleGet(request rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(AbstractGetEndpointIfc); ok {
		return handler.HandleGet(request)
	}
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement AbstractGetEndpointIfc!?",
		ctx.GetRequestId(),
	))
	return nil
}

func handlePost(request rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(AbstractPostEndpointIfc); ok {
		return handler.HandlePost(request)
	}
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement AbstractPostEndpointIfc!?",
		ctx.GetRequestId(),
	))
	return nil
}

func handlePut(request rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(AbstractPutEndpointIfc); ok {
		return handler.HandlePut(request)
	}
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement AbstractPutEndpointIfc!?",
		ctx.GetRequestId(),
	))
	return nil
}

func handleOptions(request rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(AbstractOptionsEndpointIfc); ok {
		return handler.HandleOptions(request)
	}
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement AbstractOptionsEndpointIfc!?",
		ctx.GetRequestId(),
	))
	return nil
}

func handleHead(request rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(AbstractHeadEndpointIfc); ok {
		return handler.HandleHead(request)
	}
	// The endpoint doesn't implement Head directly, but we can call GET and modify
	if handler, ok := ep.(AbstractGetEndpointIfc); ok {
		// Any endpoint with an expensive Get call should override this default
		// handling with something better tuned to skip expensive steps if possible
		response := handler.HandleGet(request)
		if nil != response {
			hdrs := response.GetHeaders()
			hdrs.Set("content-length", strconv.Itoa(len(response.GetBody())))
			response.SetBody("")
		}
		return response
	}

	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement AbstractHeadEndpointIfc!?",
		ctx.GetRequestId(),
	))
	return nil
}

func handleDelete(request rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(AbstractDeleteEndpointIfc); ok {
		return handler.HandleDelete(request)
	}
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement AbstractDeleteEndpointIfc!?",
		ctx.GetRequestId(),
	))
	return nil
}

func handlePatch(request rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(AbstractPatchEndpointIfc); ok {
		return handler.HandlePatch(request)
	}
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement AbstractPatchEndpointIfc!?",
		ctx.GetRequestId(),
	))
	return nil
}

