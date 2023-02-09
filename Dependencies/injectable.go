package dependencies

/*
DependencyInjectableIfc is meant to indicate that a given implementation is prepared to accept Dependencies through this
common interface. But wait, there's more: we also include an embeddable DependencyInjectable struct which provides any
struct that embeds it with all the basic operations of dealing with and checking Dependencies that have been injected.
This eliminates boilerplate code from the general, cross-cutting concern of Dependency Injection so that other classes
can benefit from the same capabilities without duplicate code running rampant.

ref: https://en.wikipedia.org/wiki/Dependency_injection
*/

// Implementation can consume injected DependencyInstanceIfc's
type DependencyInjectableIfc interface {
        ConsumeDependencies(depinst ...DependencyInstanceIfc) error
}

// Implementation can provide information about its declared dependencies
type DependencyDiscoveryIfc interface {
	// What are all the declared Dependecies?
	GetDeclaredDependencies() DependenciesIfc
	// What are just the required Dependencies?
	GetRequiredDependencies() DependenciesIfc
	// What Dependencies are Required that have not yet been injected?
	GetMissingDependencies() DependenciesIfc
	// What are just the optional Dependencies?
	GetOptionalDependencies() DependenciesIfc
	// What are the injected DependencyInstances?
	GetInjectedDependencies() DependenciesIfc
	// Have all the required Dependencies been injected?
	HasRequiredDependencies() bool
}

