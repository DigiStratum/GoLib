package dependencies

/*
Implementation can provide information about its declared dependencies

TODO:
* Add support for Discovery of "extra" dependencies injected, but undeclared

*/

type DependencyDiscoveryIfc interface {
	// What are all the declared Dependecies?
	GetDeclaredDependencies() readableDependenciesIfc
	// What are just the required Dependencies?
	GetRequiredDependencies() readableDependenciesIfc
	// What Dependencies are Required that have not yet been injected?
	GetMissingDependencies() readableDependenciesIfc
	// What are just the optional Dependencies?
	GetOptionalDependencies() readableDependenciesIfc
	// What are the injected DependencyInstances?
	GetInjectedDependencies() readableDependenciesIfc
	// Have all the required Dependencies been injected?
	HasAllRequiredDependencies() bool
}

