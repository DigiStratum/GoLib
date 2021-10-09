package db

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_MakeDSN_ReturnsGoodDSNString(t *testing.T) {
	// Setup
	user := "user"
	pass := "pass"
	host := "host"
	port := "port"
	name := "name"
	expected := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)

	// Test
	actual := MakeDSN(user, pass, host, port, name)

	// Verify
	ExpectString(expected, actual, t)
}

func TestThat_GetDSNHash_ReturnsHashCode(t *testing.T) {
	// Test
	actual := GetDSNHash("fakedsn")

	// Verify
	ExpectInt(32, len(actual), t)
	ExpectString("2e2e111cfb0447f58c1be469d89ea984", actual, t)
}

