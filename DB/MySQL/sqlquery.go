package mysql

/*
A SQL Query is Runnable if it implements the SQLQueryIfc interface. SQLQuery is a default implementation
of this that supports a raw query pass-through. But this paves the way to implement more advanced ORM
style query builders independent of the Query runner.
*/

import (
	"fmt"
)

type SQLQueryIfc interface {
	// Return the query as a string or else a non-nil error
	GetQuery() (string, error)
}

type SQLQuery struct {
	query		string
}

func NewSQLQuery(query string) *SQLQuery {
	return &SQLQuery{
		query:		query,
	}
}

// -------------------------------------------------------------------------------------------------
// SQLQueryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *SQLQuery) GetQuery() (string,  error) {
	if len(r.query) == 0 { return "", fmt.Errorf("Query is empty") }
	return r.query, nil
}

