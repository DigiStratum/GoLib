package mysql

/*
DB Manager for MySQL - manages connections and provides various reusable DB capabilities.

TODO: Add some sort of query builder - this will allow us to ditch writing SQL for most needs.
TODO: implement connection "pool" that allows for a named association of multiple connections
      to fail-over for retries, round-robin requests, etc.

*/

import (
	"fmt"
	"errors"
	"database/sql"
	//"net/http"

	_ "github.com/go-sql-driver/mysql"

	lib "github.com/DigiStratum/GoLib"
	db "github.com/DigiStratum/GoLib/DB"
)


// ------------------------------------------------------------------------------------------------
// Query Bits

// The spec for a prepared statement query. Single '?' substitution is handled by db.Query()
// automatically. '???' expands to include enough placeholders (as with an IN () list for any count
// of keys > min. max must be >= min unless max == 0.
type QuerySpec struct {
	Query		string	// The query to execute as prepared statement
	FieldNum	int	// How many fields are we expecting the result row to contain?
	MinKeys		int	// minimum num keys required to populate query; 0 = no min
	MaxKeys		int	// maximum num keys required to populate query; 0 = no max
}

type QueryResultRow map[string]string


// ------------------------------------------------------------------------------------------------
// DB Connection

type DBConnection struct {
	DSN		string		// Full Data Source Name for this connection
	IsConnected	bool		// Is it currently connected as far as we know?
	Conn		*sql.DB		// Read-Write DBConnection
}

// Execute a single query with varargs for substitutions
// Iterate over the results for this query and send all the QueryResultRows to a channel
// TODO: How can we return an error in place of the channel if something goes wrong?
// ref: https://ewencp.org/blog/golang-iterators/index.html
// ref: https://blog.golang.org/pipelines
// ref: https://programming.guide/go/wait-for-goroutines-waitgroup.html
// TODO: TBD: Accept a DBKey to determine which connection to execute against, or... attach this
// function to DBConnection, and return connection to consumer by key so that all the connection-
// specific operations (i.e. not related to connection management, therefore a different problem
// domain) get associated with the connection itself?
func (dbc *DBConnection) SQuery(querySpec QuerySpec, args ...string) <-chan QueryResultRow {
	// TODO: Execute the query, check the result row count
	protoQuery := querySpec.Query
	// TODO: expand querySpec.Query '???' placeholders
	finalQuery := protoQuery

	// ref: https://golang.org/pkg/database/sql/#example_DB_Query_multipleResultSets
	rows, err := dbc.Conn.Query(finalQuery, args...)
	if err != nil {
		lib.GetLogger().Error(fmt.Sprintf("Query: '%s' - Error: '%s'", finalQuery, err.Error())
		// FIXME: get out of here, don't just keep running...
	}
	defer rows.Close()

	for rows.Next() {
		// TODO: form a result set that matches the query
		var (
			id   int64
			name string
		)
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		log.Printf("id %d name is %s\n", id, name)
	}

	resultRowCount := 1

	// Make a channel the size of the result row count
	ch := make(chan QueryResultRow, resultRowCount)
	defer close(ch)
	var wg sync.WaitGroup
	wg.Add(1)

	// Fire off a go routine to fill up the channel
	go func() {
		// TODO: Iterate over and convert each Query Result to a QueryResultRow
		for k, v := range *hash {
			ch <- queryResultRow
			ch <- KeyValuePair{ Key: k, Value: v }
		}
		wg.Done()
	}()
	wg.Wait()
	return ch
}

// ------------------------------------------------------------------------------------------------
// DB Key (DBConnection identifier)

type DBKey struct {
	Key	string
}


// ------------------------------------------------------------------------------------------------
// DB Manager

// Set of connections, keyed on DSN
type DBManager struct {
	connections	map[string]DBConnection
}

// Make a new one of these!
func NewDBManager() *DBManager {
	dbm := DBManager{
		connections: make(map[string]DBConnection),
	}
	return &dbm
}

// Get DB DBConnection Key from the supplied DSN
// ref: https://en.wikipedia.org/wiki/Data_source_name
func (dbm *DBManager) Connect(dsn string) (DBKey, error) {
	dbKey := DBKey{
		key: db.GetDSNHash(dsn),
	}

	// If we already have this dbKey...
	if _, ok := dbm.connections[dbKey.Key]; ok {
		// ... then it's because we already have a good connection
		return dbKey, nil;
	}

	// Not connected yet - let's do this thing!
        conn, err := sql.Open("mysql", dsn)
        if err != nil {
		return nil, err
        }

	// Make a new connection record
	dbm.connections[dbKey.Key] = DBConnection{
		DSN:	dsn,
		IsConnected:	true,
		Conn:	conn,
	}
	return dbKey, nil
}

func (dbm *DBManager) IsConnected(dbKey DBKey) bool {
	if conn, ok := dbm.connections[dbKey.Key]; ok {
		return conn.IsConnected
	}
	return false
}

func (dbm *DBManager) GetConnection(dbKey DBKey) *DBConnection, error {
	if conn, ok := dbm.connections[dbKey.Key]; ok {
		return &conn, nil
	}
	return nil, errors.New(fmt.Sprintf("The connection for '%s' is undefined", dbKey.Key))
}

func (dbm *DBManager) Disconnect(dbKey DBKey) error {
	if ! dbm.IsDBConnectionOpen(dbKey.Key) {
		return errors.New(fmt.Sprintf("The connection for '%s' is not open", dbKey.Key))
	}
	if conn, ok := dbm.connections[dbKey.Key]; ok {
		conn.Conn.Close()
		c := conn
		c.IsConnected = false
		dbm.connections[dbKey.Key] = c
	}
	return nil
}

func (dbm *DBManager) SQuery(dbKey DBKey, querySpec QuerySpec, args ...string) <-chan QueryResultRow {
	conn, err := dbm.GetConnection(dbKey)
	return conn.SQuery(querySpec, args...)
}

