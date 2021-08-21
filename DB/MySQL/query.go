package mysql

/*

A Query is attached to a database connection.

The job of the Query interface is to abstract the MySQL interface technicalities away from the consumer.

We can add other RunReturn{type}() variants for datetime, float, etc. as needed.

Prepared statements are a good idea for even single statements for security (makes sql injection impossible):
ref: https://stackoverflow.com/questions/1849803/are-prepared-statements-a-waste-for-normal-queries-php

FIXME: we attach query to a connection upon creation... but a leased connection could go away, leaving the query and any
prepared statement attached to nothing. We should have a way to deal with this, either self-destruct, or recover a leased
connection from the pool, or cause the consumer to do the same, etc. Probably best left to the consumer so that they can
connect their own connection link to the same one... Refactored the factory function to separate attachment of the query
to a given ConnectionIfc so that the consumer can reattach a query as needed... but it still needs to receive some indicator
that this is needed.

// TODO: Consider support for literal substitutions here that could get us things like variable table names and other unescaped values for query variance
// TODO: expand query '???' placeholders
// because args is interface{}, we can pass whatever we want for this. let's declare that any arg which is an
// array will be treated as a set to expand. we will expect that there be the same number of sets as there are
// '???' placeholders. We will replace the  Nth '???' placeholder in the order that they appear with the count of
// '?,?,...' placeholders that matches the Len() of the Nth array. If the count of arrays supplied does not matche
// the count of '???' placeholders, then error

*/

import (
	"strings"
	db "database/sql"

	nullables "github.com/DigiStratum/GoLib/DB/MySQL/nullables"
)

type QueryIfc interface {
	AttachConnection(connection ConnectionIfc) error
	Run(args ...interface{}) (*Result, error)
	RunReturnValue(receiver interface{}, args ...interface{}) error
	RunReturnInt(args ...interface{}) (*int, error)
	RunReturnString(args ...interface{}) (*string, error)
	RunReturnOne(args ...interface{}) (*ResultRow, error)
	RunReturnAll(args ...interface{}) (*ResultSet, error)
	RunReturnSome(max int, args ...interface{}) (*ResultSet, error)
}

type Query struct {
	connection	ConnectionIfc
	query		string
	statement	*db.Stmt
}

// Make a new one of these!
// Returns nil+error if there is any problem setting up the query...!
func NewQuery(connection ConnectionIfc, qry string) (*Query, error) {
	q := Query{
		query:		qry,
	}

	if err := q.AttachConnection(connection); nil != err { return nil, err }

	return &q, nil
}

// -------------------------------------------------------------------------------------------------
// QueryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Query) AttachConnection(connection ConnectionIfc) error {
	// If the query does NOT contain a list for expansion ('???') then we can use a prepared statement
	// Note: a literal string value of '???' would be encoded as '\\?\\?\\?'
	// https://pkg.go.dev/database/sql#Stmt
	var statement *db.Stmt
	var err error
	if ! strings.Contains(r.query, "???") {
		statement, err = connection.Prepare(r.query)
		if nil != err { return err }
	}

	r.connection = connection
	r.statement = statement
	return nil
}

// Run this query against the supplied database Connection with the provided query arguments
func (r Query) Run(args ...interface{}) (*Result, error) {
	var result db.Result
	var err error
	if nil != r.statement {
		// Prepared statement need not specify a query (the statement is the query)
		result, err = r.connection.StmtExec(r.statement, args...)
	} else {
		// Resolve a non-prepared statement query with any of our own substitutions
		qry, err := r.resolveQuery(args...)
		if nil != err { return nil, err }
		result, err = r.connection.Exec(*qry, args...)
	}
	return NewResult(result), err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single value (any type pointed at by receiver) as the only column
// of the only row of the result
func (r Query) RunReturnValue(receiver interface{}, args ...interface{}) error {
	var row *db.Row
	var err error

	if nil != r.statement {
		// Prepared statement need not specify a query (the statement is the query)
		row = r.connection.StmtQueryRow(r.statement, args...)
	} else {
		// Resolve a non-prepared statement query with any of our own substitutions
		qry, err := r.resolveQuery(args...)
		if nil != err { return err }
		row = r.connection.QueryRow(*qry, args...)
	}

	if nil == row { return nil }

	err = row.Scan(receiver)
	if db.ErrNoRows == err { return nil }
	return err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single int value as the only column of the only row of the result
func (r Query) RunReturnInt(args ...interface{}) (*int, error) {
	var value int
	err := r.RunReturnValue(&value, args...)
	if nil == err { return &value, nil }
	return nil, err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single string value as the only column of the only row of the result
func (r Query) RunReturnString(args ...interface{}) (*string, error) {
	var value string
	err := r.RunReturnValue(&value, args...)
	if nil == err { return &value, nil }
	return nil, err
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns only a single ResultRowIfc value as the only row of the result
func (r Query) RunReturnOne(args ...interface{}) (*ResultRow, error) {
	results, err := r.RunReturnSome(1, args...)
	if nil != err { return nil, err }
	if (nil == results) || (0 == results.Len()) { return nil, nil }
	result := results.Get(0)
	return &result, nil
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns all result rows as a set
func (r Query) RunReturnAll(args ...interface{}) (*ResultSet, error) {
	return r.RunReturnSome(0, args...)
}

// Run this query against the supplied database Connection with the provided query arguments
// This variant returns a set of result rows up to the max count specified where 0=unlimited (all)
// ref: https://kylewbanks.com/blog/query-result-to-map-in-golang
func (r Query) RunReturnSome(max int, args ...interface{}) (*ResultSet, error) {
	var rows *db.Rows
	var err error

	if nil != r.statement {
		// Prepared statement need not specify a query (the statement is the query)
		rows, err = r.connection.StmtQuery(r.statement, args...)
	} else {
		// Resolve a non-prepared statement query with any of our own substitutions
		qry, err := r.resolveQuery(args...)
		if nil != err { return nil, err }
		rows, err = r.connection.Query(*qry, args...)
	}

	// Return a slice of values and pointers to those values for Scan() to map result into
	makeScanReceiver := func(size int) (*[]string, *[]interface{}) {
		columnPointers := make([]interface{}, size)
		columns := make([]string, size)
		for i, _ := range columns { columnPointers[i] = &columns[i] }
		return &columns, &columnPointers
	}

	// Create our map, and retrieve the value for each column from the pointers,
	// slice, storing it in the map with the name of the column as the key.
	// Note: names and values array len() must match. If they don't, then the Universe is off balance
	convertScanReceiverToResultRow := func(names, values *[]string) *ResultRow {
		result := NewResultRow()
		for i, name := range *names {
			nullableValue := nullables.NewNullable((*values)[i])
			result.Set(name, *nullableValue)
		}
		return result
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
		results.Add(result)
	}
	return results, nil
}

// -------------------------------------------------------------------------------------------------
// QueryIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Placeholder to support resolving magic expander tags, etc within our query
func (r Query) resolveQuery(args ... interface{}) (*string, error) {
	protoQuery := r.query
	finalQuery := protoQuery
	return &finalQuery, nil
}
