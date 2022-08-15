package objectstoremysql

/*

ObjectStore for MySQL Database

Ref: https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html

Configuration:
	* dsn	- Data Source Name (DSN) for MySQL database connection

FIXME:
 * This is not yet really implemented - we need to fill it in!

*/

import (
	"fmt"
	"errors"
	"net/url"

	cfg "github.com/DigiStratum/GoLib/Config"
	obj "github.com/DigiStratum/GoLib/Object"
	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
	cloud "github.com/DigiStratum/GoLib/Cloud"
)

// A given path specifier for this type of object store can be parsed into this logical structure
// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
type pathSpec struct {
	ObjectSpecName		string
	Keys			map[string]string
}

// The spec for a prepared statement query. Single '?' substitution is handled by db.Query()
// automatically. '???' expands to include enough placeholders (as with an IN () list for any count
// of keys > min. max must be >= min unless max == 0.
type querySpec struct {
	Query			string	// The query to execute as prepared statement
	MinKeys			int	// minimum num keys required to populate query; 0 = no min
	MaxKeys			int	// maximum num keys required to populate query; 0 = no max
}

// A given database object spec couples access queries with matching field definitions
type objectSpec struct {
	template		ObjectTemplate
	queries			map[string]mysql.QueryIfc
}

type ObjectStoreMySQL struct {
	storeConfig		cfg.ConfigIfc
	readCache		*MutableObjectStore
	awsHelper		*cloud.AWSHelper
	objectSpecs		map[string]objectSpec	// Object spec names must be part of object "path"
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObjectStoreMySQL() *ObjectStoreMySQL {
	r := ObjectStoreMySQL{
		readCache:	NewMutableObjectStore(),
		objectSpecs:	make(map[string]objectSpec),
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// Satisfies ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectStoreMySQL) Configure(config cfg.ConfigIfc) error {

	// Validate that the config has what we need for MySQL!
	requiredConfig := []string{ "dsn" }
	if ! (config.HasAll(&requiredConfig)) {
		return fmt.Errorf("Incomplete ObjectStoreMySQL configuration provided")
	}
	r.storeConfig = config

	// Light up our AWS Helper with the region from our configuration data
	//os.awsHelper = cloud.NewAWSHelper(config.Get("awsregion"))
	return nil
}

// -------------------------------------------------------------------------------------------------
// Satisfies ObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
func (r *ObjectStoreMySQL) GetObject(path string) (*obj.Object, error) {
	// Require configuration
	if nil == r.storeConfig { return nil, fmt.Errorf("Not Configured!") }

	// If it's not yet in the cache
	if ! r.readCache.HasObject(path) {
		// TODO: Read the Object from MySQL into cache
	}
	return r.readCache.GetObject(path), nil
}

// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
func (r *ObjectStoreMySQL) HasObject(path string) (bool, error) {
	// Require configuration
	if nil == r.storeConfig { return false, fmt.Errorf("Not Configured!") }

/*
	// TODO @HERE reactivate this; disabled for build
	// If it's already in the cache, then we know we have it!
	if r.readCache.HasObject(path) { return true }

	// If MySQL has a non-zero count of this record, then there's an Object!
	ps, err := os.parsePath(path)
	if nil != err {
		log.GetLogger().Warn(fmt.Sprintf(
			"Failed to retrieve requested path '%s': %s",
			path,
			err.Error(),
		))
		return false
	}

	// Get our Object Spec for this path spec...
	objectSpec, ok := r.objectSpecs[ps.ObjectSpecName]
	if ! ok {
		log.GetLogger().Warn("Failed to map requested Object Spec path '%s' (undefined!)")
		return false
	}
*/
	// TODO: look it up in the DB since it's not in the cache
	// TODO: use the objectSpec.queries["has"], prepared statement, (feed args into Query method if
	// possible? - this would prevent us from using arbitrary field ordering/spec in the path
	// query string... alpha-sort the keys and reobjectsquire same sorting in prepared query?)

	// ref: http://go-database-sql.org/prepared.html

	return false
}

// -------------------------------------------------------------------------------------------------
// Satisfies MutableObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
// blank keys to create with autoincrement; INSERT ... ON DUPLICATE KEY UPDATE syntax for create/update
func (r *ObjectStoreMySQL) PutObject(path string, object *obj.Object) error {
	// Require configuration
	if nil == r.storeConfig { return fmt.Errorf("Not Configured!") }

	// TODO: Actually implement WRITE operation to MySQL here
	return fmt.Errorf("Not Yet Implemented!")
}

// -------------------------------------------------------------------------------------------------
// ObjectStoreMySQL Private Interface
// -------------------------------------------------------------------------------------------------

// TODO: Move these supporting functions to a more generalized location such as ObjectStore
// if they don't end up with any implementation that is contextualized by ObjectStoreMySQL
// (or are otherwise generally useful)

// parse URL-Encoded path string into logical structure
// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
func (r ObjectStoreMySQL) parsePath(path string) (*pathSpec, error) {

	// Separate object spec name from keys
	pathParts := strings.Split(path, "?")
	// There should only be one '?' in the path...
	if len(pathParts) > 2 { return nil, fmt.Errorf(
		"ObjectStoreMySQL: Bad path spec (more than one '?'): '%s'",
		path,
	)}

	// Decode Object Spec Name, just in case it's been URL-encoded for some specialness
	objectSpecName, err := url.QueryUnescape(pathParts[0])
	if nil != err { return nil, fmt.Errorf(
		"ObjectStoreMySQL: Failed to unescape the Object Spec Name from path: '%s'",
		pathParts[0],
	)}

	// Make a pathspec so that we have a home for everything we find next
	ps := pathSpec{
		ObjectSpecName: objectSpecName,
		Keys:		make(map[string]string),
	}

	// More than one part means there might be keys to parse!
	if len(pathParts) > 1 {
		keyList := strings.Split(pathParts[1], "&")
		if len(keyList) > 0 {
			// Each key is a name=value - split them apart!
			for i := range keyList {
				name, value, err := os.parseURLKeyValuePair(keyList[i])
				if nil != err {
					return nil, err
				}
				ps.Keys[name] = value
			}
		}
	}
	return &ps, nil
}

// Parse a given "name=value" string such that name and.r value may be URL-encoded
// Return the name and value decoded strings or an error if there was a problem
func (r *ObjectStoreMySQL) parseURLKeyValuePair(keyValuePair string) (string, string, error) {
	keyParts := strings.Split(keyValuePair, "=")
	// There should only be one '=' in it...
	if len(keyParts) > 2 { return "", "", fmt.Errorf(
		"ObjectStoreMySQL: Bad key-value pair in path (more than one '='): '%s'",
		keyValuePair,
	)}

	name, err := url.QueryUnescape(keyParts[0])
	if nil != err { return "", "", fmt.Errorf(
		"ObjectStoreMySQL: Failed to unescape the Key Name from path: '%s'",
		keyParts[0],
	)}

	value, err := url.QueryUnescape(keyParts[1])
	if nil != err { return "", "", fmt.Errorf(
		"ObjectStoreMySQL: Failed to unescape the Key Value from path: '%s'",
		keyParts[1],
	)}
	return name, value, nil
}
