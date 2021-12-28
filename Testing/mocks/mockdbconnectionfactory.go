package mocks

import(
	"database/sql"

	"github.com/DigiStratum/GoLib/DB"
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

func (r *MockDBConnectionFactory) NewConnection(dsn db.DSNIfc) (*sql.DB, error) {
	return NewMockDBConnection("mockdriver", dsn)
}
