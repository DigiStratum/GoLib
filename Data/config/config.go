package config

/*

TODO:
 * Add Setter func to change the open/close delimiters from defaults

*/

import (
	"fmt"
	"strings"

	"GoLib/Data"
)

const DEFAULT_REFERENCE_DELIMITER_OPENER = '%'
const DEFAULT_REFERENCE_DELIMITER_CLOSER = '%'

type ConfigIfc interface {
	data.DataValueIfc

	DereferenceString(str string) *string
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

// Dereference values with %reference% selectors against referenceConfig; returns num substitutions
// Note that this is a single-pass iteration dereference; if subs comes out > 0 then an additinoal
// pass may be called for to see if more subs are possible (think of subtitutions that themselves
// contain additional keys needing deferencing). A DereferenceAll() method can run N passes up to
// some cycle limit to avoid perpetual loop scenarios.
func (r *Config) Dereference(referenceConfig ConfigIfc) int {
	// FIXME: This can only be done against an Object or Array type of Config; break out for any other primitive type
	if nil == r { return 0 }
	// If no referenceConfig is specified, just dereference against ourselves
	if nil == referenceConfig { return r.Dereference(r) }
	subs := 0
	// For each of our key/value pairs...
	it := r.GetIterator()
	switch r.GetType() {
		case data.DATA_TYPE_OBJECT:
			for kvpi := it(); nil != kvpi; kvpi = it() {
				kvp, ok := kvpi.(*data.KeyValuePair)
				if ! ok { continue } // TODO: Error/Warning warranted?
				if kvp.Value.IsString() {
					tstr := referenceConfig.DereferenceString(kvp.Value.GetString())
					// Nothing to do if nothing was done...
					if (nil == tstr) || (kvp.Value.GetString() == *tstr) { continue }
					r.SetObjectProperty(kvp.Key, data.NewString(*tstr))
					subs++
				}
			}

		case data.DATA_TYPE_ARRAY:
			for ivpi := it(); nil != ivpi; ivpi = it() {
				// FIXME: Array iterator must return the array index as the key!
				ivp, ok := ivpi.(*data.IndexValuePair)
				if ! ok { continue } // TODO: Error/Warning warranted?
				if ivp.Value.IsString() {
					tstr := referenceConfig.DereferenceString(ivp.Value.GetString())
					// Nothing to do if nothing was done...
					if (nil == tstr) || (ivp.Value.GetString() == *tstr) { continue }
					//r.SetObjectProperty(kvp.Key, data.NewString(*tstr))
					// TODO: Implement an in-place update for ARRAY data values!
					r.ReplaceArrayValue(ivp.Index, data.NewString(*tstr))
					subs++
				}
			}
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

