package main

import (
	"fmt"

	dep "github.com/DigiStratum/GoLib/Dependencies"
)

// Sample service to inject
type usefulServiceIfc interface {
	UsefulActivity()
}

type usefulService struct { }
func (r *usefulService) UsefulActivity() { fmt.Println("Useful activity output!") }

// Sample consumer that wants a simple service to be injected
type dependentLayer struct {
	*dep.DependencyInjected
	us				usefulServiceIfc
}

// If we override InjectDependencies() we can capture injected dependencies
// locally instead of asking for them every time they are needed
func (r *dependentLayer) InjectDependencies(depinst ...dep.DependencyInstanceIfc) error {

	// If DI fails, return error
	if err := r.DependencyInjected.InjectDependencies(depinst...); nil != err { return err }
	// If DI missing requirements, return error
	if err := r.DependencyInjected.ValidateRequiredDependencies(); nil != err { return err }

	// Iterate over injected dependencies; use a switch-case to map them to the correct interface assertion and member value
	// TODO: Add some variant iterator with callback function per each?
	for name, _ := range r.DependencyInjected.GetVariants() {
		switch name {
			case "usefulService":
				usdep := r.DependencyInjected.GetInstance(name)
				if us, ok := usdep.(usefulServiceIfc); ok { r.us = us }
		}
	}

	return nil
}

func (r *dependentLayer) DoWork() {
	if nil == r.us { fmt.Println("No usefulService from DI!") }
	r.us.UsefulActivity()
}

func main() {
	// Declare Dependencies
	dl := &dependentLayer{
		DependencyInjected: dep.NewDependencyInjected(
			dep.NewDependencies(
				dep.NewDependency("usefulService").SetRequired(),
			),
		),
	}

	// Inject Dependencies
	dl.InjectDependencies(
		dep.NewDependencyInstance("usefulService", &usefulService{}),
	)

	// Do something useful!
	dl.DoWork();
}

