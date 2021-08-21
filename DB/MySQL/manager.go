package mysql

/*
DB Manager for MySQL - manages a set of named (keyed) mysql database connections / pools
*/

import (
	"fmt"
	"sync"
)

type ManagerIfc interface {
	NewConnectionPool(dsn string) DBKeyIfc
	CloseConnectionPool(dbKey DBKeyIfc)
	GetConnection(dbKey DBKeyIfc) LeasedConnectionIfc
}

// Set of connections, keyed on DSN
type Manager struct {
	connectionPools		map[string]ConnectionPoolIfc
	mutex			sync.Mutex
}

// Make a new one of these!
func NewManager() *Manager {
	return &Manager{
		connectionPools: make(map[string]ConnectionPoolIfc),
	}
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Manager) NewConnectionPool(dsn string) DBKeyIfc {
	dbKey := NewDBKeyFromDSN(dsn)
	r.mutex.Lock();	defer r.mutex.Unlock()
	r.connectionPools[dbKey.GetKey()] = NewConnectionPool(dsn)
	return dbKey
}

func (r *Manager) GetConnection(dbKey DBKeyIfc) LeasedConnectionIfc {
	connPool := r.getConnectionPool(dbKey)
	if nil == connPool { return nil }
	conn, err := connPool.GetConnection()
	if nil != err { fmt.Println("error: %s", err.Error()) }
	return conn
}

func (r *Manager) CloseConnectionPool(dbKey DBKeyIfc) {
	r.mutex.Lock();	defer r.mutex.Unlock()
	connPool := r.getConnectionPool(dbKey)
	if nil == connPool { return }
	connPool.ClosePool()
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Get the connection pool for the specified key
func (r *Manager) getConnectionPool(dbKey DBKeyIfc) ConnectionPoolIfc {
	if connPool, ok := r.connectionPools[dbKey.GetKey()]; ok {
		return connPool
	}
	return nil
}
