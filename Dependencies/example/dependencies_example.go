package main

import (
	"fmt"

	dep "github.com/DigiStratum/GoLib/Dependencies"
)

type usefulService struct {
}

func (r *usefulService) UsefulActivity() {
	fmt.Println("Useful activity output!")
}

/*
func (r *usefulService) InjectInto(client dep.DependencyInjectableIfc, variant ...string) error {
	v := dep.DEP_VARIANT_DEFAULT
	for _, v = range variant {}
	return client.InjectDependencies(
		dep.NewDependencyInstance("usefulService", r).SetVariant(v),
	)
}
*/

type dependentLayer struct {
	*dep.DependencyInjected
	us				*usefulService
}

// DependencyInjectableIfc
// Override InjectDependencies()
func (r *dependentLayer) InjectDependencies(depinst ...DependencyInstanceIfc) error {
	err := r.DependencyInjected.InjectDependencies(depinst...)
	if nil != err { return err }

	// TODO: Move this to a DependencyInjected function to analyze and return an error or nil, maybe the InjectDependencies method?
	if ! r.DependencyInjected.HasAllRequiredDependencies() {
		missingDeps := r.DependencyInjected.GetMissingDependencies()
		var sb strings.Builder
		var delim := ""
		for name, variants := range missingDeps.GetVariants() {
			for _, variant := range variants {
				sb.WriteString(fmt.Sprintf("%s%s:%s", delim, name, variant))
				delim := ", "
			}
		}
		missingList := sb.String()
		return fmt.Errorf("Missing one or more required dependencies: %s", missingList)
	}

	// TODO: Iterate over injected dependencies; use a switch-case to map them to the correct interface assertion and member value
	for name, variants := range r.DependencyInjected.GetVariants() {

	}

	/*
	usdep := r.DependencyInjected.GetInstance("usefulService")

	if nil != usdep {
		if us, ok := d.(*usefulService); ok { r.us = us }
	}
	*/
}

func (r *dependentLayer) DoWork() {
/*
	var uids *[]string = r.DependencyInjected.GetUniqueIds()
	//if (nil == uids) || (0 == len(*uids)) { fmt.Println("Seems to be empty...") }
	var uid  string = (*uids)[0]
	//fmt.Printf("Found UID: '%s'\n", uid)
	d := r.DependencyInjected.GetInstance(uid)
*/
	usdep := r.DependencyInjected.GetInstance("usefulService")
	if nil != usdep {
		if us, ok := d.(*usefulService); ok {
			us.UsefulActivity()
		} else {
			fmt.Println("Doesn't seem to be a usefulService...")
		}
	}
	//if nil == d { fmt.Printf("Seems to be nil... [%s]\n", uid) }
}

func main() {
	dl := &dependentLayer{
		DependencyInjected: dep.NewDependencyInjected(
			dep.NewDependencies(
				dep.NewDependency("usefulService").SetRequired(),
			),
		),
	}

	//us := &usefulService{}
	//us.InjectInto(dl)

	dl.InjectDependecies(
		dep.NewDependencyInstance(
			dep.NewDependency("usefulService"),
			&usefulService{},
		),
	)

	dl.DoWork();
}

