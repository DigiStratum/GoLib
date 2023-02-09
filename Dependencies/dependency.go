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
	IsRequired() bool
}

type dependency struct {
	name, variant		string
	isRequired		bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependency(name, variant string, isRequired bool) *dependency {
	return &dependency{
		name:		name,
		variant:	variant,
		isRequired:	isRequired,
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyIfc
// -------------------------------------------------------------------------------------------------

func (r *dependency) GetName() string {
	return r.name
}

func (r *dependency) GetVariant() string {
	return r.variant
}

func (r *dependency) GetUniqueId() string {
	return fmt.Sprintf("%s(%s)", r.name, r.variant)
}

func (r *dependency) IsRequired() bool {
	return r.isRequired
}

