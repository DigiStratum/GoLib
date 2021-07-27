package connection

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

TODO: Consider whether this can be abstracted within reason to support any kind of database. Perhaps all it needs is
one additional layer of abstraction around the database-specific functions, but this can create an annoying requirement
on the consumer side to receive the abstracted interface and have to cast it to a DB-specific interface. We could also
use some sort of generic DB function method that accepts a "command" as a string or constant wherein the generic function
switches and dispatches based on this, however you then have this additional overhead and still would be left with
potentially returning generic/abstract structures that the consumer would have to deal with (unless the result of such
is always a ResultSetIfc (nil or 1+ rows) each with 1+ columns, and an error (nil if none)..?)
*/

import (
	"errors"

	lib "github.com/DigiStratum/GoLib"
)

// A Connection Pool to maintain a set of one or more persistent connections to a MySQL database
type ConnectionPoolIfc interface {
	GetConnection() LeasedConnectionIfc
}

type connectionPool struct {
	configured		bool
	dsn			string
	minConnections		int64
	maxConnections		int64
	maxIdle			int64
	connections		[]pooledConnectionIfc
	leasedConnections	leasedConnectionsIfc
}

const DEFAULT_MIN_CONNECTIONS = 1
const DEFAULT_MAX_CONNECTIONS = 1
const DEFAULT_MAX_IDLE = 1

// Make a new one of these
func NewConnectionPool(dsn string) ConnectionPoolIfc {
	cp := connectionPool{
		configured:		false,
		dsn:			dsn,
		minConnections:		DEFAULT_MIN_CONNECTIONS,
		maxConnections:		DEFAULT_MAX_CONNECTIONS,
		maxIdle:		DEFAULT_MAX_IDLE,
		connections:		make([]pooledConnectionIfc, 0, DEFAULT_MAX_CONNECTIONS),
		leasedConnections:	newLeasedConnectionsIfc(),
	}
	return &cp
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Optionally accept overrides for defaults in configuration
func (cp *connectionPool) Configure(config ConfigIfc) error {
	// If we have already been configured, do not accept a second configuration
	if (*cp).configured { return nil }

	// Capture optional configs
	for kvp := range config.IterateChannel() {
		switch kvp.Key {
			case "min_connections":
				// Set the new Min (cannot be < 1)
				(*cp).minConnections = math.Max(1, config.GetInt64("min_connections"))
				// If Min pushed above Max, then push Max up
				(*cp).maxConnections = math.Max((*cp).minConnections, (*cp).maxConnections)

			case "max_connections":
				// Set the new Max (cannot be < 1)
				(*cp).maxConnections = math.Max(1, config.GetInt64("max_connections"))
				// If Max dropped below Min, then push Min down
				(*cp).minConnections = math.Min((*cp).minConnections, (*cp).maxConnections)

				// If the new Max increases from default...
				if cap((*cp.connections)) < (*cp).maxConnections {
					// Increase connection pool capacity from default to the new max_connections
					// ref: https://blog.golang.org/slices-intro
					nc := make([]pooledConnectionIfc, len((*cp).connections), (*cp).maxConnections)
					copy(nc, (*cp).connections)
					(*cp).connections = nc
				}

			case "max_idle":
				// Max seconds since lastActiveAt for leased connections: 1 <= max_idle
				(*cp).maxIdle = math.Max(1, config.GetInt64("max_idle"))

			default:
				return errors.New(fmt.Sprintf("Unknown configuration key: '%s'", kvp.Key))
		}
	}
	(*cp).configured = true

	return nil
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Request a connection from the pool using multiple approaches
func (cp *connectionPool) GetConnection() LeasedConnectionIfc {
	var connection pooledConnectionIfc
	// 1) An already established connection that is available (not leased out to another consumer)
	connection = cp.findAvailableConnection()

	// 2) A newly created connection if the total number of connections is below the max
	if nil == connection { connection = cp.createNewConnection() }

	// 3) An already established connection that is leased out, but past the lease time for idle connections
	if nil == connection { connection = cp.findExpiredConnection() }

	if nil == connection { return nil }

	// Establish a lease for this connection which is ours now
	return (*cp).leasedConnections.GetLeaseForConnection(connection)
}

func (cp *connectionPool) findAvailableConnection() pooledConnectionIfc {
	for connection := range (*cp).connections {
		if ! connection.IsLeased() { return connection }
	}
	return nil
}

func (cp *connectionPool) createNewConnection() pooledConnectionIfc {
	// if we are at capacity, then we can't create a new connection
	if len((*cp).connections) >= cap((*cp).connections) { return nil }
	// TODO: Implement!
	return nil
}

func (cp *connectionPool) findExpiredConnection() pooledConnectionIfc {
	// TODO: Implement!
	return nil
}


