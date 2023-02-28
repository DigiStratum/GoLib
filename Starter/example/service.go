package main

import (
	"fmt"

	"github.com/DigiStratum/GoLib/Starter"
)

type ServiceIfc interface {
	starter.StartedIfc

	DoSomething() error
}

type Service struct {
	starter.Started
	message			string
}

func NewService() *Service {
	started := starter.NewStarted()
	return &Service{
		Started:	*started,
		message:	"",
	}
}

// ------------------------------------------------------------------------------------------------
// StartableIfc
// ------------------------------------------------------------------------------------------------

func (r *Service) Start() error {
	r.message = "Service was started! :^)"
	r.Started.SetStarted()
	return nil
}

// ------------------------------------------------------------------------------------------------
// ServiceIfc
// ------------------------------------------------------------------------------------------------

func (r *Service) DoSomething() error {
	if (! r.Started.IsStarted()) { return fmt.Errorf("Service not started! :^(\n") }
	fmt.Printf("%s\n", r.message)
	return nil
}

