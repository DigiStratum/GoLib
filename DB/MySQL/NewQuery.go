package mysql

import (
	"encoding/json"
	
	nullables "github.com/DigiStratum/GoLib/DB/MySQL/nullables"
)

type QueryIfc interface {
	Run(args ...interface{}) error
	RunReturnInt(args ...interface{}) (*int, error)
	RunReturnString(args ...interface{}) (*string, error)
	RunReturnOne(args ...interface{}) (ResultIfc, error)
	RunReturnSet(args ...interface{}) (ResultSetIfc, error)
}

type query struct {
	connection	ConnectionIfc
	query		string
}

func NewQuery(connection ConnectionIfc, query string) QueryIfc {
	q := newQuery{
		connection:	connection,
		query:		query,
	}
	return &q
}

// -------------------------------------------------------------------------------------------------
// NewQueryIfc Public Interface
// -------------------------------------------------------------------------------------------------
func (q *query) Run(args ...interface{}) error {
	// TODO: Implement
	return nil
}

func (q *query) RunReturnInt(args ...interface{}) (*int, error) {
	// TODO: Implement
	return nil, nil
}

func (q *query) RunReturnString(args ...interface{}) (*string, error) {
	// TODO: Implement
	return nil, nil
}

func (q *query) RunReturnOne(args ...interface{}) (ResultIfc, error) {
	// TODO: Implement
	return nil, nil
}

// ref: https://kylewbanks.com/blog/query-result-to-map-in-golang
func (q *query) RunReturnSet(args ...interface{}) (ResultSetIfc, error) {
	results := NewResultSet()
	rows, err := (*q).connection.GetConnection().Query((*q).query)
	if nil != err { return nil, err }
	cols, _ := rows.Columns()
	for rows.Next() {
		columnPointers := make([]interface{}, len(cols))
		columns := make([]string, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err = rows.Scan(columnPointers...); err != nil { return nil, err }

		// Create our map, and retrieve the value for each column from the pointers,
		// slice, storing it in the map with the name of the column as the key.
		result := NewResultRow()
		for i, colName := range cols {
			val := columns[i]
			result.Set(colName, nullables.NewNullable(val))
		}
		results.add(result)
	}
	return results, nil
}

// Placeholder to support resolving magic expander tags, etc within our query
func (q *query) resolveQuery(args ... interface{}) string {
	protoQuery := (*q).query
	// TODO: expand query '???' placeholders
	finalQuery := protoQuery
	return finalQuery
}
