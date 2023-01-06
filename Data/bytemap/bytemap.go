// DigiStratum GoLib - ByteMap
package bytemap

/*

This ByteMap class wraps a basic Go map with essential helper functions to make life easier for
dealing with simple key/value pair data such that values are []byte; *Should* be thread-safe.

FIXME:
 * Replace the mutex with go-routine+channel for concurrency orchestration

*/

import (
	"fmt"
	"sync"
	"strconv"
	gojson "encoding/json"

	"github.com/DigiStratum/GoLib/Data/json"
	log "github.com/DigiStratum/GoLib/Logger"
)

type KeyValuePair struct {
	Key	string
	Value	string
}

// ByteMap public interface
type ByteMapIfc interface {
	Copy() *ByteMap
	Merge(mergeHash ByteMapIfc)
	IsEmpty() bool
	Size() int
	Set(key, value []byte)
	Get(key string) *[]byte
	GetInt64(key string) *int64
	GetBool(key string) bool
	GetKeys() []string
	Has(key string) bool
	HasAll(keys *[]string) bool
	GetSubset(keys *[]string) *ByteMap
	Drop(key string) *ByteMap
	DropSet(keys *[]string) *ByteMap
	DropAll()
	GetIterator() func () interface{}
}

type ByteMap struct {
	hash		map[string]string
	mutex		sync.Mutex
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewByteMap() *ByteMap {
	return &ByteMap{
		hash:	make(map[string]string),
	}
}

func NewByteMapFromJsonString(json *string) (*ByteMap, error) {
	r := NewByteMap()
	if err := r.LoadFromJsonString(json); nil != err { return nil, err }
	return r, nil
}

func NewByteMapFromJsonFile(jsonFile string) (*ByteMap, error) {
	r := NewByteMap()
	if err := r.LoadFromJsonFile(jsonFile); nil != err { return nil, err }
	return r, nil
}

// -------------------------------------------------------------------------------------------------
// ByteMapIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get a full (deep) copy of this ByteMap
// This is so that we can give away a copy to someone else without allowing them to tamper with us
// ref: https://developer20.com/be-aware-of-coping-in-go/
func (r *ByteMap) Copy() *ByteMap {
	n := NewByteMap()
	for k, v := range (*r).hash { n.Set(k, v) }
	return n
}

// Load our hash map with JSON data from a string (or return an error)
func (r *ByteMap) LoadFromJsonString(jsonStr *string) error {
	if nil == r { return fmt.Errorf("This receiver is nil, nothing to do!") }
	r.mutex.Lock(); defer r.mutex.Unlock()
	return json.NewJson(jsonStr).Load(&r.hash)
}

// Load our hash map with JSON data from a file (or return an error)
func (r *ByteMap) LoadFromJsonFile(jsonFile string) error {
	if nil == r { return fmt.Errorf("This receiver is nil, nothing to do!") }
	r.mutex.Lock(); defer r.mutex.Unlock()
	return json.NewJsonFromFile(jsonFile).Load(&r.hash)
}

// Check whether this ByteMap is empty (has no properties)
func (r *ByteMap) IsEmpty() bool {
	return 0 == r.Size()
}

// Get the number of properties in this ByteMap
func (r *ByteMap) Size() int {
	return len(r.hash)
}

// Merge some additional data on top of our own
func (r *ByteMap) Merge(mergeHash ByteMapIfc) {
	if nil == r { return }
	r.mutex.Lock(); defer r.mutex.Unlock()
	keys := mergeHash.GetKeys()
	for _, key := range keys {
		value := mergeHash.Get(key)
		if nil != value { r.set(key, *value) }
	}
}

// Set a single data element key to the specified value
func (r *ByteMap) Set(key, value string) {
	if nil == r { return }
	r.set(key, value)
}

// Get a single data element by key name
func (r *ByteMap) Get(key string) *string {
	if val, ok := r.hash[key]; ok { return &val }
	return nil
}

// Get a single data element by key name as an int64
func (r *ByteMap) GetInt64(key string) *int64 {
	value := r.Get(key)
	if nil != value {
		if vc, err := strconv.ParseInt(*value, 0, 64); nil == err { return &vc }
	}
	return nil
}

// Get a single data element by key name as a boolean
func (r *ByteMap) GetBool(key string) bool {
	s := r.Get(key)
	if nil == s { return false }
	if ("true" == *s) || ("TRUE" == *s) || ("t" == *s) || ("T" == *s) || ("1" == *s) { return true }
	if ("on" == *s) || ("ON" == *s) || ("yes" == *s) || ("YES" == *s) { return true }
	n := r.GetInt64(key)
	if nil == n { return false }
	return *n != 0
}

// Check whether we have a data element by key name
func (r *ByteMap) Has(key string) bool {
	return r.Get(key) != nil
}

// Check whether we have configuration elements for all the key names
func (r *ByteMap) HasAll(keys *[]string) bool {
	for _, key := range *keys { if ! r.Has(key) { return false } }
	return true
}

// Make a new hashmap from a subset of the key-values from this one
func (r *ByteMap) GetSubset(keys *[]string) *ByteMap {
	n := NewByteMap()
	if nil != keys {
		for _, k := range *keys {
			if v, ok := r.hash[k]; ok { n.hash[k] = v }
		}
	}
	return n
}

// Get the set of keys currently loaded into this hashmap
func (r *ByteMap) GetKeys() []string {
	keys := make([]string, len(r.hash))
	i := 0
	for key, _ := range r.hash { keys[i] = key; i++ }
	return keys
}

// Drop a single key from the hashmap, if it's set
func (r *ByteMap) Drop(key string) *ByteMap {
	if ! r.Has(key) { return r }
	delete(r.hash, key)
	return r
}

// Drop an set of keys from the hashmap, if any are set
func (r *ByteMap) DropSet(keys *[]string) *ByteMap {
	if nil == keys { return r }
	for _, k := range *keys { r.Drop(k) }
	return r
}

// Drop all keys/values (reset to empty state)
func (r *ByteMap) DropAll() {
	r.hash = make(map[string]string)
}

// Dump JSON-like representation of our entries in readable form to supplied logger
func (r *ByteMap) ToLog(logger log.LoggerIfc, level log.LogLevel, label string) {
	if nil == logger { return }
	logger.Any(level, "\"%s\": {", label)
	for k, v := range r.hash { logger.Any(level, "\t\"%s\": \"%s\"", k, v) }
	logger.Any(level, "}")
}

// -------------------------------------------------------------------------------------------------
// IterableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Iterate over all of our items, returning each as a *KeyValuePair in the form of an interface{}
func (r *ByteMap) GetIterator() func () interface{} {
	kvps := make([]KeyValuePair, r.Size())
	var idx int = 0
	for k, v := range r.hash {
		kvps[idx] = KeyValuePair{ Key: k, Value: v }
		idx++
	}
	idx = 0
	var data_len = r.Size()
	return func () interface{} {
		// If we're done iterating, return do nothing
		if idx >= data_len { return nil }
		prev_idx := idx
		idx++
		return &kvps[prev_idx]
	}
}

// -------------------------------------------------------------------------------------------------
// JsonSerializable Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ByteMap) ToJson() (*string, error) {
	jsonBytes, err := gojson.Marshal(r.hash)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Interface Implementation
// -------------------------------------------------------------------------------------------------

func (r *ByteMap) MarshalJSON() ([]byte, error) {
	return gojson.Marshal(r.hash)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Interface Implementation
// -------------------------------------------------------------------------------------------------

func (r *ByteMap) UnmarshalJSON(value []byte) error {
	return gojson.Unmarshal(value, &(r.hash))
}

// -------------------------------------------------------------------------------------------------
// ByteMapIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Set a single data element key to the specified value
func (r *ByteMap) set(key, value string) {
	if nil == r { return }
	r.hash[key] = value
}
