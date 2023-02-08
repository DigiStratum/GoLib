// DigiStratum GoLib - Dependency Injection
package dependencies

/*
Dependencies is a Dependency set; it represents the complete collection of Dependencies as an
expression of what a client needs/wants from the provider.
*/

type DependenciesIfc interface {
	Add(dep dependency)
	Get(uniqueId string) *dependency
	GetUniqueIds() *[]string
}

type dependencies struct {
	deps		map[string]*dependency
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewDependencies(deps ...dependency) *dependencies {
	r := dependencies{
		deps:	make(map[string]dependency),
	}
	for dep in range deps... { r.Add(dep) }
	return &r
}

// -------------------------------------------------------------------------------------------------
// DependenciesIfc
// -------------------------------------------------------------------------------------------------

// Add a Dependency to the set
func (r *Dependencies) Add(dep dependency) {
	r.deps[dep.GetUniqueId()] = dep
}

// Get a dependency by name
func (r *dependencies) Get(uniqueId string) dependency {
	if d, ok := r.deps[uniqueId]; ok { return d }
	return nil
}

// Get the list of names for the currently set dependencies
func (r *dependencies) GetUniqueIds() *[]string {
	names := make([]string, 0)
	for name, _ := range r.deps { names = append(names, name) }
	return &names
}

