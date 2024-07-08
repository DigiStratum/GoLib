package metadata

/*

Metadata is a mini-hashmap with a builder to support passing of an immutable DTO of name-value pairs

TODO:
 * Add support for JSON un|marshal to de|serialze

*/

// Immutable Metadata DTO
type MetadataIfc interface {
	// Return true if the Metadata has ALL of the named values, else false
	Has(name ...string) bool
	Get(name string) string
	List() []string

	// Package private
	getMetadata() *metadata
}

type metadata struct {
	data	map[string]string
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewMetadata() *metadata {
	return &metadata{
		data:		make(map[string]string),
	}
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

func (r *metadata) List() []string {
	names := make([]string, 0)
	for name, _ := range r.data { names = append(names, name) }
	return names
}

func (r *metadata) getMetadata() *metadata {
	return r
}
