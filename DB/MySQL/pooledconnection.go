package mysql

/*

A Pooled Connection wraps a raw DB connection with additional metadata to manage leasing

TODO: Add support for restoring the state of the connection in the event that we capture changes like transaction isolation, etc.

*/

import (
	"io"
	"fmt"
	"time"
	"sync"

	"database/sql"
)

// TODO: Can we not inherit this interface from ConnectionIfc?
type PooledConnectionIfc interface {

	// Connections
	IsConnected() bool

	// Leasing
	IsLeased() bool
	MatchesLeaseKey(leaseKey int64) bool
	Lease(leaseKey int64)
	Release() error
	Touch()
	IsExpired() bool

	ConnectionCommonIfc
}

type pooledConnection struct {
	pool			ConnectionPoolIfc	// The pool that this pooled connection lives in
	connection		ConnectionIfc		// Our underlying database connection
	establishedAt		int64			// Time that this connection was established to the DB
	lastActiveAt		int64			// Last time this connection saw activity from the consumer
	lastLeasedAt		int64			// Last time this connection was leased out
	isLeased		bool			// Is this connection currently leased out?
	leaseKey		int64			// This is the lease key for the current lease holder
	mutex			sync.Mutex
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewPooledConnection(connection ConnectionIfc, connPool ConnectionPoolIfc) (*pooledConnection, error) {
	if nil == connection { return nil, fmt.Errorf("Supplied connection was nil!") }
	if nil == connPool { return nil, fmt.Errorf("Supplied connection pool was nil!") }
	now := time.Now().Unix()
	return &pooledConnection{
		pool:			connPool,
		connection:		connection,
		establishedAt:		now,
		lastActiveAt:		0,
		lastLeasedAt:		0,
		isLeased:		false,
		leaseKey:		0,
	}, nil
}

// -------------------------------------------------------------------------------------------------
// io.Closer Public Interface
// -------------------------------------------------------------------------------------------------

// Drop this connection
func (r *pooledConnection) Close() error {
	if nil == r.connection { return fmt.Errorf("Underlying connection is nil") }
	if closeableConnection, ok := r.connection.(io.Closer); ok {
		return closeableConnection.Close()
	}
	return fmt.Errorf("pooledConnection is not Closeable")
}

// -------------------------------------------------------------------------------------------------
// pooledConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Connections
func (r *pooledConnection) IsConnected() bool {
	if nil == r.connection { return false }
	return r.connection.IsConnected()
}

// Leasing
func (r *pooledConnection) IsLeased() bool { return r.isLeased }

func (r *pooledConnection) MatchesLeaseKey(leaseKey int64) bool {
	if ! r.IsLeased() { return false }
	return r.leaseKey == leaseKey
}

func (r *pooledConnection) Lease(leaseKey int64) {
	r.mutex.Lock(); defer r.mutex.Unlock()

	// Set up the lease to guarantee nobody else comes and steals this from us
	r.isLeased = true
	r.leaseKey = leaseKey
	now := time.Now().Unix()
	r.lastLeasedAt = now
	r.lastActiveAt = now
	// Just in case we evicted a previous lease holder, see if there is any connection state reset needed
	if r.InTransaction() { r.Rollback() }
}

func (r *pooledConnection) Release() error {
	r.mutex.Lock(); defer r.mutex.Unlock()
	err := r.pool.Release(r.leaseKey)
	if nil != err { return err }
	r.isLeased = false
	r.leaseKey = 0
	return nil
}

func (r *pooledConnection) Touch() {
	r.lastActiveAt = time.Now().Unix()
}

func (r *pooledConnection) IsExpired() bool {
	maxIdle := int64(r.pool.GetMaxIdle())
	now := time.Now().Unix()
	// If the last time this connection was Touch()ed, plus the max idle period is in the past, lease expired!
	return r.lastActiveAt + maxIdle < now
}

// Transactions
func (r *pooledConnection) InTransaction() bool { return r.connection.InTransaction() }
func (r *pooledConnection) Rollback() error { return r.connection.Rollback() }
func (r *pooledConnection) Begin() error { return r.connection.Begin() }
func (r *pooledConnection) Commit() error { r.Touch(); return r.connection.Commit() }

// Operations
func (r *pooledConnection) NewQuery(query SQLQueryIfc) (QueryIfc, error) { r.Touch(); return r.connection.NewQuery(query) }
func (r *pooledConnection) Exec(query SQLQueryIfc, args ...interface{}) (sql.Result, error) { r.Touch(); return r.connection.Exec(query, args...) }
func (r *pooledConnection) Query(query SQLQueryIfc, args ...interface{}) (*sql.Rows, error) { r.Touch(); return r.connection.Query(query, args...) }
func (r *pooledConnection) QueryRow(query SQLQueryIfc, args ...interface{}) *sql.Row { r.Touch(); return r.connection.QueryRow(query, args...) }
