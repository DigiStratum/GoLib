package mysql

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewResultSet_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewMySQLConnectionFactory("bogusdriver")

	// Verify
	ExpectNonNil(sut, t)
}
