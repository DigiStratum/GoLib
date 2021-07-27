package mysql

/*
DB Manager for MySQL - manages a set of named (keyed) mysql database connections / pools
*/

import (
	"errors"

	conn "github.com/DigiStratum/GoLib/DB/MySQL/connection"
)

type ManagerIfc interface {
	// Public interface
	Connect(dsn string) (DBKeyIfc, error)
	IsConnected(dbKey DBKeyIfc) bool
	NewQuery(dbKey DBKeyIfc, qry string) (QueryIfc, error)
	Disconnect(dbKey DBKeyIfc)
	// Private interface
	//getConnection(dbKey DBKeyIfc) conn.ConnectionIfc
	getConnectionPool(dbKey DBKeyIfc) conn.ConnectionPoolIfc
}

// Set of connections, keyed on DSN
type manager struct {
	//connections		map[string]conn.ConnectionIfc
	connectionPools		map[string]conn.ConnectionPoolIfc
}

// Make a new one of these!
func NewManager() ManagerIfc {
	return &manager{
		//connections: make(map[string]conn.ConnectionIfc),
		connectionPools: make(map[string]conn.ConnectionPoolIfc),
	}
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get DB Connection Key from the supplied DSN
// ref: https://en.wikipedia.org/wiki/Data_source_name
func (mgr *manager) Connect(dsn string) (DBKeyIfc, error) {

	dbKey := NewDBKeyFromDSN(dsn)
/*
	if _, ok := mgr.connections[dbKey.GetKey()]; ! ok {
		// Not connected yet - let's do this thing!
		conn, err := conn.NewConnection(dsn)
		if err != nil { return nil, err }

		// Make a new connection record
		mgr.connections[dbKey.GetKey()] = conn
	}
*/
	if _, ok := mgr.connectionPools[dbKey.GetKey()]; ! ok {
		mgr.connectionPools[dbKey.GetKey()] = NewConnectionPool(dbKey.GetKey())
	}


	return dbKey, nil
}

// Check that this connection is still established
func (mgr *manager) IsConnected(dbKey DBKeyIfc) bool {
	conn := mgr.getConnection(dbKey)
	if nil != conn { return conn.IsConnected() }
	return false
}

// Make a new Query attached to this manager session
func (mgr *manager) NewQuery(dbKey DBKeyIfc, qry string) (QueryIfc, error) {
	conn := mgr.getConnection(dbKey)
        if nil == conn { return nil, errors.New("Error getting connection") }
	return NewQuery(conn, qry), nil
}

// Close the connection with this key, if it exists, and forget about it
// (There's no value in reusing the key, just delete it)
func (mgr *manager) Disconnect(dbKey DBKeyIfc) {
	//conn := mgr.getConnection(dbKey)
	connPool := mgr.getConnectionPool(dbKey)
	if nil != conn {
		//conn.Disconnect()
		//delete(mgr.connections, dbKey.GetKey())
		// FIXME: Disconnect/close the connection pool
		delete(mgr.connectionPools, dbKey.GetKey())
	}
}

// -------------------------------------------------------------------------------------------------
// ManagerIfc Private Interface
// -------------------------------------------------------------------------------------------------

// Get the connection for the specified key
//func (mgr *manager) getConnection(dbKey DBKeyIfc) conn.ConnectionIfc {
	//if conn, ok := mgr.connections[dbKey.GetKey()]; ok {
	//	return conn
func (mgr *manager) getConnectionPool(dbKey DBKeyIfc) conn.ConnectionPoolIfc {
	if connPool, ok := mgr.connectionPools[dbKey.GetKey()]; ok {
		return connPool
	}
	return nil
}

