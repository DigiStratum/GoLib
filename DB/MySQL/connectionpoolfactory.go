package mysql

import (
	"fmt"

	cfg "github.com/DigiStratum/GoLib/Config"
	"github.com/DigiStratum/GoLib/Dependencies"
	"github.com/DigiStratum/GoLib/DB"
)

type ConnectionPoolFactoryIfc interface {
	NewConnectionPool(dsn db.DSN) *ConnectionPool
}

type ConnectionPoolFactory struct {
	// Note: intentionally captured as supplied interfaces to pass through unmodified...
	config		cfg.ConfigIfc
	dependencies	dependencies.DependenciesIfc
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewConnectionPoolFactory() *ConnectionPoolFactory {
	return &ConnectionPoolFactory{}
}

// -------------------------------------------------------------------------------------------------
// ConnectionPoolFactoryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r ConnectionPoolFactory) NewConnectionPool(dsn db.DSN) (*ConnectionPool, error) {
	// TODO: Match up supplied DSN with configured pattern validator matching this connpool
	connectionPool := NewConnectionPool(dsn)

	// Optional Configuration
	if nil != r.config {
		err := r.configure(connectionPool)
		if nil != err { return nil, err }
	}

	// Optional Dependency Injection
	if nil != r.dependencies {
		if err := r.injectDependencies(connectionPool); nil != err { return nil, err }
	}

	return connectionPool, nil
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPoolFactory) Configure(config cfg.ConfigIfc) error {
	// TODO: Check against ConfigurationPool whether this config is acceptable
	r.config = config
	return nil
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPoolFactory) InjectDependencies(deps dependencies.DependenciesIfc) error {
	// TODO: Check against ConfigurationPool whether these dependencies are acceptable
	r.dependencies = deps
	return nil
}

// -------------------------------------------------------------------------------------------------
// ConnectionPoolFactory Private Implementation
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPoolFactory) configure(connectionPool ConnectionPoolIfc) error {
	if configurableConnectionPool, ok := connectionPool.(cfg.ConfigurableIfc); ok {
		return configurableConnectionPool.Configure(r.config)
	}
	return fmt.Errorf("ConnectionPool is not Configurable? Strange!")
}

func (r *ConnectionPoolFactory) injectDependencies(connectionPool ConnectionPoolIfc) error {
	if injectableConnectionPool, ok := connectionPool.(dependencies.DependencyInjectableIfc); ok {
		return injectableConnectionPool.InjectDependencies(r.dependencies)
	}
	return fmt.Errorf("ConnectionPool is not DependencyInjectable? Strange!")
}
