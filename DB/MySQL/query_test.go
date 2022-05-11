package mysql

import(
	"fmt"
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
	mockDBConnection, _ := getGoodNewConnection()

	// Test
	sut, err := NewQuery(mockDBConnection, NewSQLQuery("bogus query"))

	// Verify
	ExpectNonNil(sut, t)
	ExpectNoError(err, t)
}

func TestThat_Run_ReturnsResult_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
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

func TestThat_Run_ReturnsError_WhenQueryExecutionFailsWithError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.Run()

	// Verify
	ExpectNil(actual, t)
	ExpectError(err2, t)
}

func TestThat_RunReturnInt_ReturnsInt_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	mock := mockDB.GetMock()

        expectedValue := 111
	expectedRows := sqlmock.NewRows([]string{"ResultID"}).AddRow(expectedValue)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.RunReturnInt()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err2, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
	ExpectInt(expectedValue, *actual, t)
}

func TestThat_RunReturnInt_ReturnsError_WhenQueryExecutionFailsWithError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	mock := mockDB.GetMock()

	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.RunReturnInt()

	// Verify
	ExpectNil(actual, t)
	ExpectError(err2, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func TestThat_RunReturnString_ReturnsString_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	mock := mockDB.GetMock()

        expectedValue := "success!"
	expectedRows := sqlmock.NewRows([]string{"ResultStr"}).AddRow(expectedValue)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.RunReturnString()

	// Verify
	ExpectNonNil(actual, t)
	ExpectNoError(err2, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
	ExpectString(expectedValue, *actual, t)
}

func TestThat_RunReturnString_ReturnsError_WhenQueryExecutionFailsWithError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	mock := mockDB.GetMock()

	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.RunReturnString()

	// Verify
	ExpectNil(actual, t)
	ExpectError(err2, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

// RunReturnOne(args ...interface{}) (*ResultRow, error)
func TestThat_RunReturnOne_ReturnsResultRow_WithoutError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	mock := mockDB.GetMock()

        expectedValueStr := "success!"
        var expectedValueInt int64 = 111
	expectedRows := sqlmock.NewRows([]string{"ResultInt", "ResultStr"}).AddRow(expectedValueInt, expectedValueStr)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.RunReturnOne()

	// Verify
	ExpectNoError(err2, t)
	ExpectNonNil(actual, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)

	actualValueStrNullable := actual.Get("ResultStr")
	ExpectNonNil(actualValueStrNullable, t)
	actualValueStr := actualValueStrNullable.GetString()
	ExpectNonNil(actualValueStr, t)
	ExpectString(expectedValueStr, *actualValueStr, t)

	actualValueIntNullable := actual.Get("ResultInt")
	ExpectNonNil(actualValueIntNullable, t)
	actualValueInt := actualValueIntNullable.GetInt64()
	ExpectNonNil(actualValueInt, t)
	ExpectInt64(expectedValueInt, *actualValueInt, t)

}

func TestThat_RunReturnOne_ReturnsError_WhenQueryExecutionFailsWithError(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := getGoodNewConnection()
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	mock := mockDB.GetMock()

	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, err1 := NewQuery(mockDBConnection, NewSQLQuery(expectedQuery))

	// Test
	ExpectNoError(err1, t)
	ExpectNonNil(sut, t)
	actual, err2 := sut.RunReturnOne()

	// Verify
	ExpectNil(actual, t)
	ExpectError(err2, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}
