package dependencies

/*
DependencyInjectableIfc is meant to indicate that a given implementation is prepared to accept Dependencies through this
common interface. But wait, there's more: we also include an embeddable DependencyInjectable struct which provides any
struct that embeds it with all the basic operations of dealing with and checking Dependencies that have been injected.
This eliminates boilerplate code from the general, cross-cutting concern of Dependency Injection so that other classes
can benefit from the same capabilities without duplicate code running rampant.

ref: https://en.wikipedia.org/wiki/Dependency_injection
*/

// Whatever implements this interface is able to publicly receive Dependencies
type DependencyInjectableIfc interface {
	GetDependencyRequirements() []DependencyIfc
	GetDependenciesInjected() []DependencyIfc
	HasRequiredDependencies() bool
        InjectDependencies(deps DependenciesIfc) error
}

// Dependency can inject itself into a provided client
type DependencyInjectorIfc interface {
        InjectInto(client DependencyInjectableIfc) error
}

