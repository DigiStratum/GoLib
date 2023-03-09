package mysql

import(
	"fmt"
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	dep "github.com/DigiStratum/GoLib/Dependencies"
	cfg "github.com/DigiStratum/GoLib/Config"

	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

func TestThat_NewConnectionPool_ReturnsSomething(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")

	// Test
	sut := NewConnectionPool(*dsn)

	// Verify
	ExpectNonNil(sut, t)
}

// StartableIfc, DependencyInjectableIfc
func TestThat_ConnectionPool_Start_ReturnsError_ForMissingDependencies(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)

	// Test
	err := sut.InjectDependencies(
		dep.NewDependencyInstance("bogusdep", "bogusstring"),
		dep.NewDependencyInstance("ConnectionFactory", "bogusstring"),
	)

	actual := sut.Start()

	// Verify
	if ! ExpectError(err, t) { return }
	if ! ExpectError(actual, t) { return }
}

func TestThat_ConnectionPool_InjectDependencies_ReturnsNoError_ForGoodDependencies(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()

	// Test
	err := sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	// Verify
	ExpectNoError(err, t)
}

func TestThat_ConnectionPool_Configure_ReturnsNoError_ForEmptyConfig(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()

	// Test
	err := sut.Configure(config)

	// Verify
	ExpectNoError(err, t)
}

func TestThat_ConnectionPool_Configure_ReturnsError_ForUnknownConfigKeys(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("boguskey", "bogusvalue")

	// Test
	err := sut.Configure(config)

	// Verify
	ExpectError(err, t)
}

func TestThat_ConnectionPool_Configure_ReturnsNoError_ForKnownConfigKeys(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")

	// Test
	err := sut.Configure(config)

	// Verify
	ExpectNoError(err, t)
}

func TestThat_ConnectionPool_GetConnection_ReturnsLeasedConnection_WhenOneAvailable(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	sut.Configure(config)
	if ! ExpectNoError(sut.Start(), t) { return }

	// Test
	conn, err := sut.GetConnection()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(conn, t)
}

func TestThat_ConnectionPool_GetConnection_ReturnsError_WhenNoConnectionsAvailable(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	sut.Configure(config)

	// Test
	conn1, _ := sut.GetConnection()
	conn2, err := sut.GetConnection()

	// Verify
	ExpectError(err, t)
	ExpectNonNil(conn1, t)
	ExpectNil(conn2, t)
}

func TestThat_ConnectionPool_GetConnection_ReturnsLeasedConnection_WhenPreviouslyReleased(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	sut.Configure(config)

	// Test
	conn1, _ := sut.GetConnection()
	err1 := conn1.Release()
	ExpectNoError(err1, t)
	conn2, err2 := sut.GetConnection()

	// Verify
	ExpectNoError(err2, t)
	ExpectNonNil(conn2, t)
}

func TestThat_ConnectionPool_GetMaxIdle_ReturnsConfiguredValue(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	expected := 33
	config.Set("max_idle", fmt.Sprintf("%d", expected))
	sut.Configure(config)

	// Test
	actual := sut.GetMaxIdle()

	// Verify
	ExpectInt(expected, actual, t)
}

func TestThat_ConnectionPool_Close_ClosesConnectionPool_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	// Test
	err := sut.Close()

	// Verify
	ExpectNoError(err, t)
}

