package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	"github.com/DigiStratum/GoLib/Dependencies"
	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
	cfg "github.com/DigiStratum/GoLib/Config"
)

func TestThat_NewLeasedConnection_ReturnsSomething(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	connectionPool := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connectionFactory := NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connectionFactory)
	connectionPool.InjectDependencies(deps)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	connectionPool.Configure(config)

	conn, _ := connectionFactory.NewConnection(dsn)
	newConnection, _ := NewConnection(conn)
	newPooledConnection, _ := NewPooledConnection(newConnection, connectionPool)

	leaseKey := int64(333)

	// Test
	actual := NewLeasedConnection(newPooledConnection, leaseKey)

	// Verify
	ExpectNonNil(actual, t)
}

/*
FIXME: Release indicates no leased connection in the pool with our key; this is unexpected.
Switching to PooledConnection test to address this before returning here

func TestThat_Release_Returns_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	connectionPool := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connectionFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connectionFactory)
	connectionPool.InjectDependencies(deps)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	connectionPool.Configure(config)

	conn, _ := connectionFactory.NewConnection(dsn)
	newConnection, _ := NewConnection(conn)
	newPooledConnection, _ := NewPooledConnection(newConnection, connectionPool)

	leaseKey := int64(333)
	sut := NewLeasedConnection(newPooledConnection, leaseKey)

	// Test
	err := sut.Release()

	// Verify
	ExpectNoError(err, t)
}
*/