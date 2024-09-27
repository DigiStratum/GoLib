package config

/*

TODO:
 * Add Setter func to change the open/close delimiters from defaults
 * Consider a configurable logger - if we wanted Config to log warnings/errors via logger, but
   logger needs Config to initialize, then a circular dependency would be formed, an anti-pattern
   which implies the need for a third resource upon which both depend. what is it? Some kind of
   separate ConfigurableLoggerIfc, a higher level construct which depends on both, but upon which
   neither depend, perhaps...
 * Add support for dereferencing other data types as they become available via Iterator, at least
   string would make sense, but Iterator only supports Array|Object currently.
*/

import (
	"fmt"
	"strings"

	"GoLib/Data"
)

// Prevent runaway processes with absurd boundaries with an absolute maximum on loop count
const MAX_REFERENCE_DEPTH			= 10
const DEFAULT_REFERENCE_DELIMITER_OPENER	= '%'
const DEFAULT_REFERENCE_DELIMITER_CLOSER	= '%'

type ConfigIfc interface {
	data.DataValueIfc

	DereferenceString(str string) *string
	Dereference(referenceConfigs ...ConfigIfc) int
}

type Config struct {
	*data.DataValue

	refDelimOpener		byte
	refDelimCloser		byte

}

func NewConfig() *Config {
	r := Config{
		DataValue:		data.NewDataValue(),
		refDelimOpener:		DEFAULT_REFERENCE_DELIMITER_OPENER,
		refDelimCloser:		DEFAULT_REFERENCE_DELIMITER_CLOSER,
	}
	return &r
}

// Dereference any %selector% references our keys in supplied string; returns dereferenced string
func (r *Config) DereferenceString(str string) *string {
	selectors, err := r.getReferenceSelectorsFromString(str)
	if nil != err {
		// TODO: Log the error or pass it back to the caller
		return &str
	}
	for _, selector := range selectors {
		value, err := r.Select(selector)
		if (nil == value) || (nil != err) {  continue }

		ref := fmt.Sprintf("%c%s%c", r.refDelimOpener, selector, r.refDelimCloser)
		str = strings.Replace(str, ref, value.ToString(), -1)
	}
	return &str
}

// Dereference values with %reference% selectors against referenceConfig(s); returns num substitutions
// Note that this is a single-pass iteration dereference; if subs comes out > 0 then an additinoal
// pass may be called for to see if more subs are possible (think of subtitutions that themselves
// contain additional keys needing deferencing). A DereferenceAll() method can run N passes up to
// some cycle limit to avoid perpetual loop scenarios.
func (r *Config) Dereference(referenceConfigs ...ConfigIfc) int {
	referenceDepth := 0
	subs := 0
	for (MAX_REFERENCE_DEPTH > referenceDepth) {
		referenceDepth++
		passSubs := r.dereferencePass(referenceConfigs...)
		if 0 == passSubs { break }
		subs += passSubs
	}
	return subs
}

func (r *Config) dereferencePass(referenceConfigs ...ConfigIfc) int {
	subs := 0
	for _, referenceConfig := range referenceConfigs {
		switch r.GetType() {
			case data.DATA_TYPE_OBJECT:
				it := r.GetIterator()
				for kvpi := it(); nil != kvpi; kvpi = it() {
					if kvp, ok := kvpi.(data.KeyValuePair); ok {
						tstr := r.dereferenceOne(kvp.Value.GetString(), referenceConfig)
						if nil == tstr { continue }
						r.SetObjectProperty(kvp.Key, data.NewString(*tstr))
						subs++
					}
				}

			case data.DATA_TYPE_ARRAY:
				it := r.GetIterator()
				for ivpi := it(); nil != ivpi; ivpi = it() {
					if ivp, ok := ivpi.(data.IndexValuePair); ok {
						tstr := r.dereferenceOne(ivp.Value.GetString(), referenceConfig)
						if nil == tstr { continue }
						r.ReplaceArrayValue(ivp.Index, data.NewString(*tstr))
						subs++
					}
				}
		}
	}
	return subs
}

func (r *Config) dereferenceOne(before string, referenceConfig ConfigIfc) *string {
	tstr := referenceConfig.DereferenceString(before)
	// If referenceConfig has/changes nothing...
	if (nil == tstr) || (before == *tstr) {
		// What if we try to dereference against ourselves?
		tstr = r.DereferenceString(before)
		// Nothing to do if nothing changed...
		if (nil == tstr) || (before == *tstr) { return nil }
	}
	return tstr
}

// -------------------------------------------------------------------------------------------------
// Config implementation
// -------------------------------------------------------------------------------------------------

func (r *Config) getReferenceSelectorsFromString(str string) ([]string, error) {
	selectors := make([]string, 0)
	inSelector := false
	var sb strings.Builder
	for _, char := range str {
		// If we found a delimiter char...
		if (char == rune(r.refDelimOpener)) || (char == rune(r.refDelimCloser)) {
			// If we're working on a selector...
			if inSelector && (char == rune(r.refDelimCloser)) {
				// This is the end!
				selectors = append(selectors, sb.String())
				inSelector = false
			} else if ! inSelector && (char == rune(r.refDelimOpener)) {
				// This is the beginning!
				sb.Reset()
				inSelector = true
			}
		} else {
			// If we're working on a key, add this byte to it
			if inSelector { sb.WriteRune(char) }
		}
	}
	var err error
	if inSelector { err = fmt.Errorf("Unmatched selector delimiter in string '%s'", str) }
	return selectors, err
}

