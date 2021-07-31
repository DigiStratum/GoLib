package mysql

import (
	"time"
	db "database/sql"
)

// A Pooled Connection wraps a raw DB connection with additional metadata to manage leasing
// We are not exporting this because it is only important to the connection package internal implementation
type PooledConnectionIfc interface {

	// Connections
	IsConnected() bool
	Connect() error
	Disconnect()
	Reconnect()

	// Leasing
	IsLeased() bool
	MatchesLeaseKey(leaseKey int64) bool
	Lease(leaseKey int64)
	Release()
	Touch()

	// Transactions
	InTransaction() bool
	Rollback() error
	Begin() error
	Commit() error

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

type pooledConnection struct {
	connection	ConnectionIfc	// Our underlying database connection
	establishedAt	int64		// Time that this connection was established to the DB
	lastActiveAt	int64		// Last time this connection saw activity from the consumer
	lastLeasedAt	int64		// Last time this connection was leased out
	isLeased	bool		// Is this connection currently leased out?
	leaseKey	int64		// This is the lease key for the current lease holder
}

func NewPooledConnection(dsn string) (PooledConnectionIfc, error) {
	connection, err := NewConnection(dsn)
	if nil != err { return nil, err }
	now := time.Now().Unix()
	pc := pooledConnection{
		connection:	connection,
		establishedAt:	now,
		lastActiveAt:	0,
		lastLeasedAt:	0,
		isLeased:	false,
		leaseKey:	0,
	}
	return &pc, nil
}

// -------------------------------------------------------------------------------------------------
// pooledConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Connections
func (pc *pooledConnection) IsConnected() bool { return (*pc).connection.IsConnected() }
func (pc *pooledConnection) Connect() error { return (*pc).connection.Connect() }
func (pc *pooledConnection) Disconnect() { (*pc).connection.Disconnect() }
func (pc *pooledConnection) Reconnect() { (*pc).connection.Reconnect() }

// Leasing
func (pc *pooledConnection) IsLeased() bool { return (*pc).isLeased }

func (pc *pooledConnection) MatchesLeaseKey(leaseKey int64) bool {
	if ! pc.IsLeased() { return false }
	return (*pc).leaseKey == leaseKey
}

func (pc *pooledConnection) Lease(leaseKey int64) {
	(*pc).isLeased = true
	(*pc).leaseKey = leaseKey
	now := time.Now().Unix()
	(*pc).lastLeasedAt = now
	(*pc).lastActiveAt = now
}

func (pc *pooledConnection) Release() {
	(*pc).isLeased = false
}

func (pc *pooledConnection) Touch() {
	(*pc).lastActiveAt = time.Now().Unix()
}

// Transactions
func (pc *pooledConnection) InTransaction() bool { return (*pc).connection.InTransaction() }
func (pc *pooledConnection) Rollback() error { return (*pc).connection.Rollback() }
func (pc *pooledConnection) Begin() error { return (*pc).connection.Begin() }
func (pc *pooledConnection) Commit() error { return (*pc).connection.Commit() }

// Operations
func (pc *pooledConnection) Prepare(query string) (*db.Stmt, error) { return (*pc).connection.Prepare(query) }
func (pc *pooledConnection) Exec(query string, args ...interface{}) (db.Result, error) { return (*pc).connection.Exec(query, args...) }
func (pc *pooledConnection) Query(query string, args ...interface{}) (*db.Rows, error) { return (*pc).connection.Query(query, args...) }
func (pc *pooledConnection) QueryRow(query string, args ...interface{}) *db.Row { return (*pc).connection.QueryRow(query, args...) }


// Statements
func (pc *pooledConnection) StmtExec(stmt *db.Stmt, args ...interface{}) (db.Result, error) { return (*pc).connection.StmtExec(stmt, args...) }
func (pc *pooledConnection) StmtQuery(stmt *db.Stmt, args ...interface{}) (*db.Rows, error) {  return (*pc).connection.StmtQuery(stmt, args...) }
func (pc *pooledConnection) StmtQueryRow(stmt *db.Stmt, args ...interface{}) *db.Row {  return (*pc).connection.StmtQueryRow(stmt, args...) }
