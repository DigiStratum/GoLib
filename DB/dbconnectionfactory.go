package db

/*
Factory for producing DBConnections; we can use Dependency Injection to enable unit testing of DB-integrations
*/

import(
	"fmt"
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

func DBConnectionFactryFromIfc(i interface{}) (DBConnectionFactoryIfc, error) {
	if ii, ok := i.(DBConnectionFactoryIfc); ok { return ii, nil }
	return nil, fmt.Errorf("Does not implement DBConnectionFactoryIfc")
}

// -------------------------------------------------------------------------------------------------
// ConnectionFactoryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *DBConnectionFactory) NewConnection(dsn DSNIfc) (*sql.DB, error) {
	return NewDBConnection(r.driver, dsn)
}
