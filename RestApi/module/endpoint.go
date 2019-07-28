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

	lib "github.com/DigiStratum/GoLib"
	rest "github.com/DigiStratum/GoLib/RestApi"
)

var supportedMethods	[]string

func init() {
	supportedMethods = []string{ "get", "post", "put", "options", "head", "patch", "delete" }
}

type EndpointIfc interface {
	// TODO: Add error return value for Configure()
	Configure(concreteEndpoint interface{}, serverConfig lib.Config, moduleConfig lib.Config, extraConfig lib.Config)
	Init(concreteEndpoint interface{})
	GetSecurityPolicy() *SecurityPolicy
	GetName() string
	GetVersion() string
	GetPattern() string
	IsDefault() bool
	GetMethods() []string
	HandleRequest(request *rest.HttpRequest, endpoint EndpointIfc) *rest.HttpResponse
}

type Endpoint struct {
        serverConfig    *lib.Config	// Server configuration copy
        moduleConfig    *lib.Config	// Module configuration copy
	endpointConfig	*lib.Config	// Endpoint configuration
        name            string		// Unique name of this Endpoint
        version         string		// Version of this Endpoint
        pattern         string		// Pattern which matches URI's to us (relative to Module)
        methods         []string	// List of HTTP request methods that we respond to
	securityPolicy	*SecurityPolicy	// Security Policy for this Endpoint
	isDefault	bool		// Is this endpoint configured as a default?
}

// Make a new one of these (typically embedded as the superclass of some subclass)
func NewEndpoint(name string, version string) Endpoint {
	return Endpoint{
		name:		name,
		version:	version,
	}
}

// Initialize
// concreteEndpoint is a sub-class of Endpoint; it needs to be passed in for inspection because
// inspecting the super-class (Endpoint) will not expose the properties of the sub-class
// TODO: Support an error response? Needed to knock out the mapping? Is having no methods enough?
func (ep *Endpoint) Init(concreteEndpoint interface{}) {
	l := lib.GetLogger()
	l.Trace(fmt.Sprintf("Endpoint{%s}: Init()", ep.name))

	// Verify that concreteEndpoint implements EndpointIfc
	if _, ok := concreteEndpoint.(EndpointIfc); ! ok {
		l.Error(fmt.Sprintf(
			"Endpoint{%s}.Init(): Object supplied is not an EndpointIfc",
			ep.name,
		))
		return
	}

	// Capture basic properties
	ep.methods = []string{}

	// Find which methods this Endpoint actually implements
	implementedMethods := make(map[string]bool)
	for _, method := range supportedMethods {
		implemented := false
		if implementsMethod(method, concreteEndpoint) {
			ep.methods = append(ep.methods, method)
			implemented = true
		}
		implementedMethods[method] = implemented
	}

	// If GET is implemented, but not HEAD, enable HEAD so we receive the default behavior
	if implementedMethods["get"] && !implementedMethods["head"] {
		ep.methods = append(ep.methods, "head")
	}

	l.Trace(fmt.Sprintf(
		"Endpoint{%s}.Init(): Methods Implemented: [%s]",
		ep.name, strings.Join(ep.methods, ","),
	))
}

// Does the supplied Endpoint implement the interface for the specified Method?
func implementsMethod(method string, endpoint interface{}) bool {
	switch (method) {
		case "get": if _, ok := endpoint.(GetEndpointIfc); ok { return true }
		case "post": if _, ok := endpoint.(PostEndpointIfc); ok { return true }
		case "put": if _, ok := endpoint.(PutEndpointIfc); ok { return true }
		case "options": if _, ok := endpoint.(OptionsEndpointIfc); ok { return true }
		case "head": if _, ok := endpoint.(HeadEndpointIfc); ok { return true }
		case "delete": if _, ok := endpoint.(PatchEndpointIfc); ok { return true }
		case "patch": if _, ok := endpoint.(DeleteEndpointIfc); ok { return true }
	}
	return false
}

// Capture the configuration data for this endpoint
// Server/Module config passed by value (copy) to prevent tampering
func (ep *Endpoint) Configure(concreteEndpoint interface{}, serverConfig lib.Config, moduleConfig lib.Config, extraConfig lib.Config) {

	// Endpoint-specific Config properties have prefix: "endpoint.{Endpoint name}."
	configPrefix := fmt.Sprintf("endpoint.%s.", ep.name)

	l := lib.GetLogger()
	l.Trace(fmt.Sprintf("Endpoint{%s}.Configure(): Prefix is: '%s'", ep.name, configPrefix))
	ep.serverConfig = &serverConfig
	ep.moduleConfig = &moduleConfig

	// The Endpoint's Config is a subset of the Module Config
	ep.endpointConfig = moduleConfig.GetSubset(configPrefix)
	requiredConfig := []string{ "version", "pattern" }
	if ! (ep.endpointConfig.HasAll(&requiredConfig)) {
		l := lib.GetLogger()
		l.Error(fmt.Sprintf("Endpoint{%s}.Configure(): Incomplete Endpoint Config provided", ep.name))
		return
	}
	ep.endpointConfig.Set("name", ep.name) // Reflect name into Module Config for reference

	ep.pattern = ep.endpointConfig.Get("pattern")
	ep.securityPolicy = NewSecurityPolicy(ep.endpointConfig.GetSubset("auth"))

	// If this Endpoint is Configurable...
	if configurableEndpoint, ok := concreteEndpoint.(ConfigurableEndpointIfc); ok {
		// Hit the Configure method!
		configurableEndpoint.ConfigureEndpoint(ep.endpointConfig)
	} else {
		l.Trace(fmt.Sprintf("Endpoint{%s}.Configure(): Not a ConfigurableEndpoint", ep.name))
	}

        // See if there are any overrides for this Endpoint hiding in extra Module Config
        overrides := extraConfig.GetSubset(configPrefix)
        if ! overrides.IsEmpty() {
		l.Trace(fmt.Sprintf(
			"Endpoint{%s}.Configure(): Applying overrides from extra Module Config",
			ep.name,
		))
		overrides.Dump()
                ep.endpointConfig.Merge(overrides)
        }

	// See if this endpoint is configured as a default
	if ep.endpointConfig.Has("isdefault") {
		isDefault := ep.endpointConfig.Get("isdefault")
		ep.isDefault = (isDefault == "true")
	} else {
		ep.isDefault = false
	}
	l.Trace(fmt.Sprintf("Endpoint{%s}.Configure(): isDefault? %t", ep.name, ep.isDefault))
	l.Crazy(fmt.Sprintf(
		"Endpoint{%s} Configuration: %s",
		ep.name,
		ep.endpointConfig.DumpString(),
	));
}

// Endpoint needs to be able to access its own Security Policy
func (ep *Endpoint) GetSecurityPolicy() *SecurityPolicy {
	return ep.securityPolicy
}

// Return our pattern
func (ep *Endpoint) GetPattern() string {
	return ep.pattern
}

// Return our version
func (ep *Endpoint) GetVersion() string {
	return ep.version
}

// Return our name
func (ep *Endpoint) GetName() string {
	return ep.name
}

// Get our defaultness
func (ep *Endpoint) IsDefault() bool {
	return ep.isDefault
}

// Return our list of methods
// This is used by the Controller to add us to the map to send us requests.
func (ep *Endpoint) GetMethods() []string {
	return ep.methods
}

// Request handler
func (ep *Endpoint) HandleRequest(request *rest.HttpRequest, endpoint EndpointIfc) *rest.HttpResponse {

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
		case "get": return handleGet(request, endpoint)
		case "head": return handleHead(request, endpoint)
		case "post": return handlePost(request, endpoint)
		case "put": return handlePut(request, endpoint)
		case "options": return handleOptions(request, endpoint)
		case "delete": return handleDelete(request, endpoint)
		case "patch": return handlePatch(request, endpoint)
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
func (endpoint *Endpoint) HandleOptions(request *rest.HttpRequest) *rest.HttpResponse {
	hdrs := rest.HttpHeaders{}
	hdrs.Set("allow", strings.Join(endpoint.methods, ","))
	hlpr := rest.GetHelper()
	return hlpr.ResponseWithHeaders(rest.STATUS_OK, nil, &hdrs)
}

type ConfigurableEndpointIfc interface {
	ConfigureEndpoint(endpointConfig *lib.Config)
}

type GetEndpointIfc interface {
	HandleGet(request *rest.HttpRequest) *rest.HttpResponse
}

type PostEndpointIfc interface {
	HandlePost(request *rest.HttpRequest) *rest.HttpResponse
}

type PutEndpointIfc interface {
	HandlePut(request *rest.HttpRequest) *rest.HttpResponse
}

type OptionsEndpointIfc interface {
	HandleOptions(request *rest.HttpRequest) *rest.HttpResponse
}

type HeadEndpointIfc interface {
	HandleHead(request *rest.HttpRequest) *rest.HttpResponse
}

type DeleteEndpointIfc interface {
	HandleDelete(request *rest.HttpRequest) *rest.HttpResponse
}

type PatchEndpointIfc interface {
	HandlePatch(request *rest.HttpRequest) *rest.HttpResponse
}

func handleImpossible(unmatchedIfc string, requestId string) *rest.HttpResponse {
	lib.GetLogger().Error(fmt.Sprintf(
		"[%s] Endpoint doesn't implement %s (should not be mapped)",
		unmatchedIfc,
		requestId,
	))
	return nil
}

func handleGet(request *rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(GetEndpointIfc); ok {
		return handler.HandleGet(request)
	}
	return handleImpossible("GetEndpointIfc", request.GetContext().GetRequestId())
}

func handlePost(request *rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(PostEndpointIfc); ok {
		return handler.HandlePost(request)
	}
	return handleImpossible("PostEndpointIfc", request.GetContext().GetRequestId())
}

func handlePut(request *rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(PutEndpointIfc); ok {
		return handler.HandlePut(request)
	}
	return handleImpossible("PutEndpointIfc", request.GetContext().GetRequestId())
}

func handleOptions(request *rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(OptionsEndpointIfc); ok {
		return handler.HandleOptions(request)
	}
	return handleImpossible("OptionsEndpointIfc", request.GetContext().GetRequestId())
}

func handleHead(request *rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(HeadEndpointIfc); ok {
		return handler.HandleHead(request)
	}
	// The endpoint doesn't implement Head directly, but we can call GET and modify
	if handler, ok := ep.(GetEndpointIfc); ok {
		// Any endpoint with an expensive Get call should override this default
		// handling with something better tuned to skip expensive steps if possible
		response := handler.HandleGet(request)
		if nil != response {
			hdrs := response.GetHeaders()
			body := response.GetBody()
			bodyLen := 0
			if nil != body { bodyLen = len(*body) }
			hdrs.Set("content-length", strconv.Itoa(bodyLen))
			response.SetBody(nil)
		}
		return response
	}
	return handleImpossible("HeadEndpointIfc", request.GetContext().GetRequestId())
}

func handleDelete(request *rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(DeleteEndpointIfc); ok {
		return handler.HandleDelete(request)
	}
	return handleImpossible("DeleteEndpointIfc", request.GetContext().GetRequestId())
}

func handlePatch(request *rest.HttpRequest, ep interface{}) *rest.HttpResponse {
	if handler, ok := ep.(PatchEndpointIfc); ok {
		return handler.HandlePatch(request)
	}
	return handleImpossible("PatchEndpointIfc", request.GetContext().GetRequestId())
}

