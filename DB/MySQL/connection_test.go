package mysql

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

func TestThat_NewConnection_ReturnsError_WhenGivenNilConnection(t *testing.T) {
	// Test
	sut, err := NewConnection(nil)

	// Verify
	ExpectNil(sut, t)
	ExpectNonNil(err, t)
}

func TestThat_NewConnection_ReturnsConnection_WhenGivenDBConnection(t *testing.T) {
	// Test
	mockDBConnection, _ := NewMockDBConnection("mockdriver", "mockdsn")
	sut, err := NewConnection(mockDBConnection)

	// Verify
	ExpectNil(err, t)
	ExpectNonNil(sut, t)
}

func TestThat_Connection_IsConnected_ReturnsTrue_WhenConnected(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection("mockdriver", "mockdsn")
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res := sut.IsConnected()

	// Verify
	ExpectTrue(res, t)
}

func TestThat_Connection_IsConnected_ReturnsFalse_WhenNotConnected(t *testing.T) {
	// Setup
	sut := Connection{}

	// Test
	res := sut.IsConnected()

	// Verify
	ExpectFalse(res, t)
}

func TestThat_Connection_IsConnected_ReturnsFalse_WhenConnectedThenClosed(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection("mockdriver", "mockdsn")
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res1 := sut.IsConnected()
	sut.Close()
	res2 := sut.IsConnected()

	// Verify
	ExpectTrue(res1, t)
	ExpectFalse(res2, t)
}

func TestThat_Connection_InTransaction_ReturnsFalse_WhenNotInTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection("mockdriver", "mockdsn")
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res := sut.InTransaction()

	// Verify
	ExpectFalse(res, t)
}

func TestThat_Connection_InTransaction_ReturnsTrue_WhenInTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection("mockdriver", "mockdsn")
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	res := sut.InTransaction()

	// Verify
	ExpectTrue(res, t)
}
