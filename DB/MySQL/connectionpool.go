package mysql

import (
	"errors"

	lib "github.com/DigiStratum/GoLib"
)

type ConnectionPoolIfc interface {
}

type pooledConnectionIfc

type pooledConnection struct {
	connection	ConnectionIfc
}

type connectionPool struct {
	dsn		string
	minConnections	int64
	maxConnections	int64
	maxIdle		int64
	connections	[]pooledConnectionIfc
}

func NewConnectionPool() ConnectionPoolIfc {
	cp := connectionPool{}
	return &cp
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------
func (cp *connectionPool) Configure(config ConfigIfc) error {
	requiredConfigs := []string{ "dsn", "min_connections", "max_connections", "max_idle" }
	if ! config.HasAll(&requiredConfigs) { return errors.New("Missing required config") }
	(*cp).dsn = config.Get("dsn")
	(*cp).minConnections = config.GetInt64("min_connections")
	(*cp).maxConnections = config.GetInt64("max_connections")
	(*cp).maxIdle = config.GetInt64("max_idle")

	// Set up the connection pool to hold max connections
	(*cp).connections = make([]pooledConnectionIfc, (*cp).maxConnections)

	// TODO: create the min connections and put them into the pool
	// TODO: track how many initialized connections we have in the pool
}
