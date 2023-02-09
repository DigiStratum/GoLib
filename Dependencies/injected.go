package dependencies

/*
Boilerplate code for DependencyInjected clients to inspect injected dependencies for completeness
and validity. Bearer must declare which dependency names are Optional and/or Required, and point
us at the injected Dependencies. Validity checking will be performed against these data points.

TODO:
 * Capture mutation vs. validity state so that IsValid() uses cached validity if not mutated,
   and mutation flag updates with changes to Set functions
 * Add support for redefinition and/or replacement of one or more Dependencies at after
   initialization to support runtime reconfigurability.

*/

// This interface may not be used, but helps for readability here nonetheless
type DependencyInjectedIfc interface {
	// This implementation supports all the Discovery functions (so embed the interface!)
	DependencyDiscoveryIfc
	// This implementation is injectable (so embed the interface!)
	DependencyInjectableIfc
	// Embed all the readableDependenciesIfc requirements
	readableDependenciesIfc

	GetInstance(uniqueId string) interface{}
}

type dependencyInjected struct {
	hasRequired	bool
	declared	DependenciesIfc
	provided	map[string]DependencyInstanceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInjected(declaredDependencies DependenciesIfc) *dependencyInjected {
	return &dependencyInjected{
		hasRequired:	false,
		declared:	declaredDependencies,
		provided:	make(map[string]DependencyInstanceIfc),
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyDiscoveryIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencyInjected) GetDeclaredDependencies() DependenciesIfc {
	declared := NewDependencies()
	for _, uniqueId := range *(r.declared.GetUniqueIds()) {
		declared.Add(r.declared.Get(uniqueId))
	}
	return declared
}

func (r *dependencyInjected) GetRequiredDependencies() DependenciesIfc {
	required := NewDependencies()
	for _, uniqueId := range *(r.declared.GetUniqueIds()) {
		dep := r.declared.Get(uniqueId)
		if dep.IsRequired() { required.Add(dep) }
	}
	return required
}

func (r *dependencyInjected) GetMissingDependencies() DependenciesIfc {
	missing := NewDependencies()
	injected := r.GetInjectedDependencies()
	required := r.GetRequiredDependencies()
	for _, uniqueId := range *(required.GetUniqueIds()) {
		if ! injected.Has(uniqueId) { missing.Add(required.Get(uniqueId)) }
	}
	return missing
}

func (r *dependencyInjected) GetOptionalDependencies() DependenciesIfc {
	optional := NewDependencies()
	for _, uniqueId := range *(r.declared.GetUniqueIds()) {
		dep := r.declared.Get(uniqueId)
		if ! dep.IsRequired() { optional.Add(dep) }
	}
	return optional
}

func (r *dependencyInjected) GetInjectedDependencies() DependenciesIfc {
	injected := NewDependencies()
	for _, instance := range r.provided {
		injected.Add(instance.GetDependency())
	}
	return injected
}

func (r *dependencyInjected) HasRequiredDependencies() bool {
	missing := r.GetMissingDependencies()
	return len(*(missing.GetUniqueIds())) == 0
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencyInjected) ConsumeDependencies(depinst ...DependencyInstanceIfc) error {
	for _, instance := range depinst {
		if nil == instance { continue }
		r.provided[instance.GetDependency().GetUniqueId()] = instance
	}
	return nil
}

// -------------------------------------------------------------------------------------------------
// readableDependenciesIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencyInjected) Get(uniqueId string) *dependency {
	if instance, ok := r.provided[uniqueId]; ok { return instance.GetDependency() }
	return nil
}

func (r *dependencyInjected) Has(uniqueId string) bool {
	_, ok := r.provided[uniqueId]
	return ok
}

func (r *dependencyInjected) GetUniqueIds() *[]string {
	uniqueIds := make([]string, len(r.provided))
	for uniqueId, _ := range r.provided { uniqueIds = append(uniqueIds, uniqueId) }
	return &uniqueIds
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectedIfc
// -------------------------------------------------------------------------------------------------

func (r *dependencyInjected) GetInstance(uniqueId string) interface{} {
	if depinst, ok := r.provided[uniqueId]; ok { return depinst.GetInstance() }
	return nil
}

