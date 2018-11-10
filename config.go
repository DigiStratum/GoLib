// DigiStratum GoLib - Config
package golib

/*

This Config class associates basic helper functions to a simple name/value string map to make life
easier for dealing with simple configuration data. There you can load configuration data from JSON
files, Get/Set individual configuration values, merge additional configuration data in, and more.

JSON Config data may only be in the form of an object with named properties with string values. We
generally pass around JSON strings as a pointer whenever we can to avoid copying potentially large
JSON strings unnecessarily. As long as we are in a trusted code/library scope, this is fine; when
we get into an untrusted code/library scope, we must revert to pass by value as needed to prevent
unauthorized tampering.

In addition to the explicit imports below, we also leverage the following classes from here:
 * Logger

*/

import(
	"strings"
	"os"
	"fmt"
	"encoding/json"
	"errors"
)

type Config map[string]string

// Make a new one of these!
func NewConfig() *Config {
	return &Config{}
}

// Check whether this Config is empty (has no properties)
func (cfg *Config) IsEmpty() bool {
	return 0 == cfg.Size()
}

// Get the number of properties in this Config
func (cfg *Config) Size() int {
	return len(*cfg)
}

// Merge some additional configuration data on top of our own
func (cfg *Config) Merge(inbound *Config) {
	for k, v := range *inbound { (*cfg)[k] = v }
}

// Set a single configuration element key to the specified value
func (cfg *Config) Set(key string, value string) {
	(*cfg)[key] = value
}

// Get a single configuration element by key name
func (cfg *Config) Get(key string) string {
	str := ""
	if val, ok := (*cfg)[key]; ok { str = val }
	return str
}

// Check whether we have a configuration element by key name
func (cfg *Config) Has(key string) bool {
	_, ok := (*cfg)[key];
	return ok
}

// Check whether we have configuration elements for all the key names
func (cfg *Config) HasAll(keys *[]string) bool {
	if nil == keys { return false }
	for _, key := range *keys {
		_, ok := (*cfg)[key];
		if ! ok { return false }
	}
	return true
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
	res := make(Config)
	for k, v := range *cfg {
		matches := strings.HasPrefix(k, prefix)
		if (matches) {
			if ! keepMatches { continue }
			strippedKey := k[len(prefix):]
			res[strippedKey] = v
		} else {
			if keepMatches { continue }
			res[k] = v
		}
	}
	return &res
}

// Get a full copy of this Config
// This is so that we can give away a copy to someone else without allowing them to tamper with us
func (cfg *Config) GetCopy() *Config {
	res := make(Config)
	for k, v := range *cfg { res[k] = v }
	return &res
}

// Dump our configuration data
func (cfg *Config) DumpConfig() {
	l := GetLogger()
	l.Crazy("Config:")
	l.Crazy("--------------------------")
	for k, v := range *cfg {
		l.Crazy(fmt.Sprintf("\t'%s': '%s'", k, v))
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

