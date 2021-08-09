// DigiStratum GoLib - Dependency Injection
package golib

/*
Dependencies - Implement a container to hold named dependencies and an interface for injection
*/
type Dependencies struct {
	deps	map[string]interface{}
}

type DependenciesIfc interface {
	Set(name string, dep interface{})
	Get(name string) interface{}
	Has(names string) bool
	HasAll(names *[]string) bool
}

// Whatever implements this interface is able to receive dependencies
type DependencyInjectableIfc interface {
        InjectDependencies(deps DependenciesIfc) error
}

// Make a new one of these!
func NewDependencies() *Dependencies {
	d := Dependencies{
		deps:	make(map[string]interface{}),
	}
	return &d
}

// -------------------------------------------------------------------------------------------------
// DependenciesIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Set a dependency by name
func (d *Dependencies) Set(name string, dep interface{}) {
	d.deps[name] = dep
}

// Get a dependency by name
func (d Dependencies) Get(name string) interface{} {
	if i, ok := d.deps[name]; ok { return i }
	return nil
}

// Check whether we have dependencies for all the names
func (d Dependencies) Has(name string) bool {
	return d.Get(name) != nil
}

// Check whether we have dependencies for all the names
func (d Dependencies) HasAll(names *[]string) bool {
	for _, name := range *names {
		if ! d.Has(name) { return false }
	}
	return true
}


