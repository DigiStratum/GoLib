package main

import (
	"fmt"

	dep "github.com/DigiStratum/GoLib/Dependencies"
)

func main() {
	app := NewApp()
	app.InjectDependencies(dep.NewDependencyInstance("Service", NewService()))
	if err := app.Start(); nil != err {
		fmt.Printf("Error starting app: %s\n\n", err.Error())
	} else {
		app.Run()
	}
}

