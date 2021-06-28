package mysql

/*
DB Manager for MySQL - manages connections and provides various reusable DB capabilities.
*/

import (
	"errors"
)

// Manager public interface
type ManagerIfc interface {
	Connect(dsn string) (DBKeyIfc, error)
	IsConnected(dbKeyIfc DBKeyIfc) bool
	Query(dbKeyIfc DBKeyIfc, query string, prototype ResultIfc, args ...interface{}) (ResultSetIfc, error)
	Run(dbKeyIfc DBKeyIfc, query QueryIfc, args ...interface{}) (ResultSetIfc, error)
	RunInt(dbKeyIfc DBKeyIfc, query QueryIfc, args ...interface{}) (*int, error)
	RunString(dbKeyIfc DBKeyIfc, query QueryIfc, args ...interface{}) (*string, error)
	Disconnect(dbKeyIfc DBKeyIfc)
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
func (mgr *manager) Connect(dsn string) (DBKeyIfc, error) {

	// If we already have this dbKey...
	dbKeyIfc := NewDBKeyFromDSN(dsn)
	if _, ok := mgr.connections[dbKeyIfc.GetKey()]; ! ok {
		// Not connected yet - let's do this thing!
		conn, err := NewConnection(dsn)
		if err != nil { return nil, err }

		// Make a new connection record
		mgr.connections[dbKeyIfc.GetKey()] = conn
	}
	return dbKeyIfc, nil
}

// Check that this connection is still established
func (mgr *manager) IsConnected(dbKeyIfc DBKeyIfc) bool {
	conn := mgr.getConnection(dbKeyIfc)
	if nil != conn { return conn.IsConnected() }
	return false
}

// Run a query against the database connection identified by the dbkey
func (mgr *manager) Query(dbKeyIfc DBKeyIfc, query string, prototype ResultIfc, args ...interface{}) (ResultSetIfc, error) {
	conn := mgr.getConnection(dbKeyIfc)
        if nil == conn { return nil, errors.New("Error getting connection") }
        return NewQuery(query, prototype).Run(conn, args...)
}

// Run a query against the database connection identified by the dbkey
func (mgr *manager) Run(dbKeyIfc DBKeyIfc, query QueryIfc, args ...interface{}) (ResultSetIfc, error) {
	conn := mgr.getConnection(dbKeyIfc)
	if nil == conn { return nil, errors.New("Error getting connection") }
	// TODO: check the result of the query and, if err, check the connection and, if fail, reconnect and try again
	return query.Run(conn, args...)
}

func (mgr *manager) RunInt(dbKeyIfc DBKeyIfc, query QueryIfc, args ...interface{}) (*int, error) {
	conn := mgr.getConnection(dbKeyIfc)
	if nil == conn { return nil, errors.New("Error getting connection") }
	// TODO: check the result of the query and, if err, check the connection and, if fail, reconnect and try again
	return query.RunInt(conn, args...)
}

func (mgr *manager) RunString(dbKeyIfc DBKeyIfc, query QueryIfc, args ...interface{}) (*string, error) {
	conn := mgr.getConnection(dbKeyIfc)
	if nil == conn { return nil, errors.New("Error getting connection") }
	// TODO: check the result of the query and, if err, check the connection and, if fail, reconnect and try again
	return query.RunString(conn, args...)
}

// Close the connection with this key, if it exists, and forget about it
// (There's no value in reusing the key, just delete it)
func (mgr *manager) Disconnect(dbKeyIfc DBKeyIfc) {
	conn := mgr.getConnection(dbKeyIfc)
	if nil != conn {
		conn.Disconnect()
		delete(mgr.connections, dbKeyIfc.GetKey())
	}
}

// ------------------------------------------------------------------------------------------------
// PRIVATE

// Get the connection for the specified key
func (mgr *manager) getConnection(dbKeyIfc DBKeyIfc) ConnectionIfc {
	if conn, ok := mgr.connections[dbKeyIfc.GetKey()]; ok {
		return conn
	}
	return nil
}

