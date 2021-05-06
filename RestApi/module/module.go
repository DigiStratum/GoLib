package module

/*
This is the integration contract for a Module implementation. From the Stratify SaaS Platform
perspective, a Module is the middle layer of HTTP concerns. It is responsible for receiving and
delegating the handling of HTTP Requests to Endpoints that are mapped to it through configuration.

Any Module implementation MUST implement the Required Module public interface. A Module MAY also
implement any of the Optional Module public interfaces as needed.
*/

import(
	lib "github.com/DigiStratum/GoLib"
)

// Required: Module public interface
type ModuleIfc interface {
	GetName() string
	GetConfig() *lib.Config
}

// Optional: Configurable Module public interface
type ConfigurableModuleIfc interface {
	Configure(moduleConfig *lib.Config) error
}

