// DigiStratum GoLib - Dependency Injection
package dependencies

/*
Dependencies is a Dependency set; it is what we use to declare a set of Dependencies and represents
the complete collection of Dependencies as an expression of what a client needs/wants from the
provider.
*/

type DependenciesIfc interface {
	// Add one or more Dependencies to the set
	Add(deps ...DependencyIfc) *dependencies

	// Get a dependency by name/variant
	Get(name string) DependencyIfc
	GetVariant(name, variant string) DependencyIfc

	// Check whether a dependency is in the set by name/variant
	Has(name string) bool
	HasVariant(name, variant string) bool

	// Get the list of currently set dependencies
	GetVariants(name string) []string
	GetAllVariants() map[string][]string
}

type dependencies struct {
	deps		map[string]map[string]DependencyIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewDependencies(deps ...DependencyIfc) *dependencies {
	r := dependencies{
		deps:	make(map[string]map[string]DependencyIfc),
	}
	r.Add(deps...)
	return &r
}

// -------------------------------------------------------------------------------------------------
// DependenciesIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencies) Add(deps ...DependencyIfc) *dependencies {
	for _, dep := range deps {
		name := dep.GetName()
		if _, ok := r.deps[name]; ! ok {
			r.deps[name] = make(map[string]DependencyIfc)
		}
		r.deps[name][dep.GetVariant()] = dep
	}
	return r
}

// Get the named Dependency; try Default, or any; use GetVariant() if specific variant needed!
func (r *dependencies) Get(name string) DependencyIfc {
	variant := ""
	if r.HasVariant(name, DEP_VARIANT_DEFAULT) { variant = DEP_VARIANT_DEFAULT }
	if variants := r.GetVariants(name); len(variants) > 0 { variant = variants[0] }
	return r.GetVariant(name, variant)
}

func (r *dependencies) GetVariant(name, variant string) DependencyIfc {
	if d, ok := r.deps[name][variant]; ok { return d }
	return nil
}

// Check if we have the named Dependency; try Default, or any; use HasVariant() if specific variant needed!
func (r *dependencies) Has(name string) bool {
	variant := ""
	if r.HasVariant(name, DEP_VARIANT_DEFAULT) { variant = DEP_VARIANT_DEFAULT }
	if variants := r.GetVariants(name); len(variants) > 0 { variant = variants[0] }
	return r.HasVariant(name, variant)
}

// Do we have the named dependency and variant combo?
func (r *dependencies) HasVariant(name, variant string) bool {
	_, ok := r.deps[name][variant]
	return ok
}

// Get all the variant names for this Dependency name, if any
func (r *dependencies) GetVariants(name string) []string {
	vstrs := []string{}
	if variants, ok := r.deps[name]; ok {
		for variant, _ := range variants { vstrs = append(vstrs, variant) }
	}
	return vstrs
}

// Get all the variant names for all Dependency names, if any
func (r *dependencies) GetAllVariants() map[string][]string {
	vmap := make(map[string][]string)
	for name, _ := range r.deps { vmap[name] = r.GetVariants(name) }
	return vmap
}

