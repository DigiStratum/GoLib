// DigiStratum GoLib - Dependency Injection
package golib

/*
Dependencies - Implement a container to hold named dependencies and an interface for injection
*/
type Dependencies struct {
	deps	map[string]interface{}
}

// Whatever implements this interface is able to receive dependencies
type DependencyInjectableIfc interface {
        InjectDependencies(deps *Dependencies) error
}

// Whatever implements this interface initializes dependencies that will be injected elsewhere
// Note: this happens AFTER Configure() stage in order to see that all configuration data is
// available in order to initialize the dependencies (some of which may need configuration!)
type DependencyInitializerIfc interface {
	InitDependencies() error
}

// Make a new one of these!
func NewDependencies() *Dependencies {
	d := Dependencies{
		deps:	make(map[string]interface{}),
	}
	return &d
}

// Set a dependency by name
func (d *Dependencies) Set(name string, dep interface{}) {
	(*d).deps[name] = dep
}

// Get a dependency by name
func (d *Dependencies) Get(name string) interface{} {
	if i, ok := (*d).deps[name]; ok { return i }
	return nil
}

// Check whether we have dependencies for all the names
func (d *Dependencies) HasAll(names *[]string) bool {
	for _, name := range *names {
		if _, ok := (*d).deps[name]; ! ok { return false }
	}
	return true
}


