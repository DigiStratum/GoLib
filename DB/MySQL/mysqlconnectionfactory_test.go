package mysql

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewResultSet_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewMySQLConnectionFactory()

	// Verify
	ExpectNonNil(sut, t)
}
