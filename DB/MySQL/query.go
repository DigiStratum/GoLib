package mysql

/*

TODO: add other RunReturn{type}() variants for datetime, float, etc. as needed

Prepared statements are a good idea for even single statements for security (sql injection is impossible):
ref: https://stackoverflow.com/questions/1849803/are-prepared-statements-a-waste-for-normal-queries-php

*/

import (
	"strings"
	db "database/sql"

	nullables "github.com/DigiStratum/GoLib/DB/MySQL/nullables"
)

type QueryIfc interface {
	// Public Interface
	Run(args ...interface{}) error
	RunReturnValue(receiver interface{}, args ...interface{}) error
	RunReturnInt(args ...interface{}) (*int, error)
	RunReturnString(args ...interface{}) (*string, error)
	RunReturnOne(args ...interface{}) (ResultRowIfc, error)
	RunReturnAll(args ...interface{}) (ResultSetIfc, error)
	RunReturnSome(max int, args ...interface{}) (ResultSetIfc, error)
	// Private interface
	resolveQuery(args ... interface{}) (*string, error)
}

type query struct {
	connection	ConnectionIfc
	query		string
	prepareOk	bool
}

// Make a new one of these!
// Returns nil if there is any problem setting up the query...!
func NewQuery(connection interface{}, qry string) QueryIfc {

	// We are going to allow multiple interfaces to be passed in here and convert to ConnectionIfc (or fail)
	var c ConnectionIfc
	if c, ok := connection.(ConnectionIfc); ! ok { return nil }

	// If the query does NOT contain a list for expansion ('???') then we can use a prepared statement
	// Note: a literal string value of '???' would be encoded as '\\?\\?\\?'
	// https://pkg.go.dev/database/sql#Stmt
	var statement *db.Stmt
	var err error
	if ! strings.Contains(qry, "???") {
		statement, err = c.Prepare(qry)
		if nil != err { return nil } // TODO: log an error! (?)
	}

	q := query{
		connection:	c,
		query:		qry,
		prepareOk:	! strings.Contains(qry, "???"),
		statement:	statement,
	}
	return &q
}

// -------------------------------------------------------------------------------------------------
// QueryIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Run this query against the supplied database Connection with the provided query arguments
// No result rows are expected or returned
func (q *query) Run(args ...interface{}) error {
	var err error
	// TODO: Capture Exec() result (swallowed into _) to get result.RowsAffected(), etc
	if nil != (*q).statement {
		// Prepared statement need not specify a query (the statement is the query)
		_, err = (*q).statement.Exec(args...)
	} else {
		// Resolve a non-prepared statement query with any of our own substitutions
		qry, err := q.resolveQuery(args...)
		if nil != err { return err }
		_, err = (*q).connection.Exec(*qry, args...)
	}
	return err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single value (any type pointed at by receiver) as the only column
// of the only row of the result
func (q *query) RunReturnValue(receiver interface{}, args ...interface{}) error {
	var row *db.Row
	var err error

	if nil != (*q).statement {
		// Prepared statement need not specify a query (the statement is the query)
		row = (*q).statement.QueryRow(args...)
	} else {
		// Resolve a non-prepared statement query with any of our own substitutions
		qry, err := q.resolveQuery(args...)
		if nil != err { return err }
		row = (*q).connection.QueryRow(*qry, args...)
	}

	if nil == row { return nil }

	err = row.Scan(receiver)
	if db.ErrNoRows == err { return nil }
	return err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single int value as the only column of the only row of the result
func (q *query) RunReturnInt(args ...interface{}) (*int, error) {
	var value int
	err := q.RunReturnValue(&value, args...)
	if nil == err { return &value, nil }
	return nil, err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single string value as the only column of the only row of the result
func (q *query) RunReturnString(args ...interface{}) (*string, error) {
	var value string
	err := q.RunReturnValue(&value, args...)
	if nil == err { return &value, nil }
	return nil, err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single ResultRowIfc value as the only row of the result
func (q *query) RunReturnOne(args ...interface{}) (ResultRowIfc, error) {
	results, err := q.RunReturnSome(1, args...)
	if nil != err { return nil, err }
	if (nil == results) || (0 == results.Len()) { return nil, nil }
	return results.Get(0), nil
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns all result rows as a set
func (q *query) RunReturnAll(args ...interface{}) (ResultSetIfc, error) {
	return q.RunReturnSome(0, args...)
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns a set of result rows up to the max count specified where 0=unlimited (all)
// ref: https://kylewbanks.com/blog/query-result-to-map-in-golang
func (q *query) RunReturnSome(max int, args ...interface{}) (ResultSetIfc, error) {
	var rows *db.Rows
	var err error

	if nil != (*q).statement {
		// Prepared statement need not specify a query (the statement is the query)
		rows, err = (*q).statement.Query(args...)
	} else {
		// Resolve a non-prepared statement query with any of our own substitutions
		qry, err := q.resolveQuery(args...)
		if nil != err { return nil, err }
		rows, err = (*q).connection.Query(*qry, args...)
	}

	// If the query returned no results, handle it specifically...
	if db.ErrNoRows == err { return nil, nil }
	if nil != err { return nil, err }
	if nil != rows { defer rows.Close() }

	results := NewResultSet()
	cols, _ := rows.Columns()
	num := 0
	for ((max == 0) || (max < num)) && rows.Next() {
		num++
		columnValues, columnPointers := makeScanReceiver(len(cols))
		if err := rows.Scan(*columnPointers...); err != nil { return nil, err }
		result := convertScanReceiverToResultRow(&cols, columnValues)
		results.add(result)
	}
	return results, nil
}

// -------------------------------------------------------------------------------------------------
// QueryIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Placeholder to support resolving magic expander tags, etc within our query
// TODO: Consider support for literal substitutions here that could get us things like variable table names and other unescaped values for query variance
func (q *query) resolveQuery(args ... interface{}) (*string, error) {
	protoQuery := (*q).query
	// TODO: expand query '???' placeholders
	// because args is interface{}, we can pass whatever we want for this. let's declare that any arg which is an
	// array will be treated as a set to expand. we will expect that there be the same number of sets as there are
	// '???' placeholders. We will replace the  Nth '???' placeholder in the order that they appear with the count of
	// '?,?,...' placeholders that matches the Len() of the Nth array. If the count of arrays supplied does not matche
	// the count of '???' placeholders, then error
	finalQuery := protoQuery
	return &finalQuery, nil
}

// -------------------------------------------------------------------------------------------------
// Private supporting functions
// -------------------------------------------------------------------------------------------------

// Return a slice of values and pointers to those values for Scan() to map result into
func makeScanReceiver(size int) (*[]string, *[]interface{}) {
	columnPointers := make([]interface{}, size)
	columns := make([]string, size)
	for i, _ := range columns { columnPointers[i] = &columns[i] }
	return &columns, &columnPointers
}

// Create our map, and retrieve the value for each column from the pointers,
// slice, storing it in the map with the name of the column as the key.
// Note: names and values array len() must match. If they don't, then the Universe is off balance
func convertScanReceiverToResultRow(names, values *[]string) ResultRowIfc {
	result := NewResultRow()
	for i, name := range *names { result.Set(name, nullables.NewNullable((*values)[i])) }
	return result
}
