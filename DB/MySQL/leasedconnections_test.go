package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	"github.com/DigiStratum/GoLib/Dependencies"
	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
	cfg "github.com/DigiStratum/GoLib/Config"
)

func TestThat_NewLeasedConnections_ReturnsSomething(t *testing.T) {
	// Setup

	// Test
	actual := NewLeasedConnections()

	// Verify
	ExpectNonNil(actual, t)
}

func TestThat_GetLeaseForConnection_ReturnsSomething(t *testing.T) {
	// Setup
	pooledConnection := getPooledConnection()
	sut := NewLeasedConnections()

	// Test
	actual := sut.GetLeaseForConnection(pooledConnection)

	// Verify
	ExpectNonNil(actual, t)
}

func TestThat_GetLeaseForConnection_ReturnsNil_ForNilConnection(t *testing.T) {
	// Setup
	sut := NewLeasedConnections()

	// Test
	actual := sut.GetLeaseForConnection(nil)

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Release_ReturnsTrue_ForGoodLeaseKey(t *testing.T) {
	// Setup
	pooledConnection := getPooledConnection()
	sut := NewLeasedConnections()
	sut.GetLeaseForConnection(pooledConnection)

	// Test
	result := sut.Release(1) // leaseKey is deerministic; first one will be 1

	// Verify
	ExpectTrue(result, t)
}

func TestThat_Release_ReturnsFalse_ForBadLeaseKey(t *testing.T) {
	// Setup
	pooledConnection := getPooledConnection()
	sut := NewLeasedConnections()
	sut.GetLeaseForConnection(pooledConnection)

	// Test
	result := sut.Release(0)

	// Verify
	ExpectFalse(result, t)
}

func getPooledConnection() *PooledConnection {
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
	conn, _ := connectionFactory.NewConnection(dsn)
	newConnection, _ := NewConnection(conn)
	newPooledConnection, _ := NewPooledConnection(newConnection, connectionPool)
	return newPooledConnection
}

