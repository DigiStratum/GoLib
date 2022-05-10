package mysql

import(
//	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/DigiStratum/GoLib/DB"
	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

func getGoodNewConnection() (*Connection, error) {
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	return NewConnection(mockDBConnection)
}

// NewQuery(connection ConnectionIfc, qry string) (*Query, error)
func TestThat_NewQuery_ReturnsError_WhenGivenNilConnection(t *testing.T) {
	// Test
	sut, err := NewQuery(nil, nil)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

func TestThat_NewQuery_ReturnsError_WhenGivenNilSQLQuery(t *testing.T) {
	// Setup
	mockDBConnection, _ := getGoodNewConnection()

	// Test
	sut, err := NewQuery(mockDBConnection, nil)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

func TestThat_NewQuery_ReturnsSomething_WhenGivenGoodParams(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	//mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnResult(expectedResult)
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.Run()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err2, t)
}


