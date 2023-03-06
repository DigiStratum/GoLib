package maps

/*
Functions and constructs for handling map objects

*/

// Get the set of keys for maps keyed with string
func Strkeys(m map[string]interface{}) []string {
	mkeys := []string{}
	for key, _ := range m { mkeys = append(mkeys, key) }
	return mkeys
}

