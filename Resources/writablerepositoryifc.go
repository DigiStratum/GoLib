package resources

/*

A WritableRepository must add the following in addition to the normal Repository IFC requirements.

*/

type WritableRepositoryIfc interface {

	// Put the supplied Resource into this Repository at the specified path
	PutResource(resource *Resource, path string) error
}

