package module

/*
Collection of Endpoints

Each Endpoint of a given Module must add itself to this collection with the AddEndpiont() function.
Use Module.GetEndpointCollectionInstance() to get the singleton instance and then call AddEndpoint().
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

