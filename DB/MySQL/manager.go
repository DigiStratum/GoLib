package mysql

/*
DB Manager for MySQL - manages a set of named (keyed) mysql database connections / pools
*/

import (
	"io"
	"fmt"
	"sync"
)

type ManagerIfc interface {
	NewConnectionPool(dsn string) DBKeyIfc
	CloseConnectionPool(dbKey DBKeyIfc)
	GetConnection(dbKey DBKeyIfc) *LeasedConnection
}

type Manager struct {
	connectionPools		map[string]*ConnectionPool // Set of connections, keyed on DSN
	mutex			sync.Mutex
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewManager() *Manager {
	return &Manager{
		connectionPools: make(map[string]*ConnectionPool),
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

func (r *Manager) GetConnection(dbKey DBKeyIfc) (*LeasedConnection, error) {
	connPool := r.getConnectionPool(dbKey)
	if nil == connPool { return nil, fmt.Errorf("No connection pool for this dbKey") }
	return connPool.GetConnection()
}

func (r *Manager) CloseConnectionPool(dbKey DBKeyIfc) {
	r.mutex.Lock();	defer r.mutex.Unlock()
	connPool := r.getConnectionPool(dbKey)
	if nil == connPool { return }
	if closeablePool, ok := connPool.(io.Closer); ok {
		closeablePool.Close()
	}
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
