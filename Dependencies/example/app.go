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
		dep.NewDependency("Service").SetRequired().CaptureWith(a.captureService),
	)

	return a
}

func (r *app) captureService(instance interface{}) error {
	if nil != instance {
		var ok bool
		if r.svc, ok = instance.(ServiceIfc); ok { return nil }
	}
	return fmt.Errorf("captureService() - Instance is not a ServiceIfc")
}

func (r *app) Run() {
	if ! r.DependencyInjectable.IsStarted() { fmt.Println("Not Started Yet...") }
	if nil == r.svc { fmt.Println("No Service from DI!") }
	r.svc.Activity()
}

