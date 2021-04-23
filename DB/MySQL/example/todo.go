package main

import (
	"fmt"

	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

type Todo struct {
	//Result
	Id	int
	Task	string
	Due	string
}

// Satisfies ResultIfc
func (t Todo) ZeroClone() (mysql.ResultIfc, mysql.PropertyPointers) {
	n := Todo{}
	npp := mysql.PropertyPointers{ &n.Id, &n.Task, &n.Due }
	return &n, npp
}


func (t Todo) Print() {
	fmt.Printf(
		"Todo: { \"id\": \"%d\", \"task\": \"%s\", \"due\": \"%s\" }\n",
		t.Id, t.Task, t.Due,
	)
}

