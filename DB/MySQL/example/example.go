package main

/*

This example demonstrates some, of many possible, functional use cases for our MySQL database package. The main()
launches multiple functions, each of which is an example of database interaction for the documented purpose. While
these examples are intended to demonstrate the steps involved with the various flows of interaction, they are written
as simple demonstrations, not necessarily as one would use directly in a fully fledged application. As such, there is
much to be desired for error handling, structure, data processing and preparation, etc.

The todolist.sql script includes the minimal statements necessary to set up the test database. Getting MySQL server
installed, configured, running, and logged into get to this point is beyond the scope of this example documentation.

mysql -u root -p < todolist.sql

Once you have the database set up, update the example configuration JSON with the appropriate connection details, then
run the example here as:

go run example.go

*/

import (
	"fmt"
	"os"
	"database/sql"

	cfg "github.com/DigiStratum/GoLib/Config"
	db "github.com/DigiStratum/GoLib/DB"
	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

func main() {
	// Load configuration
	cfg := cfg.NewConfig()
	err := cfg.LoadFromJsonFile("example.config.json")
	if nil != err {
		die("Error loading config from JSON file")
	}
	requiredConfigKeys := []string{ "User", "Passwd", "Net", "DBName" }
	if ! cfg.HasAll(&requiredConfigKeys) {
		die("Missing one or more required configuration keys")
	}
	dsnBuilder := db.BuildDSN()
	dsnBuilder.Configure(cfg)
	dsn, err := dsnBuilder.Build()
	if nil != err {
		die(fmt.Sprintf("DSN Build error: %s", err.Error()))
	}

/*
	dsn := db.MakeDSN(
		*(cfg.Get("user")),
		*(cfg.Get("pass")),
		*(cfg.Get("host")),
		*(cfg.Get("port")),
		*(cfg.Get("name")),
	)
 */
	fmt.Printf("MySQL DSN is: %s\n\n", dsn.ToString())

	//connectionFactory_example(*dsn)
	connection_example(*dsn)
//	manager_example(dsn)
}

/*
// Get the connection through a connection Manager if you want to manage multiple connections/pools to different DB's
func manager_example(dsn string) {

	fmt.Println("Manager Example")
	manager := mysql.NewManager()
	dbKey := manager.NewConnectionPool(dsn)

	// Get leased connection from pool
	conn := manager.GetConnection(dbKey)
	if nil == conn { die("Error connecting\n") }

	query, err := conn.NewQuery("SELECT id, task, due FROM todo;")

	if nil != err { die(fmt.Sprintf("Error Creating Query: %s\n", err.Error())) }

	runQueryDumpAll(query)

	err = conn.Release()
	if nil != err { die(fmt.Sprintf("Error Releasing Connection: %s\n", err.Error())) }

	manager.CloseConnectionPool(dbKey)
}
 */

func connectionFactory_example(dsn db.DSN) {
}

// Get the connection directly
func connection_example(dsn db.DSN) {
	fmt.Println("Direct Example")

	dbconn, err := sql.Open("mysql", dsn.ToString())
	if nil != err { die(fmt.Sprintf("Error getting underlying connection: %s\n", err.Error())) }

	conn, err := mysql.NewConnection(dbconn)
	if nil != err { die(fmt.Sprintf("Error getting connection wrapper: %s\n", err.Error())) }

	query, err := conn.NewQuery("SELECT id, task, due FROM todo;")
	if (nil != err ) || (nil == query) { die(fmt.Sprintf("Query Setup Error: %s\n", err)) }
	
	runQueryDumpAll(query)

	conn.Close()
}

func runQueryDumpAll(query mysql.QueryIfc) {
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
