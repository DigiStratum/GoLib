// DigiStratum GoLib - Config
package golib

/*

This Config class associates basic helper functions to a simple name/value string map to make life
easier for dealing with simple configuration data. There you can load configuration data from JSON
files, Get/Set individual configuration values, merge additional configuration data in, and more.

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

// Make a new one
func NewConfig() *Config {
	return &Config{}
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
	res := make(Config)
	for k, v := range *cfg {
		if ! strings.HasPrefix(k, prefix) { continue }
		res[k] = v[len(prefix):]
	}
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
// JSON data may only be in the form of an object with named properties with string values
func (cfg *Config) LoadFromJsonString(configJson string) {
	loadFromStringOrPanic(configJson, cfg)
	cfg.DumpConfig()
}

func loadFromStringOrPanic(configJson string, target interface{}) {
	if err := loadFromString(configJson, target); nil != err {
		l := GetLogger()
		msg := fmt.Sprintf("Config.loadFromStringOrPanic(): %s", err.Error())
		l.Fatal(msg)
		panic(msg)
	}
}

func loadFromString(configJson string, target interface{}) error {
	return json.Unmarshal([]byte(configJson), &target);
}

// Load our JSON configuration data from a file on disk
// JSON data may only be in the form of an object with named properties with string values
func (cfg *Config) LoadFromJsonFile(configFile string) {
	LoadJsonOrPanic(configFile, cfg)
	cfg.DumpConfig()
}

// TODO: DEPRECATED; replace calls with LoadFromJsonFile() above
func (cfg *Config) LoadJsonConfiguration(configFile string) {
	cfg.LoadFromJsonFile(configFile)
}

// Generic JSON load or panic
// The provided target should be a pointer to where we will dump the decoded JSON result
func LoadJsonOrPanic(jsonFile string, target interface{}) {
	if err := LoadJson(jsonFile, target); err != nil {
		l := GetLogger()
		msg := fmt.Sprintf("Config: LoadJsonOrPanic(): %s", err.Error())
		l.Fatal(msg)
		panic(msg)
	}
}

// Generic JSON load (into ANY interface)
// The provided target should be a pointer to where we will dump the decoded JSON result
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

