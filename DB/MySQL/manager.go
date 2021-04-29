package mysql

/*
DB Manager for MySQL - manages connections and provides various reusable DB capabilities.
*/

import (
	"errors"
)

// Manager public interface
type ManagerIfc interface {
	Connect(dsn string) (*DBKey, error)
	IsConnected(dbKey DBKey) bool
	Query(dbKey DBKey, query string, prototype ResultIfc, args ...interface{}) (*ResultSet, error)
	Run(dbKey DBKey, query QueryIfc, args ...interface{}) (*ResultSet, error)
	Disconnect(dbKey DBKey)
}

// Set of connections, keyed on DSN
type manager struct {
	connections	map[string]ConnectionIfc
}

// Make a new one of these!
func NewManager() ManagerIfc {
	return &manager{
		connections: make(map[string]ConnectionIfc),
	}
}

// Get DB Connection Key from the supplied DSN
// ref: https://en.wikipedia.org/wiki/Data_source_name
func (mgr *manager) Connect(dsn string) (*DBKey, error) {

	// If we already have this dbKey...
	dbKey := NewDBKeyFromDSN(dsn)
	if _, ok := mgr.connections[dbKey.GetKey()]; ! ok {
		// Not connected yet - let's do this thing!
		conn, err := NewConnection(dsn)
		if err != nil { return nil, err }

		// Make a new connection record
		mgr.connections[dbKey.GetKey()] = conn
	}
	return dbKey, nil
}

// Check that this connection is still established
func (mgr *manager) IsConnected(dbKey DBKey) bool {
	conn := mgr.getConnection(dbKey)
	if nil != conn { return conn.IsConnected() }
	return false
}

// Run a query against the dtaabase connection identified by the dbkey
func (mgr *manager) Query(dbKey DBKey, query string, prototype ResultIfc, args ...interface{}) (*ResultSet, error) {
	conn := mgr.getConnection(dbKey)
        if nil == conn { return nil, errors.New("Error getting connection") }
        return NewQuery(query, prototype).Run(conn, args...)
}

// Run a query against the dtaabase connection identified by the dbkey
func (mgr *manager) Run(dbKey DBKey, query QueryIfc, args ...interface{}) (*ResultSet, error) {
	conn := mgr.getConnection(dbKey)
	if nil == conn { return nil, errors.New("Error getting connection") }
	// TODO: check the result of the query and, if err, check the connection and, if fail, reconnect and try again
	return query.Run(conn, args...)
}

// Close the connection with this key, if it exists, and forget about it
// (There's no value in reusing the key, just delete it)
func (mgr *manager) Disconnect(dbKey DBKey) {
	conn := mgr.getConnection(dbKey)
	if nil != conn {
		conn.Disconnect()
		delete(mgr.connections, dbKey.GetKey())
	}
}

// ------------------------------------------------------------------------------------------------
// PRIVATE

// Get the connection for the specified key
func (mgr *manager) getConnection(dbKey DBKey) ConnectionIfc {
	if conn, ok := mgr.connections[dbKey.GetKey()]; ok {
		return conn
	}
	return nil
}

