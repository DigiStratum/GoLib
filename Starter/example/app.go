package main

import (
	"fmt"
)

type app struct {
	svc		ServiceIfc
}

func NewApp() *app {
	return  &app{
		svc:		NewService(),
	}
}

func (r *app) DoSomething() {
	// If we attempt to use the service without starting it, then we expect an error
	if err := r.svc.DoSomething(); nil != err { fmt.Printf("Expected Error: %s", err.Error()) }

	// Start it!
	r.svc.Start()

	// Now we expect no error
	if err := r.svc.DoSomething(); nil != err { fmt.Printf("Unexpected Error: %s", err.Error()) }
}

