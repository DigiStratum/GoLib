package config

/*

Configurable is an interface with base implementation that allows any construct to embed the data
and behaviors associated with being provided with Config data and ensuring that it is complete.

*/

import (
	"fmt"
	"strings"

	"github.com/DigiStratum/GoLib/Starter"
)

// Any type that implements ConfigurableIfc should be ready to receive configuration data one time as so:
type ConfigurableIfc interface {
	// Embedded interface(s)
	starter.StartableIfc

	// Our own interface
	Configure(config ConfigIfc) error
	GetMissingConfigs() []string
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

// Verify that all required Configs are captured
func (r *Configurable) GetMissingConfigs() []string {
	missingConfigs := []string{}
	for name, declaredConfigItem := range r.declared {
		if ! declaredConfigItem.IsRequired() { continue }
		if ! r.config.Has(name) { continue }
		missingConfigs = append(missingConfigs, name)
	}
	return missingConfigs
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *Configurable) Start() error {
	if r.Startable.IsStarted() { return nil }
	// Make sure nothing required is missing
	if missingConfigs := r.GetMissingConfigs(); len(missingConfigs) > 0 {
		return fmt.Errorf(
			"Missing required config(s) with name(s): %s",
			strings.Join(missingConfigs, ","),
		)
	}

	// Run capture funcs for all the declared configs
	for name, configItem := range r.declared {
		// If this dependency is declared and defines Capture Func...
		if ! configItem.CanCapture() { continue }
		if ! r.config.Has(name) { continue }
		// We only capture non-nil config value; TODO: is there a real scenario where we want to capture nil
		value := r.config.Get(name)
		if nil == value { continue }
		if err := configItem.Capture(*value); nil != err { return err }
	}

	return r.Startable.Start()
}

