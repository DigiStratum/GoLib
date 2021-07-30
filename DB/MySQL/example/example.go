package main

import (
	"fmt"
	"os"

	db "github.com/DigiStratum/GoLib/DB"
	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

func main() {
	// Get a database connection
	dsn := db.MakeDSN("username", "password", "localhost", "3306", "todolist")
	fmt.Printf("MySQL DSN is: %s\n", dsn)

	connection_example(dsn)
	manager_example(dsn)
}

// Get the connection through a connection Manager if you want to manage multiple connections/pools to different DB's
func manager_example(dsn string) {

	manager := mysql.NewManager()
	dbKey := manager.NewConnectionPool(dsn)

	// Get leased connection from pool
	conn := manager.GetConnection(dbKey)
	if nil == conn { die("Error connecting\n") }

	query, err := conn.NewQuery("SELECT id, task, due FROM todo;")

	if nil != err { die(fmt.Sprintf("Error Creating Query: %s\n", err.Error())) }

	runQuery(query)

	manager.DestroyConnectionPool(dbKey)
}

// Get the connection directly
func connection_example(dsn string) {
	// Get direct connection
	conn, err := mysql.NewConnection(dsn)
	if nil != err { die(fmt.Sprintf("Error getting connection: %s\n", err.Error())) }

	query, err := conn.NewQuery("SELECT id, task, due FROM todo;")
	if (nil != err ) || (nil == query) { die(fmt.Sprintf("Query Setup Error: %s\n", err)) }
	
	runQuery(query)

	conn.Disconnect()
}

func runQuery(query mysql.QueryIfc) {
	results, err := query.RunReturnAll() // No args for this example
	if nil != err { die(fmt.Sprintf("Query Error: %s\n", err.Error())) }

	// Process the results
	for index := 0; index < results.Len(); index++ {
		result := results.Get(index)
		resultJson, err := result.ToJson()
		if err != nil {
			fmt.Printf("Error converting result record to JSON: %s\n", err)
		}
		fmt.Printf("Result: %s\n\n", *resultJson)
	}
}

func die(msg string) {
	fmt.Printf("%s\n", msg)
	os.Exit(1)
}
