package module

/*

This is a service Module. It comprises a collection of related resources necessary for servicing one
or more Endpoints relative to a single base path for the module. The resources include:

* One Controller (standard interface)
* One or more Endpoints (standard interface)
* Any required static resources
* Any required functional libraries
* Configuration management, whether static, dynamic, or both

We want to standardize on a URL Pattern model for our Controller/Endpoint mappings in order to
simplify the code we manage; by establishing our own standard, we reduce the number of variations
that we would otherwise need to account for with multiple Endpoints. For example:

GET/HEAD/OPTIONS/POST                   https://hostname/controller/resources
GET/HEAD/OPTIONS/PUT/PATCH/DELETE       https://hostname/controller/resources/{id}

For the Controller with the Pattern "/controller", we can have a single endpoint with the Pattern
"/resources/*(\d+)*" such that the collection of resources managed by the endpoint can take the HTTP
verbs GET to return the list of resources, POST to create a new resource in the collection, and HEAD
or OPTIONS AS normal. The same endpoint may also support a suffix wildcard to catch everything below
that such that the suffix is treated as an individual resource ID from the collection of resources
to operate on for GET to retrieve a single resource, PUT to replace the resource record, PATCH to
modify one or more elements of the resource record, DELETE to delete it, and HEAD or OPTIONS as
normal. Using this approach, a single Endpoint may respond to all requests for "/resources/*(\d+)*"
instead of needing two Endpoints: one for "/resources" and one for "/resources/(\d+)". Thus, all
operations related to the resource collection which the Endpoint represents may be maintained in the
same place.

*/

import(
	"fmt"

	lib "github.com/DigiStratum/GoLib"
	rest "github.com/DigiStratum/GoLib/RestApi"
)

type ModulePath	string
type ModuleSet	map[ModulePath]ModuleIfc

type ModuleIfc interface {
	Configure(serverConfig lib.Config)
	GetPath() ModulePath
	GetName() string
	HandleRequest(request *rest.HttpRequest) *rest.HttpResponse
}

// TODO: Privatize all this stuff with proper accessor functions
type Module struct {
	Controller	Controller
	ServerConfig	lib.Config
	ModuleConfig	lib.Config
	securityPolicy	SecurityPolicy
}

// Make a new one!
func NewModule() Module {
	mod := Module{
		securityPolicy:	NewSecurityPolicy(),
	}
	return mod
}

// Server needs to initialize this Module with its own configuration data for reference
func (module *Module) Configure(serverConfig lib.Config) {
	l := lib.GetLogger()
	l.Trace("Module: Configure")
	// Copy over the server configuration data
	module.ServerConfig = serverConfig

	// Initialize our controller
	module.Controller.SetSecurityPolicy(*module.GetSecurityPolicy())
	module.Controller.Configure(module.ServerConfig, module.ModuleConfig)
}

// Module/Controller need to be able to access their own Security Policy
func (module *Module) GetSecurityPolicy() *SecurityPolicy {
	return &module.securityPolicy
}

// Server needs to know our module's path which it will use to map requests to us
func (module Module) GetPath() ModulePath {
	// http://hostname/server.path/module.path/endpoint.pattern
	return ModulePath(module.ModuleConfig.Get("module.path"))
}

// Server wants to know our module's name
func (module Module) GetName() string {
	return module.ModuleConfig.Get("module.name")
}

// Server wants to send us requests to be handled
func (module *Module) HandleRequest(request *rest.HttpRequest) *rest.HttpResponse {
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Trace(fmt.Sprintf(
		"[%s] Module (%s): %s %s",
		ctx.GetRequestId(),
		module.ModuleConfig.Get("module.name"),
		request.GetMethod(),
		request.GetURL(),
	))
	return module.Controller.HandleRequest(*request)
}

