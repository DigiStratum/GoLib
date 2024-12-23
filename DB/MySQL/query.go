package mysql

/*

A Query is attached to a database connection.

The job of the Query interface is to abstract the MySQL interface technicalities away from the
consumer.

We can add other RunReturn{type}() variants for datetime, float, etc. as needed.

Prepared statements are a good idea for even single statements for security (makes sql injection
impossible):
ref: https://stackoverflow.com/questions/1849803/are-prepared-statements-a-waste-for-normal-queries-php

FIXME: we attach sqlQuery to a connection upon creation... but a leased connection could go away,
leaving the sqlQuery and any prepared statement attached to nothing. We should have a way to deal with
this, either self-destruct, or recover a leased connection from the pool, or cause the consumer to
do the same, etc. Probably best left to the consumer so that they can connect their own connection
link to the same one... Refactored the factory function to separate attachment of the sqlQuery to a
given ConnectionIfc so that the consumer can reattach a sqlQuery as needed... but it still needs to
receive some indicator that this is needed.

*/

import (
	"fmt"

	db "database/sql"

	nullables "github.com/DigiStratum/GoLib/DB/MySQL/nullables"
)

type QueryIfc interface {
	Run(args ...interface{}) (*result, error)
	RunReturnValue(receiver interface{}, args ...interface{}) error
	RunReturnInt(args ...interface{}) (*int, error)
	RunReturnString(args ...interface{}) (*string, error)
	RunReturnOne(args ...interface{}) (*ResultRow, error)
	RunReturnAll(args ...interface{}) (*resultSet, error)
	RunReturnSome(max int, args ...interface{}) (*resultSet, error)
}

type query struct {
	connection	ConnectionIfc
	sqlQuery		SQLQueryIfc
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
// Returns nil+error if there is any problem setting up the sqlQuery...!
func NewQuery(connection ConnectionIfc, sqlQuery SQLQueryIfc) (*query, error) {
	if nil == connection { return nil, fmt.Errorf("Supplied connection was nil") }
	if nil == sqlQuery { return nil, fmt.Errorf("Supplied sqlQuery was nil") }
	return &query{
		connection:	connection,
		sqlQuery:		sqlQuery,
	}, nil
}

// -------------------------------------------------------------------------------------------------
// QueryIfc
// -------------------------------------------------------------------------------------------------

// Run this sqlQuery against the supplied database Connection with the provided sqlQuery arguments
func (r query) Run(args ...interface{}) (*result, error) {
	result, err := r.connection.Exec(r.sqlQuery, args...)
	return NewResult(result), err
}

// Run this sqlQuery against the supplied database Connection with the provided sqlQuery arguments
// This variant returns only a single value (any type pointed at by receiver) as the only column
// of the only row of the result
func (r query) RunReturnValue(receiver interface{}, args ...interface{}) error {
	if row := r.connection.QueryRow(r.sqlQuery, args...); nil != row {
		err := row.Scan(receiver)
		if db.ErrNoRows == err { return nil }
		return err
	}
	return nil
}

// Run this sqlQuery against the supplied database Connection with the provided sqlQuery arguments
// This variant returns only a single int value as the only column of the only row of the result
func (r query) RunReturnInt(args ...interface{}) (*int, error) {
	var value int
	err := r.RunReturnValue(&value, args...)
	if nil == err { return &value, nil }
	return nil, err
}

// Run this sqlQuery against the supplied database Connection with the provided sqlQuery arguments
// This variant returns only a single string value as the only column of the only row of the result
func (r query) RunReturnString(args ...interface{}) (*string, error) {
	var value string
	err := r.RunReturnValue(&value, args...)
	if nil == err { return &value, nil }
	return nil, err
}

// Run this sqlQuery against the supplied database Connection with the provided sqlQuery arguments
// This variant returns only a single ResultRowIfc value as the only row of the result
func (r query) RunReturnOne(args ...interface{}) (*ResultRow, error) {
	results, err := r.RunReturnSome(1, args...)
	if nil != err { return nil, err }
	if (nil == results) || (0 == results.Len()) { return nil, nil }
	return results.Get(0), nil
}

// Run this sqlQuery against the supplied database Connection with the provided sqlQuery arguments
// This variant returns all result rows as a set
func (r query) RunReturnAll(args ...interface{}) (*resultSet, error) {
	return r.RunReturnSome(0, args...)
}

// Run this sqlQuery against the supplied database Connection with the provided sqlQuery arguments
// This variant returns a set of result rows up to the max count specified where 0=unlimited (all)
// ref: https://kylewbanks.com/blog/sqlQuery-result-to-map-in-golang
func (r query) RunReturnSome(max int, args ...interface{}) (*resultSet, error) {
	// Even if max == 1 we use Query() instead of QueryRow() to leverage the same Scan converter(s)
	rows, err := r.connection.Query(r.sqlQuery, args...)
	// If the sqlQuery returned no results, handle it specifically...
	if db.ErrNoRows == err { return nil, nil }
	if nil != err { return nil, err }
	if nil != rows { defer rows.Close() }

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

	// Process all the rows of the sqlQuery result
	results := NewResultSet()
	cols, _ := rows.Columns()
	num := 0
	for rows.Next() {
		num++
		columnValues, columnPointers := makeScanReceiver(len(cols))
		if err := rows.Scan(*columnPointers...); err != nil { return nil, err }
		result := convertScanReceiverToResultRow(&cols, columnValues)
		results.Add(result)
		if (max > 0) && (num >= max) { break }
	}
	return results, nil
}
