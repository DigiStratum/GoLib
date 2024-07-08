package cli

// Package of convenience functions for Command Line Interface (CLI) programming

import (
	"os"
	"fmt"
)

type CLIIfc interface {
	Die(msg string)
}

type CLI struct {
}

func NewCLI() *CLI {
	r := CLI{}
	return &r
}

// Die with a message - DOES NOT RETURN
func (r *CLI) Die(msg string) {
        fmt.Println(msg)
        os.Exit(1)
}

