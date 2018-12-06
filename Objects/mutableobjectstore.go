package objects

/*

A MutableObjectStore must add the following in addition to the normal ObjectStore requirements.

*/

type MutableObjectStoreIfc interface {
	ObjectStoreIfc	// Inherit the requirements of the ObjectStore interface as well

	// Put the supplied Object into this ObjectStore at the specified path
	PutObject(path string, object *Object) error
}

type MutableObjectStore struct {
	ObjectStore	// Inherit the properties and functions of ObjectStore
}

// Make a new one of these!
func NewMutableObjectStore() *MutableObjectStore {
	// Ref: https://travix.io/type-embedding-in-go-ba40dd4264df
	// ObjectStore is embedded into MutableObjectStore;
	// we pass a ObjectStore in to leverage its own initializer
	store := NewObjectStore()
	objectStore := MutableObjectStore{
		*store,
	}
	return &objectStore
}

// Put the supplied Object into this ObjectStore at the specified path
func (os *MutableObjectStore) PutObject(path string, object *Object) error {
	return os.collection.PutObject(path, object)
}

