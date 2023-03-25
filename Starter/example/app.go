package main

import (
	"fmt"

	"github.com/DigiStratum/GoLib/Starter"
)

type app struct {
	*starter.Startable
	svc		ServiceIfc
	wsvc		*WrappableService
}

func NewApp() *app {
	a := app{
		svc:		NewService(),
		wsvc:		NewWrappableService(),
	}

	// Declare Startables
	a.Startable = starter.NewStartable(
		a.svc,
		starter.MakeStartable(
			func () error {
				if a.wsvc.Init() { return nil }
				return fmt.Errorf("Failed to Init() WrappableService")
			},
		),
	)
	return &a
}

func (r *app) DoSomething() error {
	if ! r.Startable.IsStarted() { return fmt.Errorf("App not Started!") }
	if err := r.svc.DoSomething(); nil != err { return err }
	if err := r.wsvc.DoSomething(); nil != err { return err }
	return nil
}

