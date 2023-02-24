package starter

/*
Embed, include, or implement this interface for constructs that require a single, initial Start
operation to initialize their state. There is also constructor function support to optionally
pass in additional Startable interfaces which must be started before we can consider ourselves to be
started.

TODO:
 * Create Stoppable counterpart; make a common base that Startable and Stoppable both use
*/

type StartableIfc interface {
	Start() error
}

// Exported to support embedding
type Startable struct {
	isStarted		bool
	startables		[]StartableIfc
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewStartable(startables ...StartableIfc) *Startable {
	subStartables := []StartableIfc{}
	for _, startable := range startables {
		subStartables = append(subStartables, startable)
	}
	return &Startable{
		startables:		subStartables,
	}
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

// Start everything; nil error indicates success
func (r *Startable) Start() error {
	if r.isStarted { return nil }
	for _, startable := range r.startables {
		if err := startable.Start(); nil != err { return err }
	}
	r.isStarted = true
	return nil
}

