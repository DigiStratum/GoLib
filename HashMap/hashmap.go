// DigiStratum GoLib - HashMap
package hashmap

/*

This HashMap class wraps a basic Go map with essential helper functions to make life easier for
dealing with simple key/value pair data. *Should* be thread-safe.

*/

import (
	"sync"
	"strconv"
	"encoding/json"
)

type keyValuePair struct {
	Key	string
	Value	string
}

// HashMap public interface
type HashMapIfc interface {
	LoadFromJsonString(json *string) error
	LoadFromJsonFile(jsonFile string) error
	IsEmpty() bool
	Size() int
	Merge(mergeHash HashMapIfc)
	Set(key string, value string)
	Get(key string) *string
	GetInt64(key string) *int64
	GetKeys() []string
	Has(key string) bool
	HasAll(keys *[]string) bool
	IterateCallback(callback func(kvp keyValuePair))
	IterateChannel() <-chan keyValuePair
	ToJson() (*string, error)
}

type HashMap struct {
	hash		map[string]string
	mutex		sync.Mutex
}

// Factory Functions
func NewHashMap() HashMap {
	return HashMap{
		hash:	make(map[string]string),
	}
}

// Get a full (deep) copy of this HashMap
// This is so that we can give away a copy to someone else without allowing them to tamper with us
// ref: https://developer20.com/be-aware-of-coping-in-go/
func CopyHashMap(source *HashMap) *HashMap {
	r := NewHashMap()
	for k, v := range (*source).hash { r.hash[k] = v }
	return &r
}

func NewHashMapFromJsonString(json *string) (*HashMap, error) {
	r := NewHashMap()
	if err := r.LoadFromJsonString(json); nil != err { return nil, err }
	return &r, nil
}

func NewHashMapFromJsonFile(jsonFile string) (*HashMap, error) {
	r := NewHashMap()
	if err := r.LoadFromJsonFile(jsonFile); nil != err { return nil, err }
	return &r, nil
}

// -------------------------------------------------------------------------------------------------
// HashMapIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Load our hash map with JSON data from a string (or return an error)
func (r *HashMap) LoadFromJsonString(json *string) error {
	r.mutex.Lock(); defer r.mutex.Unlock()
	return NewJson(json).Load(&r.hash)
}

// Load our hash map with JSON data from a file (or return an error)
func (r *HashMap) LoadFromJsonFile(jsonFile string) error {
	r.mutex.Lock(); defer r.mutex.Unlock()
	return NewJsonFromFile(jsonFile).Load(&r.hash)
}

// Check whether this HashMap is empty (has no properties)
func (r HashMap) IsEmpty() bool {
	return 0 == r.Size()
}

// Get the number of properties in this HashMap
func (r HashMap) Size() int {
	return len(r.hash)
}

// Merge some additional data on top of our own
func (r *HashMap) Merge(mergeHash HashMapIfc) {
	r.mutex.Lock(); defer r.mutex.Unlock()
	keys := mergeHash.GetKeys()
	for _, key := range keys {
		value := mergeHash.Get(key)
		if nil != value { r.Set(key, *value) }
	}
}

// Set a single data element key to the specified value
func (r *HashMap) Set(key string, value string) {
	r.mutex.Lock(); defer r.mutex.Unlock()
	r.hash[key] = value
}

// Get a single data element by key name
func (r HashMap) Get(key string) *string {
	if val, ok := r.hash[key]; ok { return &val }
	return nil
}

func (r HashMap) GetInt64(key string) *int64 {
	value := r.Get(key)
	if nil != value {
		if vc, err := strconv.ParseInt(*value, 0, 64); nil == err { return &vc }
	}
	return nil
}

// Check whether we have a data element by key name
func (r HashMap) Has(key string) bool {
	return r.Get(key) != nil
}

// Check whether we have configuration elements for all the key names
func (r HashMap) HasAll(keys *[]string) bool {
	for _, key := range *keys { if ! r.Has(key) { return false } }
	return true
}

func (r HashMap) GetKeys() []string {
	keys := make([]string, len(r.hash))
	i := 0
	for key, _ := range r.hash { keys[i] = key; i++ }
	return keys
}

// Iterate over the keys for this HashMap and call a callback for each
// ref: https://ewencp.org/blog/golang-iterators/index.html
func (r HashMap) IterateCallback(callback func(kvp keyValuePair)) {
	for k, v := range r.hash { callback(keyValuePair{ Key: k, Value: v}) }
}

// Iterate over the keys for this HashMap and send all the keyValuePairs to a channel
// ref: https://ewencp.org/blog/golang-iterators/index.html
// ref: https://blog.golang.org/pipelines
// ref: https://programming.guide/go/wait-for-goroutines-waitgroup.html
func (r HashMap) IterateChannel() <-chan keyValuePair {
	ch := make(chan keyValuePair, len(r.hash))
	defer close(ch)
	for k, v := range r.hash {
		ch <- keyValuePair{ Key: k, Value: v }
	}
	return ch
}

func (r HashMap) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(r.hash)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}