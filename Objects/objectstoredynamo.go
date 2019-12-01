package objects

/*

ObjectStore for AWS Dynamo NoSQL Database service

Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/

AWS Dynamo Object storage model adheres to our normal path model (see README.md), but with the additional
conditions that:

a) There exists an AWS / Dynamo service account
b) There exists a table within that Dynamo account
c) There exist items within that table which are keyed to the Object paths

This enables us to maintain any number of collections of Objects within a given Dynamo account by
separating the collections into different tables.

Configuration:
	* awsregion	- AWS Region identifier e.g. "us-west-1"
	* tablename	- AWS Dynamo table to retrieve content from

*/

import (
	"errors"

//	"github.com/aws/aws-sdk-go/aws"
//	"github.com/aws/aws-sdk-go/service/dynamodb"

	lib "github.com/DigiStratum/GoLib"
	"github.com/DigiStratum/GoLib/Cloud"
)

type ObjectStoreDynamo struct {
	storeConfig	*lib.Config
	readCache	*MutableObjectStore
	awsHelper	*cloud.AWSHelper
}

// Make a new one of these!
func NewObjectStoreDynamo() *ObjectStoreDynamo {
	r := ObjectStoreDynamo{
		readCache: NewMutableObjectStore(),
	}
	return &r
}

// Satisfies RespositoryIfc
func (os *ObjectStoreDynamo) Configure(config *lib.Config) error {

	// Validate that the config has what we need for AWS Dynamo!
	requiredConfig := []string{ "awsregion", "tablename" }
	if ! (config.HasAll(&requiredConfig)) {
		return errors.New("Incomplete ObjectStoreDynamo configuration provided")
	}
	os.storeConfig = config

	// Light up our AWS Helper with the region from our configuration data
	os.awsHelper = cloud.NewAWSHelper(config.Get("awsregion"))
	return nil
}

// Satisfies ObjectStoreIfc
func (os *ObjectStoreDynamo) GetObject(path string) *Object {
	// If it's not yet in the cache
	if ! os.readCache.HasObject(path) {
		// TODO: Read the Object from our Dynamo Table into cache
		// Error = no Object!
		//if nil != err { return nil }
		//os.readCache.PutObject(path, NewObjectFromString(string(buff.Bytes())))
	}
	return os.readCache.GetObject(path)
}

// Satisfies ObjectStoreIfc
func (os *ObjectStoreDynamo) HasObject(path string) bool {
	// If it's already in the cache, then we know we have it!
	if os.readCache.HasObject(path) { return true }

	// TODO: Figure out if Dynamo has this object without necessarily retrieving it
	var err error
	return nil == err
}

// Satisfies MutableObjectStoreIfc
func (os *ObjectStoreDynamo) PutObject(path string, object *Object) error {
	// TODO: Actually implement WRITE operation to Dynamo here
	return errors.New("Not Yet Implemented!")
}

