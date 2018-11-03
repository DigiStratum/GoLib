package resources

/*

Resource Repository for Compiled Resources

These will probably be hard-coded somewhere and/or code-generated and initialized right into the map
with read-only access.

*/

type RepositoryCompiledMap {
	ResourceMap	// Inherit the properties and functions of ResourceMap
}

// Satisfies ReadableRepositoryIfc
func (rcm RepositoryCompiledMap) GetResource(path string) *Resource {
	if rcm.HasResource(path) { return rcm[path] }
	return nil
}

// Satisfies ReadableRepositoryIfc
func (rcm RepositoryCompiledMap) HasResource(path string) bool {
	_, ok := rcm[path]
	return ok
}

