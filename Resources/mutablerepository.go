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
	// Ref: https://travix.io/type-embedding-in-go-ba40dd4264df
	// Repository is embedded into MutableRepository;
	// we pass a Repository in to leverage its own initializer
	repo := NewRepository()
	repository := MutableRepository{
		*repo,
	}
	return &repository
}

// Put the supplied Resource into this Repository at the specified path
func (r *MutableRepository) PutResource(path string, resource *Resource) error {
	return r.collection.PutResource(path, resource)
}

