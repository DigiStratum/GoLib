// DigiStratum GoLib - ByteMap
package bytemap

/*

This ByteMap class wraps a basic Go map with essential helper functions to make life easier for
dealing with simple key/value pair data such that values are []byte; *Should* be thread-safe.

FIXME:
 * Replace the mutex with go-routine+channel for concurrency orchestration

*/

import (
	"sync"
	"strconv"
)

type KeyValuePair struct {
	Key	string
	Value	[]byte
}

// ByteMap public interface
type ByteMapIfc interface {
	Copy() *ByteMap
	Merge(mergeHash ByteMapIfc)
	IsEmpty() bool
	Size() int
	Set(key string, value []byte)
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
	data		map[string][]byte
	mutex		sync.Mutex
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewByteMap() *ByteMap {
	return &ByteMap{
		data:	make(map[string][]byte),
	}
}

// -------------------------------------------------------------------------------------------------
// ByteMapIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get a full (deep) copy of this ByteMap
// This is so that we can give away a copy to someone else without allowing them to tamper with us
// ref: https://developer20.com/be-aware-of-coping-in-go/
func (r *ByteMap) Copy() *ByteMap {
	n := NewByteMap()
	for k, v := range (*r).data { n.Set(k, v) }
	return n
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

// Check whether this ByteMap is empty (has no properties)
func (r *ByteMap) IsEmpty() bool {
	return 0 == r.Size()
}

// Get the number of properties in this ByteMap
func (r *ByteMap) Size() int {
	if nil == r { return 0 }
	return len(r.data)
}

// Set a single data element key to the specified value
func (r *ByteMap) Set(key string, value []byte) {
	if nil == r { return }
	r.set(key, value)
}

// Get a single data element by key name
func (r *ByteMap) Get(key string) *[]byte {
	if nil == r { return nil }
	if val, ok := r.data[key]; ok { return &val }
	return nil
}

// Get a single data element by key name as an int64
func (r *ByteMap) GetInt64(key string) *int64 {
	if nil == r { return nil }
	value := r.Get(key)
	if nil != value {
		if vc, err := strconv.ParseInt(string(*value), 0, 64); nil == err { return &vc }
	}
	return nil
}

// Get a single data element by key name as a boolean
func (r *ByteMap) GetBool(key string) bool {
	if nil == r { return false }
	s := r.Get(key)
	if nil == s { return false }
	str := string(*s)
	if ("true" == str) || ("TRUE" == str) || ("t" == str) || ("T" == str) || ("1" == str) { return true }
	if ("on" == str) || ("ON" == str) || ("yes" == str) || ("YES" == str) { return true }
	n := r.GetInt64(key)
	if nil == n { return false }
	return *n != 0
}

// Check whether we have a data element by key name
func (r *ByteMap) Has(key string) bool {
	if nil == r { return false }
	return r.Get(key) != nil
}

// Check whether we have configuration elements for all the key names
func (r *ByteMap) HasAll(keys *[]string) bool {
	if nil == r { return false }
	for _, key := range *keys { if ! r.Has(key) { return false } }
	return true
}

// Make a new bytemap from a subset of the key-values from this one
func (r *ByteMap) GetSubset(keys *[]string) *ByteMap {
	if nil == r { return nil }
	n := NewByteMap()
	if nil != keys {
		for _, k := range *keys {
			if v, ok := r.data[k]; ok { n.data[k] = v }
		}
	}
	return n
}

// Get the set of keys currently loaded into this bytemap
func (r *ByteMap) GetKeys() []string {
	if nil == r { return make([]string, 0) }
	keys := make([]string, len(r.data))
	i := 0
	for key, _ := range r.data { keys[i] = key; i++ }
	return keys
}

// Drop a single key from the bytemap, if it's set
func (r *ByteMap) Drop(key string) *ByteMap {
	if nil == r { return nil }
	if ! r.Has(key) { return r }
	delete(r.data, key)
	return r
}

// Drop an set of keys from the bytemap, if any are set
func (r *ByteMap) DropSet(keys *[]string) *ByteMap {
	if nil == r { return nil }
	if (nil == keys) || (len(*keys) == 0) { return r }
	for _, k := range *keys { r.Drop(k) }
	return r
}

// Drop all keys/values (reset to empty state)
func (r *ByteMap) DropAll() {
	if nil == r { return }
	r.data = make(map[string][]byte)
}

// -------------------------------------------------------------------------------------------------
// IterableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Iterate over all of our items, returning each as a *KeyValuePair in the form of an interface{}
func (r *ByteMap) GetIterator() func () interface{} {
	if nil == r { return nil }
	kvps := make([]KeyValuePair, r.Size())
	var idx int = 0
	for k, v := range r.data {
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
// ByteMap Private Implementation
// -------------------------------------------------------------------------------------------------

// Set a single data element key to the specified value
func (r *ByteMap) set(key string, value []byte) {
	if nil == r { return }
	r.data[key] = value
}
