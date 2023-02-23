package config

type ConfiguredIfc interface {
	ConfigIfc
	ConfigurableIfc

	Configure(config ConfigIfc) error
}

// Exported to support embedding
type Configured struct {
	*Config
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
		Config:		NewConfig(),
		declared:	declared,
	}
}

// -------------------------------------------------------------------------------------------------
// ConfiguredIfc Implementation
// -------------------------------------------------------------------------------------------------

// Just capture the provided configuration by default
func (r *Configured) Configure(config ConfigIfc) error {
	r.Config.MergeConfig(config)
	return nil
}

