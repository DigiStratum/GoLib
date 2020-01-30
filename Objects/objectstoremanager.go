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

// -------------------------------------------------------------------------------------------------

// Add an unnamed ObjectStore to the set

// Accessible only by Scan method, each addition is lower in priority than the previous!
func (osm *ObjectStoreManager) AddObjectStore(objectStore interface{}) error {
	if store, ok := objectStore.(ObjectStoreIfc); ok {
		lib.GetLogger().Trace("Adding Object ObjectStore")
		osm.objectStoreCollection = append(osm.objectStoreCollection, &store)
		return nil
	}
	return lib.GetLogger().Error("Supplied ObjectStore does not satisfy ObjectStoreIfc")
}

// Accessible by Scan OR Name method!
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

// -------------------------------------------------------------------------------------------------

// Get an Object with the specified path from our set of ObjectStores

// Scan method
func (osm *ObjectStoreManager) GetObject(path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf("ObjectStoreManager.GetObject('%s')", path))
	return osm.getObject(nil, path)
}

// Name method
func (osm *ObjectStoreManager) GetNamedObjectStoreObject(objectStoreName string, path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.GetNamedObjectStoreObject('%s', '%s')",
		objectStoreName, path,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := osm.objectStoreMap[objectStoreName]; ok {
		return osm.getObject(objectStore, path)
	}
	return nil
}

// Reusable logic which supports both, called from other functions below
func (osm *ObjectStoreManager) getObject(objectStore *ObjectStoreIfc, path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf("ObjectStoreManager.getObject(*ObjectStoreIfc, '%s')", path))
	// No particular ObjectStoreIfc specified..?
	if nil == objectStore {
		// Scan UP the list of ObjectStores in the search for this Object by path
		for _, store := range osm.objectStoreCollection {
			res := (*store).GetObject(path)
			if nil != res { return res }
		}
		return nil
	} else {
		// Pluck out an Object from the specified ObjectStoreIfc with this path
		res := (*objectStore).GetObject(path)
		return res
	}
}

// -------------------------------------------------------------------------------------------------

// Find an Object relative to base path, facet on language (default="default")

// Scan method; Returns the Object or nil
func (osm *ObjectStoreManager) FindMultilingualObject(base string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindMultilingualObject('%s', [%d]string, '%s')",
		base,
		len(*languages),
		relPath,
	))
	return osm.findMultilingualObject(nil, base, languages, relPath)
}

// Name method; Returns the Object or nil
func (osm *ObjectStoreManager) FindNamedObjectStoreMultilingualObject(objectStoreName string, base string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindNamedObjectStoreMultilingualObject('%s', '%s', [%d]string, '%s')",
		objectStoreName,
		base,
		len(*languages),
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := osm.objectStoreMap[objectStoreName]; ok {
		return osm.findMultilingualObject(objectStore, base, languages, relPath)
	}
	return nil
}

// Logic; Returns the Object or nil
func (osm *ObjectStoreManager) findMultilingualObject(objectStore *ObjectStoreIfc, base string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.findMultilingualObject(*ObjectStoreIfc, '%s', [%d]string, '%s')",
		base,
		len(*languages),
		relPath,
	))
	for _, language := range *languages {
		object := osm.getObject(objectStore, fmt.Sprintf("%s/%s/%s", base, language, relPath))
		if nil != object { return object }
	}
	// One last try for "default" language
	return osm.getObject(objectStore, fmt.Sprintf("%s/%s/%s", base, "default", relPath))
}

// -------------------------------------------------------------------------------------------------

// Find a scoped (public/private), contextualized Object, facet on language (default="default")

// Scan method; Returns the Object or nil
func (osm *ObjectStoreManager) FindContextualizedObject(scope string, context string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindContextualizedObject('%s', '%s', [%d]string, '%s')",
		scope,
		context,
		len(*languages),
		relPath,
	))
	return osm.findContextualizedObject(nil, scope, context, languages, relPath)
}

// Name method; Returns the Object or nil
func (osm *ObjectStoreManager) FindNamedObjectStoreContextualizedObject(objectStoreName string, scope string, context string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindNamedObjectStoreContextualizedObject('%s', '%s', '%s', [%d]string, '%s')",
		objectStoreName,
		scope,
		context,
		len(*languages),
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := osm.objectStoreMap[objectStoreName]; ok {
		return osm.findContextualizedObject(objectStore, scope, context, languages, relPath)
	}
	return nil
}

// Logic; Returns the Object or nil
func (osm *ObjectStoreManager) findContextualizedObject(objectStore *ObjectStoreIfc, scope string, context string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindContextualizedObject(*ObjectStoreIfc, '%s', '%s', [%d]string, '%s')",
		scope,
		context,
		len(*languages),
		relPath,
	))
	base := scope
	if len(context) > 0 { base = fmt.Sprintf("%s/%s", scope, context) }
	return osm.findMultilingualObject(objectStore, base, languages, relPath)
}

// -------------------------------------------------------------------------------------------------

// Find a scoped ("private"/"public") Object for language

// Scan method; Returns the Object or nil
func (osm *ObjectStoreManager) FindScopedObject(scope string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindScopedObject('%s', '%s', '%s')",
		scope,
		language,
		relPath,
	))
	return osm.findScopedObject(nil, scope, language, relPath)
}

// Name method; Returns the Object or nil
func (osm *ObjectStoreManager) FindNamedObjectStoreScopedObject(objectStoreName string, scope string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindScopedObject('%s', '%s', '%s', '%s')",
		objectStoreName,
		scope,
		language,
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := osm.objectStoreMap[objectStoreName]; ok {
		return osm.findScopedObject(objectStore, scope, language, relPath)
	}
	return nil
}

// Logic; Returns the Object or nil
func (osm *ObjectStoreManager) findScopedObject(objectStore *ObjectStoreIfc, scope string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.findScopedObject(*ObjectStoreIfc, '%s', '%s', '%s')",
		scope,
		language,
		relPath,
	))
	languages := []string{ language }
	return osm.findContextualizedObject(objectStore, scope, "", &languages, relPath)
}

// -------------------------------------------------------------------------------------------------

// Find a private Object, facet on language (default="default")

// Scan method; Returns the Object or nil
func (osm *ObjectStoreManager) FindPrivateObject(language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindPrivateObject('%s', '%s')",
		language,
		relPath,
	))
	return osm.findScopedObject(nil, "private", language, relPath)
}

// Name method; Returns the Object or nil
func (osm *ObjectStoreManager) FindNamedObjectStorePrivateObject(objectStoreName string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindPrivateObject('%s', '%s', '%s')",
		objectStoreName,
		language,
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := osm.objectStoreMap[objectStoreName]; ok {
		return osm.findScopedObject(objectStore, "private", language, relPath)
	}
	return nil
}

// -------------------------------------------------------------------------------------------------

// Find a named (mustache) template type Object, facet on language (default="default")

// Scan method; Returns the Object or nil
func (osm *ObjectStoreManager) FindTemplate(language string, name string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindTemplate('%s', '%s')",
		language,
		name,
	))
	return osm.FindPrivateObject(language, fmt.Sprintf("templates/%s.mustache", name))
}

// Name method; Returns the Object or nil
func (osm *ObjectStoreManager) FindNamedObjectStoreTemplate(objectStoreName string, language string, name string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindTemplate('%s', '%s', '%s')",
		objectStoreName,
		language,
		name,
	))
	return osm.FindNamedObjectStorePrivateObject(objectStoreName, language, fmt.Sprintf("templates/%s.mustache", name))
}

// -------------------------------------------------------------------------------------------------


// Get an ObjectCollection filled with Objects matching the filter criteria

// Scan method; fills the collection from the first ObjectStore where a match is found
func (osm *ObjectStoreManager) GetObjectCollection(path string) *ObjectCollection {
	// TODO - implement this!
	return nil
}

// Name method; fills the collection from the named ObjectStore
func (osm *ObjectStoreManager) GetNamedObjectStoreObjectCollection(objectStoreName string, path string) *ObjectCollection {
	// TODO - implement this!
	return nil
}

