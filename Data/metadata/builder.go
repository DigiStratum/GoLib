package metadata

/*

Metadata is a mini-hashmap with a builder to support passing of an immutable DTO of name-value pairs

*/

type BuilderIfc interface {
	Set(name, value string) *builder
	Build() *metadata
}

// Metadata builder
type builder struct {
	built		bool
	md		*metadata
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewBuilder() *builder {
	return &builder{
		md:	NewMetadata(),
	}
}

// -------------------------------------------------------------------------------------------------
// MetadataBuilderIfc
// -------------------------------------------------------------------------------------------------

func (r *builder) Set(name, value string) *builder {
	if r.built { return nil }
	r.md.data[name] = value
	return r
}

func (r *builder) Build() *metadata {
	r.built = true
	return r.md
}

