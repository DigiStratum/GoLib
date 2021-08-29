package store

/*

Produce a new GoLib.Config instance from a JSON Object in an ObjectStore.

This lets us retain certain JSON configuration data within an ObjectStore. One way we use this is
to compile configuration details into the build as a generated ObjectStore which pulls the asset
right out of our compiled binary at runtime. 

*/

import (
	"fmt"

	lib "github.com/DigiStratum/GoLib"
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new Config initialized with properties from a JSON Object
func NewObjectStoreConfig(objectStore *ObjectStore, objectPath string) (lib.Config, error) {

	// Check the ObjectStore
	if nil == objectStore {
		err := fmt.Errorf("Config: ObjectStore was nil")
		return nil, err
	}

	// Request the JSON Object
	configObject := objectStore.GetObject(objectPath)
	if nil == configObject {
		err := fmt.Errorf("Config: ObjectStore does not have Object with path: '%s'", objectPath)
		return nil, err
	}

	// Get the JSON Object content
	configJson := configObject.GetContent()
	if nil == configJson {
		err := fmt.Errorf("Config: ObjectStore gave no data for Object with path: '%s'", objectPath)
		return nil, err
	}

// FIXME: validate that a non-empty string results in a non-empty config. It seems like even if Json unmarshal does not generate at error, it might not work. In this case we had a "pattern": "\d+" and it needed to be "\\d+" - withtou the double-escape, the entire JSON struct produced an empty config object, but no error.
	//lib.GetLogger().Trace(fmt.Sprintf("configJson: %s", *configJson))

	// Load up a Config structure from the JSON
	config := lib.NewConfig()
	if err := config.LoadFromJsonStringOrError(configJson); nil != err {
		return nil, fmt.Errorf("Config: Error parsing ObjectStore JSON ('%s'): %s", objectPath, err.Error())
	}
//config.Dump()

	return config, nil
}
