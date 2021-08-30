package config

// Any type that implements ConfigurableIfc should be ready to receive configuration data one time as so:
type ConfigurableIfc interface {
	Configure(config ConfigIfc) error
}
