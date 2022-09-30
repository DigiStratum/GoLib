package mysql

import (
	"sync"
)

type LeasedConnectionsIfc interface {
	// Public interface
	GetLeaseForConnection(connection PooledConnectionIfc) *LeasedConnection
	Release(leaseKey int64) bool

	// Private interface
	getNewLeaseKey() *int64
}

type LeasedConnections struct {
	leases		map[int64]*LeasedConnection
	nextLeaseKey	int64
	mutex		sync.Mutex
}

func NewLeasedConnections() *LeasedConnections {
	lc := LeasedConnections{
		leases:		make(map[int64]*LeasedConnection),
		nextLeaseKey:	0,
	}
	return &lc
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionsIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *LeasedConnections) GetLeaseForConnection(connection PooledConnectionIfc) *LeasedConnection {
	if nil == connection { return nil }
	r.mutex.Lock(); defer r.mutex.Unlock()
	// Get a new lease key...
	if ptrLeaseKey := r.getNewLeaseKey(); nil != ptrLeaseKey {
		// Set up a new lease for it...
		leasedConnection := NewLeasedConnection(connection, *ptrLeaseKey)
		r.leases[*ptrLeaseKey] = leasedConnection
		return leasedConnection
	}
	return nil
}

func (r *LeasedConnections) Release(leaseKey int64) bool {
	r.mutex.Lock(); defer r.mutex.Unlock()
	if ! r.leaseExists(leaseKey) { return false }
	delete(r.leases, leaseKey)
	return true
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionsIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Is there a Lease on record now with this key?
func (r *LeasedConnections) leaseExists(leaseKey int64) bool {
	_, ok := r.leases[leaseKey]
	return ok
}

// Return the next available Lease Key, or nil on failure
func (r *LeasedConnections) getNewLeaseKey() *int64 {
	// If we don't get this even on the first attempt, then something is wrong... but just in case...
	for attempts := 0; attempts < 100; attempts++ {
		r.nextLeaseKey++
		if ! r.leaseExists(r.nextLeaseKey) { return &r.nextLeaseKey }
	}
	return nil
}
