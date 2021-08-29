package dynamo

/*

ObjectStore for AWS Dynamo NoSQL Database service

Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/
Ref: https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/#example_DynamoDB_GetItem_shared00
Ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/dynamo-example-read-table-item.html

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

TODO: Add a check for whether we have been Configure()'d before allowing usage

*/

import (
	"fmt"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	lib "github.com/DigiStratum/GoLib"
	"github.com/DigiStratum/GoLib/Cloud"
)

type ObjectStoreDynamoIfc interface {
}

type ObjectStoreDynamo struct {
	storeConfig	lib.ConfigIfc
	readCache	*MutableObjectStore
	awsHelper	*cloud.AWSHelper
	awsDynamoDB	*dynamodb.DynamoDB
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObjectStoreDynamo() *ObjectStoreDynamo {
	r := ObjectStoreDynamo{
		readCache: NewMutableObjectStore(),
	}
	return &r
}


// -------------------------------------------------------------------------------------------------
// ObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectStoreDynamo) Configure(config lib.ConfigIfc) error {

	// Validate that the config has what we need for AWS Dynamo!
	requiredConfig := []string{ "awsregion", "tablename", "primarykey" }
	if ! (config.HasAll(&requiredConfig)) {
		return errors.New("Incomplete ObjectStoreDynamo configuration provided")
	}
	r.storeConfig = config

	// Light up our AWS Helper with the region from our configuration data
	r.awsHelper = cloud.NewAWSHelper(config.Get("awsregion"))
	return nil
}

// Note that this precludes usage of Dynamo's support for "sort keys"; who
// would have thought that two keys would be required for an object store?
// Use read-through cache which requires us to mutate state
func (r *ObjectStoreDynamo) GetObject(path string) *Object {
	// If it's not yet in the cache
	if ! r.readCache.HasObject(path) {
		// TODO: Read the Object from our Dynamo Table into cache
		key := map[string]*dynamodb.AttributeValue{
			r.storeConfig.Get("primarykey"): {
				S: aws.String(path),
			},
		}
		input := &dynamodb.GetItemInput{
			Key: key,
			TableName: aws.String(r.storeConfig.Get("tablename")),
		}
		result, err := r.awsDynamoDB.GetItem(input)
		if nil != err {
			lib.GetLogger().Error(fmt.Sprintf(
				"ObjectStoreDynamo.GetObject() : DynamoDB.GetItem() : Error: '%s'",
				err.Error(),
			))
			return nil
		}

		// Unmarshall the Dynamo result into a basic map of key=value strings
		// ref: https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
		type obj struct {
			Key	string
			Content	string
		}
		item := obj{}
		err = dynamodbattribute.UnmarshalMap(result.Item, &item)
		if err != nil {
			lib.GetLogger().Error(fmt.Sprintf(
				"ObjectStoreDynamo.GetObject() JSON UnmarshallMap() : Error: '%s'",
				err.Error(),
			))
			return nil
		}

		r.readCache.PutObject(path, NewObjectFromString(item.Content))
	}
	return r.readCache.GetObject(path)
}

func (r ObjectStoreDynamo) HasObject(path string) bool {
	// If it's already in the cache, then we know we have it!
	if r.readCache.HasObject(path) { return true }

	// TODO: Figure out if Dynamo has this object without necessarily retrieving it
	var err error
	return nil == err
}

// -------------------------------------------------------------------------------------------------
// MutableObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectStoreDynamo) PutObject(path string, object *Object) error {
	// TODO: Actually implement WRITE operation to Dynamo here
	return errors.New("Not Yet Implemented!")
}

// -------------------------------------------------------------------------------------------------
// ObjectStoreDynamoIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Get the DynamoDB service session
func (r *ObjectStoreDynamo) getDynamoService() *dynamodb.DynamoDB {
	if nil == r.awsDynamoDB {
		r.awsDynamoDB = dynamodb.New(r.awsHelper.GetSession())
	}
	return r.awsDynamoDB
}
