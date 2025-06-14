package metadata

/*

Metadata is a mini-hashmap with a builder to support passing of an immutable DTO of name-value pairs

TODO:
  * Add support for JSON un|marshal to serialze
  * Add support for Iterable interface

*/

import (
	"github.com/DigiStratum/GoLib/Data/hashmap"
)

// Immutable Metadata DTO
type MetadataIfc interface {
	// Return true if the Metadata has ALL of the named values, else false
	Has(names ...string) bool
	Get(name string) *string

	// Return a list of all names in the Metadata
	GetNames() []string
}

type metadata struct {
	data hashmap.HashMapIfc
}

// -------------------------------------------------------------------------------------------------
// MetadataIfc
// -------------------------------------------------------------------------------------------------

func (r *metadata) Has(names ...string) bool {
	return r.data.HasAll(names...)
}

func (r *metadata) Get(name string) *string {
	return r.data.Get(name)
}

func (r *metadata) GetNames() []string {
	return r.data.GetKeys()
}
