package connection

/*
A Leased Connection wraps a pooled DB connection with an internally managed lease key. If we lose our lease due to
idle/inactivity or some other reason, we either need a way to recover a working connection automatically, or our
DB access will be gone and we will need to get a new connection. Otherwise our interface here is a mirror of
connection, with each method being a pass-through based on keyed access to the underlying connection.
*/

import (
	"errors"

	query "github.com/DigiStratum/GoLib/DB/MySQL"
)

type LeasedConnectionIfc interface {
	NewQuery(query string) (QueryIfc, error)
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

func (lc *leasedConnection) Connect() error {
	// NOP - Leased connections are not allowed to tamper with the connection state
	return errors.New("Leased connection - no state changes allowed")
}

func (lc *leasedConnection) Disconnect() {
	// NOP - Leased connections are not allowed to tamper with the connection state
}

func (lc *leasedConnection) Reconnect() {
	// NOP - Leased connections are not allowed to tamper with the connection state
}

func (lc *leasedConnection) GetConnection() *db.DB {
	if ! lc.checkLease() { return nil }
	return (*lc).pooledConnection.GetConnection((*lc).leaseKey)
}

func (lc *leasedConnection) NewQuery(qry string) (QueryIfc, error) {
	if ! lc.checkLease() {
		return nil, errors.New("No Leased Connection!")
	}
	return query.NewQuery(lc, qry string)
}

// TODO: Move this to pooledConnection so that it has the power of accept/reject instead of self-policing here?
// Check whether we still hold the lease on our connection
func (lc *leasedConnection) checkLease() bool {
	if nil == (*lc).pooledConnection { return false }
	if (*lc).pooledConnection.MatchesLeaseKey((*lc).leaseKey) { return true }
	// The lease has changed! Drop the connection to eliminate any chance of further inactivity
	(*lc).pooledConnection = nil
	return false
}