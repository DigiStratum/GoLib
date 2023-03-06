package maps

/*
Functions and constructs for handling map objects

*/

// Get the set of keys for maps keyed with string
// We are expecting m to be a map[string]interface{} of some sort...
func Strkeys(m interface{}) []string {
	mkeys := []string{}
	if mi, ok := m.(map[string]interface{}); ok {
		for key, _ := range mi { mkeys = append(mkeys, key) }
	}
	return mkeys
}

