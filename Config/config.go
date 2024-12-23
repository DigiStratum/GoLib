// DigiStratum GoLib - Config
package config

/*

DEPRECATED - See GoLib/Data/config replacement

This Config class extends (embeds) our own HashMap with additional capabilities specific to the
needs of managing simple configuration data for our purposes. You can load configuration data from
JSON files, Get/Set individual configuration values, merge additional configuration data in, and
more.

JSON Config data may only be in the form of an object with named properties with string values. We
generally pass around JSON strings as a pointer whenever we can to avoid copying potentially large
JSON strings unnecessarily. As long as we are in a trusted code/library scope, this is fine; when
we get into an untrusted code/library scope, we must revert to pass by value as needed to prevent
unauthorized tampering.

Note that there is support for "dereferencing" values lurking within. If a value has a special
notation "%key%" and we call one of the Dereference member functions, if there is a matching key
within the available Config data, then the value of that matching key will replace the "%key%"
reference. This is recursive up to a depth of MAX_REFERENCE_DEPTH defined below. If there is no
match on the key, then no substitution is performed. There are some neat tricks with composition
of Config collections that can DereferenceConfig() against eachother so that we can make
references across Configs.

In addition to the explicit imports below, we use the following classes from this same package:
 * HashMap
 * Json

TODO:
 * Add support for non-string values, including nested objects, to break name=value pair limits
 * Add "configurator" Interface/boilerplate implementation(s) to fetch Config from various sources
 * Allow override for MAX_REFERENCE_DEPTH and reduce the default. By a lot.
 * Extend Data/DataValue instead of Data/hashmap
*/

import(
	"strings"
	"fmt"

	"github.com/DigiStratum/GoLib/Data/hashmap"
)

// Prevent runaway processes with absurd boundaries with an absolute maximum on loop count
const MAX_REFERENCE_DEPTH = 100

type ConfigIfc interface {
	// ref: https://www.geeksforgeeks.org/embedding-interfaces-in-golang/
	hashmap.HashMapIfc
	MergeConfig(mergeCfg ConfigIfc) *Config
	GetSubsetConfig(prefix string) *Config
	GetSubsetKeys(keys *[]string) *Config
	GetInverseSubsetConfig(prefix string) *Config
	DereferenceString(str string) *string
	Dereference(referenceConfig ConfigIfc) int
	DereferenceAll(referenceConfigs ...ConfigIfc) int
	DereferenceLoop(maxLoops int, referenceConfig ConfigIfc) bool

	//Validate(required, optional *[]string) *Config
}

// Config embeds a HashMap so that we can extend it
type Config struct {
	*hashmap.HashMap
}

// Factory Functions
func NewConfig() *Config {
	return &Config{ HashMap: hashmap.NewHashMap() }
}

// -------------------------------------------------------------------------------------------------
// ConfigIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Merge configuration data
func (r *Config) MergeConfig(mergeCfg ConfigIfc) *Config {
	r.HashMap.Merge(mergeCfg)
	return r
}

// Get configuration datum whose keys begin with the prefix...
// We also strip the prefix off leaving just the interesting parts
func (r Config) GetSubsetConfig(prefix string) *Config {
	return r.getSubset(prefix, true)
}

// Get configuration datum matching the specified keys, returned as a new Config pointer
func (r Config) GetSubsetKeys(keys *[]string) *Config {
	return &Config{ HashMap: r.HashMap.GetSubset(keys) }
}

// Get configuration datum whose keys DO NOT begin with the prefix...
// We also strip the prefix off leaving just the interesting parts
func (r Config) GetInverseSubsetConfig(prefix string) *Config {
	return r.getSubset(prefix, false)
}

// Dereference any %key% references to our own keys in the supplied string
// returns dereferenced string
func (r Config) DereferenceString(str string) *string {
	keys, err := r.getReferenceKeysFromString(str)
	if nil != err {
		// TODO: Log the error or pass it back to the caller
		return nil
	}
	for _, key := range keys {
		value := r.Get(key)
		if nil == value {  continue }

		ref := fmt.Sprintf("%%%s%%", key)
		str = strings.Replace(str, ref, *value, -1)
	}
	return &str
}

// Dereference any values we have that %reference% keys in the referenceConfig
// returns count of references substituted
func (r *Config) Dereference(referenceConfig ConfigIfc) int {
	if nil == r { return 0 }
	// If no referenceConfig is specified, just dereference against ourselves
	if nil == referenceConfig { return r.Dereference(r) }
	subs := 0
	// For each of our key/value pairs...
	it := r.GetIterator()
	for kvpi := it(); nil != kvpi; kvpi = it() {
		kvp, ok := kvpi.(*hashmap.KeyValuePair)
		if ! ok { continue } // TODO: Error/Warning warranted?
		tstr := referenceConfig.DereferenceString(kvp.Value)
		// Nothing to do if nothing was done...
		if (nil == tstr) || (kvp.Value == *tstr) { continue }
		r.Set(kvp.Key, *tstr)
		subs++
	}
	return subs
}

// Dereference against a list of other referenceConfigs
func (r *Config) DereferenceAll(referenceConfigs ...ConfigIfc) int {
	// If no reference Configs provided, just self-Dereference
	if 0 == len(referenceConfigs) { return r.Dereference(r) }
	subs := 0
	for _, referenceConfig := range referenceConfigs {
		res := r.Dereference(referenceConfig)
		subs += res
	}
	return subs
}

// Dereference until result comes back 0 or maxLoops iterations are completed
// Returns true if fully dereferenced, false, if more refereces may be hiding
func (r *Config) DereferenceLoop(maxLoops int, referenceConfig ConfigIfc) bool {
	localMax := maxLoops
	if localMax > MAX_REFERENCE_DEPTH { localMax = MAX_REFERENCE_DEPTH }
	totalSubs := 0
	for loop := 0; loop < localMax; loop++ {
		subs := r.Dereference(referenceConfig)
		totalSubs += subs
		if subs == 0 { return totalSubs > 0; }
	}
	// TODO: Figure out if new/additional references show up as a result
	// of dereferencing, otherwise, we can't be sure, so return false
	return false
}

/*
// FIXME: The idea of Validation is a good one, however this was just doing a bunch of data handling
//        and returning no useful result; revisit!

func (r *Config) Validate(required, optional *[]string) *Config {
	// TODO: Make and use a new library of string sets since all we are interested in is the keys, not values
	// Of the required keys, what do we actually have?
	requiredSubset := r.HashMap.GetSubset(required)
	requiredKeys := requiredSubset.GetKeys()

	// Of the optional keys, how many do we actually have?
	optionalSubset := r.HashMap.GetSubset(optional)
	optionalKeys := optionalSubset.GetKeys()

	// If there are extra keys remaining after removing required + optional...
	extraKeys := r.Copy().DropSet(required).DropSet(optional).GetKeys()

	return r
}
*/

// -------------------------------------------------------------------------------------------------
// Config Private Interface
// -------------------------------------------------------------------------------------------------

// Get configuration datum whose keys Do/Don't begin with the prefix...
// Return the matches if keepMatches, else return the NON-matches
func (r Config) getSubset(prefix string, keepMatches bool) *Config {
	res := NewConfig()
	it := r.GetIterator()
	for kvpi := it(); nil != kvpi; kvpi = it() {
		kvp, ok := kvpi.(*hashmap.KeyValuePair)
		if ! ok { continue } // TODO: Error/Warning warranted?

		matches := strings.HasPrefix(kvp.Key, prefix)
		if (matches) {
			if ! keepMatches { continue }
			strippedKey := kvp.Key[len(prefix):]
			res.Set(strippedKey, kvp.Value)
		} else {
			if keepMatches { continue }
			res.Set(kvp.Key, kvp.Value)
		}
	}
	return res
}

func (r Config) getReferenceKeysFromString(str string) ([]string, error) {
	runes := []rune(str)
	keys := make([]string, 0)
	inKey := false
	var keyRunes []rune
	for i := 0; i < len(runes); i++ {
		// Marker!
		if runes[i] == '%' {
			// If we're working on a key...
			if inKey {
				// This is the end!
				key := string(keyRunes)
				keys = append(keys, key)
				inKey = false
			} else {
				// This is the beginning!
				keyRunes = make([]rune, 0)
				inKey = true
			}
		} else {
			// If we're working on a key...
			if inKey {
				// Add this rune to it
				keyRunes = append(keyRunes, runes[i])
			}
		}
	}
	var err error
	if inKey { err = fmt.Errorf("Unmatched reference key marker in string '%s'", str) }
	return keys, err
}
