package config

/*

Configurable is an interface with base implementation that allows any construct to embed the data
and behaviors associated with being provided with Config data and ensuring that it is complete.

*/

import (
	"fmt"
	"strings"

	"github.com/DigiStratum/GoLib/Process/startable"
)

// Any type that implements ConfigurableIfc should be ready to receive configuration data one time as so:
type ConfigurableIfc interface {
	// Embedded interface(s)
	startable.StartableIfc

	// Our own interface
	AddConfigItems(configItems ...ConfigItemIfc) *Configurable
	Configure(config ConfigIfc) error
	HasMissingConfigs() bool
	GetMissingConfigs() []string
	GetConfig() *Config
}

// Exported to support embedding
type Configurable struct {
	*startable.Startable
	config		*Config
	declared	map[string]ConfigItemIfc	// Key is ConfigItem.name for fast lookups
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewConfigurable(configItems ...ConfigItemIfc) *Configurable {
	c := Configurable{
		Startable:	startable.NewStartable(),
		declared:	make(map[string]ConfigItemIfc),
		config:		NewConfig(),
	}
	return c.AddConfigItems(configItems...)
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc
// -------------------------------------------------------------------------------------------------

func (r *Configurable) AddConfigItems(configItems ...ConfigItemIfc) *Configurable {
	for _, configItem := range configItems { r.declared[configItem.GetName()] = configItem }
	return r
}

// Just capture the provided configuration by default
// Overrides should call this parent, and return error if this fails or for any validation problems
func (r *Configurable) Configure(config ConfigIfc) error {
	// Disallow Configure() after we've already Started
	if r.Startable.IsStarted() { return fmt.Errorf("Already started; Config is immutable now") }
	r.config = NewConfig().MergeConfig(config)
	return nil
}

// MissingConfigs as a bool!
func (r *Configurable) HasMissingConfigs() bool {
	return len(r.GetMissingConfigs()) > 0
}

// Verify that all required Configs are captured
func (r *Configurable) GetMissingConfigs() []string {
	missingConfigs := []string{}
	for name, declaredConfigItem := range r.declared {
		if ! declaredConfigItem.IsRequired() { continue }
		if r.config.Has(name) { continue }
		missingConfigs = append(missingConfigs, name)
	}
	return missingConfigs
}

func (r *Configurable) GetConfig() *Config {
	// Require Start() first to finalize Config
	if ! r.Startable.IsStarted() { return nil }
	return r.config
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

	// For all the declared Config Items...
	for name, configItem := range r.declared {

		// Get the value for this named ConfigItem
		value := r.config.Get(name)

		// If this ConfigItem has a Validation Func...
		if (nil != value) && configItem.CanValidate() {
			if ! configItem.Validate(*value) {
				return fmt.Errorf(
					"Configuration Item '%s' failed validation with value: '%s'",
					name, *value,
				)
			}
		}

		// If this ConfigItem has a Capture Func...
		if configItem.CanCapture() {
			if err := configItem.Capture(value); nil != err { return err }
		}
		if configItem.CanCaptureSubset() {
			subset := r.config.GetSubsetConfig(name)
			if err := configItem.CaptureSubset(subset); nil != err { return err }
		}
	}

	return r.Startable.Start()
}

