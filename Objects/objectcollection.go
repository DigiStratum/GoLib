package objects

/*

Collection of Objects organized by path; useful for memory-mapped ObjectStores and/or cache type
structures.

FIXME: Add some thread concurrency safety around this things accessor functions

*/

import (
	"sync"
	lib "github.com/DigiStratum/GoLib"
)

type objectMap map[string]*Object

type ObjectCollection struct {
	collection	*objectMap
}

type PathObjectPair struct {
	Path	string
	Obj	*Object
}

type ObjectFieldCondition int

const (
        OFC_NOP ObjectFieldCondition = iota	// No Operation
        OFC_EQ					// Equals
	OFC_NE					// Not Equal
	OFC_ISNULL				// Is Null
	OFC_ISNOTNULL				// Is Not NUll
	OFC_LT					// Less Than
	OFC_LTE					// Less Than or Equal
	OFC_GT					// Greater Than
	OFC_GTE					// Greater Than or Equal
	OFC_TRUTHY				// Represents true
	OFC_FALSEY				// Represents false
	OFC_SW					// string Starts With
	OFC_EW					// string Ends With
	OFC_CONTAINS				// string Contrains
	OFC_NOTCONTAIN				// string does Not Contrain
	OFC_EMPTY				// string is Empty (zero length)
	OFC_NOTEMPTY				// string is Not Empty (non-zero length)
)

type ObjectFieldRule struct {
	Condition	ObjectFieldCondition	// Must be one of the OFC_* constants
	ControlValue	string			// Significance varies with Field and Condition
}

// Make a new one of these
func NewObjectCollection() *ObjectCollection {
	om := make(objectMap)
	objectCollection := ObjectCollection {
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

// Iterate over the objects for this collectino and send all the Path-Object Pairs to a channel
func (oc *ObjectCollection) IterateChannel() <-chan PathObjectPair {
	ch := make(chan PathObjectPair, len(*oc.collection))
	defer close(ch)
	var wg sync.WaitGroup
	wg.Add(1)

	// Fire off a go routine to fill up the channel
	go func() {
		for p, o := range *oc.collection {
			ch <- PathObjectPair{ Path: p, Obj: o }
		}
		wg.Done()
	}()
	wg.Wait()
	return ch
}

