package restapi

/*

Find and retrieve static resources for a given HttpRequestContext

TODO: Convert this into an interface/struct, etc in line with the other abstract classes
TODO: Add LRU cache for this one

*/

import(
	"io/ioutil"
	"path/filepath"
	"errors"
)

func GetResourcePath(resource string) (string, error) {
	// Convert to a file path relative to our resource directory
	var resourcePath = "res/" + resource
	// Use path.filepath to make the path absolute and clean
	var absPath, err = filepath.Abs(resourcePath)
	if nil != err {
		return "", errors.New("Error converting resource path: '" + resourcePath + "'")
	}
	// Now make sure that the path still ends with templatePath
	if absPath[len(absPath) - len(resourcePath):] != resourcePath {
		// Nope! Someone is toying with us!
		return "", errors.New("Resource is not in our collection: '" + resourcePath + "'")
	}
	return resourcePath, nil
}

func ReadResourceAsString(resource string) (string, error) {
	tbuf, err := ReadResourceAsBytes(resource)
	if nil != err {
		return "", err
	}
	return string(tbuf), nil
}

func ReadResourceAsBytes(resource string) ([]byte, error) {
	var tbuf []byte
	var resourcePath, err = GetResourcePath(resource)
	if nil != err {
		return tbuf, err
	}

	tbuf, err = ioutil.ReadFile(resourcePath)
	if nil != err {
		return tbuf, errors.New("Error reading resource: '" + resourcePath + "'")
	}

	return tbuf, nil
}

