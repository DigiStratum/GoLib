package mysql

/*
TODO: Add some sort of query builder - this will allow us to ditch writing SQL for most needs.
*/

import (
	_ "github.com/go-sql-driver/mysql"
)

// The spec for a prepared statement query. Single '?' substitution is handled by db.Query()
// automatically. '???' expands to include enough placeholders (as with an IN () list for any count
// of keys > min. max must be >= min unless max == 0.
type Query struct {
	query		string          // The query to execute as prepared statement
	prototype	ResultIfc       // Object to use as a prototype to produce query Result row objects
}

// Make a new one of these
func NewQuery(query string, prototype ResultIfc) *Query {
	q := Query{
		query:		query,
		prototype:	prototype,
	}
	return &q
}

// Run this query against the supplied database Connection with the provided query arguments
func (q *Query) Run(conn *Connection, args ...interface{}) (*ResultSet, error) {
	results := ResultSet{}
	protoQuery := (*q).query
	// TODO: expand query '???' placeholders
	finalQuery := protoQuery

	// Execute the Query
	rows, err := conn.GetConnection().Query(finalQuery, args...)
	if err != nil { return nil, err }

	// Process the result rows
	for rows.Next() {
		// Make a new result object for this row
		// (... and get pointers to all the all the result object members)
		result, resultProperties := (*q).prototype.ZeroClone()

		// Read MySQL columns for this row into the result object member pointers
		err = rows.Scan(resultProperties...)
		if nil != err { return nil, err }

		results = append(results, result)
	}

	return &results, nil
}

