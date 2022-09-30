package objectstoredynamo

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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	cfg "github.com/DigiStratum/GoLib/Config"
	obj "github.com/DigiStratum/GoLib/Object"
	objs "github.com/DigiStratum/GoLib/Object/store"
	cloud "github.com/DigiStratum/GoLib/Cloud/aws"
)

type ObjectStoreDynamo struct {
	storeConfig	cfg.ConfigIfc
	readCache	objs.MutableObjectStoreIfc
	awsHelper	*cloud.AWSHelper		// TODO: change to IFC
	awsDynamoDB	*dynamodb.DynamoDB		// TODO: change to IFC
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObjectStoreDynamo() *ObjectStoreDynamo {
	r := ObjectStoreDynamo{
		readCache: objs.NewMutableObjectStore(),
	}
	return &r
}


// -------------------------------------------------------------------------------------------------
// Satisfies ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectStoreDynamo) Configure(config cfg.ConfigIfc) error {

	// Validate that the config has what we need for AWS Dynamo!
	requiredConfig := []string{ "awsregion", "tablename", "primarykey" }
	if ! (config.HasAll(&requiredConfig)) {
		return fmt.Errorf("Incomplete ObjectStoreDynamo configuration provided")
	}
	r.storeConfig = config

	// Light up our AWS Helper with the region from our configuration data
	r.awsHelper = cloud.NewAWSHelper()
	awsHelperConfigKeys := []string{ "awsregion"}
	r.awsHelper.Configure(config.GetSubsetKeys(&awsHelperConfigKeys))
	return nil
}

// -------------------------------------------------------------------------------------------------
// Satisfies ObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Note that this precludes usage of Dynamo's support for "sort keys"; who
// would have thought that two keys would be required for an object store?
// Use read-through cache which requires us to mutate state
func (r *ObjectStoreDynamo) GetObject(path string) (*obj.Object, error) {
	// Require configuration
	if nil == r.storeConfig { return nil, fmt.Errorf("Not Configured!") }

	// If it's not yet in the cache
	if ! r.readCache.HasObject(path) {
		// TODO: Read the Object from our Dynamo Table into cache
		primaryKey := r.storeConfig.Get("primarykey")
		key := map[string]*dynamodb.AttributeValue{
			*primaryKey: {
				S: aws.String(path),
			},
		}
		tableName := r.storeConfig.Get("tablename")
		input := &dynamodb.GetItemInput{
			Key: key,
			TableName: aws.String(*tableName),
		}
		result, err := r.awsDynamoDB.GetItem(input)
		if nil != err {	return nil, fmt.Errorf(
			"ObjectStoreDynamo.GetObject() : DynamoDB.GetItem() : Error: '%s'",
			err.Error(),
		)}

		// Unmarshall the Dynamo result into a basic map of key=value strings
		// ref: https://stackoverflow.com/questions/11066946/partly-json-unmarshal-into-a-map-in-go
		type tobj struct {
			Key	string
			Content	string
		}
		//item := tobj{}
		item := struct{ Key, Content string }{}
		err = dynamodbattribute.UnmarshalMap(result.Item, &item)
		if err != nil { return nil, fmt.Errorf(
			"ObjectStoreDynamo.GetObject() JSON UnmarshallMap() : Error: '%s'",
			err.Error(),
		)}

		rcObj := obj.NewObject()
		rcObj.Deserialize(&item.Content)
		r.readCache.PutObject(path, rcObj)
	}
	return r.readCache.GetObject(path)
}

func (r ObjectStoreDynamo) HasObject(path string) (bool, error) {
	// Require configuration
	if nil == r.storeConfig { return false, fmt.Errorf("Not Configured!") }

	// If it's already in the cache, then we know we have it!
	if r.readCache.HasObject(path) { return true, nil }

	// TODO: Figure out if Dynamo has this object without necessarily retrieving it
	var err error
	return nil == err, err
}

// -------------------------------------------------------------------------------------------------
// Satisfies MutableObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectStoreDynamo) PutObject(path string, object *obj.Object) error {
	// Require configuration
	if nil == r.storeConfig { return fmt.Errorf("Not Configured!") }

	// TODO: Actually implement WRITE operation to Dynamo here
	return fmt.Errorf("Not Yet Implemented!")
}

// -------------------------------------------------------------------------------------------------
// ObjectStoreDynamo Private Interface
// -------------------------------------------------------------------------------------------------

// Get the DynamoDB service session
func (r *ObjectStoreDynamo) getDynamoService() *dynamodb.DynamoDB {
	if nil == r.awsDynamoDB {
		awsSession, err := r.awsHelper.GetSession()
		if nil == err {
			r.awsDynamoDB = dynamodb.New(awsSession)
		}
	}
	return r.awsDynamoDB
}
