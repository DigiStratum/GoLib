package resources

// Map of resource path to the Resource and its properties
// TODO: Add some supporting funcs to Resource to get a list of Resources below a given path (i.e. everything in a dir)
type ResourceMap map[string]*Resource

func (rm ResourceMap) GetResource(path string) *Resource {
	if r, ok := rm[path]; ok {
		return r
	}
	return nil
}

