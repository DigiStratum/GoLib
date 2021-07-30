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
	// Embed Transaction support to this interface
	ConnectionCommonIfc
}

type leasedConnection struct {
	pooledConnection	PooledConnectionIfc
	leaseKey		int64
	errNoLease		error
}

func NewLeasedConnection(pooledConnection PooledConnectionIfc, leaseKey int64) LeasedConnectionIfc {
	pooledConnection.Lease(leaseKey)
	lc := leasedConnection{
		pooledConnection:	pooledConnection,
		leaseKey:		leaseKey,
		errNoLease:		errors.New("No Leased Connection!"),
	}
	return &lc
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

// -------------------------------------------------------------------------------------------------
// ConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (lc *leasedConnection) IsConnected() bool {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return false }
	return (*lc).pooledConnection.IsConnected()
}

//  Leased connections are not allowed to tamper with the connection lifecycle
func (lc *leasedConnection) Disconnect() { }
func (lc *leasedConnection) Reconnect() { }
func (lc *leasedConnection) Connect() error {
	return errors.New("Leased connection - no state changes allowed")
}

func (lc *leasedConnection) InTransaction() bool {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return false }
	return (*lc).pooledConnection.InTransaction()
}

func (lc *leasedConnection) Rollback() error {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return (*lc).errNoLease }
	return (*lc).pooledConnection.Rollback()
}

func (lc *leasedConnection) Begin() error {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return (*lc).errNoLease }
	return (*lc).pooledConnection.Begin()
}

func (lc *leasedConnection) Commit() error {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return (*lc).errNoLease }
	return (*lc).pooledConnection.Commit()
}

func (lc *leasedConnection) NewQuery(qry string) (QueryIfc, error) {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return nil, errors.New("No Leased Connection!") }
	// Feed NewQuery() our leasedConnection so it doesn't have direct access to underlying pooledConnection
	return NewQuery(lc, qry)
}

func (lc *leasedConnection) Prepare(query string) (*db.Stmt, error) {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return nil, (*lc).errNoLease }
	return (*lc).pooledConnection.Prepare(query)
}

func (lc *leasedConnection) Exec(query string, args ...interface{}) (db.Result, error) {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return nil, (*lc).errNoLease }
	return (*lc).pooledConnection.Exec(query, args...)
}

func (lc *leasedConnection) Query(query string, args ...interface{}) (*db.Rows, error) {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return nil, (*lc).errNoLease }
	return (*lc).pooledConnection.Query(query, args...)
}

func (lc *leasedConnection) QueryRow(query string, args ...interface{}) *db.Row {
	if ! (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return nil }
	return (*lc).pooledConnection.QueryRow(query, args...)
}
