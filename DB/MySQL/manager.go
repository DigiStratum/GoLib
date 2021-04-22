package mysql

/*
DB Manager for MySQL - manages connections and provides various reusable DB capabilities.
*/

import (
	"fmt"
	"errors"

	db "github.com/DigiStratum/GoLib/DB"
)

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
	if conn, ok := mgr.connections[dbKey.Key]; ok {
		return conn.IsConnected()
	}
	return false
}

// TODO: we should maybe get rid of this - if you want a direct connection then connect directly, no?
func (mgr *Manager) GetConnection(dbKey DBKey) (*Connection, error) {
	if conn, ok := mgr.connections[dbKey.Key]; ok {
		return conn, nil
	}
	return nil, errors.New(fmt.Sprintf("The connection for '%s' is undefined", dbKey.Key))
}

// Run a query against the dtaabase connection identified by the dbkey
func (mgr *Manager) RunQuery(dbKey DBKey, query string, prototype ResultIfc) (*ResultSet, error) {
	dbConn, err := mgr.GetConnection(dbKey)
        if nil != err {
		return nil, errors.New(fmt.Sprintf("Error getting connection: %s\n", err.Error()))
	}
        q := NewQuery(query, prototype)
        return q.Run(dbConn)
}

// Close the connection with this key, if it exists, and forget about it
// (There's no value in reusing the key, just delete it)
func (mgr *Manager) Disconnect(dbKey DBKey) {
	if conn, ok := mgr.connections[dbKey.Key]; ok {
		conn.Disconnect()
		delete(mgr.connections, dbKey.Key)
	}
}

