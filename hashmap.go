// DigiStratum GoLib - HashMap
package golib

/*

This HashMap class wraps a basic Go map with essential helper functions to make life easier for
dealing with simple key/value pair data.

We're going to let this thing panic if we are not initialized (nil). That's the caller's bad.

TODO: Put some multi-threaded protections around the accessors here

*/

import (
	"sync"
	"fmt"
	"strings"
)

type KeyValuePair struct {
	Key	string
	Value	string
}

type HashMap	map[string]string

// Make a new one of these!
func NewHashMap() *HashMap {
	return &HashMap{}
}

// Check whether this HashMap is empty (has no properties)
func (hash *HashMap) IsEmpty() bool {
	return 0 == hash.Size()
}

// Get the number of properties in this HashMap
func (hash *HashMap) Size() int {
	return len(*hash)
}

// Merge some additional configuration data on top of our own
func (hash *HashMap) Merge(inbound *HashMap) {
	for k, v := range *inbound { (*hash)[k] = v }
}

// Set a single configuration element key to the specified value
func (hash *HashMap) Set(key string, value string) {
	(*hash)[key] = value
}

// Get a single configuration element by key name
func (hash *HashMap) Get(key string) string {
	str := ""
	val, ok := (*hash)[key]
	if ok { str = val }
	return str
}

// Check whether we have a configuration element by key name
func (hash *HashMap) Has(key string) bool {
	_, ok := (*hash)[key];
	return ok
}

// Check whether we have configuration elements for all the key names
func (hash *HashMap) HasAll(keys *[]string) bool {
	for _, key := range *keys {
		_, ok := (*hash)[key];
		if ! ok { return false }
	}
	return true
}

// Get a full copy of this HashMap
// This is so that we can give away a copy to someone else without allowing them to tamper with us
func (hash *HashMap) GetCopy() *HashMap {
	if nil == hash { return nil }
	res := make(HashMap)
	for k, v := range *hash { res[k] = v }
	return &res
}

// Iterate over the keys for this HashMap and call a callback for each
// ref: https://ewencp.org/blog/golang-iterators/index.html
func (hash *HashMap) IterateCallback(callback func(kvp KeyValuePair)) {
	for k, v := range *hash { callback(KeyValuePair{ Key: k, Value: v}) }
}

// Iterate over the keys for this HashMap and send all the KeyValuePairs to a channel
// ref: https://ewencp.org/blog/golang-iterators/index.html
// ref: https://blog.golang.org/pipelines
// ref: https://programming.guide/go/wait-for-goroutines-waitgroup.html
func (hash *HashMap) IterateChannel() <-chan KeyValuePair {
	ch := make(chan KeyValuePair, len(*hash))
	defer close(ch)
	var wg sync.WaitGroup
	wg.Add(1)

	// Fire off a go routine to fill up the channel
	go func() {
		for k, v := range *hash {
			ch <- KeyValuePair{ Key: k, Value: v }
		}
		wg.Done()
	}()
	wg.Wait()
	return ch
}

// Dump the contents of this HashMap to stdout for debug purposes
func (hash *HashMap) Dump() {
	GetLogger().Info(hash.DumpString())
}

// Dump the contents of this HashMap to a string so that it can be captured/processed as needed by the caller
func (hash *HashMap) DumpString() string {
	var b strings.Builder
	for k, v := range *hash { fmt.Fprintf(&b, "\t'%s': '%s'\n", k, v) }
	return fmt.Sprintf("HashMap = {\n%s}\n", b.String())
}

