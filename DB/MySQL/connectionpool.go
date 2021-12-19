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
	GetConnection() (*LeasedConnection, error)
	Release(leaseKey int64) error
	GetMaxIdle() int
	Close() error
}

type ConnectionPool struct {
	di 			*dependencies.DependencyInjected
	configured		bool
	connectionFactory	db.ConnectionFactoryIfc
	dsn			db.DSN
	minConnections		int
	maxConnections		int
	maxIdle			int
	connections		[]*PooledConnection
	leasedConnections	*LeasedConnections
	mutex			sync.Mutex
}

const DEFAULT_MIN_CONNECTIONS = 1
const DEFAULT_MAX_CONNECTIONS = 1
const DEFAULT_MAX_IDLE = 60

// Make a new one of these
func NewConnectionPool(dsn db.DSN) *ConnectionPool {
	cp := ConnectionPool{
		configured:		false,
		dsn:			dsn,
		minConnections:		DEFAULT_MIN_CONNECTIONS,
		maxConnections:		DEFAULT_MAX_CONNECTIONS,
		maxIdle:		DEFAULT_MAX_IDLE,
		connections:		make([]*PooledConnection, 0, DEFAULT_MAX_CONNECTIONS),
		leasedConnections:	NewLeasedConnections(),
	}

	return &cp
}

// -------------------------------------------------------------------------------------------------
// DependencyInjectableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPool) InjectDependencies(deps dependencies.DependenciesIfc) error {

	// Validate Dependencies
	if nil == deps { return fmt.Errorf("No dependencies provided!") }
	r.di = dependencies.NewDependencyInjected(deps)
	requiredDeps := []string{ "connectionFactory" }
	if ! r.di.SetRequired(&requiredDeps).IsValid() { return r.di.GetValidationError() }

	// Iterate over dependencies and assign each to a local property
	it := deps.GetIterator()
	for dd := it(); nil != dd; dd = it() {
		name, dep := dd.GetDep()
		var ok bool
		switch name {
			case "connectionFactory":
				if r.connectionFactory, ok = dep.(db.ConnectionFactoryIfc); ok { continue }
			default:
				// Ignore extra dependencies - not a fatal error
		}
		return fmt.Errorf("Dependency was wrong type: %s", name )
	}

	// Set up the first resource
	r.establishMinConnections()

	return nil
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Optionally accept overrides for defaults in configuration
func (r *ConnectionPool) Configure(config cfg.ConfigIfc) error {
	if (nil == r) || (! r.di.IsValid()) { return fmt.Errorf("Not ready") }

	// If we have already been configured, do not accept a second configuration
	if r.configured { return nil }

	// Capture optional configs
	it := config.GetIterator()
	for kvpi := it(); nil != kvpi; kvpi = it() {
		kvp, ok := kvpi.(*hashmap.KeyValuePair)
		if ! ok { continue }
		switch kvp.Key {
			case "min_connections":
				if value := config.GetInt64(kvp.Key); nil != value {
					r.configureMinConnections(int(*value))
				}

			case "max_connections":
				if value := config.GetInt64(kvp.Key); nil != value {
					r.configureMaxConnections(int(*value))
				}

			case "max_idle":
				if value := config.GetInt64(kvp.Key); nil != value {
					r.configureMaxIdle(int(*value))
				}

			default:
				return fmt.Errorf("Unknown configuration key: '%s'", kvp.Key)
		}
	}
	r.configured = true

	// If the new Max increases from default...
	if cap(r.connections) < r.maxConnections {
		// Increase connection pool capacity from default to the new max_connections
		// ref: https://blog.golang.org/slices-intro
		nc := make([]*PooledConnection, len(r.connections), r.maxConnections)
		copy(nc, r.connections)
		r.connections = nc
	}

	// If the minimum resource count has gone up, fill up the difference
	r.establishMinConnections()

	return nil
}

func (r *ConnectionPool) configureMinConnections(value int) {
	r.minConnections = value
	if r.minConnections < 1 { r.minConnections = 1 }
	// If Min pushed above Max, then push Max up
	if r.maxConnections < r.minConnections { r.maxConnections = r.minConnections }
}

func (r *ConnectionPool) configureMaxConnections(value int) {
	r.maxConnections = value
	if r.maxConnections < 1 { r.maxConnections = 1 }
	// If Max dropped below Min, then push Min down
	if r.maxConnections < r.minConnections { r.minConnections = r.maxConnections }
}

func (r *ConnectionPool) configureMaxIdle(value int) {
	r.maxIdle = value
	// Max seconds since lastActiveAt for leased connections: 1 <= max_idle
	if r.maxIdle < 1 { r.maxIdle = 1 }
}

// -------------------------------------------------------------------------------------------------
// ConnectionPoolIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Request a connection from the pool using multiple approaches
func (r *ConnectionPool) GetConnection() (*LeasedConnection, error) {
	var connection *PooledConnection

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
// TODO: Denote who need to know this & why...
func (r ConnectionPool) GetMaxIdle() int {
	return r.maxIdle
}

// -------------------------------------------------------------------------------------------------
// io.Closer Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPool) Close() error {
	r.mutex.Lock(); defer r.mutex.Unlock()

	// Wipe the DSN to prevent new connections from being established
	r.dsn = db.DSN{}
	r.connectionFactory = nil

	// Drop all open leases
	r.leasedConnections = NewLeasedConnections()

	// Disconnect all open connections
	errors := 0
	for _, pooledConnection := range (*r).connections {
		err := r.closePooledConnection(pooledConnection)
		if nil != err { errors++ }
	}
	if 0 == errors { return nil }
	return fmt.Errorf("There were %d errors closing PooledConnection(s)", errors)
}

// -------------------------------------------------------------------------------------------------
// ConnectionPool (Package-Private) Implementation
// -------------------------------------------------------------------------------------------------

func (r *ConnectionPool) findAvailableConnection() *PooledConnection {
	for _, connection := range (*r).connections {
		if ! connection.IsLeased() { return connection }
	}
	return nil
}

// TODO: Pass errors back to caller and on up the chain for visibility/logging
func (r *ConnectionPool) createNewConnection() *PooledConnection {
	// if we are at capacity, then we can't create a new connection
	if len(r.connections) >= cap(r.connections) { return nil }
	// We're under capacity so should be able to add a new connection
	conn, err := r.connectionFactory.NewConnection(r.dsn)
	if nil != err { return nil }
	// Wrap the raw connection into a Connection
	newConnection, err := NewConnection(conn)
	if nil != err { return nil }
	// Wrap the new connection into a pooled connection to maintain state
	newPooledConnection, err := NewPooledConnection(newConnection, r)
	if nil == err { r.connections = append(r.connections, newPooledConnection) }
	return newPooledConnection // nil if there was an error
}

func (r ConnectionPool) findExpiredLeaseConnection() *PooledConnection {
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

func (r *ConnectionPool) closePooledConnection(pooledConnection PooledConnectionIfc) error {
	if closeableConnection, ok := pooledConnection.(io.Closer); ok {
		return closeableConnection.Close()
	}
	return fmt.Errorf("PooledConnection not closable")
}
