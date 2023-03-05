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
}

// Exported to support embedding
type DependencyInjectable struct {
	*starter.Startable
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// FIXME: @HERE We are going to structure this more similarly to Configurable and Startable
func NewDependencyInjectable(declared ...DependencyIfc) *DependencyInjectable {
	return &DependencyInjected{
		Startable:	starter.NewStartable(),
		declared:	declaredDependencies,
		injected:	make(map[string]map[string]DependencyInstanceIfc),
	}
}

