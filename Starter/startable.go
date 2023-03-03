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
        IsStarted() bool
}

// Exported to support embedding
type Startable struct {
	isStarted		bool
	startables		[]StartableIfc
	startablesStarted	[]bool
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
	// If we already started successfully, yay!
	if r.isStarted { return nil }

	// Prepare to capture start states for each of our startables
	r.startablesStarted = make([]bool, len(r.startables))

	// For each startable...
	for i, startable := range r.startables {
		// If it's already started, yay!
		if r.startablesStarted[i] { continue }

		// Try to start it, return error on failure
		if err := startable.Start(); nil != err { return err }

		// It started, yay! Make a mental note
		r.startablesStarted[i] = true
	}

	r.isStarted = true
	return nil
}

func (r *Startable) IsStarted() bool {
        return r.isStarted
}

