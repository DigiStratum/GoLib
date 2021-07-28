package mysql

/*
A Leased Connection wraps a pooled DB connection with an internally managed lease key. If we lose our lease due to
idle/inactivity or some other reason, we either need a way to recover a working connection automatically, or our
DB access will be gone and we will need to get a new connection. Otherwise our interface here is a mirror of
connection, with each method being a pass-through based on keyed access to the underlying connection.
*/

import (
	"errors"
	db "database/sql"
)

type LeasedConnectionIfc interface {
	NewQuery(query string) (QueryIfc, error)
	// Private
	checkLease() bool
	errNoLease() error
}

type leasedConnection struct {
	pooledConnection	PooledConnectionIfc
	leaseKey		int64
}

func NewLeasedConnection(pooledConnection PooledConnectionIfc, leaseKey int64) LeasedConnectionIfc {
	pooledConnection.Lease(leaseKey)
	lc := leasedConnection{
		pooledConnection:	pooledConnection,
		leaseKey:		leaseKey,
	}
	return &lc
}

// -------------------------------------------------------------------------------------------------
// ConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (lc *leasedConnection) IsConnected() bool {
	if ! lc.checkLease() { return nil }
	return (*lc).pooledConnection.IsConnected()
}

//  Leased connections are not allowed to tamper with the connection lifecycle
func (lc *leasedConnection) Disconnect() { }
func (lc *leasedConnection) Reconnect() { }
func (lc *leasedConnection) Connect() error {
	return errors.New("Leased connection - no state changes allowed")
}

// Transactions - Passthrough
func (lc *leasedConnection) InTransaction() bool {
	if ! lc.checkLease() { return false }
	return (*lc).pooledConnection.InTransaction()
}

func (lc *leasedConnection) Rollback() error {
	if ! lc.checkLease() { return lc.errNoLease() }
	return (*lc).pooledConnection.Rollback()
}

func (lc *leasedConnection) Begin() error {
	if ! lc.checkLease() { return lc.errNoLease() }
	return (*lc).pooledConnection.Begin()
}

func (lc *leasedConnection) Commit() error {
	if ! lc.checkLease() { return lc.errNoLease() }
	return (*lc).pooledConnection.Commit()
}

// Operations - Passthrough
func (lc *leasedConnection) Prepare(query string) (*db.Stmt, error) {
	if ! lc.checkLease() { return nil, lc.errNoLease() }
	return (*lc).pooledConnection.Prepare(query)
}

func (lc *leasedConnection) Exec(query string, args ...interface{}) (db.Result, error) {
	if ! lc.checkLease() { return nil, lc.errNoLease() }
	return (*lc).pooledConnection.Exec(query, args...)
}

func (lc *leasedConnection) Query(query string, args ...interface{}) (*db.Rows, error) {
	if ! lc.checkLease() { return nil, lc.errNoLease() }
	return (*lc).pooledConnection.Query(query, args...)
}

func (lc *leasedConnection) QueryRow(query string, args ...interface{}) *db.Row {
	if ! lc.checkLease() { return nil }
	return (*lc).pooledConnection.QueryRow(query, args...)
}

func (lc *leasedConnection) Stmt(stmt *db.Stmt) *db.Stmt {
	if ! lc.checkLease() { return nil }
	return (*lc).pooledConnection.Stmt(stmt)
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (lc *leasedConnection) NewQuery(qry string) (QueryIfc, error) {
	if ! lc.checkLease() { return nil, errors.New("No Leased Connection!") }
	return query.NewQuery(lc, qry)
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionIfc Private Interface
// -------------------------------------------------------------------------------------------------

// TODO: Move this to pooledConnection so that it has the power of accept/reject instead of self-policing here?
// Check whether we still hold the lease on our connection
func (lc *leasedConnection) checkLease() bool {
	if nil == (*lc).pooledConnection { return false }
	if (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return true }
	// The lease has changed! Drop the connection to eliminate any chance of further inactivity
	(*lc).pooledConnection = nil
	return false
}

func (lc *leasedConnection) errNoLease() error {
	return errors.New("No Leased Connection!")
}