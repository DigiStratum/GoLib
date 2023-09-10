package main

import (
	"fmt"

	"github.com/DigiStratum/GoLib/Process/startable"
)

type ServiceIfc interface {
	// Embedded interface(s)
	startable.StartableIfc

	// Our own interface
	DoSomething() error
}

type Service struct {
	// Embedded properties
	*startable.Startable

	// Our own properties
	message			string
}

func NewService() *Service {
	return &Service{
		Startable:	startable.NewStartable(),
		message:	"Service Not Started",
	}
}

// ------------------------------------------------------------------------------------------------
// StartableIfc
// ------------------------------------------------------------------------------------------------

func (r *Service) Start() error {
	r.message = "Service was started! :^)"
	return r.Startable.Start()
}

// ------------------------------------------------------------------------------------------------
// ServiceIfc
// ------------------------------------------------------------------------------------------------

func (r *Service) DoSomething() error {
	if ! r.Startable.IsStarted() { return fmt.Errorf("%s", r.message) }
	fmt.Printf("%s\n", r.message)
	return nil
}

