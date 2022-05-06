package mysql

import(
//	"fmt"
	"testing"

//	"github.com/DATA-DOG/go-sqlmock"

//	"github.com/DigiStratum/GoLib/DB"
	. "github.com/DigiStratum/GoLib/Testing"
//	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

// NewQuery(connection ConnectionIfc, qry string) (*Query, error)
func TestThat_NewQuery_ReturnsError_WhenGivenNilConnection(t *testing.T) {
	// Test
	sut, err := NewQuery(nil, nil)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

