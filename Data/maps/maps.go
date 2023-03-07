package maps

/*
Functions and constructs for handling map objects

*/

import (
	"reflect"
)

// Get the set of keys for maps keyed with string
// We are expecting m to be a map[string]interface{} of some sort...
func Strkeys(m interface{}) []string {
	mkeys := []string{}
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Map {
		keyValues := v.MapKeys()
		for _, kV := range keyValues {
			mkeys = append(mkeys, kV.String())
		}
	}
	return mkeys
}

