package dependencies

/*
Whatever struct embeds this struct inherit all of the behaviors attached to the struct. The

Note: if the embedding struct uses a lowercase (non-exported) alias, then this embedded struct will be private
ref: https://knight.sc/software/2018/09/20/type-embedding-in-go.html
*/

import (
	"github.com/DigiStratum/GoLib/Data/stringset"
)

// This interface may not be used, but helps for readability here nonetheless
type DependencyInjectedIfc interface {
	SetRequired(required *[]string) *DependencyInjected
	GetRequired() *[]string
	NumRequired() int
	SetOptional(optional *[]string) *DependencyInjected
	GetOptional() *[]string
	NumOptional() int
	IsValid() bool
	GetMissingRequiredDependencyNames() *[]string
	GetInvalidDependencyNames() *[]string
}

type DependencyInjected struct {
	isInstantiated	bool
	deps		DependenciesIfc
	required	[]string
	optional	[]string
	isValid		bool
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInjected(deps DependenciesIfc) *DependencyInjected {
	if nil == deps { return nil }
	return &DependencyInjected{
		isInstantiated:	true,
		deps:		deps,
		required:	make([]string, 0),
		optional:	make([]string, 0),
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyInjected Public Interface
// -------------------------------------------------------------------------------------------------

// TODO: make this additive?
func (r *DependencyInjected) SetRequired(required *[]string) *DependencyInjected {
	if r.isInstantiated { r.required = *required }
	return r
}

func (r DependencyInjected) GetRequired() *[]string {
	required := make([]string, r.NumRequired())
	for i := range r.required { required[i] = r.required[i] }
	return &required
}

func (r DependencyInjected) NumRequired() int {
	return len(r.required)
}

// TODO: make this additive?
func (r *DependencyInjected) SetOptional(optional *[]string) *DependencyInjected {
	if r.isInstantiated { r.optional = *optional }
	return r
}

func (r DependencyInjected) GetOptional() *[]string {
	optional := make([]string, r.NumOptional())
	for i := range r.optional { optional[i] = r.optional[i] }
	return &optional
}

func (r DependencyInjected) NumOptional() int {
	return len(r.optional)
}

func (r *DependencyInjected) IsValid() bool {
	if ! r.isInstantiated { return false }
	missingDeps := r.GetMissingRequiredDependencyNames()
	if (nil != missingDeps) && (len(*missingDeps) > 0) { return false }
	invalidDeps := r.GetInvalidDependencyNames()
	if (nil != invalidDeps) && (len(*invalidDeps) > 0) { return false }
	return true
}

// If some named dependencies are required, then they must all be present
func (r *DependencyInjected) GetMissingRequiredDependencyNames() *[]string {
	if ! r.isInstantiated { return nil }
	missingDeps := make([]string, 0)
	if r.NumRequired() > 0 {
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
func (r *DependencyInjected) GetInvalidDependencyNames() *[]string {
	if (! r.isInstantiated) || (len(r.optional) == 0) { return nil }
	givenNames := stringset.NewStringSet()
	givenNames.SetAll(r.deps.GetNames())
	givenNames.DropAll(&r.optional)
	if len(r.required) > 0 { givenNames.DropAll(&r.required) }
	invalidDeps := givenNames.ToArray()
	if len(*invalidDeps) == 0 { return nil }
	return invalidDeps
}
