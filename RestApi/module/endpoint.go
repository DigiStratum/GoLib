package module

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

