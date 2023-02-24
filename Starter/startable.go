package starter

type StartableIfc interface {
	Check() error
}

// Exported to support embedding
type Startable struct {
	isStarted		bool
	startChecks		[]StartableIfc
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewStartable(startChecks ...StartableIfc) *Startable {
	i := Startable{ startChecks: []StartableIfc{} }
	for _, startCheck := range startChecks { i.startChecks = append(i.startChecks, startCheck) }
	return &i
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

// Run start checks; nil error indicates success
func (r *Startable) Check() error {
	if r.isStarted { return nil }
	for _, startCheck := range r.startChecks { if err := startCheck.Check(); nil != err { return err } }
	r.isStarted = true
	return nil
}

