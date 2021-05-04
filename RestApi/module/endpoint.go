package module

import (
	lib "github.com/DigiStratum/GoLib"
	rest "github.com/DigiStratum/GoLib/RestApi"
)

// Endpoint public interface
type EndpointIfc interface {
	GetName() string
	GetVersion() string
/*
	Configure(concreteEndpoint interface{}, serverConfig lib.Config, moduleConfig lib.Config, extraConfig lib.Config)
	Init(concreteEndpoint interface{})
	GetSecurityPolicy() *SecurityPolicy
	GetId() string
	GetPattern() string
	IsDefault() bool
	GetMethods() []string
	GetRelativeURI(request *rest.HttpRequest) string
	GetRequestMatches(request *rest.HttpRequest) []string
	GetRequestURIMatches(relativeURI string) []string
	GetRequestPathParameters(request *rest.HttpRequest) (*lib.HashMap, error)
	GetRequestURIPathParameters(relativeURI string) (*lib.HashMap, error)
	HandleRequest(request *rest.HttpRequest, endpoint EndpointIfc) *rest.HttpResponse
*/
}

// Implementation-Dependent Endpoint Interface: Configurability
type AllConfigurableIfc interface {
	ConfigureAll(endpointConfig, moduleConfig, serverConfig *lib.Config)
}

// Implementation-Dependent Endpoint Interface: Configurability
type ConfigurableEndpointIfc interface {
	ConfigureEndpoint(endpointConfig *lib.Config)
}

// Implementation-Dependent Endpoint Interface: ANY METHOD request handling
type AnyEndpointIfc interface {
	HandleAny(request *rest.HttpRequest) *rest.HttpResponse
}

// Implementation-Dependent Endpoint Interface: GET request handling
type GetEndpointIfc interface {
	HandleGet(request *rest.HttpRequest) *rest.HttpResponse
}

// Implementation-Dependent Endpoint Interface: POST request handling
type PostEndpointIfc interface {
	HandlePost(request *rest.HttpRequest) *rest.HttpResponse
}

// Implementation-Dependent Endpoint Interface: PUT request handling
type PutEndpointIfc interface {
	HandlePut(request *rest.HttpRequest) *rest.HttpResponse
}

// Implementation-Dependent Endpoint Interface: OPTIONS request handling
type OptionsEndpointIfc interface {
	HandleOptions(request *rest.HttpRequest) *rest.HttpResponse
}

// Implementation-Dependent Endpoint Interface: HEAD request handling
type HeadEndpointIfc interface {
	HandleHead(request *rest.HttpRequest) *rest.HttpResponse
}

// Implementation-Dependent Endpoint Interface: DELETE request handling
type DeleteEndpointIfc interface {
	HandleDelete(request *rest.HttpRequest) *rest.HttpResponse
}

// Implementation-Dependent Endpoint Interface: PATCH request handling
type PatchEndpointIfc interface {
	HandlePatch(request *rest.HttpRequest) *rest.HttpResponse
}

