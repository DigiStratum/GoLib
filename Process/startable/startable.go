package startable

/*
Embed, include, or implement this interface for constructs that require a single, initial Start
operation to initialize their state. There is also constructor function support to optionally
pass in additional Startable interfaces which must be started before we can consider ourselves to be
started.

TODO:
 * Test coverage
 * Separate Lockable (Lockability to prevent mutations) from Startable (having a Start() func);
   Startable must be Lockable to block re-entry, but Lockable needn't be Startable...
*/

type StartableIfc interface {
	AddStartables(startables ...StartableIfc) *Startable
	Start() error
        IsStarted() bool
	Stop()
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
	s := Startable{
		startables:		make([]StartableIfc, 0),
	}
	return s.AddStartables(startables...)
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *Startable) AddStartables(startables ...StartableIfc) *Startable {
	for _, startable := range startables {
		r.startables = append(r.startables, startable)
	}
	return r
}

// Start everything; nil error indicates success
func (r *Startable) Start() error {
	// If we already started successfully, yay!
	if r.isStarted { return nil }

	// For each startable...
	for _, startable := range r.startables {
		// If it's already started, yay!
		if startable.IsStarted() { continue }

		// Try to start it, return error on failure
		if err := startable.Start(); nil != err { return err }
	}

	r.isStarted = true
	return nil
}

func (r *Startable) IsStarted() bool {
        return r.isStarted
}

func (r *Startable) Stop() {
	if ! r.IsStarted() { return }

	// For each startable...
	for _, startable := range r.startables {
		// Stop it. No-op if it wasn't started, so state doesn't matter.
		startable.Stop()
	}
}

