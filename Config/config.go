// DigiStratum GoLib - Config
package config

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
 * Json

*/

import(
	"strings"
	"fmt"

	"github.com/DigiStratum/GoLib/Data/hashmap"
)

// Prevent runaway processes with absurd boundaries with an absolute maximum on loop count
const MAX_REFERENCE_DEPTH = 100

type ConfigIfc interface {
	hashmap.HashMapIfc	// ref: https://www.geeksforgeeks.org/embedding-interfaces-in-golang/
	MergeConfig(mergeCfg ConfigIfc)
	GetSubset(prefix string) *Config
	GetInverseSubset(prefix string) *Config
	DereferenceString(str string) *string
	Dereference(referenceConfig ConfigIfc) int
	DereferenceAll(referenceConfigs ...ConfigIfc) int
	DereferenceLoop(maxLoops int, referenceConfig ConfigIfc) bool
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
func (r *Config) MergeConfig(mergeCfg ConfigIfc) {
	r.HashMap.Merge(mergeCfg)
}

// Get configuration datum whose keys begin with the prefix...
// We also strip the prefix off leaving just the interesting parts
func (r Config) GetSubset(prefix string) *Config {
	return r.getSubset(prefix, true)
}

// Get configuration datum whose keys DO NOT begin with the prefix...
// We also strip the prefix off leaving just the interesting parts
func (r Config) GetInverseSubset(prefix string) *Config {
	return r.getSubset(prefix, false)
}

// Rereference any %key% references to our own keys in the supplied string
// returns dereferenced string
func (r Config) DereferenceString(str string) *string {
/*
	// For each of our key/value pairs...
	for cpair := range r.IterateChannel() {
		// Exit if there are no references in the string
		if ! strings.ContainsRune(str, '%') { break }

		// A reference looks like '%key%'...
		reference := fmt.Sprintf("%%%s%%", cpair.Key)

		// If the referenceConfig value doesn't reference config key, move on...
		if ! strings.Contains(str, reference) { continue }

		// Replace the reference(s) in our value with the values referenced
		tmp := strings.Replace(str, reference, cpair.Value, -1)
		str = tmp
	}
*/
cfgJson, _ := r.ToJson()
fmt.Printf("json: %s\n", *cfgJson)
	keys, err := r.getReferenceKeysFromString(str)
	if nil != err {
		// TODO: Log the error or pass it back to the caller
fmt.Printf("[[ERROR]]\n")
		return nil
	}
	for _, key := range keys {
		value := r.Get(key)
		if nil == value { fmt.Printf("[[NIL]] key='%s'\n", key);continue }
/*
		ref := fmt.Sprintf("%%%s%%", key)
fmt.Printf("replacing key '%s' before: '%s'", ref, str)
		str = strings.Replace(str, ref, *value, -1)
fmt.Printf("after: '%s'\n", str)
*/
	}
	return &str
}

// Dereference any values we have that %reference% keys in the referenceConfig
// returns count of references substituted
func (r *Config) Dereference(referenceConfig ConfigIfc) int {
	if nil == referenceConfig { return 0 }
	subs := 0
	// For each of our key/value pairs...
	//for cpair := range r.IterateChannel() {
	it := r.GetIterator()
	for kvpi := it(); nil != kvpi; kvpi = it() {
		kvp, ok := kvpi.(*hashmap.KeyValuePair)
		if ! ok { continue } // TODO: Error/Warning warranted?
		tstr := referenceConfig.DereferenceString(kvp.Value)
//fmt.Printf("Dereference() '%s': '%s' -> '%s'\n", kvp.Key, kvp.Value, *tstr)
		// Nothing to do if nothing was done...
		if (nil == tstr) || (kvp.Value == *tstr) { continue }
		r.Set(kvp.Key, *tstr)
		subs++
	}
	return subs
}

// Dereference against a list of other referenceConfigs
func (r *Config) DereferenceAll(referenceConfigs ...ConfigIfc) int {
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
//fmt.Printf("Loop:%d, subs:%d, totalSubs:%d, localMax:%d\n", loop, subs, totalSubs, localMax)
		if subs == 0 { return totalSubs > 0; }
	}
	// TODO: Figure out if new/additional references show up as a result
	// of dereferencing, otherwise, we can't be sure, so return false
//fmt.Printf("totalSubs:%d, localMax:%d\n", totalSubs, localMax)
	return false
}

// -------------------------------------------------------------------------------------------------
// Config Private Interface
// -------------------------------------------------------------------------------------------------

// Get configuration datum whose keys Do/Don't begin with the prefix...
// Return the matches if keepMatches, else return the NON-matches
func (r Config) getSubset(prefix string, keepMatches bool) *Config {
	res := NewConfig()
	//for pair := range r.IterateChannel() {
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
				//fmt.Printf("found key:'%s'\n", key)
				inKey = false
			} else {
				// This is the begining!
				keyRunes = make([]rune, len(str))
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
