package main

import (
	"fmt"

	"github.com/DigiStratum/GoLib/Starter"
)

type ServiceIfc interface {
	// Embedded interface(s)
	starter.StartableIfc

	// Our own interface
	DoSomething() error
}

type Service struct {
	// Embedded properties
	*starter.Startable

	// Our own properties
	message			string
}

func NewService() *Service {
	return &Service{
		Startable:	starter.NewStartable(),
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
	if ! r.Startable.IsStarted() { return fmt.Errorf("Service not started! :^(\n") }
	fmt.Printf("%s\n", r.message)
	return nil
}

