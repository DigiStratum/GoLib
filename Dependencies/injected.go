package dependencies

/*
Boilerplate code for DependencyInjected clients to inspect injected dependencies for completeness
and validity. Bearer must declare which dependency names are Optional and/or Required, and point
us at the injected Dependencies. Validity checking will be performed against these data points.

Note that we export the DependencyInjected struct itself so that it may be embedded into other
structs that want to inherit this functionality; if it's not exported, then it can't be accessed by
another package.

TODO:
 * Cache HasAllRequiredDependencies() vs. mutation funcs so we only re-eval fully as needed
 * Add support for redefinition and/or replacement of one or more Dependencies after initialization
   to support runtime reconfigurability.
 * Add support for Discovery of "extra" dependencies injected, but undeclared

*/

import (
	"fmt"
	"strings"
)

// This interface may not be used, but helps for readability here nonetheless
type DependencyInjectedIfc interface {
	// This implementation supports all the Discovery functions (so embed the interface!)
	DependencyDiscoveryIfc
	// This implementation is injectable (so embed the interface!)
	DependencyInjectableIfc
	// Embed all the readableDependenciesIfc requirements
	readableDependenciesIfc

	GetInstance(name string) interface{}
	GetInstanceVariant(name, variant string) interface{}
	ValidateRequiredDependencies() error
}

type DependencyInjected struct {
	declared	DependenciesIfc
	injected	map[string]map[string]DependencyInstanceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInjected(declaredDependencies DependenciesIfc) *DependencyInjected {
	return &DependencyInjected{
		declared:	declaredDependencies,
		injected:	make(map[string]map[string]DependencyInstanceIfc),
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyDiscoveryIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) GetDeclaredDependencies() readableDependenciesIfc {
	declared := NewDependencies()
	// Iterate over Dependencies' names and variants for each (map[(name)][]string(variant)
	for name, variants := range r.declared.GetVariants() {
		for _, variant := range variants {
			dep := NewDependency(name).SetVariant(variant)
			if r.declared.GetVariant(name, variant).IsRequired() { dep.SetRequired() }
			declared.Add(dep)
		}
	}
	return declared
}

func (r *DependencyInjected) GetRequiredDependencies() readableDependenciesIfc {
	required := NewDependencies()
	// Iterate over Dependencies' names and variants for each (map[(name)][]string(variant)
	for name, variants := range r.declared.GetVariants() {
		for _, variant := range variants {
			// We're only interested in the required ones - skip others
			if ! r.declared.GetVariant(name, variant).IsRequired() { continue }
			required.Add(NewDependency(name).SetVariant(variant).SetRequired())
		}
	}
	return required
}

func (r *DependencyInjected) GetOptionalDependencies() readableDependenciesIfc {
	optional := NewDependencies()
	// Iterate over Dependencies' names and variants for each (map[(name)][]string(variant)
	for name, variants := range r.declared.GetVariants() {
		for _, variant := range variants {
			// We're only interested in the optional (non-required) ones - skip others
			if r.declared.GetVariant(name, variant).IsRequired() { continue }
			optional.Add(NewDependency(name).SetVariant(variant))
		}
	}
	return optional
}

func (r *DependencyInjected) GetMissingDependencies() readableDependenciesIfc {
	missing := NewDependencies()
	injected := r.GetInjectedDependencies()
	required := r.GetRequiredDependencies()
	// Iterate over required Dependencies' names and variants (map[(name)][]string(variant)
	for name, variants := range required.GetVariants() {
		for _, variant := range variants {
			// We're only interested in the missing ones - skip others
			if injected.HasVariant(name, variant) { continue }
			missing.Add(NewDependency(name).SetVariant(variant).SetRequired())
		}
	}
	return missing
}

func (r *DependencyInjected) GetInjectedDependencies() readableDependenciesIfc {
	injected := NewDependencies()
	// Iterate over Dependencies' names and variants for each (map[(name)][]string(variant)
	for name, variants := range r.injected {
		for variant, _:= range variants {
			dep := NewDependency(name).SetVariant(variant)
			decl := r.declared.GetVariant(name, variant)
			if (nil != decl) && decl.IsRequired() { dep.SetRequired() }
			injected.Add(dep)
		}
	}
	return injected
}

func (r *DependencyInjected) HasAllRequiredDependencies() bool {
	missing := r.GetMissingDependencies()
	return len(missing.GetVariants()) == 0
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) InjectDependencies(depinst ...DependencyInstanceIfc) error {
        for _, instance := range depinst {
		if nil == instance { continue }
                name := instance.GetName()
                if _, ok := r.injected[name]; ! ok {
                        r.injected[name] = make(map[string]DependencyInstanceIfc)
                }
                r.injected[name][instance.GetVariant()] = instance
        }

	return nil
}

// -------------------------------------------------------------------------------------------------
// readableDependenciesIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) Get(name string) *dependency {
	return r.GetVariant(name, DEP_VARIANT_DEFAULT)
}

func (r *DependencyInjected) GetVariant(name, variant string) *dependency {
	return r.GetInjectedDependencies().GetVariant(name, variant)
}

func (r *DependencyInjected) Has(name string) bool {
	return r.HasVariant(name, DEP_VARIANT_DEFAULT)
}

func (r *DependencyInjected) HasVariant(name, variant string) bool {
	return r.GetInjectedDependencies().HasVariant(name, variant)
}

func (r *DependencyInjected) GetVariants() map[string][]string {
	return r.GetInjectedDependencies().GetVariants()
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectedIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) GetInstance(name string) interface{} {
	return r.GetInstanceVariant(name, DEP_VARIANT_DEFAULT)
}

func (r *DependencyInjected) GetInstanceVariant(name, variant string) interface{} {
	if depinst, ok := r.injected[name][variant]; ok { return depinst.GetInstance() }
	return nil
}

func (r *DependencyInjected) ValidateRequiredDependencies() error {
	if r.HasAllRequiredDependencies() { return nil }

	missingDeps := r.GetMissingDependencies()
	var sb strings.Builder
	delim := ""
	for name, variants := range missingDeps.GetVariants() {
		for _, variant := range variants {
			sb.WriteString(fmt.Sprintf("%s%s:%s", delim, name, variant))
			delim = ", "
		}
	}
	missingList := sb.String()
	return fmt.Errorf("Missing one or more required dependencies: %s", missingList)
}

