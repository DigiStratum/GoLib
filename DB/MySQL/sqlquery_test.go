package mysql

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewSQLQuery_ReturnsSomething(t *testing.T) {
	// Test
	var sut *SQLQuery = NewSQLQuery("bogus query")

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_SQLQuery_Resolve_ReturnsError_WhenNil(t *testing.T) {
	// Setup
	var sut SQLQuery

	// Test
	_, err := sut.Resolve()

	// Verify
	ExpectError(err, t)
}

func TestThat_SQLQuery_Resolve_ReturnsQuery_NoError(t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	sut := NewSQLQuery(expectedQuery)

	// Test
	actual, err := sut.Resolve()

	// Verify
	ExpectNoError(err, t)
	ExpectString(expectedQuery, actual, t)
}

