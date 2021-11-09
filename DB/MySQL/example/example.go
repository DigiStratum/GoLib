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
	"github.com/DigiStratum/GoLib/Dependencies"
	db "github.com/DigiStratum/GoLib/DB"
	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

func main() {
	// Load configuration
	config := cfg.NewConfig()
	err := config.LoadFromJsonFile("example.config.json")
	if nil != err { die(fmt.Sprintf("Error loading config JSON: %s", err.Error())) }
	dsn, err := getDSNFromConfig(config.GetSubsetConfig("db.dsn."))
	if nil != err {
		die(fmt.Sprintf("DSN Build error: %s", err.Error()))
	}

	connection_example(*dsn)
	connectionFactory_example(*dsn)
	connectionPool_example(*dsn)
}

func getDSNFromConfig(config cfg.ConfigIfc) (*db.DSN, error) {
	requiredConfigKeys := []string{ "User", "Passwd", "Net", "DBName" }
	keys := config.GetKeys()
	if ! config.HasAll(&requiredConfigKeys) {
		for _, key := range keys { fmt.Printf("config key: %s\n", key) }
		return nil, fmt.Errorf("Missing one or more required configuration keys")
	}
	dsnBuilder := db.BuildDSN()
	dsnBuilder.Configure(config)
	dsn, err := dsnBuilder.Build()
	if nil != err {
		return nil, fmt.Errorf("DSN Build error: %s", err.Error())
	}
	fmt.Printf("MySQL DSN is: %s\n\n", dsn.ToString())
	return dsn, nil
}

func connectionPool_example(dsn db.DSN) {
	// TODO: Make a db.NewConnectionPool(), give it a connectionFactory, then Get and Release a leased connection!
	fmt.Println("ConnectionPool Example")

	// Get the connection from a MySQL connection pool
	connFactory := mysql.NewMySQLConnectionFactory()
	deps := dependencies.NewDependencies()
	deps.Set("connectionPool", connFactory)
	connPool := mysql.NewConnectionPool(dsn)
	defer connPool.Close()
	connPool.InjectDependencies(deps)

	// Lease a connection from the pool
	conn, err := connPool.GetConnection()
	if nil != err { die(fmt.Sprintf("Error getting leased connection: %s\n", err.Error())) }
	defer conn.Release()

	// Run a query through
	query, err := conn.NewQuery("SELECT id, task, due FROM todo;")
	if (nil != err ) || (nil == query) { die(fmt.Sprintf("Query Setup Error: %s\n", err)) }
	runQueryDumpAll(query)
}

func connectionFactory_example(dsn db.DSN) {
	fmt.Println("ConnectionFactory Example")

	// Get the connection from a MySQL connection factory
	connFactory := mysql.NewMySQLConnectionFactory()
	dbconn, err := connFactory.NewConnection(dsn)
	if nil != err { die(fmt.Sprintf("Error getting underlying connection: %s\n", err.Error())) }

	// Wrap the raw connection
	conn, err := mysql.NewConnection(dbconn)
	if nil != err { die(fmt.Sprintf("Error getting connection wrapper: %s\n", err.Error())) }
	defer conn.Close()

	// Run a query through
	query, err := conn.NewQuery("SELECT id, task, due FROM todo;")
	if (nil != err ) || (nil == query) { die(fmt.Sprintf("Query Setup Error: %s\n", err)) }
	runQueryDumpAll(query)
}

// Get the connection directly
func connection_example(dsn db.DSN) {
	fmt.Println("Connection Example")

	// Get the connection directly from SQL driver
	dbconn, err := sql.Open("mysql", dsn.ToString())
	if nil != err { die(fmt.Sprintf("Error getting underlying connection: %s\n", err)) }

	// Wrap the raw connection
	conn, err := mysql.NewConnection(dbconn)
	if nil != err { die(fmt.Sprintf("Error getting connection wrapper: %s\n", err)) }
	defer conn.Close()

	// Run a query through
	query, err := conn.NewQuery("SELECT id, task, due FROM todo;")
	if (nil != err ) || (nil == query) { die(fmt.Sprintf("Query Setup Error: %s\n", err)) }
	runQueryDumpAll(query)
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
