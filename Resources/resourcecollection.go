package resources

/*

Collection of Resources organized by path

*/

import (
	"errors"

	lib "github.com/DigiStratum/GoLib"
)

type ResourceCollection struct {
	collection	*map[string]*Resource
}

// Make a new one of these
func NewResourceCollection() *ResourceCollection {
	rc := make(map[string]*Resource)
	resourceCollection := ResourceCollection{
		collection:	&rc,
	}
	return &resourceCollection
}

// Get a Resource out of the Collection by path
func (rc *ResourceCollection) GetResource(path string) *Resource {
	if rc.HasResource(path) { return (*rc.collection)[path] }
	return nil
}

// Check whether a Resource is in the Collection by path
func (rc *ResourceCollection) HasResource(path string) bool {
	_, ok := (*rc.collection)[path]
	return ok
}

// Put a Resource into the Collection by path
func (rc *ResourceCollection) PutResource(path string, resource *Resource) error {
	if nil == resource {
		l := lib.GetLogger()
		msg := "RepositoryWritableMap.PutResource() - resource can't be nil"
		l.Error(msg)
		return errors.New(msg)
	}
	(*rc.collection)[path] = resource
	return nil
}

