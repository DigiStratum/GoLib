package main

import (
	"fmt"
	"os"

	db "github.com/DigiStratum/GoLib/DB"
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

func main() {
	// Get a database connection
	dsn := db.MakeDSN("username", "password", "localhost", "3306", "todolist")
	fmt.Printf("MySQL DSN is: %s\n", dsn)

	dbConn, err := mysql.NewConnection(dsn)
	if nil != err { die(fmt.Sprintf("Error getting connection: %s\n", err.Error())) }

	// Run the query
	query := mysql.NewQuery(
		"SELECT id, task, due FROM todo;",
		Todo{},
	)
	results, err := query.Run(dbConn)
	if nil != err { die(fmt.Sprintf("Query Error: %s\n", err.Error())) }

	// Process the results
	for index, result := range *results {
		if todoResult, ok := result.(*Todo); ok {
			todoResult.Print()
		} else {
			fmt.Printf("Error converting result record to Todo{%d}\n", index)
		}
	}

	dbConn.Disconnect()
}

func die(msg string) {
	fmt.Printf("%s\n", msg)
	os.Exit(1)
}

