package starter

// Interface to determine whether this one thing has been startialized
type StartedIfc interface {
	// Embedder needs to implement this:
	StartableIfc

	// We implement this:
	SetStarted()
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

func (r *Started) SetStarted() {
	r.isStarted = true
}

func (r *Started) IsStarted() bool {
	return r.isStarted
}

