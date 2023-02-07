// DigiStratum GoLib - Dependency Injection
package dependencies

/*
Dependencies - Implement a container to hold named dependencies and an interface for injection
*/

type DependenciesIfc interface {
	Add(dependency DependencyIfc)
	Get(uniqueId string) DependencyIfc
	GetUniqueIds() *[]string
	Has(uniqueId string) bool
	HasAll(uniqueIds *[]string) bool
	GetIterator() func () *Dependency
}

type dependencies struct {
	deps		map[string]DependencyIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewDependencies(deps ...DependencyIfc) *Dependencies {
	r := Dependencies{
		deps:	make(map[string]DependencyIfc),
	}
	for dep in range deps... { r.Add(dep) }
	return &r
}

// -------------------------------------------------------------------------------------------------
// DependenciesIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Add a Dependency to the set
func (r *Dependencies) Add(dep DependencyIfc) {
	r.deps[dep.GetUniqueId()] = dep
}

// Get a dependency by name
func (r Dependencies) Get(uniqueId string) DependencyIfc {
	if d, ok := r.deps[uniqueId]; ok { return d }
	return nil
}

// Check whether we have dependencies for all the names
func (r Dependencies) Has(uniqueId string) bool {
	return r.Get(uniqueId) != nil
}

// Check whether we have dependencies for all the names
func (r Dependencies) HasAll(uniqueIds *[]string) bool {
	for _, name := range *names {
		if ! r.Has(name) { return false }
	}
	return true
}

// Get the list of names for the currently set dependencies
func (r Dependencies) GetUniqueIds() *[]string {
	names := make([]string, 0)
	for name, _ := range r.deps {
		names = append(names, name)
	}
	return &names
}

// -------------------------------------------------------------------------------------------------
// IterableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Dependencies) GetIterator() func () *Dependency {
	idx := 0
	names := r.GetNames()
	// Return a Dependency(name/dep) or nil when done iterating
	return func () *Dependency {
		// If we're done iterating, return do nothing
		if idx >= len(*names) { return nil }
		name := (*names)[idx]
		dep, ok := r.deps[name]
		if ! ok { return nil } // This can only happen if someone tampers with deps while iterating
		idx++
		return NewDependency(name, dep)
	}
}
