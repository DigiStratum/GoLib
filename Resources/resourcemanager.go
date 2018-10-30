package resources

/*
A Resource Manager to home the logic needed to locate a given Resource within one or more configured
respositories in a given sequence. By scanning UP the list of repositories, starting at 0, whichever
repository matches the Resource path first wins; this is an override model whereby the closest match
to 0 will override everything higher.
*/

type ResourceManager struct {
	// Ordered list of Resource Repositories to find resources within:
	resourceRepositories	[]ResourceRepository,
}

// Make a new one of these!
func NewResourceManager() *ResourceManager {
	rm := ResourceManager{
		resourceRepositories:	make([]ResourceRepository, 0),
	}
	return &rm
}

// Add a Resource repository to the set
// Remember: each addition is lower in priority than the previous!
func (rm *ResourceManager) AddResourceRepository(rr *ResourceRepository) {
	*rm.resourceRepositories = append(*rm.resourceRepositories, *rr)
}

// Get a Resource with the specified path from our set of Resource repositories
func (rm *ResourceManager) GetResource(path string) *Resource {
	// Scan UP the list of Resource repositories in the search for this Resource by path
	for _, repo := range rm.resourceRepositories {
		res := repo.GetResource(path)
		if nil != res { return res }
	}
	return nil
}

