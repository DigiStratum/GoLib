package mysql

/*

DB Connection

*/

import (
	"database/sql"
)

type Connection struct {
	dsn	string          // Full Data Source Name for this connection
	conn	*sql.DB         // Read-Write Connection
}

// Make a new one of these and connect!
func NewConnection(dsn string) (*Connection, error) {
	connection := Connection{
		dsn:	dsn,
	}
	return &connection, connection.Connect()
}

// Check whether this connection is established
func (c *Connection) IsConnected() bool {
	if nil == (*c).conn { return false }
	return nil == (*c).conn.Ping()
}

// Establish the connection using the suplied DSN
func (c *Connection) Connect() error {
	// If we're already connected, nothing to do
	if c.IsConnected() { return nil }
	var err error
	(*c).conn, err = sql.Open("mysql", (*c).dsn)
	return err
}

// Drop this connection
func (c * Connection) Disconnect() {
	// If we're not connected, nothing to do
	if ! c.IsConnected() { return }
	(*c).conn.Close()
}

// Cycle this connection, or establish a new connection if we're not connected
func (c *Connection) Reconnect() {
	if c.IsConnected() { c.Disconnect() }
	c.Connect()
}

// Get the underlying connection for the caller to put it to work!
func (c *Connection) GetConnection() *sql.DB {
	return (*c).conn
}

