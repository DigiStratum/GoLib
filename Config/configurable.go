package config

import (
	"fmt"
	"github.com/DigiStratum/GoLib/Starter"
)

// Any type that implements ConfigurableIfc should be ready to receive configuration data one time as so:
type ConfigurableIfc interface {
	// Embedded interface(s)
	starter.StartedIfc

	// Our own interface
	Configure(config ConfigIfc) error
}

// Exported to support embedding
type Configurable struct {
	config		*Config
	startable	*starter.Startable
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
	configurable := Configurable{
		declared:	declared,
		config:		NewConfig(),
	}
	configurable.startable = starter.NewStartable(
		starter.MakeStartable(configurable.start),
	)
	return &configurable
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

func (r *Configurable) start() error {
	// Verify that all required Configs are captured
	for name, declaredConfigItem := range r.declared {
		if declaredConfigItem.IsRequired() {
			if ! r.config.Has(name) {
				return fmt.Errorf("Missing required config with name '%s'", name)
			}
		}
	}
	return nil
}

