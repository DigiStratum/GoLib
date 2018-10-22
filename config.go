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
func NewConfig() Config {
	return Config{}
}

// Merge some additional configuration data on top of our own
func (cfg *Config) Merge(inbound Config) {
	for k, v := range inbound { (*cfg)[k] = v }
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

//  Get configuration datum whose keys begin with the base string...
func (cfg *Config) GetSubset(base string) Config {
	res := make(Config)
	for k, v := range *cfg {
		if ! strings.HasPrefix(k, base) { continue }
		res[k] = v
	}
	return res
}

// Load our JSON configuration data from a file on disk
// JSON data may only be in the form of an object with named properties with string values
func (cfg *Config) LoadJsonConfiguration(configFile string) {
	LoadJsonOrPanic(configFile, cfg)
	cfg.DumpConfig()
}

// Dump our configuration data
func (cfg *Config) DumpConfig() {
	l := GetLogger()
	l.Trace("Config:")
	l.Trace("--------------------------")
	for k, v := range *cfg {
		l.Trace(fmt.Sprintf("\t'%s': '%s'", k, v))
	}
	l.Trace("--------------------------")
}

// Generic JSON load or panic
func LoadJsonOrPanic(jsonFile string, target interface{}) {
	err := LoadJson(jsonFile, target)
        if err == nil { return }
	l := GetLogger()
	msg := fmt.Sprintf("Config: LoadJsonOrPanic(): %s", err.Error())
	l.Fatal(msg)
	panic(msg)
}

// Generic JSON load (into ANY interface)
// Ref: https://blogger.golang.org/json-and-go
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

