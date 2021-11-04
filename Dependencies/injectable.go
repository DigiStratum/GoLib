package dependencies

// Whatever implements this interface is able to receive dependencies
type DependencyInjectableIfc interface {
        InjectDependencies(deps DependenciesIfc) error
}

type DependencyInjectable struct {
	deps		DependenciesIfc
	required	[]string
	optional	[]string
	isValid		bool
}

func NewDependencyInjectable(deps DependenciesIfc) *DependencyInjectable {
	if nil == deps { return nil }
	return &DependencyInjectable{
		deps:		deps,
		required:	make([]string, 0),
		optional:	make([]string, 0),
	}
}

func (r *DependencyInjectable) SetRequired(required []string) *DependencyInjectable {
	if nil != r { r.required = required }
	return r
}

func (r *DependencyInjectable) SetOptional(optional []string) *DependencyInjectable {
	if nil != r { r.optional = optional }
	return r
}

func (r *DependencyInjectable) IsValid() bool {
	if nil == r { return false }
	// If required dependencies are specified...
	if len(r.required) > 0 {
		// All required dependency names must have been provided and non-nil
		for i, name := range (*r).required {
			if ! r.deps.Has(name) { return false }
		}
		if ! r.deps.HasAll(&(r.required)) { return false }
	}
	depNames := r.deps.GetNames()
}

// If some named dependencies are required, then they must all be present
func (r *DependencyInjectable) GetMissingRequiredDependencyNames() *[]string {
	if nil == r { return nil }
	missingDeps := make([]string, 0)
	if len(r.required) > 0 {
		// For each of the required dependency names...
		for _, name := range (*r).required {
			// ... is this named dependency present...?
			if r.deps.Has(name) {
				// ... and non-nil?
				dep := r.deps.Get(name)
				if nil != dep { continue }
			}
			// Missing or nil!
			missingDeps = append(missingDeps, name)
		}
	}
	return &missingDeps
}

// If some named dependencies are optional, then all present must be valid (either required or optional)
func (r *DependencyInjectable) GetInvalidDependencyNames() *[]string {
	if nil == r { return nil }
	invalidDeps := make([]string, 0)
	if len(r.optional) > 0 {
		// Collect up the valid names
		validNames := r.optional
		if len(r.required) > 0 {
			for _, name := range r.required {
				validNames = append(validNames, name)
			}
		}

		// Make sure all the supplied dependency names are in the set of valid names
		depNames := r.deps.GetNames()
		for _, name := range *depNames {
			for _, validName := range validNames {
				if name == validName { continue nextDepName }
			}
			invalidDeps = append(invalidDeps, name)
			nextDepName:
		}
	}
	return &invalidDeps
}
