// DigiStratum GoLib - Dependency Injection
package dependencies

/*
Dependencies is a Dependency set; it represents the complete collection of Dependencies as an
expression of what a client needs/wants from the provider.
*/

type readableDependenciesIfc interface {
	// Get a dependency by name/variant
	Get(name string) *dependency
	GetVariant(name, variant string) *dependency

	// Check whether a dependency is in the set by name/variant
	Has(name string) bool
	HasVariant(name, variant string) bool

	// Get the list of currently set dependencies
	GetVariants() map[string][]string
}

type DependenciesIfc interface {
	// Embed all the readableDependenciesIfc requirements
	readableDependenciesIfc
	// Add a Dependency to the set
	Add(dep *dependency)
}

type dependencies struct {
	deps		map[string]map[string]*dependency	// map[name][variant] -> DependencyIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewDependencies(deps ...*dependency) *dependencies {
	r := dependencies{
		deps:	make(map[string]map[string]*dependency),
	}
	for _, dep := range deps { r.Add(dep) }
	return &r
}

// -------------------------------------------------------------------------------------------------
// DependenciesIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencies) Add(dep *dependency) {
	r.deps[dep.GetName()][dep.GetVariant()] = dep
}

func (r *dependencies) Get(name string) *dependency {
	return r.GetVariant(name, DEP_VARIANT_DEFAULT)
}

func (r *dependencies) GetVariant(name, variant string) *dependency {
	if d, ok := r.deps[name][variant]; ok { return d }
	return nil
}

func (r *dependencies) Has(name string) bool {
	return r.HasVariant(name, DEP_VARIANT_DEFAULT)
}

func (r *dependencies) HasVariant(name, variant string) bool {
	_, ok := r.deps[name][variant]
	return ok
}

func (r *dependencies) GetVariants() map[string][]string {
	vmap := make(map[string][]string)
	for name, variants := range r.deps {
		vstrs := []string{}
		for variant, _ := range variants {
			vstrs = append(vstrs, variant)
		}
		vmap[name]=vstrs
	}
	return vmap
}

