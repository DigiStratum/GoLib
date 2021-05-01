package module

import(
	lib "github.com/DigiStratum/GoLib"
)

// Required: Module public interface
type ModuleIfc interface {
	GetName() string
	GetConfig() *lib.Config
}

// Optional: Configurable Module public interface
// Note: if a concrete module optionally implements this interface, then it can receive the Module config (from ModuleIfc.GetConfig())
type ConfigurableModuleIfc interface {
	Configure(moduleConfig *lib.Config) error
}

