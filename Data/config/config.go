package config

/*

TODO:
 * Add Setter func to change the open/close delimiters from defaults
 * Consider a configurable logger - if we wanted Config to log warnings/errors via logger, but
   logger needs Config to initialize, then a circular dependency would be formed, an anti-pattern
   which implies the need for a third resource upon which both depend. what is it? Some kind of
   separate ConfigurableLoggerIfc, a higher level construct which depends on both, but upon which
   neither depend, perhaps...
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
		return nil
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
	if nil == r { return 0 }
	var subs int
	referenceDepth := 0
	for true {
		subs = 0
		referenceDepth++
		if MAX_REFERENCE_DEPTH >= referenceDepth { break }

		// Dereference against ourselves; allows for interesting combos like we have a value
		// and reference config back-references a value that we have to sub it back in, etc.
		subs = r.Dereference(r) // <- Beware, recursion!

		for _, referenceConfig := range referenceConfigs {
			// Add support for dereferencing other data types as they become available via Iterator, at
			// least string would make sense, but Iterator only supports Array|Object currently.
			switch r.GetType() {
				case data.DATA_TYPE_OBJECT:
					it := r.GetIterator()
					for kvpi := it(); nil != kvpi; kvpi = it() {
						if kvp, ok := kvpi.(*data.KeyValuePair); (nil != kvp) && ok {
							if kvp.Value.IsString() {
								tstr := referenceConfig.DereferenceString(kvp.Value.ToString())
								// Nothing to do if nothing changed...
								if (nil == tstr) || (kvp.Value.GetString() == *tstr) { continue }
								r.SetObjectProperty(kvp.Key, data.NewString(*tstr))
								subs++
							}
						}
					}

				case data.DATA_TYPE_ARRAY:
					it := r.GetIterator()
					for ivpi := it(); nil != ivpi; ivpi = it() {
						if ivp, ok := ivpi.(*data.IndexValuePair); (nil != ivp) && ok {
							if ivp.Value.IsString() {
								tstr := referenceConfig.DereferenceString(ivp.Value.ToString())
								// Nothing to do if nothing changed...
								if (nil == tstr) || (ivp.Value.GetString() == *tstr) { continue }
								r.ReplaceArrayValue(ivp.Index, data.NewString(*tstr))
								subs++
							}
						}
					}
			}
		}
		if 0 == subs { break }
	}
	return subs
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

