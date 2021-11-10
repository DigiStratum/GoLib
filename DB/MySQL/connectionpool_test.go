package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewConnectionPool_ReturnsSomething(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")

	// Test
	sut := NewConnectionPool(*dsn)

	// Verify
	ExpectNonNil(sut, t)
}
