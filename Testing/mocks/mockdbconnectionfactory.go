package mockdb

import(
	"database/sql"
)

type MockDBConnectionFactory struct {
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMockDBConnectionFactory() *MockDBConnectionFactory {
	return &MockDBConnectionFactory{}
}

// -------------------------------------------------------------------------------------------------
// db.DBConnectionFactoryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *MockDBConnectionFactory) NewConnection(dsn string) (*sql.DB, error) {
	return NewMockDBConnection("mockdriver", dsn)
}
