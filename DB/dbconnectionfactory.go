package db

/*
Factory for producing DBConnections; we can use Dependency Injection to enable unit testing of DB-integrations
*/

type DBConnectionFactoryIfc interface {
	NewConnection() (*DBConnection, error)
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

func (r *DBConnectionFactory) NewConnection(dsn string) (*DBConnection, error) {
	return NewDBConnection(r.driver, dsn)
}
