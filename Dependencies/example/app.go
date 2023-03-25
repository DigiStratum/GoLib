package main

import (
	"fmt"

	dep "github.com/DigiStratum/GoLib/Dependencies"
)

type app struct {
	*dep.DependencyInjectable
	svc				ServiceIfc
}

func NewApp() *app {
	a := &app{}

	// Declare Dependencies
	a.DependencyInjectable = dep.NewDependencyInjectable(
		dep.NewDependency("Service").SetRequired().CaptureWith(
			func (instance interface{}) error {
				var ok bool
				if a.svc, ok = instance.(ServiceIfc); ok { return nil }
				return fmt.Errorf("captureService() - Instance is not a ServiceIfc")
			},
		),
	)

	return a
}

func (r *app) Run() {
	if ! r.DependencyInjectable.IsStarted() { fmt.Println("Not Started Yet...") }
	if nil == r.svc { fmt.Println("No Service from DI!") }
	r.svc.Activity()
}

