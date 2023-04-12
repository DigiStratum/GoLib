package metadata

/*

Metadata is a mini-hashmap with a builder to support passing of an immutable DTO of name-value pairs

*/

// Immutable Metadata DTO
type MetadataIfc interface {
	// Return true if the Metadata has ALL of the named values, else false
	Has(name ...string) bool
	Get(name string) string
	GetNames() []string
}

type MetadataBuilderIfc interface {
	Set(name, value string) *MetadataBuilder
	Build() *metadata
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
// MetadataBuilderIfc
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

func (r *metadata) Has(name ...string) bool {
	for _, n := range name {
		_, ok := r.data[n]
		if ! ok { return false }
	}
	return true
}

func (r *metadata) Get(name string) string {
	if value, ok := r.data[name]; ok { return value }
	return ""
}

func (r *metadata) GetNames() []string {
	names := make([]string, 0)
	for name, _ := range r.data { names = append(names, name) }
	return names
}

