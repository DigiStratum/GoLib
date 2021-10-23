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
	// Embed Transaction support to this interface
	ConnectionCommonIfc

	Release() error
}

type LeasedConnection struct {
	pooledConnection	PooledConnectionIfc
	leaseKey		int64
	errNoLease		error
}

func NewLeasedConnection(pooledConnection PooledConnectionIfc, leaseKey int64) *LeasedConnection {
	pooledConnection.Lease(leaseKey)
	lc := LeasedConnection{
		pooledConnection:	pooledConnection,
		leaseKey:		leaseKey,
		errNoLease:		fmt.Errorf("No Leased Connection!"),
	}
	return &lc
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *LeasedConnection) Release() error {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return r.errNoLease }
	if err := r.pooledConnection.Release(); nil != err { return err }
	r.leaseKey = 0
	return nil
}

// -------------------------------------------------------------------------------------------------
// ConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r LeasedConnection) IsConnected() bool {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return false }
	return r.pooledConnection.IsConnected()
}

//  Leased connections are not allowed to tamper with the connection lifecycle
func (r *LeasedConnection) Disconnect() { }
func (r *LeasedConnection) Reconnect() { }
func (r *LeasedConnection) Connect() error {
	return fmt.Errorf("Leased connection - no state changes allowed")
}

func (r LeasedConnection) InTransaction() bool {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return false }
	return r.pooledConnection.InTransaction()
}

func (r *LeasedConnection) Begin() error {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return r.errNoLease }
	return r.pooledConnection.Begin()
}

func (r *LeasedConnection) NewQuery(qry string) (QueryIfc, error) {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil, fmt.Errorf("No Leased Connection!") }
	// Feed NewQuery() our LeasedConnection so it doesn't have direct access to underlying pooledConnection
	return NewQuery(r, qry)
}

func (r LeasedConnection) Commit() error {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return r.errNoLease }
	return r.pooledConnection.Commit()
}

func (r LeasedConnection) Rollback() error {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return r.errNoLease }
	return r.pooledConnection.Rollback()
}

func (r LeasedConnection) Exec(query string, args ...interface{}) (db.Result, error) {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil, r.errNoLease }
	return r.pooledConnection.Exec(query, args...)
}

func (r LeasedConnection) Query(query string, args ...interface{}) (*db.Rows, error) {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil, r.errNoLease }
	return r.pooledConnection.Query(query, args...)
}

func (r LeasedConnection) QueryRow(query string, args ...interface{}) *db.Row {
	if ! r.pooledConnection.MatchesLeaseKey(r.leaseKey) { return nil }
	return r.pooledConnection.QueryRow(query, args...)
}
