package init

type InitializableIfc interface {
	Check() error
}

// Exported to support embedding
type Initializable struct {
	isInitialized		bool
	initChecks		[]InitializableIfc
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewInitializable(initChecks ...InitializableIfc) *Initializable {
	i := Initializable{ initChecks: []InitializableIfc{} }
	for _, initCheck := range initChecks { i.initChecks = append(i.initChecks, initCheck) }
	return &i
}

// -------------------------------------------------------------------------------------------------
// InitializableIfc
// -------------------------------------------------------------------------------------------------

// Run init checks; nil error indicates success
func (r *Initializable) Check() error {
	if r.isInitialized { return nil }
	for _, initCheck := range r.initChecks { if err := initCheck.Check(); nil != err { return err } }
	r.isInitialized = true
	return nil
}

