package dependencies

/*
Boilerplate code for DependencyInjected clients to inspect injected dependencies for completeness
and validity. Bearer must declare which dependency names are Optional and/or Required, and point
us at the injected Dependencies. Validity checking will be performed against these data points.

TODO:
 * Cache HasRequired() vs. mutation funcs so that we only re-eval HasRequired() as needed
 * Add support for redefinition and/or replacement of one or more Dependencies after initialization
   to support runtime reconfigurability.
 * Add support for Discovery of "extra" dependencies injected, but undeclared

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

type DependencyInjected struct {
	declared	DependenciesIfc
	injected	map[string]DependencyInstanceIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDependencyInjected(declaredDependencies DependenciesIfc) *DependencyInjected {
	return &DependencyInjected{
		declared:	declaredDependencies,
		injected:	make(map[string]DependencyInstanceIfc),
	}
}

// -------------------------------------------------------------------------------------------------
// DependencyDiscoveryIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) GetDeclaredDependencies() DependenciesIfc {
	declared := NewDependencies()
	for _, uniqueId := range *(r.declared.GetUniqueIds()) {
		declared.Add(r.declared.Get(uniqueId))
	}
	return declared
}

func (r *DependencyInjected) GetRequiredDependencies() DependenciesIfc {
	required := NewDependencies()
	for _, uniqueId := range *(r.declared.GetUniqueIds()) {
		dep := r.declared.Get(uniqueId)
		if dep.IsRequired() { required.Add(dep) }
	}
	return required
}

func (r *DependencyInjected) GetMissingDependencies() DependenciesIfc {
	missing := NewDependencies()
	injected := r.GetInjectedDependencies()
	required := r.GetRequiredDependencies()
	for _, uniqueId := range *(required.GetUniqueIds()) {
		if ! injected.Has(uniqueId) { missing.Add(required.Get(uniqueId)) }
	}
	return missing
}

func (r *DependencyInjected) GetOptionalDependencies() DependenciesIfc {
	optional := NewDependencies()
	for _, uniqueId := range *(r.declared.GetUniqueIds()) {
		dep := r.declared.Get(uniqueId)
		if ! dep.IsRequired() { optional.Add(dep) }
	}
	return optional
}

func (r *DependencyInjected) GetInjectedDependencies() DependenciesIfc {
	injected := NewDependencies()
	for _, instance := range r.injected {
		injected.Add(NewDependency(instance.GetName(), instance.GetVariant(), false))
	}
	return injected
}

func (r *DependencyInjected) HasRequiredDependencies() bool {
	missing := r.GetMissingDependencies()
	return len(*(missing.GetUniqueIds())) == 0
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) InjectDependencies(depinst ...DependencyInstanceIfc) error {
	for _, instance := range depinst {
		if nil == instance { continue }
		r.injected[instance.GetUniqueId()] = instance
	}
	return nil
}

// -------------------------------------------------------------------------------------------------
// readableDependenciesIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) Get(uniqueId string) *dependency {
	if instance, ok := r.injected[uniqueId]; ok { return NewDependency(instance.GetName(), instance.GetVariant(), false) }
	return nil
}

func (r *DependencyInjected) Has(uniqueId string) bool {
	_, ok := r.injected[uniqueId]
	return ok
}

func (r *DependencyInjected) GetUniqueIds() *[]string {
	uniqueIds := make([]string, len(r.injected))
	for uniqueId, _ := range r.injected { uniqueIds = append(uniqueIds, uniqueId) }
	return &uniqueIds
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectedIfc
// -------------------------------------------------------------------------------------------------

func (r *DependencyInjected) GetInstance(uniqueId string) interface{} {
	if depinst, ok := r.injected[uniqueId]; ok { return depinst.GetInstance() }
	return nil
}

