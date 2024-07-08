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
	if ! ExpectNonNil(sut, t) { return }
}

// StartableIfc, DependencyInjectableIfc
func TestThat_ConnectionPool_Start_ReturnsError_ForMissingDependencies(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)

	// Test
	err1 := sut.InjectDependencies(
		dep.NewDependencyInstance("bogusdep", "bogusstring"),
	)
	actualErr1 := sut.Start()
	err2 := sut.InjectDependencies(
		// Right name, but wrong interface, should throw an error!
		dep.NewDependencyInstance("ConnectionFactory", "bogusstring"),
	)

	actualErr2 := sut.Start()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectError(err2, t) { return }
	if ! ExpectError(actualErr1, t) { return }
	if ! ExpectError(actualErr2, t) { return }
}

func TestThat_ConnectionPool_InjectDependencies_ReturnsNoError_ForGoodDependencies(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)

	// Test
	err := sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
	)

	// Verify
	if ! ExpectNoError(err, t) { return }
}

func TestThat_ConnectionPool_Configure_ReturnsNoError_ForEmptyConfig(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
	)

	config := cfg.NewConfig()

	// Test
	err := sut.Configure(config)

	// Verify
	if ! ExpectNoError(err, t) { return }
}

func TestThat_ConnectionPool_Configure_ReturnsNoError_ForKnownConfigKeys(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")

	// Test
	err := sut.Configure(config)

	// Verify
	if ! ExpectNoError(err, t) { return }
}

func TestThat_ConnectionPool_GetConnection_ReturnsLeasedConnection_WhenOneAvailable(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
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
	if ! ExpectNoError(err, t) { return }
	if ! ExpectNonNil(conn, t) { return }
}

func TestThat_ConnectionPool_GetConnection_ReturnsError_WhenNoConnectionsAvailable(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	sut.Configure(config)
	sut.Start()

	// Test
	conn1, err1 := sut.GetConnection()
	conn2, err2 := sut.GetConnection()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNonNil(conn1, t) { return }
	if ! ExpectError(err2, t) { return }
	if ! ExpectNil(conn2, t) { return }
}

func TestThat_ConnectionPool_GetConnection_ReturnsLeasedConnection_WhenPreviouslyReleased(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	sut.Configure(config)
	sut.Start()

	// Test
	conn1, _ := sut.GetConnection()
	err1 := conn1.Release()
	if ! ExpectNoError(err1, t) { return }
	conn2, err2 := sut.GetConnection()

	// Verify
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectNonNil(conn2, t) { return }
}

func TestThat_ConnectionPool_GetMaxIdle_ReturnsConfiguredValue(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
	)

	config := cfg.NewConfig()
	expected := 33
	config.Set("max_idle", fmt.Sprintf("%d", expected))
	sut.Configure(config)
	sut.Start()

	// Test
	actual := sut.GetMaxIdle()

	// Verify
	if ! ExpectInt(expected, actual, t) { return }
}

func TestThat_ConnectionPool_Close_ClosesConnectionPool_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	sut.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", NewMockDBConnectionFactory()),
	)
	sut.Start()

	// Test
	err := sut.Close()

	// Verify
	if ! ExpectNoError(err, t) { return }
}

