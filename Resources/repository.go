package resources

/*

Repository for Resources (immutable)

TODO: Add some supporting funcs to Resource to get a list of Resources below a given path (i.e. everything in a dir)

*/

import (
	lib "github.com/DigiStratum/GoLib"
)

type RepositoryIfc interface {

	// Configure Repository after it exists (properties are implementation-specific)
	Configure(repoConfig *lib.Config) error

	// Get the Resource located at this path, or nil if none
	GetResource(path string) *Resource

	// Check whether there is a Resource located at this path, true if so
	HasResource(path string) bool
}

type Repository struct {
	collection	*ResourceCollection
}

// Make a new one of these
func NewRepository() *Repository {
	repository := Repository{}
	return &repository
}

// Make a new one of these, preloaded with a ResourceCollection
func NewRepositoryPreloaded(collection *ResourceCollection) *Repository {
	repository := Repository{
		collection: collection,
	}
	return &repository
}

// Satisfies RepositoryIfc
// Any Repository implementation should override this as needed
func (r Repository) Configure(repoConfig *lib.Config) error {
	// There is no configuration data required for this repository type
	return nil
}

// Satisfies RepositoryIfc
func (r Repository) GetResource(path string) *Resource {
	return r.collection.GetResource(path)
}

// Satisfies RepositoryIfc
func (r Repository) HasResource(path string) bool {
	return r.collection.HasResource(path)
}

