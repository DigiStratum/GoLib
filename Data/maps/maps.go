package maps

/*
Functions and constructs for handling map objects

*/

func strkeys(m map[string]interface{}) []string {
	mkeys := []string{}
	for key, _ := range m { mkeys = append(mkeys, key) }
	return mkeys
}

