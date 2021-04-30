package mysql

/*
TODO: Add some sort of query builder - this will allow us to ditch writing SQL for most needs.
*/

import (
	_ "github.com/go-sql-driver/mysql"
)

// Query public interface
type QueryIfc interface {
	Run(conn ConnectionIfc, args ...interface{}) (ResultSetIfc, error)
}

// The spec for a prepared statement query. Single '?' substitution is handled by db.Query()
// automatically. '???' expands to include enough placeholders (as with an IN () list for any count
// of keys > min. max must be >= min unless max == 0.
type qry struct {
	query		string          // The query to execute as prepared statement
	prototype	ResultIfc       // Object to use as a prototype to produce query Result row objects
}

// Make a new one of these
func NewQuery(query string, prototype ResultIfc) QueryIfc {
	q := qry{
		query:		query,
		prototype:	prototype,
	}
	return &q
}

// Run this query against the supplied database Connection with the provided query arguments
func (q *qry) Run(conn ConnectionIfc, args ...interface{}) (ResultSetIfc, error) {
	protoQuery := (*q).query
	// TODO: expand query '???' placeholders
	finalQuery := protoQuery

	// Execute the Query
	rows, err := conn.GetConnection().Query(finalQuery, args...)
	if err != nil { return nil, err }

	// Process the result rows
	results := NewResultSet()
	for rows.Next() {
		// Make a new result object for this row
		// (... and get pointers to all the all the result object members)
		result, resultProperties := (*q).prototype.ZeroClone()

		// Read MySQL columns for this row into the result object member pointers
		// ref: https://kylewbanks.com/blog/query-result-to-map-in-golang
		err = rows.Scan(resultProperties...)
		if nil != err { return nil, err }

		results = results.Append(result)
	}

	return results, nil
}

