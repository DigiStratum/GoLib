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
	GetUniqueId() string
	GetInstance() interface{}
}

type dependencyInstance struct {
	dep		*dependency
	instance	interface{}
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInstance(name, variant string, instance interface{}) *dependencyInstance {
	return &dependencyInstance{
		dep:		NewDependency(name, variant, false),
		instance:	instance,
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyInstanceIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencyInstance) GetDependency() *dependency {
	return r.dep
}

func (r *dependencyInstance) GetName() string {
	return r.dep.GetName()
}

func (r *dependencyInstance) GetVariant() string {
	return r.dep.GetVariant()
}

func (r *dependencyInstance) GetUniqueId() string {
	return r.dep.GetUniqueId()
}

func (r *dependencyInstance) GetInstance() interface{} {
	return r.instance
}

