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
			func (instance interface{}) bool {
				var ok bool
				a.svc, ok = instance.(ServiceIfc)
				return ok
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

