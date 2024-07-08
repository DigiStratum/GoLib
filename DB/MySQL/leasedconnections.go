package mysql

import (
	"sync"
)

type LeasedConnectionsIfc interface {
	// Public interface
	GetLeaseForConnection(connection PooledConnectionIfc) *leasedConnection
	Release(leaseKey int64) bool

	// Private interface
	getNewLeaseKey() *int64
}

type leasedConnections struct {
	leases		map[int64]*leasedConnection
	nextLeaseKey	int64
	mutex		sync.Mutex
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewLeasedConnections() *leasedConnections {
	lc := leasedConnections{
		leases:		make(map[int64]*leasedConnection),
		nextLeaseKey:	0,
	}
	return &lc
}

// -------------------------------------------------------------------------------------------------
// LeasedConnectionsIfc
// -------------------------------------------------------------------------------------------------

func (r *leasedConnections) GetLeaseForConnection(connection PooledConnectionIfc) *leasedConnection {
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

func (r *leasedConnections) Release(leaseKey int64) bool {
	r.mutex.Lock(); defer r.mutex.Unlock()
	if ! r.leaseExists(leaseKey) { return false }
	delete(r.leases, leaseKey)
	return true
}

// -------------------------------------------------------------------------------------------------
// leasedConnections
// -------------------------------------------------------------------------------------------------

// Is there a Lease on record now with this key?
func (r *leasedConnections) leaseExists(leaseKey int64) bool {
	_, ok := r.leases[leaseKey]
	return ok
}

// Return the next available Lease Key, or nil on failure
func (r *leasedConnections) getNewLeaseKey() *int64 {
	// If we don't get this even on the first attempt, then something is wrong... but just in case...
	for attempts := 0; attempts < 100; attempts++ {
		r.nextLeaseKey++
		if ! r.leaseExists(r.nextLeaseKey) { return &r.nextLeaseKey }
	}
	return nil
}
