package resources

/*

Repository interface for Resources; supports all the minimum operations needed to read from it.

TODO: Add some supporting funcs to Resource to get a list of Resources below a given path (i.e. everything in a dir)

*/

import (
	"errors"

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

