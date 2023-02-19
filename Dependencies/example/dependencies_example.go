package main

import (
	"fmt"
	"strings"

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
	err := r.DependencyInjected.InjectDependencies(depinst...)
	if nil != err { return err }

	// TODO: Move this to a DependencyInjected function to analyze and return an error or nil, maybe the InjectDependencies method?
	if ! r.DependencyInjected.HasAllRequiredDependencies() {
		missingDeps := r.DependencyInjected.GetMissingDependencies()
		var sb strings.Builder
		delim := ""
		for name, variants := range missingDeps.GetVariants() {
			for _, variant := range variants {
				sb.WriteString(fmt.Sprintf("%s%s:%s", delim, name, variant))
				delim = ", "
			}
		}
		missingList := sb.String()
		return fmt.Errorf("Missing one or more required dependencies: %s", missingList)
	}

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
	dl := &dependentLayer{
		DependencyInjected: dep.NewDependencyInjected(
			dep.NewDependencies(
				dep.NewDependency("usefulService").SetRequired(),
			),
		),
	}

	dl.InjectDependencies(
		dep.NewDependencyInstance("usefulService", &usefulService{}),
	)

	dl.DoWork();
}

