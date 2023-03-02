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

	"github.com/DigiStratum/GoLib/Starter"
)

// This interface may not be used, but helps for readability here nonetheless
type DependencyInjectedIfc interface {
	DependencyDiscoveryIfc
	DependencyInjectableIfc
	readableDependenciesIfc
	starter.StartedIfc

	GetInstance(name string) interface{}
	GetInstanceVariant(name, variant string) interface{}
}

type DependencyInjected struct {
	*starter.Started
	declared	DependenciesIfc
	injected	map[string]map[string]DependencyInstanceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInjected(declaredDependencies DependenciesIfc) *DependencyInjected {
	return &DependencyInjected{
		Started:	starter.NewStarted(),
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
	return 0 == len(r.GetMissingDependencies().GetVariants())
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) InjectDependencies(depinst ...DependencyInstanceIfc) error {
	for _, instance := range depinst {
		if nil == instance { continue }
		name := instance.GetName()

		// Only capture declared dependencies, ignore extras
		if ! r.declared.Has(name) { continue }

		// Capture into our map for basic access
		if _, ok := r.injected[name]; ! ok {
			r.injected[name] = make(map[string]DependencyInstanceIfc)
		}
		r.injected[name][instance.GetVariant()] = instance

		// If this declared dependency defines Capture Func...
		declaredDep := r.declared.Get(name)
		if declaredDep.CanCapture() {
			err := declaredDep.Capture(instance.GetInstance())
			if nil != err { return err }
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------
// readableDependenciesIfc
// -------------------------------------------------------------------------------------------------

// Note: Naming of these does not reflect that it inspects the INJECTED dependencies, not DECLARED

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
	// Try default variant first; if we find it, great...
	if res := r.GetInstanceVariant(name, DEP_VARIANT_DEFAULT); nil != res { return res }

	// No default variant - take the first match (if there are any); correlates to "the" variant for this name
	if variants, ok := r.injected[name]; ok {
		for _, variant := range variants {
			if res := r.GetInstanceVariant(name, variant.GetVariant()); nil != res { return res }
		}
	}

	// We got nothing
	return nil
}

func (r *DependencyInjected) GetInstanceVariant(name, variant string) interface{} {
	if depinst, ok := r.injected[name][variant]; ok { return depinst.GetInstance() }
	return nil
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) Start() error {
	if r.Started.IsStarted() { return nil }
	missingDepVariants := r.GetMissingDependencies().GetVariants()
	if 0 == len(missingDepVariants) {
		r.Started.SetStarted()
		return nil
	}
	mdvs := []string{}
	for name, variants := range missingDepVariants {
		for _, variant := range variants {
			mdvs = append(mdvs, fmt.Sprintf("%s:%s", name, variant))
		}
	}
	return fmt.Errorf("Missing one or more required dependencies: %s", strings.Join(mdvs, ","))
}

