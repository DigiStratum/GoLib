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
 * Why is this in the MySQL package? There doesn't seem to be anything MySQL-specific. Move up to DB package if possible!
 * Refactor errors to use Logger.Error()

 FIXME:
  * Double check that this does not return a connection from GetConnection() if not IsStarted(); there were indicators of
    this in at least one unit test that was missing a Start() operation, but the error was buried when we added Start()
*/

import (
	"io"
	"fmt"
	"errors"
	"sync"
	"strconv"

	"github.com/DigiStratum/GoLib/Process/startable"
	cfg "github.com/DigiStratum/GoLib/Config"
	dep "github.com/DigiStratum/GoLib/Dependencies"
	"github.com/DigiStratum/GoLib/DB"
)

// A Connection Pool to maintain a set of one or more persistent connections to a MySQL database
type ConnectionPoolIfc interface {
	// Embedded interface(s)
	startable.StartableIfc
	dep.DependencyInjectableIfc
	cfg.ConfigurableIfc

	// Our own interface
	GetConnection() (*leasedConnection, error)
	Release(leaseKey int64) error
	GetMaxIdle() int
	Close() error
}

type connectionPool struct {
	// Embedded struct(s)
	*startable.Startable
	*dep.DependencyInjectable
	*cfg.Configurable

	// Our own properties
	connectionFactory		db.ConnectionFactoryIfc
	dsn				db.DSN
	minConnections			int
	maxConnections			int
	maxIdle				int
	connections			[]*pooledConnection
	leasedConnections		LeasedConnectionsIfc
	mutex				sync.Mutex
}

const DEFAULT_MIN_CONNECTIONS = 1
const DEFAULT_MAX_CONNECTIONS = 1
const DEFAULT_MAX_IDLE = 60

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewConnectionPool(dsn db.DSN) *connectionPool {
	cp := connectionPool{
		// Config Defaults
		minConnections:		DEFAULT_MIN_CONNECTIONS,
		maxConnections:		DEFAULT_MAX_CONNECTIONS,
		maxIdle:		DEFAULT_MAX_IDLE,

		dsn:			dsn,
		connections:		make([]*pooledConnection, 0, DEFAULT_MAX_CONNECTIONS),
		leasedConnections:	NewLeasedConnections(),
	}

	// Declare Dependencies
	cp.DependencyInjectable = dep.NewDependencyInjectable(
		dep.NewDependency("ConnectionFactory").SetRequired().CaptureWith(
			func (instance interface{}) bool {
				var ok bool
				cp.connectionFactory, ok = instance.(db.ConnectionFactoryIfc)
				return ok
			},
		),
	)

	// Declare Configuration
	cp.Configurable = cfg.NewConfigurable(
		cfg.NewConfigItem("min_connections").CaptureWith(cp.captureConfigMinConnections),
		cfg.NewConfigItem("max_connections").CaptureWith(cp.captureConfigMaxConnections),
		cfg.NewConfigItem("max_idle").CaptureWith(cp.captureConfigMaxIdle),
	)

	// Starters
	cp.Startable = startable.NewStartable(
		cp.DependencyInjectable,
		cp.Configurable,
	)

	return &cp
}

func ConnectionPoolFromIfc(i interface{}) (ConnectionPoolIfc, error) {
	if ii, ok := i.(ConnectionPoolIfc); ok { return ii, nil }
	return nil, fmt.Errorf("Does not implement ConnectionPoolIfc")
}


// -------------------------------------------------------------------------------------------------
// ConnectionPoolIfc
// -------------------------------------------------------------------------------------------------

// Request a connection from the pool using multiple approaches
func (r *connectionPool) GetConnection() (*leasedConnection, error) {
	if ! r.Startable.IsStarted() { return nil, fmt.Errorf("ConnectionPool.GetConnection() : Not Started") }
	var connection *pooledConnection

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

func (r *connectionPool) Release(leaseKey int64) error {
	if ! r.Startable.IsStarted() { return fmt.Errorf("ConnectionPool.GetConnection() : Not Started") }
	r.mutex.Lock(); defer r.mutex.Unlock()
	if ! r.leasedConnections.Release(leaseKey) {
		return errors.New(fmt.Sprintf("Pool contains no lease key = '%d'", leaseKey))
	}
	return nil
}

// Max Idle has a default, but may be overridden by configuration; this gets access to the current setting value
// TODO: Denote who need to know this & why...
func (r *connectionPool) GetMaxIdle() int {
	return r.maxIdle
}

// -------------------------------------------------------------------------------------------------
// io.Closer
// -------------------------------------------------------------------------------------------------

// TODO: Tie this into some new StoppableIfc
func (r *connectionPool) Close() error {
	if ! r.Startable.IsStarted() { return fmt.Errorf("ConnectionPool.GetConnection() : Not Started") }
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
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *connectionPool) Start() error {
	if r.Startable.IsStarted() { return nil }

	// If the configured Max increases from default...
	if cap(r.connections) < r.maxConnections {
		// Increase connection pool capacity from default to the new max_connections
		// ref: https://blog.golang.org/slices-intro
		nc := make([]*pooledConnection, len(r.connections), r.maxConnections)
		copy(nc, r.connections)
		r.connections = nc
	}

	// If the minimum resource count has gone up, fill up the difference
	r.establishMinConnections()

	return r.Startable.Start()
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc
// -------------------------------------------------------------------------------------------------

// Optionally accept overrides for defaults in configuration
func (r *connectionPool) Configure(config cfg.ConfigIfc) error {
	// If we have already been configured, do not accept a second configuration
	if r.Startable.IsStarted() { return nil }

	return r.Configurable.Configure(config)
}

// -------------------------------------------------------------------------------------------------
// ConnectionPool
// -------------------------------------------------------------------------------------------------

// Config Capture Funcs
// -----------------------------------------------

func (r *connectionPool) captureConfigMinConnections(value string) error {
	v, err := strconv.Atoi(value)
	if nil != err { return err }
	r.minConnections = v
	if r.minConnections < 1 { r.minConnections = 1 }
	// If Min pushed above Max, then push Max up
	if r.maxConnections < r.minConnections { r.maxConnections = r.minConnections }
	return nil
}

func (r *connectionPool) captureConfigMaxConnections(value string) error {
	v, err := strconv.Atoi(value)
	if nil != err { return err }
	r.maxConnections = v
	if r.maxConnections < 1 { r.maxConnections = 1 }
	// If Max dropped below Min, then push Min down
	if r.maxConnections < r.minConnections { r.minConnections = r.maxConnections }
	return nil
}

func (r *connectionPool) captureConfigMaxIdle(value string) error {
	v, err := strconv.Atoi(value)
	if nil != err { return err }
	r.maxIdle = v
	// Max seconds since lastActiveAt for leased connections: 1 <= max_idle
	if r.maxIdle < 1 { r.maxIdle = 1 }
	return nil
}

// Supporting funcs
// -----------------------------------------------

// TODO: Pass errors back to caller and on up the chain for visibility/logging
func (r *connectionPool) createNewConnection() *pooledConnection {
	if ! r.Startable.IsStarted() { return nil }
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

func (r *connectionPool) findExpiredLeaseConnection() *pooledConnection {
	for _, connection := range r.connections {
		if connection.IsLeased() && connection.IsExpired() { return connection }
	}
	return nil
}

func (r *connectionPool) establishMinConnections() {
	// If the minimum resource count has gone up, fill up the difference
	for ci := len(r.connections); ci < r.minConnections; ci++ {
		_ = r.createNewConnection()
	}
}

func (r *connectionPool) closePooledConnection(pooledConnection PooledConnectionIfc) error {
	if ! r.Startable.IsStarted() { return fmt.Errorf("ConnectionPool.GetConnection() : Not Started") }
	if closeableConnection, ok := pooledConnection.(io.Closer); ok {
		return closeableConnection.Close()
	}
	return fmt.Errorf("PooledConnection not closable")
}

func (r *connectionPool) findAvailableConnection() *pooledConnection {
	for _, connection := range (*r).connections {
		if ! connection.IsLeased() { return connection }
	}
	return nil
}

