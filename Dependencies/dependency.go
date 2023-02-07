package dependencies

/*

A Dependency is metadata to uniquely identify some DependencyInstance. It is convenient to pass
around the Dependency metadata for purposes of identifying, organizing, and assessing dependency
requirements. Thus, an Injectable (which consumes DependencyInstances) can declaratively identify
the dependencies that it requires as a set of Dependencies such that the Injector (which injects
DependencyInstances) can ensure that the Dependency requirements are met.

*/

import (
	"fmt"
)

type DependencyIfc interface {
	GetName() string
	GetVariant() string
	GetUniqueId() string
}

type dependency struct {
	name, variant		string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependency(name, variant string) *dependency {
	return &dependency{
		name:		name,
		variant:	variant,
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyIfc
// -------------------------------------------------------------------------------------------------

func (r *dependency) GetName() string {
	return r.name
}

func r *dependency) GetVariant() string {
	return r.variant
}

func (r *dependency) GetUniqueId() string {
	return fmt.Sprintf("%s(%s)", r.name, r.variant)
}

