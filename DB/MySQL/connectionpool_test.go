package mysql

import(
	"fmt"
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	"github.com/DigiStratum/GoLib/Dependencies"
	. "github.com/DigiStratum/GoLib/Testing"
	"github.com/DigiStratum/GoLib/Testing/mocks"
	cfg "github.com/DigiStratum/GoLib/Config"
)

func TestThat_NewConnectionPool_ReturnsSomething(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")

	// Test
	sut := NewConnectionPool(*dsn)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_InjectDependencies_ReturnsError_ForNilDependencies(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)

	// Test
	err := sut.InjectDependencies(nil)

	// Verify
	ExpectError(err, t)
}

func TestThat_InjectDependencies_ReturnsError_ForWrongDependencies(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps1 := dependencies.NewDependencies()
	deps1.Set("bogusdep", "bogusstring")
	deps2 := dependencies.NewDependencies()
	deps2.Set("connectionFactory", "bogusstring")

	// Test
	err1 := sut.InjectDependencies(deps1)
	err2 := sut.InjectDependencies(deps2)

	// Verify
	ExpectError(err1, t)
	ExpectError(err2, t)
}

func TestThat_InjectDependencies_ReturnsNoError_ForGoodDependencies(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)

	// Test
	err := sut.InjectDependencies(deps)

	// Verify
	ExpectNoError(err, t)
}

func TestThat_Configure_ReturnsNoError_ForEmptyConfig(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

	config := cfg.NewConfig()

	// Test
	err := sut.Configure(config)

	// Verify
	ExpectNoError(err, t)
}

func TestThat_Configure_ReturnsError_ForUnknownConfigKeys(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

	config := cfg.NewConfig()
	config.Set("boguskey", "bogusvalue")

	// Test
	err := sut.Configure(config)

	// Verify
	ExpectError(err, t)
}

func TestThat_Configure_ReturnsNoError_ForKnownConfigKeys(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")

	// Test
	err := sut.Configure(config)

	// Verify
	ExpectNoError(err, t)
}

func TestThat_GetConnection_ReturnsLeasedConnection_WhenOneAvailable(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	sut.Configure(config)

	// Test
	conn, err := sut.GetConnection()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(conn, t)
}

func TestThat_GetConnection_ReturnsError_WhenNoConnectionsAvailable(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

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

func TestThat_GetConnection_ReturnsLeasedConnection_WhenPreviouslyReleased(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

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

func TestThat_GetMaxIdle_ReturnsConfiguredValue(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

	config := cfg.NewConfig()
	expected := 33
	config.Set("max_idle", fmt.Sprintf("%d", expected))
	sut.Configure(config)

	// Test
	actual := sut.GetMaxIdle()

	// Verify
	ExpectInt(expected, actual, t)
}

func TestThat_Close_ClosesConnectionPool_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	sut := NewConnectionPool(*dsn)
	deps := dependencies.NewDependencies()
	connecitonFactory := mockdb.NewMockDBConnectionFactory()
	deps.Set("connectionFactory", connecitonFactory)
	sut.InjectDependencies(deps)

	// Test
	err := sut.Close()

	// Verify
	ExpectNoError(err, t)
}
