package starter

/*
Embed or include this into any construct that needs to be started once so that it can track its
own state as well as make that state known to consumers.
*/

type StartedIfc interface {
	// Embedded construct must implement this
	StartableIfc

	// We implement this
	SetStarted() *Started
	IsStarted() bool
}

// Exported to support embedding
type Started struct {
	isStarted		bool
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewStarted() *Started {
	return &Started{}
}

// -------------------------------------------------------------------------------------------------
// StartedIfc
// -------------------------------------------------------------------------------------------------

func (r *Started) SetStarted() *Started {
	r.isStarted = true
	return r
}

func (r *Started) IsStarted() bool {
	return r.isStarted
}
