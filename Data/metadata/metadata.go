package metadata

/*

Metadata is a mini-hashmap with a builder to support passing of an immutable DTO of name-value pairs

*/

// Immutable Metadata DTO
type MetadataIfc interface {
	Get(name string) string
}

type metadata struct {
	data	map[string]string
}

// Metadata builder
type MetadataBuilder struct {
	built		bool
	md		*metadata
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func newMetadata() *metadata {
	return &metadata{
		data:		make(map[string]string),
	}
}

func NewMetadataBuilder() *MetadataBuilder {
	return &MetadataBuilder{
		md:	newMetadata(),
	}
}

// -------------------------------------------------------------------------------------------------
// MetadataBuilder
// -------------------------------------------------------------------------------------------------

func (r *MetadataBuilder) Set(name, value string) *MetadataBuilder {
	if r.built { return nil }
	r.md.data[name] = value
	return r
}

func (r *MetadataBuilder) Build() *metadata {
	r.built = true
	return r.md
}

// -------------------------------------------------------------------------------------------------
// MetadataIfc
// -------------------------------------------------------------------------------------------------

func (r *metadata) Get(name string) string {
	if value, ok := r.data[name]; ok { return value }
	return ""
}

