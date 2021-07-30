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
	Rollback() error
	Begin() error
	Commit() error
	NewQuery(query string) (QueryIfc, error)
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

type connection struct {
	dsn		string		// Full Data Source Name for this connection
	conn		*db.DB		// Read-Write Connection
	transaction	*db.Tx		// Our transaction, if we're in the middle of one
}

// Make a new one of these and connect!
func NewConnection(dsn string) (ConnectionIfc, error) {
	connection := connection{
		dsn:	dsn,
	}
	return &connection, connection.Connect()
}

// -------------------------------------------------------------------------------------------------
// ConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

// ------------
// Connections
// ------------

// Check whether this connection is established
func (c *connection) IsConnected() bool {
	if nil == (*c).conn { return false }
	return nil == (*c).conn.Ping()
}

// Establish the connection using the suplied DSN
func (c *connection) Connect() error {
	// If we're already connected, nothing to do
	if c.IsConnected() { return nil }
	var err error
	(*c).conn, err = db.Open("mysql", (*c).dsn)
	return err
}

// Drop this connection
func (c * connection) Disconnect() {
	// If we're not connected, nothing to do
	if ! c.IsConnected() { return }
	(*c).conn.Close()
}

// Cycle this connection, or establish a new connection if we're not connected
func (c *connection) Reconnect() {
	if c.IsConnected() { c.Disconnect() }
	c.Connect()
}

// Get the underlying connection for the caller to put it to work!
func (c *connection) GetConnection() *db.DB {
	return (*c).conn
}

// ------------
// Transactions
// ------------

func (c *connection) InTransaction() bool {
	return nil != (*c).transaction
}

func (c *connection) Rollback() error {
	// Not in the middle of a Transaction? no-op, no-error!
	if ! c.InTransaction() { return nil }
	err := (*c).transaction.Rollback()
	(*c).transaction = nil
	return err
}

func (c *connection) Begin() error {
	// If we're already in a Transaction...
	if c.InTransaction() {
		// Assume that the app has lost track of the Transaction, maybe lost the connection lease: reset!
		err := c.Rollback()
		if nil != err { return err }
	}
	var err error
	(*c).transaction, err = (*c).conn.Begin()
	return err
}

func (c *connection) Commit() error {
	if ! c.InTransaction() { return errors.New("No active transaction!") }
	err := (*c).transaction.Commit()
	(*c).transaction = nil
	return err
}

func (c *connection) NewQuery(query string) (QueryIfc, error) {
	return NewQuery(c, query)
}

// ------------
// Operations
// ------------

func (c *connection) Prepare(query string) (*db.Stmt, error) {
	if c.InTransaction() { return (*c).transaction.Prepare(query) }
	return (*c).conn.Prepare(query)
}

func (c *connection) Exec(query string, args ...interface{}) (db.Result, error) {
	if c.InTransaction() { return (*c).transaction.Exec(query, args...) }
	return (*c).conn.Exec(query, args...)
}

func (c *connection) Query(query string, args ...interface{}) (*db.Rows, error) {
	if c.InTransaction() { return (*c).transaction.Query(query, args...) }
	return (*c).conn.Query(query, args...)
}

func (c *connection) QueryRow(query string, args ...interface{}) *db.Row {
	if c.InTransaction() { return (*c).transaction.QueryRow(query, args...) }
	return (*c).conn.QueryRow(query, args...)
}

// ------------
// Statements
// ------------

func (c *connection) StmtExec(stmt *db.Stmt, args ...interface{}) (db.Result, error) {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if c.InTransaction() { return (*c).transaction.Stmt(stmt).Exec(args...) }
	return stmt.Exec(args...)
}

func (c *connection) StmtQuery(stmt *db.Stmt, args ...interface{}) (*db.Rows, error) {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if c.InTransaction() { return (*c).transaction.Stmt(stmt).Query(args...) }
	return stmt.Query(args...)
}

func (c *connection) StmtQueryRow(stmt *db.Stmt, args ...interface{}) *db.Row {
	// If we're in a transaction, attach the statement and invoke, otherwise invoke directly
	if c.InTransaction() { return (*c).transaction.Stmt(stmt).QueryRow(args...) }
	return stmt.QueryRow(args...)
}
