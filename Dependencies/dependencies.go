// DigiStratum GoLib - Dependency Injection
package dependencies

/*
Dependencies - Implement a container to hold named dependencies and an interface for injection
*/
type Dependencies struct {
	deps		map[string]interface{}
}

type DependenciesIfc interface {
	Set(name string, dep interface{})
	Get(name string) interface{}
	Has(name string) bool
	GetNames() *[]string
	HasAll(names *[]string) bool
	GetIterator() func () *Dependency
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewDependencies() *Dependencies {
	r := Dependencies{
		deps:	make(map[string]interface{}),
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// DependenciesIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Set a dependency by name
func (r *Dependencies) Set(name string, dep interface{}) {
	r.deps[name] = dep
}

// Get a dependency by name
func (r Dependencies) Get(name string) interface{} {
	if i, ok := r.deps[name]; ok { return i }
	return nil
}

// Check whether we have dependencies for all the names
func (r Dependencies) Has(name string) bool {
	return r.Get(name) != nil
}

// Check whether we have dependencies for all the names
func (r Dependencies) HasAll(names *[]string) bool {
	for _, name := range *names {
		if ! r.Has(name) { return false }
	}
	return true
}

// Get the list of names for the currently set dependencies
func (r Dependencies) GetNames() *[]string {
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
