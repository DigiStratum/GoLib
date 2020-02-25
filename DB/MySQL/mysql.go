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
	Query		string		// The query to execute as prepared statement
	FieldNum	int		// How many fields are we expecting the result row to contain?
	MinKeys		int		// minimum num keys required to populate query; 0 = no min
	MaxKeys		int		// maximum num keys required to populate query; 0 = no max
}

// A query always results in a row of column data where each column has a name and a value as a map
type QueryResultRow map[string]string

// A query result may have an error during processing; this wrapper lets us combine data and error
type QueryResult struct {
	Row		*QueryResultRow
	Err		error
}

// ------------------------------------------------------------------------------------------------
// DB Connection

type DBConnection struct {
	DSN		string		// Full Data Source Name for this connection
	IsConnected	bool		// Is it currently connected as far as we know?
	Conn		*sql.DB		// Read-Write DBConnection
}

// Execute a query with varargs for substitutions
// Iterate over the results for this query and send all the QueryResult rows to a channel
// TODO: How can we return an error in place of the channel if something goes wrong?
// ref: https://ewencp.org/blog/golang-iterators/index.html
// ref: https://blog.golang.org/pipelines
// ref: https://programming.guide/go/wait-for-goroutines-waitgroup.html
// ref: https://golang.org/pkg/database/sql/#example_DB_Query_multipleResultSets
// ref: https://kylewbanks.com/blog/query-result-to-map-in-golang
func (dbc *DBConnection) Query(querySpec QuerySpec, args ...string) <-chan QueryResult {
	protoQuery := querySpec.Query
	// TODO: expand querySpec.Query '???' placeholders
	finalQuery := protoQuery

	var err error = nil

	// Make a result channel
	ch := make(chan QueryResult)
	defer close(ch)

	rows, err := dbc.Conn.Query(finalQuery, args...)
	defer rows.Close()
	if (nil == err) {
		resultCols, err := rows.Columns()
		if (nil == err) {
			for rows.Next() {
				// Create a slice of interface{}'s to represent each column,
				// and a second slice to contain pointers to each item in the columns slice.
				columns := make([]interface{}, len(resultCols))
				columnPointers := make([]interface{}, len(resultCols))

				// Scan the result into the column pointers...
				for i, _ := range columns { columnPointers[i] = &columns[i] }
				if err = rows.Scan(columnPointers...); err == nil {

					// Create our map, and retrieve the value for each column from the pointers
					// slice, storing it in the map with the name of the column as the key.
					qrr := make(QueryResultRow)
					for i, colName := range cols {
						qrr[colName] = fmt.Sprintf("%v", columnPointers[i].(*interface{}))
					}
					ch <- QueryResult{ Row: qrr, Err: nil }
				}
			}
		}
	}

	// Any errors above get dumped into the channel as a single result here
	if (nil != err) {
		qrErr := lib.GetLogger().Error(fmt.Sprintf("Query: '%s' - Error: '%s'", finalQuery, err.Error()))
		ch <- QueryResult{ Row:	nil, Err: qrErr }
	}

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

func (dbm *DBManager) GetConnection(dbKey DBKey) (*DBConnection, error) {
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

