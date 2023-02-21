package main

import (
	"fmt"
)

type ServiceIfc interface {
	Activity()
}

type service struct { }

func NewService() *service {
	return &service{}
}

func (r *service) Activity() {
	fmt.Println("Activity output!")
}

