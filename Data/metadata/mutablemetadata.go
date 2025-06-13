package metadata

/*

TODO:
  * Add support for JSON un|marshal to deserialze
  * Convert this into a builder, not a mutable object
*/

type MetadataBuilderIfc interface {
	// Embedded interface(s)
	MetadataIfc

	// Our own interface
	Set(name, value string) *metadataBuilder
	GetMetadata() *metadata
}

type metadataBuilder struct {
	// Embedded struct+interface(s)
	metadata *metadata
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewMetadataBuilder() *metadataBuilder {
	return &metadataBuilder{
		metadata: NewMetadata(),
	}
}

// -------------------------------------------------------------------------------------------------
// MetadataBuilderIfc
// -------------------------------------------------------------------------------------------------

func (r *metadataBuilder) Set(name, value string) *metadataBuilder {
	r.metadata.data[name] = value
	return r
}

func (r *metadataBuilder) GetMetadata() *metadata {
	return r.metadata
}
