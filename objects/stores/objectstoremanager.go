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

type ObjectStoreManagerIfc {
	AddObjectStore(objectStore ObjectStoreIfc)
	AddNamedObjectStore(name string, objectStore ObjectStoreIfc) error
	GetNamedObjectStoreObject(objectStoreName string, path string) *Object
	FindMultilingualObject(base string, languages *[]string, relPath string) *Object
	FindObject(scope string, possibleContexts *[]string, languages *[]string, relPath string) *Object
	FindContextualizedObject(scope string, context string, languages *[]string, relPath string) *Object
	FindNamedObjectStoreContextualizedObject(objectStoreName string, scope string, context string, languages *[]string, relPath string) *Object
	FindScopedObject(scope string, language string, relPath string) *Object
	FindNamedObjectStoreScopedObject(objectStoreName string, scope string, language string, relPath string) *Object
	FindPrivateObject(language string, relPath string) *Object
	FindNamedObjectStorePrivateObject(objectStoreName string, language string, relPath string) *Object
	FindTemplate(language string, name string) *Object
	FindNamedObjectStoreTemplate(objectStoreName string, language string, name string) *Object
	GetObjectCollection(path string) *ObjectCollection
	GetNamedObjectStoreObjectCollection(objectStoreName string, path string) *ObjectCollection
}

type ObjectStoreManager struct {
	objectStoreCollection	[]*ObjectStoreIfc		// objectStoreCollection[N] -> *ObjectStoreIfc
	objectStoreMap		map[string]*ObjectStoreIfc	// objectStoreMap[name] -> *ObjectStoreIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObjectStoreManager() *ObjectStoreManager {
	osm := ObjectStoreManager{
		objectStoreCollection:	make([]*ObjectStoreIfc, 0),
	}
	return &osm
}

// -------------------------------------------------------------------------------------------------
// ObjectStoreManagerIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Add an unnamed ObjectStore to the set

// Accessible only by Scan method, each addition is lower in priority than the previous!
func (r *ObjectStoreManager) AddObjectStore(objectStore ObjectStoreIfc) {
	lib.GetLogger().Trace("Adding Object ObjectStore")
	r.objectStoreCollection = append(r.objectStoreCollection, &objectStore)
}

// Accessible by Scan OR Name method!
func (r *ObjectStoreManager) AddNamedObjectStore(name string, objectStore ObjectStoreIfc) error {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.AddNamedObjectStore('%s', ObjectStoreIfc)",
		name,
	))

	// If we already have an ObjectStore with this name...
	if _, ok := r.objectStoreMap[name]; ok {
		return lib.GetLogger().Error(fmt.Sprintf(
			"There is already an ObjectStore in the collection with the name '%s'",
			name,
		))
	}

	// Try to add it to the colletion...
	r.AddObjectStore(objectStore)

	// Capture a reference to it by name into the map!
	r.objectStoreMap[name] = &objectStore
}

// Get an Object with the specified path from our set of ObjectStores

// Scan method
func (r ObjectStoreManager) GetObject(path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf("ObjectStoreManager.GetObject('%s')", path))
	return r.getObject(nil, path)
}

// Name method
func (r ObjectStoreManager) GetNamedObjectStoreObject(objectStoreName string, path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.GetNamedObjectStoreObject('%s', '%s')",
		objectStoreName, path,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := r.objectStoreMap[objectStoreName]; ok {
		return r.getObject(objectStore, path)
	}
	return nil
}

// Find an Object relative to base path, facet on language (default="default")

// Scan method; Returns the Object or nil
func (r ObjectStoreManager) FindMultilingualObject(base string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindMultilingualObject('%s', [%d]string, '%s')",
		base,
		len(*languages),
		relPath,
	))
	return r.findMultilingualObject(nil, base, languages, relPath)
}

// Find a scoped (public/private), contextualized Object, facet on language (default="default")

func (r ObjectStoreManager) FindObject(scope string, possibleContexts *[]string, languages *[]string, relPath string) *Object {
	for _, context := range *possibleContexts {
                lib.GetLogger().Trace(fmt.Sprintf("Trying path: %s/{language}/%s", context, relPath))
                object := r.FindContextualizedObject(
                        scope,
                        context,
                        languages,
                        relPath,
                )
                if nil == object { continue }
		return object
        }
	return nil
}

// Scan method; Returns the Object or nil
func (r ObjectStoreManager) FindContextualizedObject(scope string, context string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindContextualizedObject('%s', '%s', [%d]string, '%s')",
		scope,
		context,
		len(*languages),
		relPath,
	))
	return r.findContextualizedObject(nil, scope, context, languages, relPath)
}

// Name method; Returns the Object or nil
func (r ObjectStoreManager) FindNamedObjectStoreContextualizedObject(objectStoreName string, scope string, context string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindNamedObjectStoreContextualizedObject('%s', '%s', '%s', [%d]string, '%s')",
		objectStoreName,
		scope,
		context,
		len(*languages),
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := r.objectStoreMap[objectStoreName]; ok {
		return r.findContextualizedObject(objectStore, scope, context, languages, relPath)
	}
	return nil
}

// Find a scoped ("private"/"public") Object for language

// Scan method; Returns the Object or nil
func (r ObjectStoreManager) FindScopedObject(scope string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindScopedObject('%s', '%s', '%s')",
		scope,
		language,
		relPath,
	))
	return r.findScopedObject(nil, scope, language, relPath)
}

// Name method; Returns the Object or nil
func (r ObjectStoreManager) FindNamedObjectStoreScopedObject(objectStoreName string, scope string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindScopedObject('%s', '%s', '%s', '%s')",
		objectStoreName,
		scope,
		language,
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := r.objectStoreMap[objectStoreName]; ok {
		return r.findScopedObject(objectStore, scope, language, relPath)
	}
	return nil
}

// Find a private Object, facet on language (default="default")

// Scan method; Returns the Object or nil
func (r ObjectStoreManager) FindPrivateObject(language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindPrivateObject('%s', '%s')",
		language,
		relPath,
	))
	return r.findScopedObject(nil, "private", language, relPath)
}

// Name method; Returns the Object or nil
func (r ObjectStoreManager) FindNamedObjectStorePrivateObject(objectStoreName string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindPrivateObject('%s', '%s', '%s')",
		objectStoreName,
		language,
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := r.objectStoreMap[objectStoreName]; ok {
		return r.findScopedObject(objectStore, "private", language, relPath)
	}
	return nil
}

// Find a named (mustache) template type Object, facet on language (default="default")

// Scan method; Returns the Object or nil
func (r ObjectStoreManager) FindTemplate(language string, name string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindTemplate('%s', '%s')",
		language,
		name,
	))
	return r.FindPrivateObject(language, fmt.Sprintf("templates/%s.mustache", name))
}

// Name method; Returns the Object or nil
func (r rObjectStoreManager) FindNamedObjectStoreTemplate(objectStoreName string, language string, name string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindTemplate('%s', '%s', '%s')",
		objectStoreName,
		language,
		name,
	))
	return r.FindNamedObjectStorePrivateObject(objectStoreName, language, fmt.Sprintf("templates/%s.mustache", name))
}

// Get an ObjectCollection filled with Objects matching the filter criteria

// Scan method; fills the collection from the first ObjectStore where a match is found
func (r ObjectStoreManager) GetObjectCollection(path string) *ObjectCollection {
	// TODO - implement this!
	return nil
}

// Name method; fills the collection from the named ObjectStore
func (r ObjectStoreManager) GetNamedObjectStoreObjectCollection(objectStoreName string, path string) *ObjectCollection {
	// TODO - implement this!
	return nil
}

// -------------------------------------------------------------------------------------------------
// ObjectStoreManagerIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Reusable logic which supports both, called from other functions below
func (r ObjectStoreManager) getObject(objectStore *ObjectStoreIfc, path string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf("ObjectStoreManager.getObject(*ObjectStoreIfc, '%s')", path))
	// No particular ObjectStoreIfc specified..?
	if nil == objectStore {
		// Scan UP the list of ObjectStores in the search for this Object by path
		for _, store := range r.objectStoreCollection {
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

// Name method; Returns the Object or nil
func (r ObjectStoreManager) FindNamedObjectStoreMultilingualObject(objectStoreName string, base string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindNamedObjectStoreMultilingualObject('%s', '%s', [%d]string, '%s')",
		objectStoreName,
		base,
		len(*languages),
		relPath,
	))
	// If we can find an ObjectStore with this name
	if objectStore, ok := r.objectStoreMap[objectStoreName]; ok {
		return r.findMultilingualObject(objectStore, base, languages, relPath)
	}
	return nil
}

// Logic; Returns the Object or nil
func (r ObjectStoreManager) findMultilingualObject(objectStore *ObjectStoreIfc, base string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.findMultilingualObject(*ObjectStoreIfc, '%s', [%d]string, '%s')",
		base,
		len(*languages),
		relPath,
	))
	for _, language := range *languages {
		object := r.getObject(objectStore, fmt.Sprintf("%s/%s/%s", base, language, relPath))
		if nil != object { return object }
	}
	// One last try for "default" language
	return r.getObject(objectStore, fmt.Sprintf("%s/%s/%s", base, "default", relPath))
}

// Logic; Returns the Object or nil
func (r ObjectStoreManager) findContextualizedObject(objectStore *ObjectStoreIfc, scope string, context string, languages *[]string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.FindContextualizedObject(*ObjectStoreIfc, '%s', '%s', [%d]string, '%s')",
		scope,
		context,
		len(*languages),
		relPath,
	))
	base := scope
	if len(context) > 0 { base = fmt.Sprintf("%s/%s", scope, context) }
	return r.findMultilingualObject(objectStore, base, languages, relPath)
}

// Logic; Returns the Object or nil
func (r *ObjectStoreManager) findScopedObject(objectStore *ObjectStoreIfc, scope string, language string, relPath string) *Object {
	lib.GetLogger().Trace(fmt.Sprintf(
		"ObjectStoreManager.findScopedObject(*ObjectStoreIfc, '%s', '%s', '%s')",
		scope,
		language,
		relPath,
	))
	languages := []string{ language }
	return r.findContextualizedObject(objectStore, scope, "", &languages, relPath)
}
