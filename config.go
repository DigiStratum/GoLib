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

In addition to the explicit imports below, we also leverage the following classes from here:
 * HashMap
 * Logger

*/

import(
	"strings"
	"os"
	"fmt"
	"encoding/json"
	"errors"
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

// Dump our configuration data
func (cfg *Config) DumpConfig() {
	l := GetLogger()
	l.Crazy("Config:")
	l.Crazy("--------------------------")
	for pair := range cfg.IterateChannel() {
		l.Crazy(fmt.Sprintf("\t'%s': '%s'", pair.Key, pair.Value))
	}
	l.Crazy("--------------------------")
}

// Load our JSON configuration data from a string
func (cfg *Config) LoadFromJsonString(configJson *string) {
	loadFromJsonStringOrPanic(configJson, cfg)
	cfg.DumpConfig()
}

func loadFromJsonStringOrPanic(configJson *string, target interface{}) {
	if err := loadFromJsonString(configJson, target); nil != err { panic(err.Error()) }
}

// Load our JSON configuration data from a string (or return an error)
func (cfg *Config) LoadFromJsonStringOrError(configJson *string) error {
	if err := loadFromJsonString(configJson, cfg); nil != err {
		return err
	}
	cfg.DumpConfig()
	return nil
}

func loadFromJsonString(configJson *string, target interface{}) error {
	if nil == configJson {
		msg := "Config.loadFromJsonString(): We were given nil for the Config JSON"
		GetLogger().Error(msg)
		return  errors.New(msg)
	}
	if err := json.Unmarshal([]byte(*configJson), &target); err != nil {
		msg := fmt.Sprintf("Config.loadFromJsonString(): Failed to unmarshall JSON: %s", err.Error())
		GetLogger().Error(msg)
		return errors.New(msg)
	}
	return nil
}

// Load our JSON configuration data from a file on disk
func (cfg *Config) LoadFromJsonFile(configFile string) {
	LoadJsonOrPanic(configFile, cfg)
	cfg.DumpConfig()
}

// FIXME: DEPRECATED; replace calls with LoadFromJsonFile() above
func (cfg *Config) LoadJsonConfiguration(configFile string) {
	cfg.LoadFromJsonFile(configFile)
}

// Dereference any values we have that %reference% keys in the referenceConfig
func (cfg *Config) Dereference(referenceConfig *Config) {
	GetLogger().Trace("Config.Dereference()")
	// For each of our key/value pairs...
	for cpair := range cfg.IterateChannel() {
		// For each of the referenceConfig's key/value pairs...
		for rcpair := range cfg.IterateChannel() {
			// A reference looks like '%key%'...
			reference := fmt.Sprintf("%%%s%%", rcpair.Key)

			// And if our value doesn't have the reference, move on...
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

// Generic JSON load or panic
// The provided target should be a pointer to where we will dump the decoded JSON result
func LoadJsonOrPanic(jsonFile string, target interface{}) {
	if err := LoadJson(jsonFile, target); err != nil {
		msg := fmt.Sprintf("Config.LoadJsonOrPanic(): %s", err.Error())
		GetLogger().Fatal(msg)
		panic(msg)
	}
}

// Generic JSON load (into ANY interface)
// The provided target should be a pointer to where we will dump the decoded JSON result
// TODO: relocate this to a general purpose JSON library as it is not Config-specific
func LoadJson(jsonFile string, target interface{}) error {
        file, err := os.Open(jsonFile)
        if nil == err {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(target)
		file.Close()
		if nil == err { return nil }
	}
	// Decorate the errror with a little more context
	msg := fmt.Sprintf("LoadJson(): file='%s': '%s'", jsonFile, err.Error())
	return errors.New(msg)
}

