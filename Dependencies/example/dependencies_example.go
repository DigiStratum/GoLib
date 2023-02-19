package main

import (
	"fmt"

	dep "github.com/DigiStratum/GoLib/Dependencies"
)

type usefulServiceIfc interface {
	UsefulActivity()
}

type usefulService struct {
}

func (r *usefulService) UsefulActivity() {
	fmt.Println("Useful activity output!")
}

type dependentLayer struct {
	*dep.DependencyInjected
	us				usefulServiceIfc
}

// DependencyInjectableIfc
// Override InjectDependencies()
func (r *dependentLayer) InjectDependencies(depinst ...dep.DependencyInstanceIfc) error {
	if err := r.DependencyInjected.InjectDependencies(depinst...); nil != err { return err }
	if err := r.DependencyInjected.ValidateRequiredDependencies(); nil != err { return err }

	// Iterate over injected dependencies; use a switch-case to map them to the correct interface assertion and member value
	// TODO: Add some variant iterator with callback function per each?
	for name, variants := range r.DependencyInjected.GetVariants() {
		switch name {
			case "usefulService":
				// Any variant wins
				for _, variant := range variants {
					usdep := r.DependencyInjected.GetInstanceVariant(name, variant)
					if nil != usdep {
						if us, ok := usdep.(usefulServiceIfc); ok { r.us = us }
					}
				}
		}
	}

	return nil
}

func (r *dependentLayer) DoWork() {
	if nil != r.us {
		r.us.UsefulActivity()
	} else {
		fmt.Println("Doesn't seem to be a usefulService...")
	}
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

