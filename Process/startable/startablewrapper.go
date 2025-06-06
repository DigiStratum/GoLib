package startable

/*

This sneaky little wrapper allows us to pass a non-exported initializtion func as a StartableIfc so
that we can enjoy the benefits of Startability without having to publicly Exported so that everyone
who has us can access these properties. Check out the example(s) below this package.

*/

import (
	"fmt"
)

// Wrap some other func as a Startable
type startableWrapper struct {
	startFunc	func () error
	stopFunc	func ()
	isStarted	bool
}

func MakeStartable(startFunc func() error, stopFunc func()) *startableWrapper {
	return &startableWrapper{
		startFunc:		startFunc,
		stopFunc:		stopFunc,
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

// This is a no-op for startableWrapper which need not support nested sub-startables
func (r *startableWrapper) AddStartables(startables ...StartableIfc) *Startable {
	return nil
}

func (r *startableWrapper) Stop() {
	r.stopFunc()
}

