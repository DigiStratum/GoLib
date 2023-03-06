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
	"github.com/DigiStratum/GoLib/Data/maps"
)

// Implementation can consume injected DependencyInstanceIfc's
type DependencyInjectableIfc interface {
	// Embedded interface(s)
	starter.StartableIfc

	// Our own interface

	// Injection & Retrieval
	// Receive dependency injection from external source
	InjectDependencies(depinst ...DependencyInstanceIfc) error
	// Get the injected dependency matching name (default or any variant)
	GetInstance(name string) interface{}
	// Get the injected dependency matching name and specific variant
	GetInstanceVariant(name, variant string) interface{}
	// Have all required dependencies have been injected?
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

// Make a new one of these
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
	// If there are any missing required dependencies...
	if missingDeps := r.GetMissingDependencies().GetAllVariants(); 0 < len(missingDeps) {
		// Make a comprehensive list of them and return that as the error
		mdvs := []string{}
		for name, variants := range missingDeps {
			for _, variant := range variants {
				mdvs = append(mdvs, fmt.Sprintf("%s:%s", name, variant))
			}
		}
		return fmt.Errorf("Missing one or more required dependencies: %s", strings.Join(mdvs, ","))
	}
	// Nothing missing, let's start!
	return r.Startable.Start()
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc
// -------------------------------------------------------------------------------------------------

// What are all the declared Dependecies?
func (r *DependencyInjectable) GetDeclaredDependencies() DependenciesIfc {
	return r.declared
}

// What are just the required Dependencies?
func (r *DependencyInjectable) GetRequiredDependencies() DependenciesIfc {
	return r.getFilteredDependencies(true)
}

// What are just the optional Dependencies?
func (r *DependencyInjectable) GetOptionalDependencies() DependenciesIfc {
	return r.getFilteredDependencies(false)
}

func (r *DependencyInjectable) getFilteredDependencies(matchRequired bool) DependenciesIfc {
	matches := NewDependencies()
	// For each of the declared dependency names...
	for name, variants := range r.declared.GetAllVariants() {
		// For each of the variants...
		for _, variant := range variants {
			// Note the ones matching matchRequired
			if matchRequired == r.declared.GetVariant(name, variant).IsRequired() {
				matches.Add(NewDependency(name).SetVariant(variant))
			}
		}
	}
	return matches
}

// What Dependencies are Required that have not yet been injected?
func (r *DependencyInjectable) GetMissingDependencies() DependenciesIfc {
	missing := NewDependencies()
	injected := r.GetInjectedDependencies()
	required := r.GetRequiredDependencies()
	// For each of the required dependency names...
	for name, variants := range required.GetAllVariants() {
		// For each of the variants...
		for _, variant := range variants {
			// Note the ones not yet injected as missing, skip others
			if injected.HasVariant(name, variant) { continue }
			missing.Add(NewDependency(name).SetVariant(variant).SetRequired())
		}
	}
	return missing
}

// What are the injected DependencyInstances?
func (r *DependencyInjectable) GetInjectedDependencies() DependenciesIfc {
	injected := NewDependencies()
	// For each of the injected DependencyInstance names...
	for name, variants := range r.injected {
		// For each of the variants...
		for variant, _:= range variants {
			// Make Dependency matching name and vairant
			dep := NewDependency(name).SetVariant(variant)
			// If it's a declared Dependency, match required state
			decl := r.declared.GetVariant(name, variant)
			if (nil != decl) && decl.IsRequired() { dep.SetRequired() }
			injected.Add(dep)
		}
	}
	return injected
}

// What Dependencies are injected, but unknown (undeclared) to us?
func (r *DependencyInjectable) GetUnknownDependencies() DependenciesIfc {
	unknown := NewDependencies()
	injected := r.GetInjectedDependencies()
	declared := r.GetDeclaredDependencies()
	// For each of the injected dependency names...
	for name, variants := range r.injected {
		// For each of the variants...
		for variant, _:= range variants {
			// If it is declared, skip it
			if declared.HasVariant(name, variant) { continue }
			// Otherwise, it is unknown!
			unknown.Add(injected.GetVariant(name, variant))
		}
	}
	return unknown
}

// Receive dependency injection from external source
func (r *DependencyInjectable) InjectDependencies(depinst ...DependencyInstanceIfc) error {
	// For each of the injected DependencyInstances...
	for _, instance := range depinst {
		if nil == instance { continue }
		name := instance.GetName()

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

// Get the injected dependency matching name (default or any variant)
func (r *DependencyInjectable) GetInstance(name string) interface{} {
        variant := ""
	if _, ok := r.injected[name][DEP_VARIANT_DEFAULT]; ok { variant = DEP_VARIANT_DEFAULT }
        if vars, ok := r.injected[name]; ok && (len(vars) > 0) { variant = maps.Strkeys(vars)[0] }
        return r.GetInstanceVariant(name, variant)
}

// Get the injected dependency matching name and specific variant
func (r *DependencyInjectable) GetInstanceVariant(name, variant string) interface{} {
	if depinst, ok := r.injected[name][variant]; ok { return depinst.GetInstance() }
	return nil
}

// Have all required dependencies have been injected?
// TODO: Do we need this, or is StartableIfc.IsStarted() sufficient?
func (r *DependencyInjectable) HasAllRequiredDependencies() bool {
	return 0 == len(r.GetMissingDependencies().GetAllVariants())
}

