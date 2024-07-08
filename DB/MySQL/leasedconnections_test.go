package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	dep "github.com/DigiStratum/GoLib/Dependencies"
	cfg "github.com/DigiStratum/GoLib/Config"

	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

func TestThat_NewLeasedConnections_ReturnsSomething(t *testing.T) {
	// Setup

	// Test
	var actual LeasedConnectionsIfc = NewLeasedConnections()	// <- Ensures that we satisfy our interface

	// Verify
	ExpectNonNil(actual, t)
}

func TestThat_LeasedConnections_GetLeaseForConnection_ReturnsSomething(t *testing.T) {
	// Setup
	pooledConnection := getPooledConnection()
	sut := NewLeasedConnections()

	// Test
	actual := sut.GetLeaseForConnection(pooledConnection)

	// Verify
	ExpectNonNil(actual, t)
}

func TestThat_LeasedConnections_GetLeaseForConnection_ReturnsNil_ForNilConnection(t *testing.T) {
	// Setup
	sut := NewLeasedConnections()

	// Test
	actual := sut.GetLeaseForConnection(nil)

	// Verify
	ExpectNil(actual, t)
}

func TestThat_LeasedConnections_Release_ReturnsTrue_ForGoodLeaseKey(t *testing.T) {
	// Setup
	pooledConnection := getPooledConnection()
	sut := NewLeasedConnections()
	sut.GetLeaseForConnection(pooledConnection)

	// Test
	result := sut.Release(1) // leaseKey is deerministic; first one will be 1

	// Verify
	ExpectTrue(result, t)
}

func TestThat_LeasedConnections_Release_ReturnsFalse_ForBadLeaseKey(t *testing.T) {
	// Setup
	pooledConnection := getPooledConnection()
	sut := NewLeasedConnections()
	sut.GetLeaseForConnection(pooledConnection)

	// Test
	result := sut.Release(0)

	// Verify
	ExpectFalse(result, t)
}

func getPooledConnection() *pooledConnection {
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	connectionPool := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	connectionPool.InjectDependencies(
		dep.NewDependencyInstance("connectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	conn, _ := connectionFactory.NewConnection(dsn)
	newConnection, _ := NewConnection(conn)
	newPooledConnection, _ := NewPooledConnection(newConnection, connectionPool)
	return newPooledConnection
}

