package connection

import (
	"time"
)

// A Pooled Connection wraps a raw DB connection with additional metadata to manage leasing
// We are not exporting this because it is only important to the connection package internal implementation
type PooledConnectionIfc interface {
	GetConnection(leaseKey int64) ConnectionIfc
	IsLeased() bool
	MatchesLeaseKey(leaseKey int64) bool
	Lease(leaseKey int64)
	Release()
	Touch()
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
}

// -------------------------------------------------------------------------------------------------
// pooledConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------
func (pc *pooledConnection) GetConnection(leaseKey int64) ConnectionIfc {
	if ! pc.MatchesLeaseKey(leaseKey) { return nil }
	return (*pc).connection
}

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