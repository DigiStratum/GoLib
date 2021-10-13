// DigiStratum GoLib - Dependency Injection
package dependencies

/*
Dependencies - Implement a container to hold named dependencies and an interface for injection
*/
type Dependencies struct {
	deps	map[string]interface{}
}

type DependenciesIfc interface {
	Set(name string, dep interface{})
	Get(name string) interface{}
	Has(name string) bool
	HasAll(names *[]string) bool
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
