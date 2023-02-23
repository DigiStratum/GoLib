package init

type InitializableIfc interface {
	Check() error
}

// Exported to support embedding
type Initializable struct {
	isInitialized		bool
	initChecks		[]func() error
}

func NewInitializable(initChecks ...func() error) *Initializable {
	i := Initializable{ initChecks: []func() error{} }
	for _, initCheck := range initChecks { i.initChecks = append(i.initChecks, initCheck) }
	return &i
}

// Run init checks; nil error indicates success
func (r *Initializable) Check() error {
	if r.isInitialized { return nil }
	for _, initCheck := range r.initChecks { if err := initCheck(); nil != err { return err } }
	r.isInitialized = true
	return nil
}

