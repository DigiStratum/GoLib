// DigiStratum GoLib - Dependency Injection
package golib

/*
Dependencies - Implement a container to hold named dependencies and an interface for injection
*/
type dependencies struct {
	deps	map[string]interface{}
}

// Dependencies public interface
type DependenciesIfc interface {
	Set(name string, dep interface{})
	Get(name string) interface{}
	HasAll(names *[]string) bool
}

// Whatever implements this interface is able to receive dependencies
type DependencyInjectableIfc interface {
        InjectDependencies(deps DependenciesIfc) error
}

// Make a new one of these!
func NewDependencies() DependenciesIfc {
	d := dependencies{
		deps:	make(map[string]interface{}),
	}
	return &d
}

// Set a dependency by name
func (d *dependencies) Set(name string, dep interface{}) {
	(*d).deps[name] = dep
}

// Get a dependency by name
func (d *dependencies) Get(name string) interface{} {
	if i, ok := (*d).deps[name]; ok { return i }
	return nil
}

// Check whether we have dependencies for all the names
func (d *dependencies) HasAll(names *[]string) bool {
	for _, name := range *names {
		if _, ok := (*d).deps[name]; ! ok { return false }
	}
	return true
}


