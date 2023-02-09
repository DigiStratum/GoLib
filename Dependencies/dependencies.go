// DigiStratum GoLib - Dependency Injection
package dependencies

/*
Dependencies is a Dependency set; it represents the complete collection of Dependencies as an
expression of what a client needs/wants from the provider.
*/

type readableDependenciesIfc interface {
	// Get a dependency by uniqueId
	Get(uniqueId string) *dependency
	// Check whether a dependency is in the set by uniqueId
	Has(uniqueId string) bool
	// Get the list of uniqueIds for the currently set dependencies
	GetUniqueIds() *[]string
}

type DependenciesIfc interface {
	// Embed all the readableDependenciesIfc requirements
	readableDependenciesIfc
	// Add a Dependency to the set
	Add(dep *dependency)
}

type dependencies struct {
	deps		map[string]*dependency
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewDependencies(deps ...*dependency) *dependencies {
	r := dependencies{
		deps:	make(map[string]*dependency),
	}
	for _, dep := range deps { r.Add(dep) }
	return &r
}

// -------------------------------------------------------------------------------------------------
// DependenciesIfc
// -------------------------------------------------------------------------------------------------

// Add a Dependency to the set
func (r *dependencies) Add(dep *dependency) {
	r.deps[dep.GetUniqueId()] = dep
}

// Get a dependency by uniqueId
func (r *dependencies) Get(uniqueId string) *dependency {
	if d, ok := r.deps[uniqueId]; ok { return d }
	return nil
}

// Check whether a dependency is in the set by uniqueId
func (r *dependencies) Has(uniqueId string) bool {
	_, ok := r.deps[uniqueId]
	return ok
}

// Get the list of uniqueIds for the currently set dependencies
func (r *dependencies) GetUniqueIds() *[]string {
	names := make([]string, 0)
	for name, _ := range r.deps { names = append(names, name) }
	return &names
}

