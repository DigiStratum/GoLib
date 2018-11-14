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
	res "github.com/DigiStratum/GoLib/Resources"
	rest "github.com/DigiStratum/GoLib/RestApi"
)

type ModuleSet	map[string]ModuleIfc

type ModuleIfc interface {
	Configure(serverConfig lib.Config, extraConfig lib.Config)
	GetPath() string
	GetName() string
	HandleRequest(request *rest.HttpRequest) *rest.HttpResponse
}

type Module struct {
	moduleName	string
	serverConfig	*lib.Config	// Server Config passed to us
	moduleConfig	*lib.Config	// Our own Config that we initialize with
	extraConfig	*lib.Config	// Extra data from our own Config to pass on to Endpoints
	repository	*res.Repository
}

// Make a new one of these!
// repository is where we can retrieve all our Module-specific assets (like configuration data)
// name is the unique name of this Module which allows the Server to separate it from others
// TODO: Validate name; non-empty, prefer [a-zA-Z0-9_-.]+ (not starting or ending with '.'!)
func NewModule(repository *res.Repository, name string) *Module {
	return &Module{
		moduleName:	name,
		repository:	repository,
	}
}

// Server needs to initialize this Module with its own configuration data for reference
// Config is passed by value so that we can have a copy, but not tamper with original
func (module *Module) Configure(serverConfig lib.Config, extraConfig lib.Config) {
	l := lib.GetLogger()
	l.Trace(fmt.Sprintf("Module{%s}.Configure()", module.moduleName))

	// Copy Server configuration data for reference
	module.serverConfig = &serverConfig

	// Load Module Config from Resource Repository
	config, err := res.NewRepositoryConfig(module.repository, "config/config.json")
	if nil != err {
		l.Error(fmt.Sprintf(
			"Module{%s}.Configure(): Error loading JSON Config from Repository: %s",
			module.moduleName,
			err.Error(),
		))
		return
	}

	// Validate that the Config has what we need for a Module!
	configPrefix := "module." + module.moduleName + "."
	module.moduleConfig = config.GetSubset(configPrefix)
	requiredConfig := []string{ "version", "path" }
	if ! (module.moduleConfig.HasAll(&requiredConfig)) {
		l.Error(fmt.Sprintf(
			"Module{%s}.Configure(): Incomplete Module Config provided",
			module.moduleName,
		))
		return
	}
	module.moduleConfig.Set("name", module.moduleName) // Reflect name into Module Config

	// See if there are any overrides for this Module hiding in extra Server Config
	overrides := extraConfig.GetSubset(configPrefix)
	if ! overrides.IsEmpty() {
		module.moduleConfig.Merge(overrides)
	}

	// Capture any extra configuration
	module.extraConfig = config.GetInverseSubset(configPrefix)

	// Initialize our controller
	controller := GetController()
	controller.Configure(module.serverConfig, module.moduleConfig, module.extraConfig)
	controller.SetSecurityPolicy(NewSecurityPolicy(config.GetSubset("auth")))
}

// Server needs to know our module's path which it will use to map requests to us
func (module Module) GetPath() string {
	// http://hostname/server.path/module.path/endpoint.pattern
	return module.moduleConfig.Get("path")
}

// Server wants to know our name
func (module Module) GetName() string {
	return module.moduleName
}

// Server wants to send us requests to be handled
// TODO: Eliminate this hop: got from Server directly to Module Controller
func (module *Module) HandleRequest(request *rest.HttpRequest) *rest.HttpResponse {
	ctx := request.GetContext()
	l := lib.GetLogger()
	l.Trace(fmt.Sprintf(
		"[%s] Module (%s): %s %s",
		ctx.GetRequestId(),
		module.moduleName,
		request.GetMethod(),
		request.GetURL(),
	))
	return GetController().HandleRequest(request)
}

