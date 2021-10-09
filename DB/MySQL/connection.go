package mysql

/*

DB Connection - All the low-level nitty-gritty interacting with the sql driver

ref: https://github.com/go-sql-driver/mysql#interpolateparams
ref: https://pkg.go.dev/database/sql#Tx.Stmt

*/

import (
	"errors"
	db "database/sql"
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
	Connect() error
	Disconnect()
	Reconnect()

	// Transactions
	ConnectionCommonIfc

	// Operations
	Prepare(query string) (*db.Stmt, error)
	Exec(query string, args ...interface{}) (db.Result, error)
	Query(query string, args ...interface{}) (*db.Rows, error)
	QueryRow(query string, args ...interface{}) *db.Row

	// Statements
	StmtExec(stmt *db.Stmt, args ...interface{}) (db.Result, error)
	StmtQuery(stmt *db.Stmt, args ...interface{}) (*db.Rows, error)
	StmtQueryRow(stmt *db.Stmt, args ...interface{}) *db.Row
}

type Connection struct {
	dsn		string		// Full Data Source Name for this connection
	conn		*db.DB		// Read-Write Connection
	transaction	*db.Tx		// Our transaction, if we're in the middle of one
}

// Make a new one of these and connect!
func NewConnection(dsn string) (*Connection, error) {
	connection := Connection{
		dsn:	dsn,
	}
	err := connection.Connect()
	if nil != err { return nil, err }
	return &connection, nil
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

// Establish the connection using the suplied DSN
func (r *Connection) Connect() error {
	// If we're already connected, nothing to do
	if r.IsConnected() { return nil }
	var err error
	r.conn, err = db.Open("mysql", r.dsn)
	return err
}

// Drop this connection
func (r *Connection) Disconnect() {
	// If we're not connected, nothing to do
	if ! r.IsConnected() { return }
	r.conn.Close()
}

// Cycle this connection, or establish a new connection if we're not connected
func (r *Connection) Reconnect() {
	if r.IsConnected() { r.Disconnect() }
	r.Connect()
}

// Get the underlying connection for the caller to put it to work!
func (r *Connection) GetConnection() *db.DB {
	return r.conn
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
	if ! r.InTransaction() { return errors.New("No active transaction!") }
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

func (r Connection) Prepare(query string) (*db.Stmt, error) {
	if r.InTransaction() { return r.transaction.Prepare(query) }
	return r.conn.Prepare(query)
}

func (r Connection) Exec(query string, args ...interface{}) (db.Result, error) {
	if r.InTransaction() { return r.transaction.Exec(query, args...) }
	return r.conn.Exec(query, args...)
}

func (r Connection) Query(query string, args ...interface{}) (*db.Rows, error) {
	if r.InTransaction() { return r.transaction.Query(query, args...) }
	return r.conn.Query(query, args...)
}

func (r Connection) QueryRow(query string, args ...interface{}) *db.Row {
	if r.InTransaction() { return r.transaction.QueryRow(query, args...) }
	return r.conn.QueryRow(query, args...)
}

// ------------
// Statements
// ------------

func (r Connection) StmtExec(stmt *db.Stmt, args ...interface{}) (db.Result, error) {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if r.InTransaction() { return r.transaction.Stmt(stmt).Exec(args...) }
	return stmt.Exec(args...)
}

func (r Connection) StmtQuery(stmt *db.Stmt, args ...interface{}) (*db.Rows, error) {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if r.InTransaction() { return r.transaction.Stmt(stmt).Query(args...) }
	return stmt.Query(args...)
}

func (r Connection) StmtQueryRow(stmt *db.Stmt, args ...interface{}) *db.Row {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if r.InTransaction() { return r.transaction.Stmt(stmt).QueryRow(args...) }
	return stmt.QueryRow(args...)
}
