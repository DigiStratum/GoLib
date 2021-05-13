package main

import (
	"fmt"
	"os"

	db "github.com/DigiStratum/GoLib/DB"
	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

func main() {
	connection_example()
	manager_example()
}

func manager_example() {
	// Get a database connection
	dsn := db.MakeDSN("username", "password", "localhost", "3306", "todolist")
	fmt.Printf("MySQL DSN is: %s\n", dsn)

	// Get the connection through a connection Manager if you want to manage multiple connections to different DB's
	manager := mysql.NewManager()
	dbKey, err := manager.Connect(dsn)
	if nil != err { die(fmt.Sprintf("Error connecting: %s\n", err.Error())) }

	// Run the query
	results, err := manager.Query(
		*dbKey,
		"SELECT id, task, due FROM todo;",
		Todo{},
	)
	if nil != err { die(fmt.Sprintf("Query Error: %s\n", err.Error())) }

	// Process the results
	for index := 0; index < results.Len(); index++ {
		result := results.Get(index)
		if todoResult, ok := result.(*Todo); ok {
			todoResult.Print()
		} else {
			fmt.Printf("Error converting result record to todo{%d}\n", index)
		}
	}

	manager.Disconnect(*dbKey)
}

func connection_example() {
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
	//for index, result := range *results {
	for index := 0; index < results.Len(); index++ {
		result := results.Get(index)
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

