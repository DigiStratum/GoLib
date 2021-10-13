package mysql

/*
This Database Connection Pool establishes one or more persistent connections to a MySQL database given a configured DSN.
When a consumer requests a connection from the pool, we will attempt to provide one using multiple approaches, in the
following order of priority:

1) An already established connection that is available (not leased out to another consumer)
2) A newly created connection if the total number of connections is below the max
3) An already established connection that is leased out, but past the lease time for idle connections

When a consumer returns a leased connection, we should decide whether we should just close the connection, or close and
reopen it, or just cycle it back into the pool of available connections. This can be based on any number of factors,
including overall activity, maximum age of established connections, any sort of flag indicating that the connection is
"dirty" (e.g some change has been made to transaction isolation mode, etc.) We could also take this opportunity to audit
all of the open connections to see if any others have been sitting open and idle too long and need similar treatment.

TODO:
 * Change mutex to go-routine+channel for multithreaded orchestration
 * Look more closely at sql.DB which is a connection pool natively; would it give us enough control/visibility over state?
   * ref: https://pkg.go.dev/database/sql#DB

*/

import (
	"io"
	"fmt"
	"errors"
	"sync"
	cfg "github.com/DigiStratum/GoLib/Config"
	"github.com/DigiStratum/GoLib/Dependencies"
	"github.com/DigiStratum/GoLib/Data/hashmap"
	"github.com/DigiStratum/GoLib/DB"
)

// A Connection Pool to maintain a set of one or more persistent connections to a MySQL database
type ConnectionPoolIfc interface {
	GetConnection() (LeasedConnectionIfc, error)
	Release(leaseKey int64) error
	GetMaxIdle() int
}

type ConnectionPool struct {
	configured		bool
	dbConnectionFactory	db.DBConnectionFactoryIfc
	dsn			string
	minConnections		int
	maxConnections		int
	maxIdle			int
	connections		[]PooledConnectionIfc
	leasedConnections	LeasedConnectionsIfc
	mutex			sync.Mutex
}

const DEFAULT_MIN_CONNECTIONS = 1
const DEFAULT_MAX_CONNECTIONS = 1
const DEFAULT_MAX_IDLE = 60

// Make a new one of these
func NewConnectionPool(dsn string) *ConnectionPool {
	cp := ConnectionPool{
		configured:		false,
		dsn:			dsn,
		minConnections:		DEFAULT_MIN_CONNECTIONS,
		maxConnections:		DEFAULT_MAX_CONNECTIONS,
		maxIdle:		DEFAULT_MAX_IDLE,
		connections:		make([]PooledConnectionIfc, 0, DEFAULT_MAX_CONNECTIONS),
		leasedConnections:	NewLeasedConnections(),
	}
	// Set up the first resource
	cp.establishMinConnections()

	return &cp
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
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Optionally accept overrides for defaults in configuration
func (r *ConnectionPool) Configure(config cfg.ConfigIfc) error {
	// If we have already been configured, do not accept a second configuration
	if r.configured { return nil }

	// Capture optional configs
	it := config.GetIterator()
	for kvpi := it(); nil != kvpi; kvpi = it() {
		kvp, ok := kvpi.(hashmap.KeyValuePair)
		if ! ok { continue }
		switch kvp.Key {
			case "min_connections":
				value := config.GetInt64("min_connections")
				if nil == value { break }
				// Set the new Min (cannot be < 1)
				r.minConnections = int(*value)
				if r.minConnections < 1 { r.minConnections = 1 }
				// If Min pushed above Max, then push Max up
				if r.maxConnections < r.minConnections { r.maxConnections = r.minConnections }

			case "max_connections":
				value := config.GetInt64("max_connections")
				if nil == value { break }
				// Set the new Max (cannot be < 1)
				r.maxConnections = int(*value)
				if r.maxConnections < 1 { r.maxConnections = 1 }
				// If Max dropped below Min, then push Min down
				if r.maxConnections < r.minConnections { r.minConnections = r.maxConnections }


			case "max_idle":
				value := config.GetInt64("max_idle")
				if nil == value { break }
				r.maxIdle = int(*value)
				// Max seconds since lastActiveAt for leased connections: 1 <= max_idle
				if r.maxIdle < 1 { r.maxIdle = 1 }

			default:
				return errors.New(fmt.Sprintf("Unknown configuration key: '%s'", kvp.Key))
		}
	}
	r.configured = true

	// If the new Max increases from default...
	if cap(r.connections) < r.maxConnections {
		// Increase connection pool capacity from default to the new max_connections
		// ref: https://blog.golang.org/slices-intro
		nc := make([]PooledConnectionIfc, len(r.connections), r.maxConnections)
		copy(nc, r.connections)
		r.connections = nc
	}

	// If the minimum resource count has gone up, fill up the difference
	r.establishMinConnections()

	return nil
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Request a connection from the pool using multiple approaches
func (r *ConnectionPool) GetConnection() (LeasedConnectionIfc, error) {
	var connection PooledConnectionIfc

	r.mutex.Lock(); defer r.mutex.Unlock()

	// 1) An already established connection that is available (not leased out to another consumer)
	connection = r.findAvailableConnection()

	// 2) A newly created connection if the total number of connections is below the max
	if nil == connection { connection = r.createNewConnection() }

	// 3) An already established connection that is leased out, but past the lease time for idle connections
	if nil == connection { connection = r.findExpiredLeaseConnection() }

	if nil == connection { return nil, errors.New("No available pooled connections!") }

	// Establish a lease for this connection which is ours now
	return r.leasedConnections.GetLeaseForConnection(connection), nil
}

func (r *ConnectionPool) Release(leaseKey int64) error {
	r.mutex.Lock(); defer r.mutex.Unlock()
	if ! r.leasedConnections.Release(leaseKey) {
		return errors.New(fmt.Sprintf("Pool contains no lease key = '%d'", leaseKey))
	}
	return nil
}

// Max Idle has a default, but may be overridden by configuration; this gets access to the current setting value
func (r ConnectionPool) GetMaxIdle() int {
	return r.maxIdle
}

// -------------------------------------------------------------------------------------------------
// io.Closer Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPool) Close() error {
	r.mutex.Lock(); defer r.mutex.Unlock()

	// Wipe the DSN to prevent new connections from being established
	r.dsn = ""
	r.dbConnectionFactory = nil

	// Drop all open leases
	r.leasedConnections = NewLeasedConnections()

	// Disconnect all open connections
	for _, pooledConnection := range (*r).connections {
		if closeableConnection, ok := pooledConnection.(io.Closer); ok {
			closeableConnection.Close()
		}
	}

	return nil
}

// -------------------------------------------------------------------------------------------------
// Configurable (Package-Private) Implementation
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPool) findAvailableConnection() PooledConnectionIfc {
	for _, connection := range (*r).connections {
		if ! connection.IsLeased() { return connection }
	}
	return nil
}

// TODO: Pass errors back to caller and on up the chain for visibility/logging
func (r *ConnectionPool) createNewConnection() PooledConnectionIfc {
	// if we are at capacity, then we can't create a new connection
	if len(r.connections) >= cap(r.connections) { return nil }
	// We're under capacity so should be able to add a new connection
	conn, err := r.dbConnectionFactory.NewConnection(r.dsn)
	if nil != err { return nil }
	// Wrap the raw connection into a Connection
	newConnection, err := NewConnection(conn)
	if nil != err { return nil }
	// Wrap the new connection into a pooled connection to maintain state
	newPooledConnection, err := NewPooledConnection(newConnection, r)
	if nil == err { r.connections = append(r.connections, newPooledConnection) }
	return newPooledConnection // nil if there was an error
}

func (r ConnectionPool) findExpiredLeaseConnection() PooledConnectionIfc {
	for _, connection := range r.connections {
		if connection.IsLeased() && connection.IsExpired() { return connection }
	}
	return nil
}

func (r *ConnectionPool) establishMinConnections() {
	// If the minimum resource count has gone up, fill up the difference
	for ci := len(r.connections); ci < r.minConnections; ci++ {
		_ = r.createNewConnection()
	}
}
