package db

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDSN_ReturnsError(t *testing.T) {
	// Test
	sut, err := NewDSN("bogusdsn")

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

func TestThat_NewDSN_ReturnsDSNObject(t *testing.T) {
	// Setup
	dsn := "user:pass@tcp(host:port)/name"

	// Test
	sut, err := NewDSN(dsn)

	// Verify
	ExpectNonNil(sut, t)
	ExpectNoError(err, t)
}

func TestThat_DSN_GetDSNHash_ReturnsHashCode(t *testing.T) {
	// Setup
	dsn := "user:pass@tcp(host:port)/name"

	// Test
	sut, _ := NewDSN(dsn)
	hash := sut.ToHash()

	// Verify
	ExpectInt(32, len(hash), t)
	ExpectString("a802d6a3bd91e0d67d39bbb5ce03c153", hash, t)
}

