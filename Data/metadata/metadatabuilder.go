package metadata

/*

TODO:
  * Add support for JSON un|marshal to deserialze
  * Add Factory Function to derive from existing *metadata
  * Add Factory Function to derive from existing HashMapIfc

*/

import (
	"github.com/DigiStratum/GoLib/Data/hashmap"
)

type MetadataBuilderIfc interface {
	// Our own interface
	Set(name, value string) *metadataBuilder
	GetMetadata() *metadata
}

type metadataBuilder struct {
	metadata *metadata
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewMetadataBuilder() *metadataBuilder {
	return &metadataBuilder{
		metadata: &metadata{
			data: hashmap.NewHashMap(),
		},
	}
}

// -------------------------------------------------------------------------------------------------
// MetadataBuilderIfc
// -------------------------------------------------------------------------------------------------

func (r *metadataBuilder) Set(name, value string) *metadataBuilder {
	r.metadata.data.Set(name, value)
	return r
}

func (r *metadataBuilder) GetMetadata() *metadata {
	return r.metadata
}
