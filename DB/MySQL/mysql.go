package mysql

import (
	"fmt"
	"errors"
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	lib "github.com/DigiStratum/GoLib"
	"github.com/DigiStratum/GoLib/DB"
)

type connection struct {
	Dsn	string
	IsOpen	bool
	Conn	*sql.DB
}

// Set of connections, keyed on DSN
type DBManagerMySQL struct {
	connections	map[string]Connection
}

// Make a new one of these!
func NewDBManagerMySQL() *DBManagerMySQL {
	dbm := DBManagerMySQL{
	}
}

// Get DB Connection Key from the supplied DSN
// ref: https://en.wikipedia.org/wiki/Data_source_name
func (dbm *DBManagerMySQL) Connect(dsn string) string, error {
	key := db.GetDSNHash(dsn)

	// If we already have this key...
	if _, ok := dbm.connections[key]; ok {
		// ... then it's because we already have a good connection
		return key, nil;
	}

	// Not connected yet - let's do this thing!
        conn, err := sql.Open("mysql", connectString)
        if err != nil {
		return "", err
        }

	// Make a new connection record
	newConnection := connection{
		Dsn:	key,
		IsOpen:	true,
		Conn:	conn,
	}
	return key, nil
}

func (dbm *DBManagerMySQL) IsConnectionOpen(key string) bool {
	if conn, ok := dbm.connections[key]; ok {
		return conn.IsOpen
	}
	return false
}

func (dbm *DBManagerMySQL) Disconnect(key string) error {
	if ! dbm.IsConnectionOpen(key) {
		return errors.Error(fmt.Sprintf("The connection for '%s' is not open", key))
	}
	if conn, ok := dbm.connections[key]; ok {
		conn.Conn.Close()
		c := conn
		c.Open = false
		dbm.connections[key] = c
	}
	return nil
}

