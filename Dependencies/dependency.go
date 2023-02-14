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
	SetVariant(variant string) *dependency
	SetRequired() *dependency
}

type dependency struct {
	name, variant		string
	isRequired		bool
}

const DEP_VARIANT_DEFAULT = "default"

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependency(name string) *dependency {
	return &dependency{
		name:		name,
		variant:	DEP_VARIANT_DEFAULT,	// Optionally override with SetVariant()
		isRequired:	false,			// Optionally override with SetRequired()
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

func (r *dependency) SetVariant(variant string) *dependency {
	r.variant = variant
	return r
}

func (r *dependency) SetRequired() *dependency {
	r.isRequired = true
	return r
}

