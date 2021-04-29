package mysql

/*

DB Connection

*/

import (
	"database/sql"
)

// Connection public interface
type ConnectionIfc interface {
	IsConnected() bool
	Connect() error
	Disconnect()
	Reconnect()
	GetConnection() *sql.DB
}

type connection struct {
	dsn	string          // Full Data Source Name for this connection
	conn	*sql.DB         // Read-Write Connection
}

// Make a new one of these and connect!
func NewConnection(dsn string) (ConnectionIfc, error) {
	connection := connection{
		dsn:	dsn,
	}
	return &connection, connection.Connect()
}

// Check whether this connection is established
func (c *connection) IsConnected() bool {
	if nil == (*c).conn { return false }
	return nil == (*c).conn.Ping()
}

// Establish the connection using the suplied DSN
func (c *connection) Connect() error {
	// If we're already connected, nothing to do
	if c.IsConnected() { return nil }
	var err error
	(*c).conn, err = sql.Open("mysql", (*c).dsn)
	return err
}

// Drop this connection
func (c * connection) Disconnect() {
	// If we're not connected, nothing to do
	if ! c.IsConnected() { return }
	(*c).conn.Close()
}

// Cycle this connection, or establish a new connection if we're not connected
func (c *connection) Reconnect() {
	if c.IsConnected() { c.Disconnect() }
	c.Connect()
}

// Get the underlying connection for the caller to put it to work!
func (c *connection) GetConnection() *sql.DB {
	return (*c).conn
}

