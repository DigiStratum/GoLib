package mysql

/*

DB Connection - All the low-level nitty-gritty interacting with the sql driver

ref: https://github.com/go-sql-driver/mysql#interpolateparams
ref: https://pkg.go.dev/database/sql#Tx.Stmt

TODO:
 * Should we force a single-statement transaction begin/commit for one-off queries run outside of
   an explicit transaction?

*/

import (
	"fmt"
	"database/sql"
)

type ConnectionCommonIfc interface {
	InTransaction() bool
	Begin() error
	NewQuery(query SQLQueryIfc) (QueryIfc, error)
	Commit() error
	Rollback() error
	Exec(query SQLQueryIfc, args ...interface{}) (sql.Result, error)
	Query(query SQLQueryIfc, args ...interface{}) (*sql.Rows, error)
	QueryRow(query SQLQueryIfc, args ...interface{}) *sql.Row
}

type ConnectionIfc interface {
	IsConnected() bool
	ConnectionCommonIfc
}

type Connection struct {
	conn			*sql.DB			// Read-Write Connection
	transaction		*sql.Tx			// Our transaction, if we're in the middle of one
	// Local Cache structures, keyed on resolved SQL query text
	transactionStatements	map[string]*sql.Stmt	// retains transaction-specific prepared statements
	statements		map[string]*sql.Stmt	// retains non-transaction prepared statements
}

// Make a new one of these and connect!
func NewConnection(conn *sql.DB) (*Connection, error) {
	if nil == conn { return nil, fmt.Errorf("Cannot wrap nil connection") }
	connection := Connection{
		conn:			conn,
		statements:		make(map[string]*sql.Stmt),
	}
	return &connection, nil
}

// -------------------------------------------------------------------------------------------------
// io.Closer Public Interface
// -------------------------------------------------------------------------------------------------

// Drop this connection
func (r *Connection) Close() error {
	// If we're not connected, nothing to do
	if ! r.IsConnected() { return nil }
	r.conn.Close()
	r.conn = nil
	return nil
}

// -------------------------------------------------------------------------------------------------
// ConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Check whether this connection is established
func (r Connection) IsConnected() bool {
	if nil == r.conn { return false }
	return nil == r.conn.Ping()
}

func (r Connection) InTransaction() bool {
	return nil != r.transaction
}

func (r *Connection) Begin() error {
	// If we're already in a Transaction...
	if r.InTransaction() {
		// Assume that the app has lost track of the Transaction, maybe lost the connection lease: reset!
		err := r.Rollback()
		if nil != err { return err }
	}
	var err error
	r.transaction, err = r.conn.Begin()
	// Reset the prepared statements for a new transaction
	if nil == err { r.transactionStatements = make(map[string]*sql.Stmt) }
	return err
}

func (r *Connection) NewQuery(query SQLQueryIfc) (QueryIfc, error) {
	return NewQuery(r, query)
}

func (r *Connection) Commit() error {
	if ! r.InTransaction() { return fmt.Errorf("No active transaction!") }
	err := r.transaction.Commit()
	r.transaction = nil
	return err
}

func (r *Connection) Rollback() error {
	// Not in the middle of a Transaction? no-op, no-error!
	if ! r.InTransaction() { return nil }
	err := r.transaction.Rollback()
	r.transaction = nil
	return err
}

func (r Connection) Exec(query SQLQueryIfc, args ...interface{}) (sql.Result, error) {
	stmt, err := r.prepare(query)
	if nil != err { return nil, err }
	return stmt.Exec(args...)
}

func (r Connection) Query(query SQLQueryIfc, args ...interface{}) (*sql.Rows, error) {
	stmt, err := r.prepare(query)
	if nil != err { return nil, err }
	return stmt.Query(args...)
}

// Note: DB.(Stmt.)QueryRow always returns a non-nil value.
func (r Connection) QueryRow(query SQLQueryIfc, args ...interface{}) *sql.Row {
	stmt, err := r.prepare(query)
	if nil != err { return nil }
	return stmt.QueryRow(args...)
}

// -------------------------------------------------------------------------------------------------
// Connection Private Implementation
// -------------------------------------------------------------------------------------------------

func (r Connection) prepare(query SQLQueryIfc) (*sql.Stmt, error) {
	// Resolve the query
	sql, err := query.Resolve()
	if nil != err { return nil, err }

	if r.InTransaction() {
		// If this query is already in the transaction's prepared statements...
		if stmt, ok := r.transactionStatements[sql]; ok {
			return stmt, nil
		}
		stmt, err := r.transaction.Prepare(sql)
		if nil == err { r.transactionStatements[sql] = stmt }
		return stmt, err
	}

	// If this query is already in the non-transaction prepared statements...
	if stmt, ok := r.statements[sql]; ok {
		return stmt, nil
	}
	stmt, err := r.conn.Prepare(sql)
	if nil == err { r.statements[sql] = stmt }
	return stmt, err
}

