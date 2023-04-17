package metadata

/*

Mutable extension of otherwise immutable Metadata

*/

type MutableMetadataIfc interface {
	// Embedded interface(s)
	MetadataIfc

	// Our own interface
	Set(name, value string) *mutableMetadata
}

type mutableMetadata struct {
	// Embedded struct+interface(s)
	MetadataIfc
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewMutableMetadata() *mutableMetadata {
	return &mutableMetadata{
		MetadataIfc:		NewMetadata(),
	}
}

// -------------------------------------------------------------------------------------------------
// MutableMetadataIfc
// -------------------------------------------------------------------------------------------------

func (r *mutableMetadata) Set(name, value string) *mutableMetadata {
	//r.MetadataIfc.data[name] = value
	metadata := r.MetadataIfc.getMetadata()
	metadata.data[name] = value
	return r
}

