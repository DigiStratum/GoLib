package db

import(
	"testing"
	"database/sql"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDBConnectionFactory_ReturnsSomething(t *testing.T) {
	// Setup
	var sut *DBConnectionFactory

	// Test
	sut = NewDBConnectionFactory("bogusdriver")

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewConnectionFactory_Returns_Connection_ForGoodFactoryDriver(t *testing.T) {
	// Setup
	sut := NewDBConnectionFactory("mysql")
	var actual *sql.DB
	var err error
	dsn, _ := NewDSN("user:pass@tcp(host:port)/name")

	// Test
	actual, err = sut.NewConnection(dsn)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actual, t)
}

func TestThat_NewConnectionFactory_Returns_Error_ForBadFactoryDriver(t *testing.T) {
	// Setup
	sut := NewDBConnectionFactory("baddriver")
	var actual *sql.DB
	var err error
	dsn, _ := NewDSN("user:pass@tcp(host:port)/name")

	// Test
	actual, err = sut.NewConnection(dsn)

	// Verify
	ExpectError(err, t)
	ExpectNil(actual, t)
}

