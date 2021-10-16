package db

/*
Factory for producing DBConnections; we can use Dependency Injection to enable unit testing of DB-integrations
*/

import(
	"database/sql"
)

type DBConnectionFactoryIfc interface {
	// Return interface instead of struct so that other implementations may satisfy
	NewConnection(dsn string) (*sql.DB, error)
}

type DBConnectionFactory struct {
	driver		string
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDBConnectionFactory(driver string) *DBConnectionFactory {
	return &DBConnectionFactory{
		driver:		driver,
	}
}

// -------------------------------------------------------------------------------------------------
// DBConnectionFactoryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *DBConnectionFactory) NewConnection(dsn string) (*sql.DB, error) {
	return NewDBConnection(r.driver, dsn)
}
