package resources

/*

A MutableRepository must add the following in addition to the normal Repository requirements.

*/

type MutableRepositoryIfc interface {
	RepositoryIfc	// Inherit the requirements of the Repository interface as well

	// Put the supplied Resource into this Repository at the specified path
	PutResource(path string, resource *Resource) error
}

type MutableRepository struct {
	Repository	// Inherit the properties and functions of Repository
}

// Make a new one of these!
func NewMutableRepository() *MutableRepository {
	repository := MutableRepository{}
	return &repository
}

// Put the supplied Resource into this Repository at the specified path
func (r *MutableRepository) PutResource(path string, resource *Resource) error {
	return r.collection.PutResource(path, resource)
}

