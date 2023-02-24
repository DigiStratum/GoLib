package init

// Interface to determine whether this one thing has been initialized
type InitializedIfc interface {
	// Embedder needs to implement this:
	InitializableIfc

	// We implement this:
	SetInitialized()
	IsInitialized() bool
}

// Exported to support embedding
type Initialized struct {
	isInitialized		bool
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewInitialized() *Initialized {
	return &Initialized{}
}

// -------------------------------------------------------------------------------------------------
// InitializedIfc
// -------------------------------------------------------------------------------------------------

func (r *Initialized) SetInitialized() {
	r.isInitialized = true
}

func (r *Initialized) IsInitialized() bool {
	return r.isInitialized
}

