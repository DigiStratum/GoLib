package mysql

/*

DB Connection - All the low-level nitty-gritty interacting with the sql driver

ref: https://github.com/go-sql-driver/mysql#interpolateparams
ref: https://pkg.go.dev/database/sql#Tx.Stmt

*/

import (
	"fmt"
	"database/sql"

	"github.com/DigiStratum/GoLib/DB/dbwrapper"
)

type ConnectionCommonIfc interface {
	InTransaction() bool
	Begin() error
	NewQuery(query string) (QueryIfc, error)
	Commit() error
	Rollback() error
}

type ConnectionIfc interface {
	// Connections
	IsConnected() bool

	// Transactions
	ConnectionCommonIfc

	// Operations
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row

	// Statements
	StmtExec(stmt *sql.Stmt, args ...interface{}) (sql.Result, error)
	StmtQuery(stmt *sql.Stmt, args ...interface{}) (*sql.Rows, error)
	StmtQueryRow(stmt *sql.Stmt, args ...interface{}) *sql.Row
}

type Connection struct {
	conn			dbwrapper.DBConnectionIfc	// Read-Write Connection
	transaction		*sql.Tx			// Our transaction, if we're in the middle of one
}

// Make a new one of these and connect!
func NewConnection(conn dbwrapper.DBConnectionIfc) (*Connection, error) {
	if nil == conn { return nil, fmt.Errorf("Cannot wrap nil connection") }
	connection := Connection{
		conn:			conn,
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

// ------------
// Connections
// ------------

// Check whether this connection is established
func (r Connection) IsConnected() bool {
	if nil == r.conn { return false }
	return nil == r.conn.Ping()
}

// ------------
// Transactions
// ------------

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
	return err
}

func (r *Connection) NewQuery(query string) (QueryIfc, error) {
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

// ------------
// Operations
// ------------

func (r Connection) Prepare(query string) (*sql.Stmt, error) {
	if r.InTransaction() { return r.transaction.Prepare(query) }
	return r.conn.Prepare(query)
}

func (r Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	if r.InTransaction() { return r.transaction.Exec(query, args...) }
	return r.conn.Exec(query, args...)
}

func (r Connection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if r.InTransaction() { return r.transaction.Query(query, args...) }
	return r.conn.Query(query, args...)
}

func (r Connection) QueryRow(query string, args ...interface{}) *sql.Row {
	if r.InTransaction() { return r.transaction.QueryRow(query, args...) }
	return r.conn.QueryRow(query, args...)
}

// ------------
// Statements
// ------------

func (r Connection) StmtExec(stmt *sql.Stmt, args ...interface{}) (sql.Result, error) {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if r.InTransaction() { return r.transaction.Stmt(stmt).Exec(args...) }
	return stmt.Exec(args...)
}

func (r Connection) StmtQuery(stmt *sql.Stmt, args ...interface{}) (*sql.Rows, error) {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if r.InTransaction() { return r.transaction.Stmt(stmt).Query(args...) }
	return stmt.Query(args...)
}

func (r Connection) StmtQueryRow(stmt *sql.Stmt, args ...interface{}) *sql.Row {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if r.InTransaction() { return r.transaction.Stmt(stmt).QueryRow(args...) }
	return stmt.QueryRow(args...)
}
