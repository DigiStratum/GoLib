package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	dep "github.com/DigiStratum/GoLib/Dependencies"
	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
	cfg "github.com/DigiStratum/GoLib/Config"
)

func TestThat_NewPooledConnection_ReturnsSomething_WithoutError(t *testing.T) {
	// Setup
	var actual PooledConnectionIfc
	var err error

	// Test
	actual, err = getGoodNewPooledConnection()	// <- Ensures that we satisfy our interface

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err, t)
}

func TestThat_PooledConnection_Close_Returns_WithoutError(t *testing.T) {
	// Setup
	sut, _ := getGoodNewPooledConnection()

	// Test
	err := sut.Close()

	// Verify
	ExpectNoError(err, t)
}

func TestThat_PooledConnection_Close_Returns_Error_ForBadUnderlyingConnection(t *testing.T) {
	// Setup
	sut, _ := getNilNewPooledConnection()

	// Test
	err := sut.Close()

	// Verify
	ExpectError(err, t)
}

func TestThat_PooledConnection_IsConnected_ReturnsTrue_ForGoodConnection(t *testing.T) {
	// Setup
	sut, _ := getGoodNewPooledConnection()

	// Test
	actual := sut.IsConnected()

	// Verify
	ExpectTrue(actual, t)
}

func TestThat_PooledConnection_IsConnected_ReturnsFalse_ForBadConnection(t *testing.T) {
	// Setup
	sut, _ := getNilNewPooledConnection()

	// Test
	actual := sut.IsConnected()

	// Verify
	ExpectFalse(actual, t)
}

func getGoodNewPooledConnection() (*pooledConnection, error) {
	// Make a new ConnectionPool
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	connectionPool := NewConnectionPool(*dsn)

	// Inject ConnectionFactoryIfc
	connectionFactory := NewMockDBConnectionFactory()
	connectionPool.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	connectionPool.Configure(config)

	conn, _ := connectionFactory.NewConnection(dsn)
	newConnection, _ := NewConnection(conn)

	return NewPooledConnection(newConnection, connectionPool)
}

func getNilNewPooledConnection() (*pooledConnection, error) {
	return &pooledConnection{}, nil
}

