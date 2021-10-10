package mysql

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewConnection_ReturnsNil_WithError(t *testing.T) {
	// Test
	conn, err := NewConnection("fakedsn")

	// Verify
	ExpectNil(conn, t)
	ExpectNonNil(err, t)
}

func TestThat_Connection_IsConnected_ReturnsFalse_WhenNotConnected(t *testing.T) {
	// Setup
	sut := Connection{}

	// Test
	res := sut.IsConnected()

	// Verify
	ExpectFalse(res, t)
}
