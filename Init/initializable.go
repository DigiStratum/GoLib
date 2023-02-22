package init

type InitializableIfc interface {
	IsInitialized() bool
}

// Exported to support embedding
type Initializable struct {
	isInitialized		bool
	initChecks		[]func() bool
}

func NewInitializable(initChecks ...func() bool) *Initializable {
	i := Initializable{ initChecks: []func() bool{} }
	for _, initCheck := range initChecks { i.initChecks = append(i.initChecks, initCheck) }
	return &i
}

func (r *Initializable) IsInitialized() bool {
	if r.isInitialized { return true }
	for _, initCheck := range r.initChecks { if ! initCheck() { return false } }
	r.isInitialized = true
	return true
}

