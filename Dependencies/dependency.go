package dependencies

type Dependency struct {
	name	string
	dep	interface{}
}

func NewDependency(name string, dep interface{}) *Dependency {
	return &Dependency{
		name:	name,
		dep:	dep,
	}
}

func (r Dependency) GetDep() (string, interface{}) {
	return r.name, r.dep
}
