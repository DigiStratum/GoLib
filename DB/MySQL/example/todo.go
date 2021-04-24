package main

import (
	"fmt"
	"encoding/json"

	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

type Todo struct {
	Id	int	`json:"id"`
	Task	string	`json:"task"`
	Due	string	`json:"due"`
}

// Satisfies ResultIfc
func (t Todo) ZeroClone() (mysql.ResultIfc, mysql.PropertyPointers) {
	n := Todo{}
	npp := mysql.PropertyPointers{ &n.Id, &n.Task, &n.Due }
	return &n, npp
}


func (t Todo) Print() {
	jbytes, _ := json.Marshal(t)
	fmt.Println(string(jbytes))
}

