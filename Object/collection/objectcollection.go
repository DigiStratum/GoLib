package collection

/*

Collection of Objects organized by path; useful for memory-mapped ObjectStores and/or cache type
structures.

FIXME:
 * Add some thread concurrency safety around this thing's accessor functions

TODO:
 * Add support to get the list of Object pats (keys)

*/

import (
	"fmt"

	obj "github.com/DigiStratum/GoLib/Object"
)

type ObjectCollectionIfc interface {
	GetObject(path string) *obj.Object
	HasObject(path string) bool
	PutObject(path string, object *obj.Object) error
}

type ObjectCollection struct {
	collection	map[string]*obj.Object
}

/*
// TODO: Move this to a feature branch related to field validity checking
// TODO: Add support to pass in a validation function which will receive the Object and return true/false for validity
// TODO: Add support for wrapping these rules/validators into some reusable/templatized form to avoid repetitive redefinition

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

*/

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these
func NewObjectCollection() *ObjectCollection {
	objectCollection := ObjectCollection {
		collection:	make(map[string]*obj.Object),
	}
	return &objectCollection
}

// -------------------------------------------------------------------------------------------------
// ObjectCollectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get a Object out of the Collection by path
func (r ObjectCollection) GetObject(path string) *obj.Object {
	if r.HasObject(path) { return r.collection[path] }
	return nil
}

// Check whether a Object is in the Collection by path
func (r ObjectCollection) HasObject(path string) bool {
	_, ok := r.collection[path]
	return ok
}

// Put a Object into the Collection by path
func (r *ObjectCollection) PutObject(path string, object *obj.Object) error {
	if nil == object {
		return fmt.Errorf("ObjectCollection.PutObject() - object can't be nil")
	}
	r.collection[path] = object
	return nil
}

// -------------------------------------------------------------------------------------------------
// IterableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Reciever can cast/assert this struct type against the iterator function's return value
type PathObjectPair struct {
	Path	string
	Obj	*obj.Object
}

// Iterate over all of our items, returning each as a *KeyValuePair in the form of an interface{}
func (r ObjectCollection) GetIterator() func () interface{} {
	data_len := len(r.collection)
	keys := make([]string, data_len)
	var idx int = 0
	for k, _ := range r.collection {
		keys[idx] = k
		idx++
	}
	idx = 0
	return func () interface{} {
		// If we're done iterating, return do nothing
		if idx >= data_len { return nil }
		prev_idx := idx
		idx++
		return PathObjectPair{ Path: keys[prev_idx], Obj: r.GetObject(keys[prev_idx]) }
	}
}

