package config

/*

General purpose structured configuration data

This next generation Configuration object supports structured data by embedding GoLib/Data/DataValue
as the underlying storage model.

We support a string dereferencing model to pull other values into the current value for string
building. Dereferencing only works for strings at this time because DataValue does not support
this intrinsically - you can't have, for example, an Integer data type that is represented as a
"%intvalue%" selector which is a String data type (see related TODO below). Dereferencing is
based on delimiter encapsulated identifiers which are handled as DataValue selectors. These can
be nested such that multiple reference Configs can cross-reference each other up to a maximum
reference depth to prevent runaway recursion.

TODO:
 * Consider a configurable logger - if we wanted Config to log warnings/errors via logger, but
   logger needs Config to initialize, then a circular dependency would be formed, an anti-pattern
   which implies the need for a third resource upon which both depend. what is it? Some kind of
   separate ConfigurableLoggerIfc, a higher level construct which depends on both, but upon which
   neither depend, perhaps...
 * Add support for dereferencing other data types as they become available via Iterator, at least
   string would make sense, but Iterator only supports Array|Object currently.
 * Add support for casting dereferenced values to a non-String. e.g. %intvalue:integer% to cause
   the result of the dereferenced string to be stored as a NewInteger({parsed intvalue}) instead
   of storing back as a string.

*/

import (
	"fmt"
	"strings"

	"github.com/DigiStratum/GoLib/Data"
)

// Sane defaults
const DEFAULT_MAX_REFERENCE_DEPTH		= 10
const DEFAULT_REFERENCE_DELIMITER_OPENER	= '%'
const DEFAULT_REFERENCE_DELIMITER_CLOSER	= '%'

type ConfigIfc interface {
	data.DataValueIfc

	SetMaxDepth(max int) *Config
	SetDelimiters(opener, closer byte) *Config

	DereferenceString(str string) (*string, int)
	Dereference(referenceConfigs ...ConfigIfc) int
	MergeConfig(config ConfigIfc) *Config
	CloneConfig() *Config
}

type Config struct {
	*data.DataValue

	refDepthMax		int
	refDelimOpener		byte
	refDelimCloser		byte
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewConfig() *Config {
	return FromDataValue(data.NewDataValue())
}

// Create instance from DataValue (which can have its own factories from various data sources, like
// JSON, YAML, XML etc., string/stream/file/environment, etc)
func FromDataValue(dataValue *data.DataValue) *Config {
	r := &Config{ DataValue: dataValue.Clone() }

	return r.
		SetMaxDepth(DEFAULT_MAX_REFERENCE_DEPTH).
		SetDelimiters(
			DEFAULT_REFERENCE_DELIMITER_OPENER,
			DEFAULT_REFERENCE_DELIMITER_CLOSER,
		)
}

// -------------------------------------------------------------------------------------------------
// ConfigIfc
// -------------------------------------------------------------------------------------------------

func (r *Config) SetMaxDepth(max int) *Config {
	r.refDepthMax = max
	return r
}

func (r *Config) SetDelimiters(opener, closer byte) *Config {
	r.refDelimOpener = opener
	r.refDelimCloser = closer
	return r
}

// Dereference any %selector% references our keys in supplied string; returns dereferenced string
// and num substitutions
func (r *Config) DereferenceString(str string) (*string, int) {
	selectors, err := r.getReferenceSelectorsFromString(str)
	if nil != err {
		// TODO: Log the error or pass it back to the caller
		return &str, 0
	}
	subs := 0
	for _, selector := range selectors {
		value := r.Select(selector)
		if nil == value {  continue }

		ref := fmt.Sprintf("%c%s%c", r.refDelimOpener, selector, r.refDelimCloser)
		str = strings.Replace(str, ref, value.ToString(), -1)
		subs++
	}
	return &str, subs
}

// Dereference values with %reference% selectors against referenceConfig(s); returns num substitutions
// This is a multple-pass iteration dereference; if subs comes out > 0 then an additional pass may
// be called for to see if more subs are possible (think of subtitutions that themselves contain
// additional keys needing deferencing), so we make another pass up to a configured max depth.
// TODO: Technically max passes is not "depth" since this is non-recursive; rename to "derefPassMax"
// TODO: It doesn't seem like the return value int actually provides any utility value. Maybe just
// return self and set immutable - should only need to call this once. Perform any mutations/merges
// needed before Dereferencing, and then it's baked, no more changes!
func (r *Config) Dereference(referenceConfigs ...ConfigIfc) int {
	referenceDepth := 0
	subs := 0
	for (r.refDepthMax > referenceDepth) {
		referenceDepth++
		passSubs := r.dereferencePass(referenceConfigs...)
		if 0 == passSubs { break }
		subs += passSubs
	}
	return subs
}

// Merge properties of passed config into our own embedded data
func (r *Config) MergeConfig(config ConfigIfc) *Config {
	r.DataValue.Merge(config)
	return r
}

func (r *Config) CloneConfig() *Config {
	return &Config{
		DataValue:		r.DataValue.Clone(),
		refDepthMax:		r.refDepthMax,
		refDelimOpener:		r.refDelimOpener,
		refDelimCloser:		r.refDelimCloser,
	}
}

// -------------------------------------------------------------------------------------------------
// Config implementation
// -------------------------------------------------------------------------------------------------

func (r *Config) dereferencePass(referenceConfigs ...ConfigIfc) int {
	subs := 0
	for _, referenceConfig := range referenceConfigs {
		switch r.GetType() {
			case data.DATA_TYPE_OBJECT:
				it := r.GetIterator()
				for kvpi := it(); nil != kvpi; kvpi = it() {
					if kvp, ok := kvpi.(data.KeyValuePair); ok {
						tstr, drsubs := r.dereferenceOne(kvp.Value.GetString(), referenceConfig)
						if (nil == tstr) || (0 == drsubs) { continue }
						r.SetObjectProperty(kvp.Key, data.NewString(*tstr))
						subs = subs + drsubs
					}
				}

			case data.DATA_TYPE_ARRAY:
				it := r.GetIterator()
				for ivpi := it(); nil != ivpi; ivpi = it() {
					if ivp, ok := ivpi.(data.IndexValuePair); ok {
						tstr, drsubs := r.dereferenceOne(ivp.Value.GetString(), referenceConfig)
						if (nil == tstr) || (0 == drsubs) { continue }
						r.ReplaceArrayValue(ivp.Index, data.NewString(*tstr))
						subs = subs + drsubs
					}
				}
		}
	}
	return subs
}

func (r *Config) dereferenceOne(before string, referenceConfig ConfigIfc) (*string, int) {
	tstr, subs := referenceConfig.DereferenceString(before)
	// If referenceConfig has/changes nothing...
	if (nil == tstr) || (before == *tstr) || (0 == subs) {
		// What if we try to dereference against ourselves?
		tstr, subs = r.DereferenceString(before)
		// Nothing to do if nothing changed...
		if (nil == tstr) || (before == *tstr) || (0 == subs) { return nil, 0 }
	}
	return tstr, subs
}

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

