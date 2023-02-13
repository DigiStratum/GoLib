package dependencies

// Implementation can consume injected DependencyInstanceIfc's
type DependencyInjectableIfc interface {
        InjectDependencies(depinst ...DependencyInstanceIfc) error
}

