package db

import(
        "testing"
        "database/sql"

        . "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDBConnection_Returns_DBConnection_ForGoodDriver(t *testing.T) {
        // Setup
        var actual *sql.DB
	var err error
        dsn, _ := NewDSN("user:pass@tcp(host:port)/name")

        // Test
        actual, err = NewDBConnection("mysql", dsn)

        // Verify
        ExpectNoError(err, t)
        ExpectNonNil(actual, t)
}

func TestThat_NewDBConnection_ReturnsError_ForBadDriver(t *testing.T) {
        // Setup
        var actual *sql.DB
	var err error
        dsn, _ := NewDSN("user:pass@tcp(host:port)/name")

        // Test
        actual, err = NewDBConnection("bogusdriver", dsn)

        // Verify
        ExpectError(err, t)
        ExpectNil(actual, t)
}

