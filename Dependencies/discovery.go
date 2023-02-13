package dependencies

/*
Implementation can provide information about its declared dependencies

TODO:
* Add support for Discovery of "extra" dependencies injected, but undeclared

*/

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

