package store

/*

ObjectStore for Objects (immutable)

TODO:
 * Add some supporting funcs to ObjectStore to get a list of Objects below a given path (i.e. everything in a dir)

*/

import (
	obj "github.com/DigiStratum/GoLib/Object"
	objc "github.com/DigiStratum/GoLib/Object/collection"
)

type ObjectStoreIfc interface {

	// Get the Object located at this path, or nil if none
	GetObject(path string) (*obj.Object, error)

	// Check whether there is a Object located at this path, true if so
	HasObject(path string) bool
}

type ObjectStore struct {
	collection	*objc.ObjectCollection
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these
func NewObjectStore() *ObjectStore {
	objectStore := ObjectStore{
		collection: objc.NewObjectCollection(),
	}
	return &objectStore
}

// Make a new one of these, preloaded with a ObjectCollection
func NewObjectStorePreloaded(collection *objc.ObjectCollection) *ObjectStore {
	objectStore := ObjectStore{
		collection: collection,
	}
	return &objectStore
}

// -------------------------------------------------------------------------------------------------
// Satisfies ObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r ObjectStore) GetObject(path string) *obj.Object {
	return r.collection.GetObject(path)
}

func (r ObjectStore) HasObject(path string) bool {
	return r.collection.HasObject(path)
}
