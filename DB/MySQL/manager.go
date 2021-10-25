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
	GetConnection(dbKey DBKeyIfc) LeasedConnectionIfc
}

type Manager struct {
	connectionPools		map[string]ConnectionPoolIfc // Set of connections, keyed on DSN
	mutex			sync.Mutex
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

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
	if nil != err { fmt.Printf("error: %s\n", err.Error()) }
	return conn
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
