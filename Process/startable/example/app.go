package main

import (
	"fmt"

	"github.com/DigiStratum/GoLib/Process/startable"
)

type app struct {
	*startable.Startable
	svc		ServiceIfc
	wsvc		*WrappableService
}

func NewApp() *app {
	a := app{
		svc:		NewService(),
		wsvc:		NewWrappableService(),
	}

	// Declare Startables
	a.Startable = startable.NewStartable(
		a.svc,
		startable.MakeStartable(
			// Start func
			func () error {
				if a.wsvc.Init() {
					fmt.Println("Started!")
					return nil
				}
				return fmt.Errorf("Failed to Init() WrappableService")
			},
			// Stop func
			func () {
				fmt.Println("Stopped!")
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

