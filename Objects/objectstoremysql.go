package objects

/*

ObjectStore for MySQL Database

Ref: https://www.golangprograms.com/example-of-golang-crud-using-mysql-from-scratch.html

Configuration:
	* dsn	- Data Source Name (DSN) for MySQL database connection


*/

import (
	"fmt"
	"errors"
	"strings"
	"net/url"

	lib "github.com/DigiStratum/GoLib"
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
	getQuery		mysql.QuerySpec	// e.g. SELECT * FROM tablename WHERE ... LIMIT 1
	hasQuery		mysql.QuerySpec	// e.g. SELECT COUNT(*) FROM tablename WHERE ... LIMIT 1
	putQuery		mysql.QuerySpec	// e.g. INSERT INTO tablename SET name=value ... WHERE ... ON DUPLICATE KEY UPDATE tablename ...
}

type ObjectStoreMySQL struct {
	storeConfig		*lib.Config
	readCache		*MutableObjectStore
	awsHelper		*cloud.AWSHelper
	objectSpecs		map[string]objectSpec	// Object spec names must be part of object "path"
}

// Make a new one of these!
func NewObjectStoreMySQL() *ObjectStoreMySQL {
	r := ObjectStoreMySQL{
		readCache:	NewMutableObjectStore(),
		objectSpecs:	make(map[string]objectSpec),
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
	//os.awsHelper = cloud.NewAWSHelper(config.Get("awsregion"))
	return nil
}

// Satisfies ObjectStoreIfc
// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
func (os *ObjectStoreMySQL) GetObject(path string) *Object {
	// If it's not yet in the cache
	if ! os.readCache.HasObject(path) {
		// TODO: Read the Object from MySQL into cache
	}
	return os.readCache.GetObject(path)
}

// Satisfies ObjectStoreIfc
// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
func (os *ObjectStoreMySQL) HasObject(path string) bool {
	// If it's already in the cache, then we know we have it!
	if os.readCache.HasObject(path) { return true }

	// If MySQL has a non-zero count of this record, then there's an Object!
	ps, err := os.parsePath(path)
	if nil != err {
		lib.GetLogger().Warn(fmt.Sprintf(
			"Failed to retrieve requested path '%s': %s",
			path,
			err.Error(),
		))
		return false
	}

	// Get our Object Spec for this path spec...
	if objectSpec, ok := os.objectSpecs[ps.ObjectSpecName]; ok {
	} else {
		lib.GetLogger().Warn("Failed to map requested Object Spec path '%s' (undefined!)")
		return false
	}

	// TODO: look it up in the DB since it's not in the cache
	// TODO: use the objectSpec.hasQuery, prepared statement, (feed args into Query method if
	// possible? - this would prevent us from using arbitrary field ordering/spec in the path
	// query string... alpha-sort the keys and require same sorting in prepared query?)

	// ref: http://go-database-sql.org/prepared.html

	return false
}

// Satisfies MutableObjectStoreIfc
// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
// blank keys to create with autoincrement; INSERT ... ON DUPLICATE KEY UPDATE syntax for create/update
func (os *ObjectStoreMySQL) PutObject(path string, object *Object) error {
	// TODO: Actually implement WRITE operation to MySQL here
	return errors.New("Not Yet Implemented!")
}

// parse URL-Encoded path string into logical structure
// path format: "objectspecname?key1=value1&key2=value2&keyN=valueN
func (os *ObjectStoreMySQL) parsePath(path string) (*pathSpec, error) {

	// Separate object spec name from keys
	pathParts := strings.Split(path, "?")
	if len(pathParts) > 2 {
		// Hmm - there should only be one '?' in the path...
		return nil, errors.New(fmt.Sprintf(
			"ObjectStoreMySQL: Bad path spec (more than one '?'): '%s'",
			path,
		))
	}

	// Decode Object Spec Name, just in case it's been URL-encoded for some specialness
	objectSpecName, err := url.QueryUnescape(pathParts[0])
	if nil != err {
		return nil, errors.New(fmt.Sprintf(
			"ObjectStoreMySQL: Failed to unescape the Object Spec Name from path: '%s'",
			pathParts[0],
		))
	}

	// Okay, make a pathspec so that we have a home for everything we find next
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

func (os *ObjectStoreMySQL) parseURLKeyValuePair(keyValuePair string) (string, string, error) {
	keyParts := strings.Split(keyValuePair, "=")
	if len(keyParts) > 2 {
		// Hmm - there should only be one '=' in it...
		return "", "", errors.New(fmt.Sprintf(
			"ObjectStoreMySQL: Bad key-value pair in path (more than one '='): '%s'",
			keyValuePair,
		))
	}
	name, err := url.QueryUnescape(keyParts[0])
	if nil != err {
		return "", "", errors.New(fmt.Sprintf(
			"ObjectStoreMySQL: Failed to unescape the Key Name from path: '%s'",
			keyParts[0],
		))
	}
	value, err := url.QueryUnescape(keyParts[1])
	if nil != err {
		return "", "", errors.New(fmt.Sprintf(
			"ObjectStoreMySQL: Failed to unescape the Key Value from path: '%s'",
			keyParts[1],
		))
	}
	return name, value, nil
}

