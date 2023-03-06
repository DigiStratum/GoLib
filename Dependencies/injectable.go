package dependencies

/*

DependencyInjectable is an interface with base implementation that allows any construct to embed the data
and behaviors associated with being able to declare, receive, inspect, validate, discover, and utilize
injected Dependencies.

*/

import (
	"fmt"
	"strings"

	"github.com/DigiStratum/GoLib/Starter"
)

// Implementation can consume injected DependencyInstanceIfc's
type DependencyInjectableIfc interface {
	// Embedded interfaces
	starter.StartableIfc

	// Our own interface

	// Injection & Retrieval
	InjectDependencies(depinst ...DependencyInstanceIfc) error
	GetInstance(name string) interface{}
	GetInstanceVariant(name, variant string) interface{}
	HasAllRequiredDependencies() bool

	// Discovery (Declared)
	// What are all the declared Dependecies?
	GetDeclaredDependencies() DependenciesIfc
	// What are just the required Dependencies?
	GetRequiredDependencies() DependenciesIfc
	// What are just the optional Dependencies?
	GetOptionalDependencies() DependenciesIfc

	// Discovery (Injected)
	// What are the injected DependencyInstances?
	GetInjectedDependencies() DependenciesIfc
	// What Dependencies are Required that have not yet been injected?
	GetMissingDependencies() DependenciesIfc
	// What Dependencies are injected, but unknown (undeclared) to us?
	GetUnknownDependencies() DependenciesIfc
}

// Exported to support embedding
type DependencyInjectable struct {
	*starter.Startable

	declared		*dependencies
	injected		map[string]map[string]DependencyInstanceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// FIXME: @HERE We are going to structure this more similarly to Configurable and Startable
func NewDependencyInjectable(deps ...DependencyIfc) *DependencyInjectable {
	return &DependencyInjectable{
		Startable:	starter.NewStartable(),
		declared:	NewDependencies(deps...),
		injected:	make(map[string]map[string]DependencyInstanceIfc),
	}
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

// Starter func; ensure that all required dependencies are injected
func (r *DependencyInjectable) Start() error {
	if r.Startable.IsStarted() { return nil }
	if missingDeps := r.GetMissingDependencies().GetAllVariants(); 0 < len(missingDeps) {
		mdvs := []string{}
		for name, variants := range missingDeps {
			for _, variant := range variants {
				mdvs = append(mdvs, fmt.Sprintf("%s:%s", name, variant))
			}
		}
		return fmt.Errorf("Missing one or more required dependencies: %s", strings.Join(mdvs, ","))
	}
	return r.Startable.Start()
}

func (r *DependencyInjectable) GetDeclaredDependencies() DependenciesIfc {
	declared := NewDependencies()
	// Iterate over Dependencies' names and variants for each (map[(name)][]string(variant)
	for name, variants := range r.declared.GetAllVariants() {
		for _, variant := range variants {
			dep := NewDependency(name).SetVariant(variant)
			if r.declared.GetVariant(name, variant).IsRequired() { dep.SetRequired() }
			declared.Add(dep)
		}
	}
	return declared
}

func (r *DependencyInjectable) GetRequiredDependencies() DependenciesIfc {
	required := NewDependencies()
	// Iterate over Dependencies' names and variants for each (map[(name)][]string(variant)
	for name, variants := range r.declared.GetAllVariants() {
		for _, variant := range variants {
			// We're only interested in the required ones - skip others
			if ! r.declared.GetVariant(name, variant).IsRequired() { continue }
			required.Add(NewDependency(name).SetVariant(variant).SetRequired())
		}
	}
	return required
}

func (r *DependencyInjectable) GetOptionalDependencies() DependenciesIfc {
	optional := NewDependencies()
	// Iterate over Dependencies' names and variants for each (map[(name)][]string(variant)
	for name, variants := range r.declared.GetAllVariants() {
		for _, variant := range variants {
			// We're only interested in the optional (non-required) ones - skip others
			if r.declared.GetVariant(name, variant).IsRequired() { continue }
			optional.Add(NewDependency(name).SetVariant(variant))
		}
	}
	return optional
}

func (r *DependencyInjectable) GetMissingDependencies() DependenciesIfc {
	missing := NewDependencies()
	injected := r.GetInjectedDependencies()
	required := r.GetRequiredDependencies()
	// Iterate over required Dependencies' names and variants (map[(name)][]string(variant)
	for name, variants := range required.GetAllVariants() {
		for _, variant := range variants {
			// We're only interested in the missing ones - skip others
			if injected.HasVariant(name, variant) { continue }
			missing.Add(NewDependency(name).SetVariant(variant).SetRequired())
		}
	}
	return missing
}

func (r *DependencyInjectable) GetInjectedDependencies() DependenciesIfc {
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

func (r *DependencyInjectable) HasAllRequiredDependencies() bool {
	return 0 == len(r.GetMissingDependencies().GetAllVariants())
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjectable) InjectDependencies(depinst ...DependencyInstanceIfc) error {
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

func (r *DependencyInjectable) GetInstance(name string) interface{} {
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

func (r *DependencyInjectable) GetInstanceVariant(name, variant string) interface{} {
	if depinst, ok := r.injected[name][variant]; ok { return depinst.GetInstance() }
	return nil
}

