package main

/*

This is an example of a service that might come from a third party package which has some sort of
init function, but which does not implement the Startable interface directly. We can use a
StartableWrapper in order to trivially adapt this to Startable.

*/

import (
	"fmt"
)

type WrappableService struct {
	// Our own properties
	message			string
	initialized		bool
}

func NewWrappableService() *WrappableService {
	return &WrappableService{
		message:	"WrappableService Not Initialized",
	}
}

// ------------------------------------------------------------------------------------------------
// StartableIfc
// ------------------------------------------------------------------------------------------------

func (r *WrappableService) Init() bool {
	r.message = "WrappableService was initialized! :^)"
	r.initialized = true
	return r.initialized
}

// ------------------------------------------------------------------------------------------------
// WrappableService
// ------------------------------------------------------------------------------------------------

func (r *WrappableService) DoSomething() error {
	if ! r.initialized { return fmt.Errorf("WrappableService not initialized! :^(") }
	fmt.Printf("%s\n", r.message)
	return nil
}

