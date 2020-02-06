package mysql

/*
DB Manager for MySQL - manages connections and provides various reusable DB capabilities.

TODO: Add some sort of query builder - this will allow us to ditch writing SQL for most needs.

*/

import (
	"fmt"
	"errors"
	"database/sql"
	//"net/http"

	_ "github.com/go-sql-driver/mysql"

	//lib "github.com/DigiStratum/GoLib"
	db "github.com/DigiStratum/GoLib/DB"
)

type connection struct {
	Dsn	string
	IsOpen	bool
	Conn	*sql.DB
}

// Set of connections, keyed on DSN
type DBManagerMySQL struct {
	connections	map[string]connection
}

// Make a new one of these!
func NewDBManagerMySQL() *DBManagerMySQL {
	dbm := DBManagerMySQL{
		connections: make(map[string]connection),
	}
	return &dbm
}

// Get DB Connection Key from the supplied DSN
// ref: https://en.wikipedia.org/wiki/Data_source_name
func (dbm *DBManagerMySQL) Connect(dsn string) (string, error) {
	key := db.GetDSNHash(dsn)

	// If we already have this key...
	if _, ok := dbm.connections[key]; ok {
		// ... then it's because we already have a good connection
		return key, nil;
	}

	// Not connected yet - let's do this thing!
        conn, err := sql.Open("mysql", dsn)
        if err != nil {
		return "", err
        }

	// Make a new connection record
	dbm.connections[key] = connection{
		Dsn:	dsn,
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
		return errors.New(fmt.Sprintf("The connection for '%s' is not open", key))
	}
	if conn, ok := dbm.connections[key]; ok {
		conn.Conn.Close()
		c := conn
		c.IsOpen = false
		dbm.connections[key] = c
	}
	return nil
}

// The spec for a prepared statement query. Single '?' substitution is handled by db.Query()
// automatically. '???' expands to include enough placeholders (as with an IN () list for any count
// of keys > min. max must be >= min unless max == 0.
type QuerySpec struct {
	Query	string	// The query to execute as prepared statement
	MinKeys	int	// minimum num keys required to populate query; 0 = no min
	MaxKeys	int	// maximum num keys required to populate query; 0 = no max
}

