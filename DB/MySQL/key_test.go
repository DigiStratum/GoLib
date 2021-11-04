package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDBKey_ReturnsSomething(t *testing.T) {
	// Test
	sut := NewDBKey("boguskey")

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewDBKeyFromDSN_ReturnsSomething(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("bogusdsn")

	// Test
	sut := NewDBKeyFromDSN(dsn)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Key_GetKey_ReturnsExpectedKey(t *testing.T) {
	// Setup
	expectedKey := "boguskey"
	sut := NewDBKey(expectedKey)

	// Test
	actualKey := sut.GetKey()

	// Verify
	ExpectString(expectedKey, actualKey, t)
}

