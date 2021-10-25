package mysql

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewManager_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewManager()

	// Verify
	ExpectNonNil(sut, t)
}


func TestThat_Manager_GetConnection_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewManager()
	dbKey := NewDBKeyFromDSN("bogusdsn")

	// Test
	leasedConnection, err := sut.GetConnection(dbKey)

	// Verify
	ExpectError(err, t)
	ExpectNil(leasedConnection, t)
}