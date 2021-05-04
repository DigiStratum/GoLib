package module

/*
Collection for Endpoints

This is a simple mechanism for Endpoints to be able to automatically register themselves from their
own init() functions without having to be called to do so. It reduces boilerplate code, thus making
the programming interface a simpler one to understand.

This is implemented as a singleton which is globally accessible across all Modules. The functional
design is such that, because each Module initializes one at a time, in sequence, at the time of a
given Module's load as a plugin, the init() functions for all the Endpoints found can register here,
then, when the Module initializes its own Controller, the Controller can pull the Endpoints out of
here. When the Endpoints are pulled out by the Controller, the Collection is once again empty, ready
for the next Module.
*/

// Collection public interface
type EndpointCollectionIfc interface {
	AddEndpoint(endpoint EndpointIfc)
	DrainEndpoints() []EndpointIfc
}

type endpointCollection struct {
	endpoints	[]EndpointIfc
}

var endpointCollectionInstance *endpointCollection

// Get our singleton Collection instance
func GetEndpointCollectionInstance() EndpointCollectionIfc {
	if nil == endpointCollectionInstance {
		endpointCollectionInstance = &endpointCollection{ }
		endpointCollectionInstance.resetEndpoints()
	}
	return endpointCollectionInstance
}

// Add an Endpoint to our Collection
func (ec *endpointCollection) AddEndpoint(endpoint EndpointIfc) {
	(*ec).endpoints = append((*ec).endpoints, endpoint)
}

// Drain all Endpoints out of the Collection, returning them to the caller
func (ec *endpointCollection) DrainEndpoints() []EndpointIfc {
	endpoints := (*ec).endpoints
	ec.resetEndpoints()
	return endpoints
}

// Reset the Collection Endpoint Collection
func (ec *endpointCollection) resetEndpoints() {
	(*ec).endpoints = make([]EndpointIfc, 0)
}

