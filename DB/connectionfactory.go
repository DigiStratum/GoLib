package db

/*
Factory for producing DBConnections; we can use Dependency Injection to enable unit testing of DB-integrations
*/

import(
	"database/sql"
)

type ConnectionFactoryIfc interface {
	// Return interface instead of struct so that other implementations may satisfy
	NewConnection(dsn DSNIfc) (*sql.DB, error)
}
