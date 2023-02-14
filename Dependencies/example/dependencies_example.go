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

func (r *usefulService) InjectInto(client dep.DependencyInjectableIfc, variant ...string) error {
	v := dep.DEP_VARIANT_DEFAULT
	for _, v = range variant {}
	return client.InjectDependencies(
		dep.NewDependencyInstance("usefulService", r).SetVariant(v),
	)
}

type dependentLayer struct {
	*dep.DependencyInjected
}

func (r *dependentLayer) DoWork() {
	var uids *[]string = r.DependencyInjected.GetUniqueIds()
	//if (nil == uids) || (0 == len(*uids)) { fmt.Println("Seems to be empty...") }
	var uid  string = (*uids)[0]
	//fmt.Printf("Found UID: '%s'\n", uid)
	d := r.DependencyInjected.GetInstance(uid)
	if us, ok := d.(*usefulService); ok {
		us.UsefulActivity()
	} else {
		fmt.Println("Doesn't seem to be a usefulService...")
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

	us := &usefulService{}
	us.InjectInto(dl)

	dl.DoWork();
}

