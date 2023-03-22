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
	gojson "encoding/json"

	cfg "github.com/DigiStratum/GoLib/Config"
	dep "github.com/DigiStratum/GoLib/Dependencies"
	db "github.com/DigiStratum/GoLib/DB"
	mysql "github.com/DigiStratum/GoLib/DB/MySQL"
)

func main() {
	// Load configuration
	config := cfg.NewConfig()
	err := config.LoadFromJsonFile("example.config.json")
	if nil != err { dief("Error loading config JSON: %s", err) }

	dsn, err := getDSNFromConfig(config.GetSubsetConfig("db.dsn."))
	if nil != err { dief("DSN Build error: %s", err) }

	example_Connection(*dsn)
	example_ConnectionFactory(*dsn)
	example_ConnectionPool(*dsn)
	example_ResourceObject(*dsn)
}

// Examples
// -----------------------------------------------

// For long-running processes, use ConnectionPool when you want a pool of persistent connections with
// all the conveniences.
func example_ConnectionPool(dsn db.DSN) {
	fmt.Println("ConnectionPool Example")

	// Get the connection from a MySQL connection pool
	connFactory := mysql.NewMySQLConnectionFactory()
	connPool := mysql.NewConnectionPool(dsn)
	defer connPool.Close()
	err := connPool.InjectDependencies(dep.NewDependencyInstance("ConnectionFactory", connFactory))
	if nil != err { dief("Error injecting dependencies: %s\n", err) }
	if err = connPool.Start(); nil != err { dief("Error starting ConnectionPool(): %s", err.Error()) }

	// Lease a connection from the pool
	conn, err := connPool.GetConnection()
	if nil != err { dief("Error getting leased connection: %s\n", err) }
	defer conn.Release()

	// Run a query through
	query, err := conn.NewQuery(mysql.NewSQLQuery("SELECT id, task, due FROM todo;"))
	if (nil != err ) || (nil == query) { dief("Query Setup Error: %s\n", err) }

	// Dump results
	fmt.Printf("%s\n\n", runQueryReturnJson(query))
}

// Use ConnectionFactory to get Connection which can be replaced with a mock for unit test coverage
func example_ConnectionFactory(dsn db.DSN) {
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

	// Dump results
	fmt.Printf("%s\n\n", runQueryReturnJson(query))
}

// Use a Connection to wrap the sql/driver primitives with intrinsic handling for transactions and
// prepared statements.
func example_Connection(dsn db.DSN) {
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

	// Dump results
	fmt.Printf("%s\n\n", runQueryReturnJson(query))
}

type ResourceTask struct {
	Id			int	`json:"id"`
	Task			string	`json:"task"`
	Due			string	`json:"due"`
}

type resourceTask struct {
	resource		ResourceTask
	connPool		mysql.ConnectionPoolIfc
}

func (r resourceTask) GetAllTasks() []resourceTask {
	conn, err := r.connPool.GetConnection()
	if nil != err { dief(fmt.Sprintf("GetConnection Error: %s\n", err.Error())) }
	defer conn.Release()
	query, err := conn.NewQuery(mysql.NewSQLQuery("SELECT id, task, due FROM todo;"))
	results, err := query.RunReturnAll() // No args for this example
	if nil != err { dief(fmt.Sprintf("Query Error: %s\n", err.Error())) }

	// Collect results into []resourceTask
	resourceTasks := make([]resourceTask, 0)
	it := results.GetIterator()
	for rr := it(); nil != rr; rr = it() {
		if resultRow, ok := rr.(mysql.ResultRowIfc); ok {
			rt := resourceTask{
				resource:	ResourceTask{
					Id: int(resultRow.Get("id").GetInt64Default(0)),
					Task: resultRow.Get("task").GetStringDefault("null"),
					Due: resultRow.Get("due").GetStringDefault("1970-01-01 00:00:00"),
				},
				connPool: r.connPool,
			}
			resourceTasks = append(resourceTasks, rt)
		}
	}
	return resourceTasks
}

func example_ResourceObject(dsn db.DSN) {
	fmt.Println("Resource Object Example")

	// Get the connection from a MySQL connection pool
	connFactory := mysql.NewMySQLConnectionFactory()
	connPool := mysql.NewConnectionPool(dsn)
	defer connPool.Close()
	err := connPool.InjectDependencies(dep.NewDependencyInstance("ConnectionFactory", connFactory))
	if nil != err { dief("Error injecting dependencies: %s\n", err) }
	if err = connPool.Start(); nil != err { dief("Error starting ConnectionPool(): %s", err.Error()) }

	// Use the Resource Object to access data, get all the records from the connection pool
	rt := resourceTask{
		resource:	ResourceTask{},
		connPool:	connPool,
	}
	resourceTasks := rt.GetAllTasks()

	// Convert to JSON
	jsonBytes, err := gojson.Marshal(resourceTasks)
	if nil != err { dief("Error Marshaling JSON: %s\n", err) }
        jsonString := string(jsonBytes[:])

	// Dump results
	fmt.Printf("%s\n\n", jsonString)
}

// Supporting Functions
// -----------------------------------------------

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

func runQueryReturnJson(query mysql.QueryIfc) string {
	results, err := query.RunReturnAll() // No args for this example
	if nil != err { dief(fmt.Sprintf("Query Error: %s\n", err.Error())) }

	json, err := results.ToJson()
	if nil != err { dief(fmt.Sprintf("JSON Marshaler Error: %s\n", err.Error())) }

	return *json
}

func dief(msg string, args ...interface{}) {
	fmt.Printf(msg, args...)
	os.Exit(1)
}
