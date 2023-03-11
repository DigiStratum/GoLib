package main

import (
	"fmt"
)

func main() {
	app := NewApp()
	// We expect an error if we ask it to do something without Start()ing App first
	if err := app.DoSomething(); nil != err { fmt.Printf("ERROR: %s\n\n", err.Error()) }
	// Start() the App properly...
	if err := app.Start(); nil != err { fmt.Printf("ERROR: %s\n\n", err.Error()) }
	// Now we should see no errors...
	if err := app.DoSomething(); nil != err { fmt.Printf("ERROR: %s\n\n", err.Error()) }
	fmt.Println("Seems ok!")
}

