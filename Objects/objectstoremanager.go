package objects

/*

ObjectStore Manager provides an interface to a collection of ObjectStores to find and access Objects
by path. There are two methods by which an ObjectStore in the collection may be accessed:

A) Scan Method: By scanning UP the list of ObjectStores, starting at 0, whichever ObjectStore matches
the Object path first wins; this is a short-circuit model whereby the match closest to index 0 will
return immediately without consideration for anything deeper

B) Name Method: By supplying a unique name to the ObjectStore such that it can be addressed by name
directly and without scanning as above

*/

import (
	"fmt"

	lib "github.com/DigiStratum/GoLib"
)

type ObjectStoreManager struct {
	objectStoreCollection	[]*ObjectStoreIfc		// objectStoreCollection[N] -> *ObjectStoreIfc
	objectStoreMap		map[string]*ObjectStoreIfc	// objectStoreMap[name] -> *ObjectStoreIfc
}

// Make a new one of these!
func NewObjectStoreManager() *ObjectStoreManager {
	osm := ObjectStoreManager{
		objectStoreCollection:	make([]*ObjectStoreIfc, 0),
	}
	return &osm
}

// Scan Method functionality
// ------------------------------------------------------------------------------------------------

// Add an unnamed ObjectStore to the set
// The only way to access this ObjectStore is using the Scan Method
// Remember: each addition is lower in priority than the previous!
// objectStore parameter must be a pointer to a concrete implementation of an ObjectStoreIfc
// Ref: https://stackoverflow.com/questions/24422810/golang-convert-struct-pointer-to-interface#
func (osm *ObjectStoreManager) AddObjectStore(objectStore interface{}) error {
	if store, ok := objectStore.(ObjectStoreIfc); ok {
		lib.GetLogger().Trace("Adding Object ObjectStore")
		osm.objectStoreCollection = append(osm.objectStoreCollection, &store)
		return nil
	}
	return lib.GetLogger().Error("Supplied ObjectStore does not satisfy ObjectStoreIfc")
}

// Get an Object with the specified path from our set of ObjectStores
func (osm *ObjectStoreManager) GetObject(path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf("ObjectStoreManager.GetObject('%s')", path))
	// Scan UP the list of ObjectStores in the search for this Object by path
	for _, store := range osm.objectStoreCollection {
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

// Name Method functionality
// ------------------------------------------------------------------------------------------------

// Add an ObjectStore to the Collection and name Map with a unique name as the key
func (osm *ObjectStoreManager) AddNamedObjectStore(name string, objectStore interface{}) error {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.AddNamedObjectStore('%s', ObjectStoreIfc)",
		name,
	))

	// If we already have an ObjectStore with this name...
	if _, ok := osm.objectStoreMap[name]; ok {
		return lib.GetLogger().Error(fmt.Sprintf(
			"There is already an ObjectStore in the collection with the name '%s'",
			name,
		))
	}

	// Try to add it to the colletion...
	err := osm.AddObjectStore(objectStore)

	// If there was no problem adding it to the collection...
	if nil == err {
		// And if it casts correctly...
		if store, ok := objectStore.(ObjectStoreIfc); ok {
			// Capture a reference to it by name into the map!
			osm.objectStoreMap[name] = &store
		} else {
			return lib.GetLogger().Error("Supplied ObjectStore does not satisfy ObjectStoreIfc")
		}
	}
	return err
}

// Get an Object with the specified path from named ObjectStore
func (osm *ObjectStoreManager) GetNamedObjectStoreObject(objectStoreName string, path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.GetNamedObjectStoreObject('%s', '%s')",
		objectStoreName, path,
	))
	// If we can find an ObjectStore with this name
	if store, ok := osm.objectStoreMap[objectStoreName]; ok {
		// Pluck out an Object from it with this path
		res := (*store).GetObject(path)
		return res
	}
	return nil
}

