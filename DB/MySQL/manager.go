package mysql

/*
DB Manager for MySQL - manages a set of named (keyed) mysql database connections / pools
*/

import (
	"io"
	"fmt"
	"sync"

	"github.com/DigiStratum/GoLib/Dependencies"
	"github.com/DigiStratum/GoLib/DB"
)

type ManagerIfc interface {
	CloseConnectionPool(dbKey DBKeyIfc)
	GetConnection(dbKey DBKeyIfc) *LeasedConnection
}

type Manager struct {
	connectionPoolFactory	ConnectionPoolFactoryIfc
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

func (r *Manager) InjectDependencies(deps dependencies.DependenciesIfc) error {
	if nil == deps { return fmt.Errorf("Dependencies were nil") }

	depName := "connectionPoolFactory"
	if ! deps.Has(depName) { return fmt.Errorf("Missing Dependency: %s", depName) }
	dep := deps.Get(depName)
	if nil == dep { return fmt.Errorf("Dependency was nil: %s", depName) }
	connectionPoolFactory, ok := dep.(ConnectionPoolFactoryIfc)
	if ! ok { return fmt.Errorf("Dependency was nil: %s", depName) }
	r.connectionPoolFactory = connectionPoolFactory
	return nil
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Manager) OpenConnectionPool(dsn db.DSN) error {
	if nil == r.connectionPoolFactory { return fmt.Errorf("Missing Dependency: connectionPoolFactory") }
	return nil
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
	key := dbKey.GetKey()
	// If we already have a ConnectionPool for this DBKey...
	if connPool, ok := r.connectionPools[key]; ok { return connPool }
	return nil
}
