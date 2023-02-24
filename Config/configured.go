package config

import (
	"fmt"
)

type ConfiguredIfc interface {
	ConfigIfc
	InitializableIfc

	Capture(config ConfigIfc) error
}

// Exported to support embedding
type Configured struct {
	Config
	init		*Initialized
	declared	map[string]ConfigItemIfc	// Key is ConfigItem.name for fast lookups
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewConfigured(configItems ...ConfigItemIfc) *Configured {
	declared := make(map[string]ConfigItemIfc)
	for _, configItem := range configItems {
		declared[configItem.GetName()] = configItem
	}
	return &Configured{
		Config:		*(NewConfig()),
		init:		NewInitialized(),
		declared:	declared,
	}
}

// -------------------------------------------------------------------------------------------------
// ConfiguredIfc
// -------------------------------------------------------------------------------------------------

// Just capture the provided configuration by default
func (r *Configured) Capture(config ConfigIfc) error {
	r.Config.MergeConfig(config)
	return nil
}

// -------------------------------------------------------------------------------------------------
// InitializableIfc
// -------------------------------------------------------------------------------------------------

func (r *Configured) Check() error {
	// Verify that all required Configs are captured
	for name, declaredConfigItem := range r.declared {
		if declaredConfigItem.IsRequired() {
			if ! r.Config.Has(name) {
				return fmt.Errorf("Missing required config with name '%s'", name)
			}
		}
	}
	return nil
}

