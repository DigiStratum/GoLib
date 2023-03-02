package config

import (
	"fmt"
	"github.com/DigiStratum/GoLib/Starter"
)

type ConfiguredIfc interface {
	ConfigurableIfc
	ConfigIfc
	starter.StartedIfc
}

// Exported to support embedding
type Configured struct {
	Config
	started		*starter.Started
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
		started:	starter.NewStarted(),
		declared:	declared,
	}
}

// -------------------------------------------------------------------------------------------------
// ConfiguredIfc
// -------------------------------------------------------------------------------------------------

// Just capture the provided configuration by default
func (r *Configured) Configure(config ConfigIfc) error {
	r.Config.MergeConfig(config)
	return nil
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *Configured) Start() error {
	if r.started.IsStarted() { return nil }
	// Verify that all required Configs are captured
	for name, declaredConfigItem := range r.declared {
		if declaredConfigItem.IsRequired() {
			if ! r.Config.Has(name) {
				return fmt.Errorf("Missing required config with name '%s'", name)
			}
		}
	}
	r.started.SetStarted()
	return nil
}

