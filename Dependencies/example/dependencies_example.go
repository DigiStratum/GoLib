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
	v := "default"
	for _, v = range variant {}
	return client.InjectDependencies(
		dep.NewDependencyInstance("usefulService", v, r),
	)
}

type dependentLayer struct {
	*dep.DependencyInjected
}

func (r *dependentLayer) DoWork() {
	uids := r.DependencyInjected.GetUniqueIds()
	if 0 == len(*uids) { fmt.Println("Seems to be empty...") }
	d := r.DependencyInjected.GetInstance((*uids)[0])
	if nil == d { fmt.Printf("Seems to be nil... [%s]\n", (*uids)[0]) }
	if us, ok := d.(*usefulService); ok {
		us.UsefulActivity()
	} else {
		fmt.Println("Doesn't seem to be a usefulService...")
	}
}

func main() {
	dl := &dependentLayer{
		DependencyInjected: dep.NewDependencyInjected(
			dep.NewDependencies(
				dep.NewDependency("usefulService", "default", true),
			),
		),
	}

	us := &usefulService{}
	us.InjectInto(dl)

	dl.DoWork();
}

