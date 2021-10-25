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

func TestThat_Manager_GetConnection_ReturnsNothing_WhenNoConnectionPoolSetup(t *testing.T) {
	// Setup
	sut := NewManager()
	dbKey := NewDBKeyFromDSN("bogusdsn")

	// Test
	leasedConnection, err := sut.GetConnection(dbKey)

	// Verify
	ExpectError(err, t)
	ExpectNil(leasedConnection, t)
}

func TestThat_Manager_NewConnectionPool_ReturnsDBKey(t *testing.T) {
	// Setup
	sut := NewManager()

	// Test
	dbKey := NewConnectionPool("bogusdsn")

	// Verify
	ExpectTrue(len(dbKey) > 0, t)
}

func TestThat_Manager_GetConnection_ReturnsSomething_WhenConnectionPoolSetup(t *testing.T) {
	// Setup
	sut := NewManager()
	dbKey := NewConnectionPool("bogusdsn")

	// Test
	leasedConnection, err := sut.GetConnection(dbKey)

	// Verify
	ExpectError(err, t)
	ExpectNil(leasedConnection, t)
}
