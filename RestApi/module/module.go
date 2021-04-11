package module

/*

This is a service Module. It comprises a collection of related objects necessary for servicing one
or more Endpoints relative to a single base path for the module. The objects include:

* One Controller (standard interface)
* One or more Endpoints (standard interface)
* Any required static objects
* Any required functional libraries
* Configuration management, whether static, dynamic, or both

We want to standardize on a URL Pattern model for our Controller/Endpoint mappings in order to
simplify the code we manage; by establishing our own standard, we reduce the number of variations
that we would otherwise need to account for with multiple Endpoints. For example:

GET/HEAD/OPTIONS/POST                   https://hostname/controller/objects
GET/HEAD/OPTIONS/PUT/PATCH/DELETE       https://hostname/controller/objects/{id}

For the Controller with the Pattern "/controller", we can have a single endpoint with the Pattern
"/objects/*(\d+)*" such that the collection of objects managed by the endpoint can take the HTTP
verbs GET to return the list of objects, POST to create a new object in the collection, and HEAD
or OPTIONS AS normal. The same endpoint may also support a suffix wildcard to catch everything below
that such that the suffix is treated as an individual object ID from the collection of objects
to operate on for GET to retrieve a single object, PUT to replace the object record, PATCH to
modify one or more elements of the object record, DELETE to delete it, and HEAD or OPTIONS as
normal. Using this approach, a single Endpoint may respond to all requests for "/objects/*(\d+)*"
instead of needing two Endpoints: one for "/objects" and one for "/objects/(\d+)". Thus, all
operations related to the object collection which the Endpoint represents may be maintained in the
same place.

*/

import(
	"fmt"

	lib "github.com/DigiStratum/GoLib"
	obj "github.com/DigiStratum/GoLib/Objects"
)

type ModuleSet map[string]*Module

type ModuleIfc interface {
	Configure(serverConfig lib.Config, extraConfig lib.Config) error
	GetPath() string
	GetName() string
}

type Module struct {
	name		string
	serverConfig	*lib.Config	// Server Config passed to us
	moduleConfig	*lib.Config	// Our own Config that we initialize with
	extraConfig	*lib.Config	// Extra data from our own Config to pass on to Endpoints
	objectStore	*obj.ObjectStore
	controller	*Controller
}

// Make a new one of these!
// objectStore is where we can retrieve all our Module-specific assets (like configuration data)
// name is the unique name of this Module which allows the Server to separate it from others
// TODO: Validate name; non-empty, prefer [a-zA-Z0-9_-.]+ (not starting or ending with '.'!)
func NewModule(objectStore *obj.ObjectStore, name string) *Module {
	return &Module{
		name:		name,
		objectStore:	objectStore,
	}
}

// Server needs to initialize this Module with its own configuration data for reference
// Config is passed by value so that we can have a copy, but not tamper with original
// TODO: Break this into smaller, testable functions
func (module *Module) Configure(serverConfig lib.Config, extraConfig lib.Config) error {
	l := lib.GetLogger()
	l.Trace(fmt.Sprintf("Module{%s}.Configure()", module.name))

	// Copy Server configuration data for reference
	module.serverConfig = &serverConfig

	// Load Module Config from Object ObjectStore
	config, err := obj.NewObjectStoreConfig(module.objectStore, "config/config.json")
	if nil != err {
		return l.Error(fmt.Sprintf(
			"Module{%s}.Configure(): Error loading JSON Config from ObjectStore: %s",
			module.name,
			err.Error(),
		))
	}

//config.Dump()

	// Validate that the Config has what we need for a Module!
	configPrefix := "module." + module.name + "."
	l.Trace(fmt.Sprintf(
		"Module{%s}.Cofigure(): Looking for config subset with prefix '%s'",
		module.name,
		configPrefix,
	))
	module.moduleConfig = config.GetSubset(configPrefix)
	requiredConfig := []string{ "version", "path" }
	if ! (module.moduleConfig.HasAll(&requiredConfig)) {
		return l.Error(fmt.Sprintf(
			"Module{%s}.Configure(): Incomplete Module Config provided",
			module.name,
		))
	}
	module.moduleConfig.Set("name", module.name) // Reflect name into Module Config

	// See if there are any overrides for this Module hiding in extra Server Config
	overrides := extraConfig.GetSubset(configPrefix)
	if ! overrides.IsEmpty() {
		l.Trace(fmt.Sprintf(
			"Module{%s}.Configure(): Applying overrides from extra Server Config",
			module.name,
		))
		overrides.Dump()
		module.moduleConfig.Merge(overrides)
	}

	// Dereference any Config references, descending from highest level config
	module.moduleConfig.DereferenceAll(module.serverConfig, module.moduleConfig)
	l.Crazy(fmt.Sprintf(
		"Module{%s} Configuration: %s",
		module.name,
		module.moduleConfig.DumpString(),
	));

	// Capture any extra configuration
	module.extraConfig = config.GetInverseSubset(configPrefix)
	module.extraConfig.DereferenceAll(module.serverConfig, module.moduleConfig, module.extraConfig)
	l.Crazy(fmt.Sprintf(
		"Module{%s} Extra Configuration: %s",
		module.name,
		module.extraConfig.DumpString(),
	));


	// Initialize our Controller
	module.controller = NewController()
	module.controller.Configure(module.serverConfig, module.moduleConfig, module.extraConfig)
	module.controller.SetSecurityPolicy(NewSecurityPolicy(config.GetSubset("auth")))
	return nil
}

// Server needs to know our Module's path which it will use to map Requests to us
func (module Module) GetPath() string {
	// http://hostname/server.path/module.path/endpoint.pattern
	return module.moduleConfig.Get("path")
}

// Server wants to know our name
func (module Module) GetName() string {
	return module.name
}

// Server wants to send Requests to our Controller
func (module Module) GetController() *Controller {
	return module.controller
}

