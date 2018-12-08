package objects

/*

ObjectStore Manager provides an interface to one or more ObjectStores, in sequence, to find and
access Objects by path. By scanning UP the list of ObjectStores, starting at 0, whichever
ObjectStore matches the Object path first wins; this is a short-circuit model whereby the match
closest to index 0 will return immediately without consideration for anything deeper.

*/

import (
	"fmt"
	"errors"

	lib "github.com/DigiStratum/GoLib"
)

type ObjectStoreManager struct {
	// Ordered list of ObjectStores to find Objects within:
	objectStores	[]*ObjectStoreIfc
}

// Make a new one of these!
func NewObjectStoreManager() *ObjectStoreManager {
	osm := ObjectStoreManager{
		objectStores:	make([]*ObjectStoreIfc, 0),
	}
	return &osm
}

// Add an ObjectStore to the set
// Remember: each addition is lower in priority than the previous!
// objectStore parameter must be a pointer to a concrete implementation of an ObjectStoreIfc
// Ref: https://stackoverflow.com/questions/24422810/golang-convert-struct-pointer-to-interface#
func (osm *ObjectStoreManager) AddObjectStore(objectStore interface{}) error {
	l := lib.GetLogger()
	if store, ok := objectStore.(ObjectStoreIfc); ok {
		l.Trace("Adding Object ObjectStore")
		osm.objectStores = append(osm.objectStores, &store)
		return nil
	}
	msg := "Supplied ObjectStore does not satisfy ObjectStoreIfc"
	l.Error(msg)
	return errors.New(msg)
}

// Get an Object with the specified path from our set of ObjectStores
func (osm *ObjectStoreManager) GetObject(path string) *Object {
	// Scan UP the list of ObjectStores in the search for this Object by path
	for _, store := range osm.objectStores {
		res := (*store).GetObject(path)
		if nil != res { return res }
	}
	return nil
}

// Find a scoped ("private"/"public") Object, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) GetMultilingualObject(base string, languages *[]string, relPath string) *Object {
	for _, language := range *languages {
		object := osm.GetObject(fmt.Sprintf("%s/%s/%s", base, language, relPath))
		if nil != object { return object }
	}
	return nil
}

// Find a contextualized Object, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) GetContextualizedObject(context string, languages *[]string, relPath string) *Object {
	return osm.GetMultilingualObject(fmt.Sprintf("public/%s", context), languages, relPath)
}

// Find a scoped ("private"/"public") Object for language or default
// Returns the Object or nil
func (osm *ObjectStoreManager) GetScopedObject(scope string, language string, relPath string) *Object {
	languages := [2]string{ language, "default" }
	langSlice := languages[:]
	return osm.GetMultilingualObject(scope, &langSlice, relPath)
}

// Find a private Object, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) GetPrivateObject(language string, relPath string) *Object {
	return osm.GetScopedObject("private", language, relPath)
}

// Find a (mustache) template type Object, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) GetTemplate(name string, language string) *Object {
	return osm.GetPrivateObject(language, fmt.Sprintf("templates/%s.mustache", name))
}

