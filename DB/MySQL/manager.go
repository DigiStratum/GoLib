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
	NewConnectionPool(dsn string, config cfg.ConfigIfc) *DBKey
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
// DependencyInjectableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPool) InjectDependencies(deps dependencies.DependenciesIfc) error {
	if nil == deps { return fmt.Errorf("Dependencies were nil") }

	depName := "dbConnectionFactory"
	if ! deps.Has(depName) { return fmt.Errorf("Missing Dependency: %s", depName) }
	dep := deps.Get(depName)
	if nil == dep { return fmt.Errorf("Dependency was nil: %s", depName) }
	dbConnectionFactory, ok := dep.(db.DBConnectionFactoryIfc)
	if ! ok { return fmt.Errorf("Dependency was nil: %s", depName) }
	r.dbConnectionFactory = dbConnectionFactory
	return nil
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Manager) NewConnectionPool(dsn string, config cfg.ConfigIfc) *DBKey {
	dbKey := NewDBKeyFromDSN(dsn)
	r.mutex.Lock();	defer r.mutex.Unlock()
	connectionPool := NewConnectionPool(dsn)
	if nil != config {
		err := connectionPool.Configure(config)
		if nil != err { return nil, err }
	}
	r.connectionPools[dbKey.GetKey()] = connectionPool
	return dbKey, nil
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
