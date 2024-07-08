package mysql

/*
A Leased Connection wraps a pooled DB connection with an internally managed lease key. If we lose our lease due to
idle/inactivity or some other reason, we either need a way to recover a working connection automatically, or our
DB access will be gone and we will need to get a new connection. Otherwise our interface here is a mirror of
connection, with each method being a pass-through based on keyed access to the underlying connection.
*/

import (
	"fmt"
	db "database/sql"
)

type LeasedConnectionIfc interface {
	ConnectionCommonIfc

	Release() error
}

type leasedConnection struct {
	pooledConnection	PooledConnectionIfc
	leaseKey		int64
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewLeasedConnection(pooledConnection PooledConnectionIfc, leaseKey int64) *leasedConnection {
	pooledConnection.Lease(leaseKey)
	return &leasedConnection{
		pooledConnection:	pooledConnection,
		leaseKey:		leaseKey,
	}
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionIfc
// -------------------------------------------------------------------------------------------------

func (r *leasedConnection) Release() error {
	if nil == r { return fmt.Errorf("leasedConnection is nil") }
	if nil == r.pooledConnection { return fmt.Errorf("leasedConnection.pooledConnection is nil") }
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return errNoLease() }
	if err := r.pooledConnection.Release(); nil != err { return err }
	r.leaseKey = 0
	return nil
}

// -------------------------------------------------------------------------------------------------
// ConnectionIfc
// -------------------------------------------------------------------------------------------------

func (r *leasedConnection) IsConnected() bool {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return false }
	return r.pooledConnection.IsConnected()
}

//  Leased connections are not allowed to tamper with the connection lifecycle
func (r *leasedConnection) Disconnect() { }
func (r *leasedConnection) Reconnect() { }
func (r *leasedConnection) Connect() error {
	return fmt.Errorf("Leased connection - no state changes allowed")
}

func (r *leasedConnection) InTransaction() bool {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return false }
	return r.pooledConnection.InTransaction()
}

func (r *leasedConnection) Begin() error {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return errNoLease() }
	return r.pooledConnection.Begin()
}

func (r *leasedConnection) NewQuery(query SQLQueryIfc) (QueryIfc, error) {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil, fmt.Errorf("No Leased Connection!") }
	// Feed NewQuery() our *leasedConnection so it doesn't have direct access to underlying pooledConnection
	return NewQuery(r, query)
}

func (r *leasedConnection) Commit() error {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return errNoLease() }
	return r.pooledConnection.Commit()
}

func (r *leasedConnection) Rollback() error {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return errNoLease() }
	return r.pooledConnection.Rollback()
}

func (r *leasedConnection) Exec(query SQLQueryIfc, args ...interface{}) (db.Result, error) {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil, errNoLease() }
	return r.pooledConnection.Exec(query, args...)
}

func (r *leasedConnection) Query(query SQLQueryIfc, args ...interface{}) (*db.Rows, error) {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil, errNoLease() }
	return r.pooledConnection.Query(query, args...)
}

func (r *leasedConnection) QueryRow(query SQLQueryIfc, args ...interface{}) *db.Row {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil }
	return r.pooledConnection.QueryRow(query, args...)
}

// -------------------------------------------------------------------------------------------------
// leasedConnection
// -------------------------------------------------------------------------------------------------

func errNoLease() error {
	return fmt.Errorf("No Leased Connection!")
}

