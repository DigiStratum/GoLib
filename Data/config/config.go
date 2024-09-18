package config


const DEFAULT_REFERENCE_DELIMITER_OPENER = "%%"
const DEFAULT_REFERENCE_DELIMITER_CLOSER = "%%"

type ConfigIfc interface {
	DereferenceString(str string) *string
}

type Config struct {
	refDelimOpener		string
	refDelimCloser		string

}

func NewConfig() *Config {
	r := Config{
		refDelimOpener:		DEFAULT_REFERENCE_DELIMITER_OPENER,
		refDelimCloser:		DEFAULT_REFERENCE_DELIMITER_CLOSER,
	}
}

// Dereference any %selector% references to our own keys in the supplied string
// returns dereferenced string
func (r *Config) DereferenceString(str string) *string {
	keys, err := getReferenceKeysFromString(str)
	if nil != err {
		// TODO: Log the error or pass it back to the caller
		return nil
	}
	for _, key := range keys {
		value := r.Get(key)
		if nil == value {  continue }

		ref := fmt.Sprintf("%s%s%s", r.refDelimOpener, key, r.refDelimCloser)
		str = strings.Replace(str, ref, *value, -1)
	}
	return &str
}

// Dereference any values we have that %reference% selectors in the referenceConfig
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

// -------------------------------------------------------------------------------------------------
// Helper Functions
// -------------------------------------------------------------------------------------------------

func getReferenceKeysFromString(str string) ([]string, error) {
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

