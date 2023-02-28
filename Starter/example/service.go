package main

import (
	"fmt"

	"github.com/DigiStratum/GoLib/Starter"
)

type ServiceIfc interface {
	starter.StartedIfc
	DoSomething()
}

type Service struct {
	starter.Started
}

func NewService() *Service {
	started := starter.NewStarted()
	return &Service{
		Started:	*started,
	}
}

// StartedIfc
func (r *Service) Start() error {
	r.Started.SetStarted()
	return nil
}

func (r *Service) DoSomething() {
	if (! r.Started.IsStarted()) {
		fmt.Printf("Service Not Started!\n");
	} else {
		fmt.Printf("Service Did Something!\n")
	}
}

