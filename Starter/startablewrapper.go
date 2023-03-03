package starter

import (
	"fmt"
)

// Wrap some other func as a Startable
type startableWrapper struct {
	startFunc	func () error
	isStarted	bool
}

func MakeStartable(startFunc func() error) *startableWrapper {
	return &startableWrapper{
		startFunc:		startFunc,
	}
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *startableWrapper) Start() error {
	if r.IsStarted() { return nil }
	if nil == r.startFunc { return fmt.Errorf("start func is nil!") }
	if err := r.startFunc(); nil != err { return err }
	r.isStarted = true
	return nil
}

func (r *startableWrapper) IsStarted() bool {
	return r.isStarted
}

