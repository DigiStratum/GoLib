package objects

/*

Collection of Objects organized by path; useful for memory-mapped ObjectStores and/or cache type
structures.

FIXME: Add some thread concurrency safety around this things accessor functions

*/

import (
	lib "github.com/DigiStratum/GoLib"
)

type objectMap map[string]*Object

type ObjectCollection struct {
	collection	*objectMap
}

// Make a new one of these
func NewObjectCollection() *ObjectCollection {
	om := make(objectMap)
	objectCollection := ObjectCollection{
		collection:	&om,
	}
	return &objectCollection
}

// Get a Object out of the Collection by path
func (oc *ObjectCollection) GetObject(path string) *Object {
	if oc.HasObject(path) { return (*oc.collection)[path] }
	return nil
}

// Check whether a Object is in the Collection by path
func (oc *ObjectCollection) HasObject(path string) bool {
	_, ok := (*oc.collection)[path]
	return ok
}

// Put a Object into the Collection by path
func (oc *ObjectCollection) PutObject(path string, object *Object) error {
	if nil == object {
		return lib.GetLogger().Warn("ObjectCollection.PutObject() - object can't be nil")
	}
	(*oc.collection)[path] = object
	return nil
}

