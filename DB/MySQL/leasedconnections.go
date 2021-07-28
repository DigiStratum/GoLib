package mysql

type LeasedConnectionsIfc interface {
	// Public interface
	GetLeaseForConnection(connection PooledConnectionIfc) LeasedConnectionIfc
	// Private interface
	leaseExists(leaseKey int64) bool
	getNewLeaseKey() *int64
}

type leasedConnections struct {
	leases		map[int64]LeasedConnectionIfc
	nextLeaseKey	int64
}

func NewLeasedConnectionsIfc() LeasedConnectionsIfc {
	lc := leasedConnections{
		leases:		make(map[int64]LeasedConnectionIfc),
		nextLeaseKey:	0,
	}
	return &lc
}

func (lc *leasedConnections) GetLeaseForConnection(connection PooledConnectionIfc) LeasedConnectionIfc {
	// TODO: Implement!
	if ptrLeaseKey := lc.getNewLeaseKey(); nil != ptrLeaseKey {
		leasedConnection := NewLeasedConnection(connection, *ptrLeaseKey)
		if nil != leasedConnection {
			(*lc).leases[*ptrLeaseKey] = leasedConnection
			return leasedConnection
		}
	}
	return nil
}

// Is there a Lease on record now with this key?
func (lc *leasedConnections) leaseExists(leaseKey int64) bool {
	_, ok := (*lc).leases[leaseKey]
	return ok
}

// Return the next available Lease Key, or nil on failure
func (lc *leasedConnections) getNewLeaseKey() *int64 {
	attempts := 0
	// If we don't get this even on the first attempt, then something is wrong... but just in case...
	for attempts := 0; attempts < 100; attempts++ {
		(*lc).nextLeaseKey++
		if ! lc.leaseExists((*lc).nextLeaseKey) { return &(*lc).nextLeaseKey }
	}
	return nil
}
