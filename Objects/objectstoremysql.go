package objects

/*

ObjectStore for MySQL Database

Ref: https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html

Configuration:
	* dsn	- Data Source Name (DSN) for MySQL database connection


*/

import (
	"errors"

	lib "github.com/DigiStratum/GoLib"
	cloud "github.com/DigiStratum/GoLib/Cloud"
)

type ObjectStoreMySQL struct {
	storeConfig		*lib.Config
	readCache		*MutableObjectStore
	awsHelper		*cloud.AWSHelper
}

// Make a new one of these!
func NewObjectStoreMySQL() *ObjectStoreMySQL {
	r := ObjectStoreMySQL{
		readCache: NewMutableObjectStore(),
	}
	return &r
}

// Satisfies RespositoryIfc
func (os *ObjectStoreMySQL) Configure(config *lib.Config) error {

	// Validate that the config has what we need for MySQL!
	requiredConfig := []string{ "dsn" }
	if ! (config.HasAll(&requiredConfig)) {
		return errors.New("Incomplete ObjectStoreMySQL configuration provided")
	}
	os.storeConfig = config

	// Light up our AWS Helper with the region from our configuration data
	os.awsHelper = cloud.NewAWSHelper(config.Get("awsregion"))
	return nil
}

// Satisfies ObjectStoreIfc
// Ref: https://stackoverflow.com/questions/41645377/golang-s3-download-to-buffer-using-s3manager-downloader
func (os *ObjectStoreMySQL) GetObject(path string) *Object {
	// If it's not yet in the cache
	if ! os.readCache.HasObject(path) {
		// TODO: Read the Object from MySQL bucket into cache
	}
	return os.readCache.GetObject(path)
}

// Satisfies ObjectStoreIfc
func (os *ObjectStoreMySQL) HasObject(path string) bool {
	// If it's already in the cache, then we know we have it!
	if os.readCache.HasObject(path) { return true }

	// If MySQL has a non-zero count of this record, then there's an Object!
	// TODO: look it up in the DB since it's not in the cache
	return false
}

// Satisfies MutableObjectStoreIfc
func (os *ObjectStoreMySQL) PutObject(path string, object *Object) error {
	// TODO: Actually implement WRITE operation to MySQL here
	return errors.New("Not Yet Implemented!")
}

