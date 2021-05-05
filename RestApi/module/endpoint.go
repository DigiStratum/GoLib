package module

/*

This is the integration contract for an Endpoint implementation. From the Stratify SaaS Platform
perspective, an Endpoint is the lowest layer of HTTP concerns that adapts to business logic
concerns. It is responsible for receiving and handling HTTP Requests that are mapped to it through
configuration, and returning a well-formed HTTP Response that the higher layers can return to the
client.

Any Endpoint implementation MUST implement the Required Endpoint public interface. As well, at least
one Optional Endpoint public interface SHOULD be implemented, otherwise the Endpoint would not be
able to receive any HTTP Requests to handle. All of the Optional Endpoint public interfaces MAY be
implemented with one caveat: the AnyEndpointIfc is only used as a fallback if any of the other
Optional Endpoint public interfaces has NOT been implemented. If ALL of the other Optional Endpoint
public interfaces have been implemented, then the AnyEndpointIfc will not receive requests. If the
AnyEndpointIfc is implemented, then it will receive any HTTP Requests for any supported HTTP method
where there is not an implementation provided for that method.

*/

import (
	lib "github.com/DigiStratum/GoLib"
	rest "github.com/DigiStratum/GoLib/RestApi"
)

// Required: Endpoint public interface
type EndpointIfc interface {
	GetName() string
	GetVersion() string
	GetId() string
}

// Optional: Endpoint public interface: Configurability
type ConfigurableEndpointIfc interface {
	Configure(endpointConfig *lib.Config)
}

// Optional: Endpoint public interface: ANY METHOD request handling
type AnyEndpointIfc interface {
	HandleAny(request *rest.HttpRequest) *rest.HttpResponse
}

// Optional: Endpoint public interface: GET request handling
type GetEndpointIfc interface {
	HandleGet(request *rest.HttpRequest) *rest.HttpResponse
}

// Optional: Endpoint public interface: POST request handling
type PostEndpointIfc interface {
	HandlePost(request *rest.HttpRequest) *rest.HttpResponse
}

// Optional: Endpoint public interface: PUT request handling
type PutEndpointIfc interface {
	HandlePut(request *rest.HttpRequest) *rest.HttpResponse
}

// Optional: Endpoint public interface: OPTIONS request handling
type OptionsEndpointIfc interface {
	HandleOptions(request *rest.HttpRequest) *rest.HttpResponse
}

// Optional: Endpoint public interface: HEAD request handling
type HeadEndpointIfc interface {
	HandleHead(request *rest.HttpRequest) *rest.HttpResponse
}

// Optional: Endpoint public interface: DELETE request handling
type DeleteEndpointIfc interface {
	HandleDelete(request *rest.HttpRequest) *rest.HttpResponse
}

// Optional: Endpoint public interface: PATCH request handling
type PatchEndpointIfc interface {
	HandlePatch(request *rest.HttpRequest) *rest.HttpResponse
}

