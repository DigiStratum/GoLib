package main

/*

This example demonstrates some, of many possible, functional use cases for our MySQL database
package. The main() launches multiple functions, each of which is an example of database interaction
for the documented purpose. While these examples are intended to demonstrate the steps involved with
the various flows of interaction, they are written as simple demonstrations, not necessarily as one
would use directly in a fully fledged application. As such, there is much to be desired for error
handling, structure, data processing and preparation, etc.

The todolist.sql script includes the minimal statements necessary to set up the test database.
Getting MySQL server installed, configured, running, and logged into get to this point is beyond the
scope of this example documentation.

mysql -u root -p < todolist.sql

Once you have the database set up, update the example configuration JSON with the appropriate
connection details, then run the example here as:

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
	if nil != err { dief("Error loading config JSON: %s", err) }
	dsn, err := getDSNFromConfig(config.GetSubsetConfig("db.dsn."))
	if nil != err {
		dief("DSN Build error: %s", err)
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
		return nil, fmt.Errorf("DSN Build error: %s", err)
	}
	fmt.Printf("MySQL DSN is: %s\n\n", dsn.ToString())
	return dsn, nil
}

/*
For long-running processes, use a ConnectionPool when you want a pool of persistent connections with
all the conveniences.
*/
func connectionPool_example(dsn db.DSN) {
	fmt.Println("ConnectionPool Example")

	// Get the connection from a MySQL connection pool
	connFactory := mysql.NewMySQLConnectionFactory()
	deps := dependencies.NewDependencies()
	deps.Set("connectionFactory", connFactory)
	connPool := mysql.NewConnectionPool(dsn)
	defer connPool.Close()
	err := connPool.InjectDependencies(deps)
	if nil != err { dief("Error injecting dependencies: %s\n", err) }

	// Lease a connection from the pool
	conn, err := connPool.GetConnection()
	if nil != err { dief("Error getting leased connection: %s\n", err) }
	defer conn.Release()

	// Run a query through
	query, err := conn.NewQuery(mysql.NewSQLQuery("SELECT id, task, due FROM todo;"))
	if (nil != err ) || (nil == query) { dief("Query Setup Error: %s\n", err) }
	runQueryDumpAll(query)
}

/*
Use a ConnectionFactory to get a Connection which can be replaced with a mock for unit test coverage
*/
func connectionFactory_example(dsn db.DSN) {
	fmt.Println("ConnectionFactory Example")

	// Get the connection from a MySQL connection factory
	connFactory := mysql.NewMySQLConnectionFactory()
	dbconn, err := connFactory.NewConnection(dsn)
	if nil != err { dief("Error getting underlying connection: %s\n", err) }

	// Wrap the raw connection
	conn, err := mysql.NewConnection(dbconn)
	if nil != err { dief("Error getting connection wrapper: %s\n", err) }
	defer conn.Close()

	// Run a query through
	query, err := conn.NewQuery(mysql.NewSQLQuery("SELECT id, task, due FROM todo;"))
	if (nil != err ) || (nil == query) { dief("Query Setup Error: %s\n", err) }
	runQueryDumpAll(query)
}

/*
Use a Connection to wrap the sql/driver primitives with intrinsic handling for transactions and
prepared statements.
*/
func connection_example(dsn db.DSN) {
	fmt.Println("Connection Example")

	// Get the connection directly from SQL driver
	dbconn, err := sql.Open("mysql", dsn.ToString())
	if nil != err { dief("Error getting underlying connection: %s\n", err) }

	// Wrap the raw connection
	conn, err := mysql.NewConnection(dbconn)
	if nil != err { dief("Error getting connection wrapper: %s\n", err) }
	defer conn.Close()

	// Run a query through
	query, err := conn.NewQuery(mysql.NewSQLQuery("SELECT id, task, due FROM todo;"))
	if (nil != err ) || (nil == query) { dief("Query Setup Error: %s\n", err) }
	runQueryDumpAll(query)
}

func runQueryDumpAll(query mysql.QueryIfc) {
	results, err := query.RunReturnAll() // No args for this example
	if nil != err { dief(fmt.Sprintf("Query Error: %s\n", err.Error())) }

	// Output the results
	fmt.Printf("Result: [\n")
	for index := 0; index < results.Len(); index++ {
		result := results.Get(index)
		resultJson, err := result.ToJson()
		if err != nil {
			dief("Error converting result record to JSON: %s\n", err)
		}
		comma := ","
		if index == (results.Len() - 1) { comma = "" }
		fmt.Printf("\t%s%s\n", *resultJson, comma)
	}
	fmt.Printf("]\n\n")
}

func dief(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}
