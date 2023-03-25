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

const DEP_VARIANT_DEFAULT = "default"

type DependencyCaptureFunc func (instance interface{}) error

type DependencyIfc interface {
	GetName() string

	SetVariant(variant string) *dependency
	GetVariant() string

	SetRequired() *dependency
	IsRequired() bool

	CanCapture() bool
	CaptureWith(captureFunc DependencyCaptureFunc) *dependency
	Capture(instance interface{}) error
}

type dependency struct {
	name, variant		string
	isRequired		bool
	captureFunc		DependencyCaptureFunc
}


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

func (r *dependency) SetVariant(variant string) *dependency {
	r.variant = variant
	return r
}

func (r *dependency) GetVariant() string {
	return r.variant
}

func (r *dependency) SetRequired() *dependency {
	r.isRequired = true
	return r
}

func (r *dependency) IsRequired() bool {
	return r.isRequired
}

func (r *dependency) CanCapture() bool {
	return nil != r.captureFunc
}

func (r *dependency) CaptureWith(captureFunc DependencyCaptureFunc) *dependency {
	r.captureFunc = captureFunc
	return r
}

func (r *dependency) Capture(instance interface{}) error {
	if ! r.CanCapture() {
		return fmt.Errorf(
			"No Capture function is set for dependency: %s:%s",
			r.GetName(), r.GetVariant(),
		)
	}

	return r.captureFunc(instance)
}

