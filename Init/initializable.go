package init

type InitializableIfc interface {
	IsInitialized() bool
	CheckWith(func() bool) *Initializable
}

// Exported to support embedding
type Initializable struct {
	initialized		bool
	initChecks		[]func() bool
}

func NewInitializable() *Initializable {
	return &Initializable{}
}

func (r *Initializable) IsInitialized() bool {
	if r.initialized { return true }
	for _, initCheck := range r.initChecks { if ! initCheck() { return false } }
	r.initialized = true
	return true
}

func (r *Initializable) CheckWith(initCheck func() bool) *Initializable {
	r.initChecks = append(r.initChecks, initCheck)
	return r
}

