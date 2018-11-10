package resources

import (
	"fmt"
	"errors"

	lib "github.com/DigiStratum/GoLib"
)

// Make a new Config initialized with properties from a JSON Repository Resource
func NewRepositoryConfig(repository *Repository, resourcePath string) (*lib.Config, error) {

	// Check the Repository
	if nil == repository {
		err := errors.New("Config: Repository was nil")
		lib.GetLogger().Error(err.Error())
		return nil, err
	}

	// Request the JSON Resource
	configResource := repository.GetResource(resourcePath)
	if nil == configResource {
		err := errors.New(fmt.Sprintf("Config: Repository does not have Resource with path: '%s'", resourcePath))
		lib.GetLogger().Error(err.Error())
		return nil, err
	}

	// Get the JSON Resource content
	configJson := configResource.GetContent()
	if nil == configJson {
		err := errors.New(fmt.Sprintf("Config: Repository gave no data for Resource with path: '%s'", resourcePath))
		lib.GetLogger().Error(err.Error())
		return nil, err
	}

	// Load up a Config structure from the JSON
	config := lib.NewConfig()
	if err := config.LoadFromJsonStringOrError(configJson); nil != err {
		lib.GetLogger().Error(err.Error())
		return nil, err
	}

	return config, nil
}

