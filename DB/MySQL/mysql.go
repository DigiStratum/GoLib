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
	"reflect"

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
	Template	interface{}	// Structure template that each row result is expected to match; makes FieldNum obsolete
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

// Execute a query with varargs for substitution and Structured results
// ref: https://appliedgo.net/generics/
// ref: https://stackoverflow.com/questions/37851500/how-to-copy-an-interface-value-in-go
// ref: https://forum.golangbridge.org/t/database-rows-scan-unknown-number-of-columns-json/7378/2
// ref: https://stackoverflow.com/questions/26744873/converting-map-to-struct/26746461
// ref: https://stackoverflow.com/questions/29184933/golang-reflect-get-pointer-to-a-struct-field-value
func (dbc *DBConnection) SQuery(querySpec *QuerySpec, args ...string) (*[]interface{}, error) {
	results := []interface{}{}

	// Ref: https://stackoverflow.com/questions/18926303/iterate-through-the-fields-of-a-struct-in-go
	numFields := reflect.TypeOf(querySpec.Template).NumField()
	fmt.Printf("SQuery() QuerySpec.Template has %d Fields\n", numFields)
	values := make([]interface{}, numFields)
	templateValue := reflect.ValueOf(querySpec.Template)
	tvType := templateValue.Type()
	//templateValue := reflect.ValueOf(querySpec.Template).Elem()
	for i := 0; i < numFields; i++ {
		//fi := templateValue.Field(i).Interface()
		//values[i] = &fi
		// ref: https://stackoverflow.com/questions/27992821/how-get-pointer-of-structs-member-from-interface
		//values[i] = templateValue.Field(i).Addr().Interface()
		//values[i] = templateValue.Elem().FieldByIndex(i).Addr().Interface()
		//values[i] = templateValue.Field(i).Addr()
		//values[i] = templateValue.Field(i).Interface().Addr()
		//values[i] = templateValue.Field(i).Interface()

		// ref: https://stackoverflow.com/questions/29184933/golang-reflect-get-pointer-to-a-struct-field-value
		//valueField := templateValue.Field(i)
		//values[i] = valueField.Addr().Interface()

		// ref: https://samwize.com/2015/03/20/how-to-use-reflect-to-set-a-struct-field/
		fieldName := tvType.Field(i).Name
		field:= templateValue.Field(i)
		fmt.Printf("Field name: '%s', type: '%s'\n", fieldName, field.Type())
		//v := templateValue.FieldByName(fieldName).Interface()
		v := field.Interface()
		switch field.Type().String() {
			case "*int":
				values[i] = v.(*int)
			case "*string":
				values[i] = v.(*string)
		}
		//values[i] = templateValue.Elem().FieldByName(fieldName).Addr().Interface()
	}

	protoQuery := querySpec.Query
	// TODO: expand querySpec.Query '???' placeholders
	finalQuery := protoQuery

	// Convert string args to interface{} for Query()
	iArgs := make([]interface{}, len(args))
	for i, v := range args { iArgs[i] = v }

	// Execute the Query
	rows, err := dbc.Conn.Query(finalQuery, iArgs...)
	if err != nil { return nil, err }

	// Process the result rows
	for rows.Next() {

		// figure out what columns were returned
		// the column names will be the JSON object field keys
		//columns, err := rows.ColumnTypes()
		//if err != nil { return nil, err }

		// Scan needs an array of pointers to the values it is setting
		// This creates the object and sets the values correctly
		//values := make([]interface{}, len(columns))
		//object := map[string]interface{}{}
		//for i, column := range columns {
		//	object[column.Name()] = reflect.New(column.ScanType()).Interface()
		//	values[i] = object[column.Name()]
		//}
		err = rows.Scan(values...)
		if err != nil { return nil, err }
		results = append(results, querySpec.Template)
	}

	return &results, nil
}

// Execute a query with varargs for substitution and Mapped results
// Iterate over the results for this query and send all the QueryResult rows to a channel
// TODO: How can we return an error in place of the channel if something goes wrong?
// ref: https://ewencp.org/blog/golang-iterators/index.html
// ref: https://blog.golang.org/pipelines
// ref: https://programming.guide/go/wait-for-goroutines-waitgroup.html
// ref: https://golang.org/pkg/database/sql/#example_DB_Query_multipleResultSets
// ref: https://kylewbanks.com/blog/query-result-to-map-in-golang
func (dbc *DBConnection) MQuery(querySpec QuerySpec, args ...string) <-chan QueryResult {
	protoQuery := querySpec.Query
	// TODO: expand querySpec.Query '???' placeholders
	finalQuery := protoQuery
	// Make a result channel
	ch := make(chan QueryResult)
	defer close(ch)

	// Convert string args to interface{} for Query()
	iArgs := make([]interface{}, len(args))
	for i, v := range args { iArgs[i] = v }

	rows, err := dbc.Conn.Query(finalQuery, iArgs...)
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
					for i, colName := range resultCols {
						qrr[colName] = fmt.Sprintf("%v", columnPointers[i].(*interface{}))
					}
					ch <- QueryResult{ Row: &qrr, Err: nil }
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
func (dbm *DBManager) Connect(dsn string) (*DBKey, error) {

	// If we already have this dbKey...
	dbKey := DBKey{ Key: db.GetDSNHash(dsn) }
	if _, ok := dbm.connections[dbKey.Key]; ! ok {
		// Not connected yet - let's do this thing!
		conn, err := sql.Open("mysql", dsn)
		if err != nil { return nil, err }

		// Make a new connection record
		dbm.connections[dbKey.Key] = DBConnection{
			DSN:		dsn,
			IsConnected:	true,
			Conn:		conn,
		}
	}
	return &dbKey, nil
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
	if ! dbm.IsConnected(dbKey) {
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

func (dbm *DBManager) Query(dbKey DBKey, querySpec QuerySpec, args ...string) <-chan QueryResult {
	conn, err := dbm.GetConnection(dbKey)
	if (nil != err) {
		qrErr := lib.GetLogger().Error(fmt.Sprintf("GetConnection(): '%s' - Error: '%s'", dbKey.Key, err.Error()))
		// Make a result channel
		ch := make(chan QueryResult)
		defer close(ch)
		ch <- QueryResult{ Row:	nil, Err: qrErr }
		return ch
	}
	return conn.MQuery(querySpec, args...)
}

