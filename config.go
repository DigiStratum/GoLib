// DigiStratum GoLib - Config
package golib

/*

This Config class extends (embeds) our own HashMap with additional capabilities specific to the
needs of managing simple configuration data for our purposes. You can load configuration data from
JSON files, Get/Set individual configuration values, merge additional configuration data in, and
more.

JSON Config data may only be in the form of an object with named properties with string values. We
generally pass around JSON strings as a pointer whenever we can to avoid copying potentially large
JSON strings unnecessarily. As long as we are in a trusted code/library scope, this is fine; when
we get into an untrusted code/library scope, we must revert to pass by value as needed to prevent
unauthorized tampering.

In addition to the explicit imports below, we use the following classes from this same package:
 * HashMap
 * Logger
 * Json

*/

import(
	"strings"
	"fmt"
)

// Config embeds a HashMap so that we can extend it
// ref: https://stackoverflow.com/questions/28800672/how-to-add-new-methods-to-an-existing-type-in-go
type Config struct {
	HashMap
}

// Make a new one of these!
func NewConfig() *Config {
	hash := NewHashMap()
	return &Config{ HashMap: *hash }
}

// Merge configuration data
// Just because we embed HashMap doesn't mean data type casting works out for us
func (cfg *Config) Merge(mergeCfg *Config) {
	cfg.HashMap.Merge(&(mergeCfg.HashMap))
}

// Get configuration datum whose keys begin with the prefix...
// We also strip the prefix off leaving just the interesting parts
func (cfg *Config) GetSubset(prefix string) *Config {
	return cfg.getSubset(prefix, true)
}

// Get configuration datum whose keys DO NOT begin with the prefix...
// We also strip the prefix off leaving just the interesting parts
func (cfg *Config) GetInverseSubset(prefix string) *Config {
	return cfg.getSubset(prefix, false)
}

// Get configuration datum whose keys Do/Don't begin with the prefix...
// Return the matches if keepMatches, else return the NON-matches
func (cfg *Config) getSubset(prefix string, keepMatches bool) *Config {
	res := NewConfig()
	for pair := range cfg.IterateChannel() {
		matches := strings.HasPrefix(pair.Key, prefix)
		if (matches) {
			if ! keepMatches { continue }
			strippedKey := pair.Key[len(prefix):]
			res.Set(strippedKey, pair.Value)
		} else {
			if keepMatches { continue }
			res.Set(pair.Key, pair.Value)
		}
	}
	return res
}

// Load our JSON configuration data from a string
func (cfg *Config) LoadFromJsonString(configJson *string) {
	NewJson(configJson).LoadOrPanic(&cfg.HashMap)
	//cfg.Dump()
}

// Load our JSON configuration data from a string (or return an error)
func (cfg *Config) LoadFromJsonStringOrError(configJson *string) error {
	if err := NewJson(configJson).Load(&cfg.HashMap); nil == err { return err }
	//cfg.Dump()
	return nil
}

// Load our JSON configuration data from a file on disk
func (cfg *Config) LoadFromJsonFile(configFile string) {
	NewJsonFromFile(configFile).LoadOrPanic(&cfg.HashMap)
	//cfg.Dump()
}

// Dereference any values we have that %reference% keys in the referenceConfig
func (cfg *Config) Dereference(referenceConfig *Config) {
	GetLogger().Trace("Config.Dereference()")
	GetLogger().Crazy(fmt.Sprintf(
		"Dereferencing against Config: %s",
		referenceConfig.DumpString(),
	));
	// For each of our key/value pairs...
	for cpair := range cfg.IterateChannel() {
		// For each of the referenceConfig's key/value pairs...
		for rcpair := range referenceConfig.IterateChannel() {
			// A reference looks like '%key%'...
			reference := fmt.Sprintf("%%%s%%", rcpair.Key)
			GetLogger().Crazy(fmt.Sprintf(
				"Config.Dereference() -> config['%s'] = '%s' value has '%s' ... ?",
				cpair.Key,
				cpair.Value,
				reference,
			));

			// If the referenceConfig value doesn't reference config key, move on...
			if ! strings.Contains(cpair.Value, reference) { continue }

			// Replace the reference(s) in our value with the values referenced
			GetLogger().Trace(fmt.Sprintf(
				"\tReplaced '%s' with '%s' in '%s'",
				reference, rcpair.Value, cpair.Value,
			))
			cfg.Set(
				cpair.Key,
				strings.Replace(cpair.Value, reference, rcpair.Value, -1),
			)
		}
	}
}

// Dereference against a list of other referenceConfigs
func (cfg *Config) DereferenceAll(referenceConfigs ...*Config) {
	for _, referenceConfig := range referenceConfigs {
		GetLogger().Crazy(fmt.Sprintf(
			"DereferenceAll against Config: %s",
			referenceConfig.DumpString(),
		));
		cfg.Dereference(referenceConfig)
	}
}

