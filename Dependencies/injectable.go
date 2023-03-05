package dependencies

/*

DependencyInjectable is an interface with base implementation that allows any construct to embed the data
and behaviors associated with being able to declare, receive, inspect, validate, discover, and utilize
injected Dependencies.

*/

import (
	"github.com/DigiStratum/GoLib/Starter"
)

// Implementation can consume injected DependencyInstanceIfc's
type DependencyInjectableIfc interface {
	// Embedded interfaces
	starter.StartableIfc

	// Our own interface
        InjectDependencies(depinst ...DependencyInstanceIfc) error

	// Discovery
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
	HasAllRequiredDependencies() bool

}

// Exported to support embedding
type DependencyInjectable struct {
	*starter.Startable

	declared		*dependencies
	injected		map[string]map[string]DependencyInstanceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// FIXME: @HERE We are going to structure this more similarly to Configurable and Startable
func NewDependencyInjectable(deps ...DependencyIfc) *DependencyInjectable {
	return &DependencyInjectable{
		Startable:	starter.NewStartable(),
		declared:	NewDependencies(deps...),
		injected:	make(map[string]map[string]DependencyInstanceIfc),
	}
}

