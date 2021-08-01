package mysql

/*
DB Manager for MySQL - manages a set of named (keyed) mysql database connections / pools
*/

import (
	"fmt"
	"sync"
)

type ManagerIfc interface {
	// Public interface
	NewConnectionPool(dsn string) DBKeyIfc
	CloseConnectionPool(dbKey DBKeyIfc)
	GetConnection(dbKey DBKeyIfc) LeasedConnectionIfc

	// Private interface
	getConnectionPool(dbKey DBKeyIfc) ConnectionPoolIfc
}

// Set of connections, keyed on DSN
type manager struct {
	connectionPools		map[string]ConnectionPoolIfc
	mutex			sync.Mutex
}

// Make a new one of these!
func NewManager() ManagerIfc {
	return &manager{
		connectionPools: make(map[string]ConnectionPoolIfc),
	}
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (mgr *manager) NewConnectionPool(dsn string) DBKeyIfc {
	dbKey := NewDBKeyFromDSN(dsn)
	(*mgr).mutex.Lock()
	defer (*mgr).mutex.Unlock()
	(*mgr).connectionPools[dbKey.GetKey()] = NewConnectionPool(dsn)
	return dbKey
}

func (mgr *manager) GetConnection(dbKey DBKeyIfc) LeasedConnectionIfc {
	connPool := mgr.getConnectionPool(dbKey)
	if nil == connPool { return nil }
	conn, err := connPool.GetConnection()
	if nil != err { fmt.Println("error: %s", err.Error()) }
	return conn
}

func (mgr *manager) CloseConnectionPool(dbKey DBKeyIfc) {
	(*mgr).mutex.Lock()
	defer (*mgr).mutex.Unlock()
	connPool := mgr.getConnectionPool(dbKey)
	if nil == connPool { return }
	connPool.ClosePool()
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Get the connection pool for the specified key
func (mgr *manager) getConnectionPool(dbKey DBKeyIfc) ConnectionPoolIfc {
	if connPool, ok := mgr.connectionPools[dbKey.GetKey()]; ok {
		return connPool
	}
	return nil
}

