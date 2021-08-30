package dependencies

// Whatever implements this interface is able to receive dependencies
type DependencyInjectableIfc interface {
        InjectDependencies(deps DependenciesIfc) error
}
