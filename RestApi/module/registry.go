package module

/*
Registry for Endpoints

This is a simple mechanism for Endpoints to be able to automatically register themselves from their
own init() functions without having to be called to do so. It reduces boilerplate code, thus making
the programming interface a simpler one to understand.

This is implemented as a singleton which is globally accessible across all Modules. The functional
design is such that, because each Module initializes one at a time, in sequence, at the time of a
given Module's load as a plugin, the init() functions for all the Endpoints found can register here,
then, when the Module initializes its own Controller, the Controller can pull the Endpoints out of
here. When the Endpoints are pulled out by the Controller, the Registry is once again empty, ready
for the next Module.
*/

type registeredEndpoint struct {
	endpoint	interface{}
}

type registry struct {
	endpoints	*[]interface{}
}

var registryInstance *registry

// Get our singleton Registry instance
func GetRegistry() *registry {
	if nil == registryInstance {
		eps := make([]interface{}, 0)
		registryInstance = &registry{
			endpoints:	&eps,
		}
	}
	return registryInstance
}

// Add an Endpoint to our Registry
func (reg *registry) AddEndpoint(endpoint interface{}) {
	*reg.endpoints = append(*reg.endpoints, endpoint)
}

// Drain all Endpoints out of the Registry, returning them to the caller
func (reg *registry) DrainEndpoints() *[]interface{} {
	endpoints := reg.endpoints
	eps := make([]interface{}, 0)
	reg.endpoints = &eps
	return endpoints
}

