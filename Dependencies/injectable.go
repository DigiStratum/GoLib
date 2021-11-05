package dependencies

import (
	"github.com/DigiStratum/GoLib/Data/stringset"
)

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
	missingDeps := r.GetMissingRequiredDependencyNames()
	if (nil != missingDeps) && (len(*missingDeps) > 0) { return false }
	invalidDeps := r.GetInvalidDependencyNames()
	if (nil != invalidDeps) && (len(*invalidDeps) > 0) { return false }
	return true
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
	if (nil == r) || (len(r.optional) == 0) { return nil }
	givenNames := stringset.NewStringSet()
	givenNames.SetAll(r.deps.GetNames())
	givenNames.DropAll(&r.optional)
	if len(r.required) > 0 { givenNames.DropAll(&r.required) }
	invalidDeps := givenNames.ToArray()
	if len(*invalidDeps) == 0 { return nil }
	return invalidDeps
}
