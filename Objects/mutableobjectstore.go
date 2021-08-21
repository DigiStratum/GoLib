package objects

/*

A MutableObjectStore extends an ObjectStore with write capability

TODO:
 * Add support for Delete/Drop object
 * Add support to synchronize with source if we put/delete in memory and it is now different from
   source (write-through? implementation dependent?)

*/

type MutableObjectStoreIfc interface {
	ObjectStoreIfc	// Embed ObjectStore interface

	// Put the supplied Object into this ObjectStore at the specified path
	PutObject(path string, object ObjectIfc) error
}

type MutableObjectStore struct {
	ObjectStore	// Embed ObjectStore properties
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
func (r *MutableObjectStore) PutObject(path string, object ObjectIfc) error {
	return r.collection.PutObject(path, object)
}
