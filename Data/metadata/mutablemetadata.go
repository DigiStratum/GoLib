package metadata

/*

TODO:
 * Add support for JSON un|marshal to deserialze

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

func NewMutableMetadata() *metadataBuilder {
	return &metadataBuilder{
		metadata: NewMetadata(),
	}
}

// -------------------------------------------------------------------------------------------------
// MutableMetadataIfc
// -------------------------------------------------------------------------------------------------

func (r *metadataBuilder) Set(name, value string) *metadataBuilder {
	//r.MetadataIfc.data[name] = value
	r.metadata.data[name] = value
	return r
}

func (r *metadataBuilder) GetMetadata() *metadata {
	return r.metadata
}
