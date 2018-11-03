package resources

/*

Repository Manager provides an interface to one or more Resource Repositories, in sequence, to find
and access Resources by path. By scanning UP the list of repositories, starting at 0, whichever
repository matches the Resource path first wins; this is an override model whereby the closest match
to 0 will override everything higher.

*/

type RepositoryManager struct {
	// Ordered list of Resource Repositories to find resources within:
	repositories	[]RepositoryIfc
}

// Make a new one of these!
func NewRepositoryManager() *RepositoryManager {
	rm := RepositoryManager{
		repositories:	make([]RepositoryIfc, 0),
	}
	return &rm
}

// Add a Resource repository to the set
// Remember: each addition is lower in priority than the previous!
func (rm *RepositoryManager) AddRepository(repository *RepositoryIfc) {
	rm.repositories = append(rm.repositories, *repository)
}

// Get a Resource with the specified path from our set of Resource repositories
func (rm *RepositoryManager) GetResource(path string) *Resource {
	// Scan UP the list of Resource repositories in the search for this Resource by path
	for _, repo := range rm.repositories {
		res := repo.GetResource(path)
		if nil != res { return res }
	}
	return nil
}

// Find a scoped ("private"/"public") resource, facet on language (default="default")
// Returns the Resource or nil
func (rm *RepositoryManager) GetScopedResource(scope string, relPath string, language string) *Resource {
	possibilities := [...]string{ language, "default" }
	for _, possibility := range possibilities {
		resource := rm.GetResource(scope + "/" + possibility + "/" + relPath)
		if nil != resource { return resource }
	}
	return nil
}

// Find a private resource, facet on language (default="default")
// Returns the Resource or nil
func (rm *RepositoryManager) GetPrivateResource(relPath string, language string) *Resource {
	return rm.GetScopedResource("private", relPath, language)
}

// Find a public resource, facet on language (default="default")
// Returns the Resource or nil
func (rm *RepositoryManager) GetPublicResource(relPath string, language string) *Resource {
	return rm.GetScopedResource("public", relPath, language)
}

// Find a (mustache) template type resource, facet on language (default="default")
// Returns the Resource or nil
func (rm *RepositoryManager) GetTemplate(name string, language string) *Resource {
	return rm.GetPrivateResource("templates/" + name + ".mustache", language)
}

