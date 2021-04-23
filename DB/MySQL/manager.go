package mysql

/*
DB Manager for MySQL - manages connections and provides various reusable DB capabilities.
*/

import (
	"errors"

	db "github.com/DigiStratum/GoLib/DB"
)

// ------------------------------------------------------------------------------------------------
// PUBLIC

// DB Key (Connection identifier)
type DBKey struct {
	Key	string
}

// DB Manager

// Set of connections, keyed on DSN
type Manager struct {
	connections	map[string]*Connection
}

// Make a new one of these!
func NewManager() *Manager {
	mgr := Manager{
		connections: make(map[string]*Connection),
	}
	return &mgr
}

// Get DB Connection Key from the supplied DSN
// ref: https://en.wikipedia.org/wiki/Data_source_name
func (mgr *Manager) Connect(dsn string) (*DBKey, error) {

	// If we already have this dbKey...
	dbKey := DBKey{ Key: db.GetDSNHash(dsn) }
	if _, ok := mgr.connections[dbKey.Key]; ! ok {
		// Not connected yet - let's do this thing!
		conn, err := NewConnection(dsn)
		if err != nil { return nil, err }

		// Make a new connection record
		mgr.connections[dbKey.Key] = conn
	}
	return &dbKey, nil
}

// Check that this connection is still established
func (mgr *Manager) IsConnected(dbKey DBKey) bool {
	conn := mgr.getConnection(dbKey)
	if nil != conn { return conn.IsConnected() }
	return false
}

// Run a query against the dtaabase connection identified by the dbkey
func (mgr *Manager) Query(dbKey DBKey, query string, prototype ResultIfc, args ...interface{}) (*ResultSet, error) {
	dbConn := mgr.getConnection(dbKey)
        if nil == dbConn { return nil, errors.New("Error getting connection") }
        q := NewQuery(query, prototype)
        return q.Run(dbConn, args...)
}

// Close the connection with this key, if it exists, and forget about it
// (There's no value in reusing the key, just delete it)
func (mgr *Manager) Disconnect(dbKey DBKey) {
	conn := mgr.getConnection(dbKey)
	if nil != conn {
		conn.Disconnect()
		delete(mgr.connections, dbKey.Key)
	}
}

// ------------------------------------------------------------------------------------------------
// PRIVATE

// Get the connection for the specified key
func (mgr *Manager) getConnection(dbKey DBKey) *Connection {
	if conn, ok := mgr.connections[dbKey.Key]; ok {
		return conn
	}
	return nil
}

