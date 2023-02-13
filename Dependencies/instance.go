package dependencies

/*

A DependencyInstance is a pairing of a Dependency definition and matching interface{} instance. 

*/

type DependencyInstanceIfc interface {
	GetDependency() *dependency
	GetInstance() interface{}
}

type dependencyInstance struct {
	dep		*dependency
	instance	interface{}
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInstance(dep *dependency, instance interface{}) *dependencyInstance {
	return &dependencyInstance{
		dep:		dep,
		instance:	instance,
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyInstanceIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencyInstance) GetDependency() *dependency {
	return r.dep
}

func (r *dependencyInstance) GetInstance() interface{} {
	return r.instance
}

