package db

/*
Factory for producing DBConnections; we can use Dependency Injection to enable unit testing of DB-integrations
*/

import(
	"database/sql"
)

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
// ConnectionFactoryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *DBConnectionFactory) NewConnection(dsn DSNIfc) (*sql.DB, error) {
	return NewDBConnection(r.driver, dsn)
}
