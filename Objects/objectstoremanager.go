package objects

/*

ObjectStore Manager provides an interface to one or more ObjectStores, in sequence, to find and
access Objects by path. By scanning UP the list of ObjectStores, starting at 0, whichever
ObjectStore matches the Object path first wins; this is a short-circuit model whereby the match
closest to index 0 will return immediately without consideration for anything deeper.

*/

import (
	"fmt"

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
	if store, ok := objectStore.(ObjectStoreIfc); ok {
		lib.GetLogger().Trace("Adding Object ObjectStore")
		osm.objectStores = append(osm.objectStores, &store)
		return nil
	}
	return lib.GetLogger().Error("Supplied ObjectStore does not satisfy ObjectStoreIfc")
}

// Get an Object with the specified path from our set of ObjectStores
func (osm *ObjectStoreManager) GetObject(path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf("ObjectStoreManager.GetObject('%s')", path))
	// Scan UP the list of ObjectStores in the search for this Object by path
	for _, store := range osm.objectStores {
		res := (*store).GetObject(path)
		if nil != res { return res }
	}
	return nil
}

// Find an Object relative to base path, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) FindMultilingualObject(base string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindMultilingualObject('%s', [%d]string, '%s')",
		base,
		len(*languages),
		relPath,
	))
	for _, language := range *languages {
		object := osm.GetObject(fmt.Sprintf("%s/%s/%s", base, language, relPath))
		if nil != object { return object }
	}
	// One last try for "default" language
	return osm.GetObject(fmt.Sprintf("%s/%s/%s", base, "default", relPath))
	return nil
}

// Find a scoped (public/private), contextualized Object, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) FindContextualizedObject(scope string, context string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindContextualizedObject('%s', '%s', [%d]string, '%s')",
		scope,
		context,
		len(*languages),
		relPath,
	))
	base := scope
	if len(context) > 0 { base = fmt.Sprintf("%s/%s", scope, context) }
	return osm.FindMultilingualObject(base, languages, relPath)
}

// Find a scoped ("private"/"public") Object for language
// Returns the Object or nil
func (osm *ObjectStoreManager) FindScopedObject(scope string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindScopedObject('%s', '%s', '%s')",
		scope,
		language,
		relPath,
	))
	languages := []string{ language }
	return osm.FindContextualizedObject(scope, "", &languages, relPath)
}

// Find a private Object, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) FindPrivateObject(language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindPrivateObject('%s', '%s')",
		language,
		relPath,
	))
	return osm.FindScopedObject("private", language, relPath)
}

// Find a named (mustache) template type Object, facet on language (default="default")
// Returns the Object or nil
func (osm *ObjectStoreManager) FindTemplate(language string, name string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindTemplate('%s', '%s')",
		language,
		name,
	))
	return osm.FindPrivateObject(language, fmt.Sprintf("templates/%s.mustache", name))
}

