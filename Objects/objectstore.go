package objects

/*

ObjectStore for Objects (immutable)

TODO: Add some supporting funcs to Object to get a list of Objects below a given path (i.e. everything in a dir)

*/

import (
	lib "github.com/DigiStratum/GoLib"
)

type ObjectStoreIfc interface {

	// Configure ObjectStore after it exists (properties are implementation-specific)
	Configure(storeConfig *lib.Config) error

	// Get the Object located at this path, or nil if none
	GetObject(path string) *Object

	// Check whether there is a Object located at this path, true if so
	HasObject(path string) bool
}

type ObjectStore struct {
	collection	*ObjectCollection
}

// Make a new one of these
func NewObjectStore() *ObjectStore {
	objectStore := ObjectStore{
		collection: NewObjectCollection(),
	}
	return &objectStore
}

// Make a new one of these, preloaded with a ObjectCollection
func NewObjectStorePreloaded(collection *ObjectCollection) *ObjectStore {
	objectStore := ObjectStore{
		collection: collection,
	}
	return &objectStore
}

// Satisfies ObjectStoreIfc
// Any ObjectStore implementation should override this as needed
func (os *ObjectStore) Configure(storeConfig *lib.Config) error {
	// There is no configuration data required for this objectStore type
	return nil
}

// Satisfies ObjectStoreIfc
func (os *ObjectStore) GetObject(path string) *Object {
	return os.collection.GetObject(path)
}

// Satisfies ObjectStoreIfc
func (os *ObjectStore) HasObject(path string) bool {
	return os.collection.HasObject(path)
}

