package config

/*

Configurable is an interface with base implementation that allows any construct to embed the data
and behaviors associated with being provided with Config data and ensuring that it is complete.

*/

import (
	"fmt"

	"github.com/DigiStratum/GoLib/Starter"
)

// Any type that implements ConfigurableIfc should be ready to receive configuration data one time as so:
type ConfigurableIfc interface {
	// Embedded interface(s)
	starter.StartableIfc

	// Our own interface
	Configure(config ConfigIfc) error
}

// Exported to support embedding
type Configurable struct {
	*starter.Startable
	config		*Config
	declared	map[string]ConfigItemIfc	// Key is ConfigItem.name for fast lookups
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewConfigurable(configItems ...ConfigItemIfc) *Configurable {
	declared := make(map[string]ConfigItemIfc)
	for _, configItem := range configItems {
		declared[configItem.GetName()] = configItem
	}
	return &Configurable{
		Startable:	starter.NewStartable(),
		declared:	declared,
		config:		NewConfig(),
	}
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc
// -------------------------------------------------------------------------------------------------

// Just capture the provided configuration by default
func (r *Configurable) Configure(config ConfigIfc) error {
	r.config = NewConfig().MergeConfig(config)
	return nil
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *Configurable) Start() error {
	// Verify that all required Configs are captured
	for name, declaredConfigItem := range r.declared {
		if declaredConfigItem.IsRequired() {
			if ! r.config.Has(name) {
				return fmt.Errorf("Missing required config with name '%s'", name)
			}
		}
	}
	return r.Startable.Start()
}

