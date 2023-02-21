package main

import (
	"fmt"

	dep "github.com/DigiStratum/GoLib/Dependencies"
)

type app struct {
	*dep.DependencyInjected
	svc				ServiceIfc
}

func NewApp() *app {
	return  &app{
		// Declare Dependencies
		DependencyInjected: dep.NewDependencyInjected(
			dep.NewDependencies(
				dep.NewDependency("Service").SetRequired(),
			),
		),
	}
}

// If we override InjectDependencies() we can capture injected dependencies
// locally instead of asking for them every time they are needed
func (r *app) InjectDependencies(depinst ...dep.DependencyInstanceIfc) error {

	// If DI fails, return error
	if err := r.DependencyInjected.InjectDependencies(depinst...); nil != err { return err }
	// If DI missing requirements, return error
	if err := r.DependencyInjected.ValidateRequiredDependencies(); nil != err { return err }

	// Iterate over injected dependencies; use a switch-case to map
	// them to the correct interface assertion and member value
	for name, _ := range r.DependencyInjected.GetVariants() {
		switch name {
			case "Service":
				svcdep := r.DependencyInjected.GetInstance(name)
				if svc, ok := svcdep.(ServiceIfc); ok { r.svc = svc }
		}
	}

	return nil
}

func (r *app) Run() {
	if nil == r.svc { fmt.Println("No Service from DI!") }
	r.svc.Activity()
}

