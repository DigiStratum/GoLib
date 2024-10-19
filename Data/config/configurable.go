package config

/*

Configurable is an interface with base implementation that allows any construct to embed the data
and behaviors associated with being provided with Config data and ensuring that it is complete.

Note that while some of what is covered here could also be covered with general purpose dependency
injection, however configurability and config details are such fundamental building blocks that we
believe they deserve to be first class citizens of our framework for improved clarity and config
specific operations that are not to be confused or comingled with other dependencies (which may
also have need of some or all of the config data).

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
	DeclareConfigItems(configItems ...ConfigItemIfc) *Configurable
	Configure(config ConfigIfc) error
	GetMissingConfigs() []string
	HasMissingConfigs() bool
	GetConfig() ConfigIfc
}

// Exported to support embedding
type Configurable struct {
	*startable.Startable
	config		ConfigIfc
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
	return c.DeclareConfigItems(configItems...)
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc
// -------------------------------------------------------------------------------------------------

func (r *Configurable) DeclareConfigItems(configItems ...ConfigItemIfc) *Configurable {
	for _, configItem := range configItems { r.declared[configItem.GetSelector()] = configItem }
	return r
}

// Just capture the provided configuration by default
// Overrides should call this parent, and return error if this fails or for any validation problems
func (r *Configurable) Configure(config ConfigIfc) error {
	// Disallow Configure() after we've already Started
	if r.Startable.IsStarted() { return fmt.Errorf("Already started; Config is immutable now") }
	if nil == config {
		r.config = NewConfig()
	} else { r.config = config }
	return nil
}

// Verify that all required Configs are captured
func (r *Configurable) GetMissingConfigs() []string {
	requiredConfigs := []string{}
	for selector, declaredConfigItem := range r.declared {
		if ! declaredConfigItem.IsRequired() { continue }
		requiredConfigs = append(requiredConfigs, selector)
	}
	return r.config.GetMissing(requiredConfigs...)
}

// MissingConfigs as a bool!
func (r *Configurable) HasMissingConfigs() bool {
	return len(r.GetMissingConfigs()) > 0
}

func (r *Configurable) GetConfig() ConfigIfc {
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
			"Missing required config(s): '%s'",
			strings.Join(missingConfigs, "','"),
		)
	}

	// For all the declared Config Items...
	for selector, configItem := range r.declared {

		// Get the value for this config selector
		configDataValue := r.config.Select(selector)

		// If this ConfigItem has a Validation Func...
		if (nil != configDataValue) && configItem.CanValidate() {
			if err := configItem.Validate(configDataValue); nil != err {
				return fmt.Errorf(
					"Config Item '%s' failed validation with value: (%s) '%s': %s",
					selector,
					configDataValue.GetType().ToString(),
					configDataValue.ToString(),
					err.Error(),
				)
			}
		}

		// If this ConfigItem has a Capture Func...
		if configItem.CanCapture() {
			if err := configItem.Capture(configDataValue); nil != err { return err }
		}
	}

	return r.Startable.Start()
}

