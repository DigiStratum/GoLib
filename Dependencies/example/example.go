package main

import (
	dep "github.com/DigiStratum/GoLib/Dependencies"
)

func main() {
	app := NewApp()
	app.InjectDependencies(dep.NewDependencyInstance("Service", NewService()))
	app.Run();
}

