package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewMySQLConnectionFactory_ReturnsSomething(t *testing.T) {
	// Test
	var sut *db.DBConnectionFactory = NewMySQLConnectionFactory()

	// Verify
	ExpectNonNil(sut, t)
}
