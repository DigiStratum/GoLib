// DigiStratum GoLib - HashMap
package hashmap

/*

This HashMap class wraps a basic Go map with essential helper functions to make life easier for
dealing with simple key/value pair data. *Should* be thread-safe.

*/

import (
	"fmt"
	"sync"
	"strconv"
	gojson "encoding/json"

	"github.com/DigiStratum/GoLib/Data/json"
)

type KeyValuePair struct {
	Key	string
	Value	string
}

// HashMap public interface
type HashMapIfc interface {
	Copy() *HashMap
	LoadFromJsonString(jsonStr *string) error
	LoadFromJsonFile(jsonFile string) error
	IsEmpty() bool
	Size() int
	Merge(mergeHash HashMapIfc)
	Set(key, value string)
	Get(key string) *string
	GetInt64(key string) *int64
	GetBool(key string) bool
	GetKeys() []string
	Has(key string) bool
	HasAll(keys *[]string) bool
	GetSubset(keys *[]string) *HashMap
	Drop(key string) *HashMap
	DropSet(keys *[]string) *HashMap
	GetIterator() func () interface{}
	ToJson() (*string, error)
}

type HashMap struct {
	hash		map[string]string
	mutex		sync.Mutex
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewHashMap() *HashMap {
	return &HashMap{
		hash:	make(map[string]string),
	}
}

func NewHashMapFromJsonString(json *string) (*HashMap, error) {
	r := NewHashMap()
	if err := r.LoadFromJsonString(json); nil != err { return nil, err }
	return r, nil
}

func NewHashMapFromJsonFile(jsonFile string) (*HashMap, error) {
	r := NewHashMap()
	if err := r.LoadFromJsonFile(jsonFile); nil != err { return nil, err }
	return r, nil
}

// -------------------------------------------------------------------------------------------------
// HashMapIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get a full (deep) copy of this HashMap
// This is so that we can give away a copy to someone else without allowing them to tamper with us
// ref: https://developer20.com/be-aware-of-coping-in-go/
func (r *HashMap) Copy() *HashMap {
	n := NewHashMap()
	for k, v := range r.hash { n.hash[k] = v }
	return n
}

// Load our hash map with JSON data from a string (or return an error)
func (r *HashMap) LoadFromJsonString(jsonStr *string) error {
	if nil == r { return fmt.Errorf("This receiver is nil, nothing to do!") }
	r.mutex.Lock(); defer r.mutex.Unlock()
	return json.NewJson(jsonStr).Load(&r.hash)
}

// Load our hash map with JSON data from a file (or return an error)
func (r *HashMap) LoadFromJsonFile(jsonFile string) error {
	if nil == r { return fmt.Errorf("This receiver is nil, nothing to do!") }
	r.mutex.Lock(); defer r.mutex.Unlock()
	return json.NewJsonFromFile(jsonFile).Load(&r.hash)
}

// Check whether this HashMap is empty (has no properties)
func (r *HashMap) IsEmpty() bool {
	return 0 == r.Size()
}

// Get the number of properties in this HashMap
func (r *HashMap) Size() int {
	return len(r.hash)
}

// Merge some additional data on top of our own
func (r *HashMap) Merge(mergeHash HashMapIfc) {
	if nil == r { return }
	r.mutex.Lock(); defer r.mutex.Unlock()
	keys := mergeHash.GetKeys()
	for _, key := range keys {
		value := mergeHash.Get(key)
		if nil != value { r.set(key, *value) }
	}
}

// Set a single data element key to the specified value
func (r *HashMap) Set(key, value string) {
	if nil == r { return }
	r.set(key, value)
}

// Get a single data element by key name
func (r *HashMap) Get(key string) *string {
	if val, ok := r.hash[key]; ok { return &val }
	return nil
}

// Get a single data element by key name as an int64
func (r *HashMap) GetInt64(key string) *int64 {
	value := r.Get(key)
	if nil != value {
		if vc, err := strconv.ParseInt(*value, 0, 64); nil == err { return &vc }
	}
	return nil
}

// Get a single data element by key name as a boolean
func (r *HashMap) GetBool(key string) bool {
	s := r.Get(key)
	if nil == s { return false }
	if ("true" == *s) || ("TRUE" == *s) || ("t" == *s) || ("T" == *s) || ("1" == *s) { return true }
	if ("on" == *s) || ("ON" == *s) || ("yes" == *s) || ("YES" == *s) { return true }
	n := r.GetInt64(key)
	if nil == n { return false }
	return *n != 0
}

// Check whether we have a data element by key name
func (r *HashMap) Has(key string) bool {
	return r.Get(key) != nil
}

// Check whether we have configuration elements for all the key names
func (r *HashMap) HasAll(keys *[]string) bool {
	for _, key := range *keys { if ! r.Has(key) { return false } }
	return true
}

// Make a new hashmap from a subset of the key-values from this one
func (r *HashMap) GetSubset(keys *[]string) *HashMap {
	n := NewHashMap()
	if nil != keys {
		for _, k := range *keys {
			if v, ok := r.hash[k]; ok { n.hash[k] = v }
		}
	}
	return n
}

// Get the set of keys currently loaded into this hashmap
func (r *HashMap) GetKeys() []string {
	keys := make([]string, len(r.hash))
	i := 0
	for key, _ := range r.hash { keys[i] = key; i++ }
	return keys
}

// Drop a single key from the hashmap, if it's set
func (r *HashMap) Drop(key string) *HashMap {
	if ! r.Has(key) { return r }
	delete(r.hash, key)
	return r
}

// Drop an set of keys from the hashmap, if any are set
func (r *HashMap) DropSet(keys *[]string) *HashMap {
	if nil == keys { return r }
	for _, k := range *keys { r.Drop(k) }
	return r
}

// -------------------------------------------------------------------------------------------------
// IterableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Iterate over all of our items, returning each as a *KeyValuePair in the form of an interface{}
func (r *HashMap) GetIterator() func () interface{} {
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

func (r *HashMap) ToJson() (*string, error) {
	jsonBytes, err := gojson.Marshal(r.hash)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// HashMapIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Set a single data element key to the specified value
func (r *HashMap) set(key, value string) {
	if nil == r { return }
	r.hash[key] = value
}
