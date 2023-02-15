package dependencies

/*

A DependencyInstance is a pairing of a Dependency definition and matching interface{} instance. 

Note: we use DependencyIfc under the covers here, but only partially, so we don't expose it; we
don't use the IsRequired functionality which is only meaningful in the context of declared
dependencies.
*/

type DependencyInstanceIfc interface {
	GetName() string
	GetVariant() string
	GetInstance() interface{}
	SetVariant(variant string) *dependencyInstance
}

type dependencyInstance struct {
	*dependency	// Implements GetName() and GetVariant()
	instance	interface{}
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInstance(name string, instance interface{}) *dependencyInstance {
	return &dependencyInstance{
		dependency:	NewDependency(name),
		instance:	instance,
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyInstanceIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencyInstance) GetInstance() interface{} {
	return r.instance
}

func (r *dependencyInstance) SetVariant(variant string) *dependencyInstance {
	r.dependency.SetVariant(variant)
	return r
}

