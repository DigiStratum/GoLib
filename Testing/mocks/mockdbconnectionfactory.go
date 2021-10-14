package mocks

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

func (r *MockDBConnectionFactory) NewConnection(dsn string) (*MockDBConnection, error) {
	return NewMockDBConnection("mockdriver", dsn)
}
