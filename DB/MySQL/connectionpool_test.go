package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	"github.com/DigiStratum/GoLib/Dependencies"
	. "github.com/DigiStratum/GoLib/Testing"
	"github.com/DigiStratum/GoLib/Testing/mocks"
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
	//deps.Set("bogusdep", "bogusstring")

	// Test
	err := sut.InjectDependencies(deps)

	// Verify
	ExpectNoError(err, t)
}
