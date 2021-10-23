package mysql

import(
	"fmt"
	"testing"

        "github.com/DATA-DOG/go-sqlmock"

	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

var driverName string = "mockdriver"
var dataSourceName string = "mockdsn"

func TestThat_NewConnection_ReturnsError_WhenGivenNilConnection(t *testing.T) {
	// Test
	sut, err := NewConnection(nil)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

func TestThat_NewConnection_ReturnsConnection_WhenGivenDBConnection(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)

	// Test
	sut, err := NewConnection(mockDBConnection)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(sut, t)
}

func TestThat_Connection_IsConnected_ReturnsTrue_WhenConnected(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
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
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
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
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res := sut.InTransaction()

	// Verify
	ExpectFalse(res, t)
}

func TestThat_Connection_InTransaction_ReturnsTrue_WhenInTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err := sut.Begin()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err, t)
	ExpectTrue(res, t)
}

func TestThat_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenRollback(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	(*mockInfo.Mock).ExpectRollback()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err1 := sut.Begin()
	err2 := sut.Rollback()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectFalse(res, t)
}

func TestThat_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenCommit(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	(*mockInfo.Mock).ExpectCommit()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err1 := sut.Begin()
	err2 := sut.Commit()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectFalse(res, t)
}

func TestThat_Connection_Begin_ReturnsNoError_WhenCalledTwice(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	(*mockInfo.Mock).ExpectRollback()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err1 := sut.Begin()
	(*mockInfo.Mock).ExpectBegin()
	err2 := sut.Begin()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectTrue(res, t)
}

// TODO: Figure out if this is passing because the mock connection is closed, or because of
// ExpectBegin(), or both or neither; We want the result of Begin() on a Closed connection
func TestThat_Connection_Begin_ReturnsError_WhenConnectionClosed(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	mockDBConnection.Close()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err := sut.Begin()
	res := sut.InTransaction()

	// Verify
	ExpectError(err, t)
	ExpectFalse(res, t)
}

func TestThat_Connection_NewQuery_ReturnsQueryNoError(t *testing.T) {
	// Setup
	query := "bogus query"
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectPrepare(query)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res, err := sut.NewQuery("bogus query")

	// Verify
	ExpectNonNil(res, t)
	ExpectNoError(err, t)
}

func TestThat_Connection_Commit_ReturnsError_WhenNotInTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err := sut.Commit()

	// Verify
	ExpectError(err, t)
}

func TestThat_Connection_Commit_ReturnsError_WhenErrorOnTransactionCommit(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedMsg := "bogus failure"
	(*mockInfo.Mock).ExpectCommit().WillReturnError(fmt.Errorf(expectedMsg))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	err := sut.Commit()

	// Verify
	ExpectError(err, t)
	ExpectString(expectedMsg, err.Error(), t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Commit_ReturnsNoError_WhenTransactionCommits(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	(*mockInfo.Mock).ExpectCommit()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	err := sut.Commit()

	// Verify
	ExpectNoError(err, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Rollback_ReturnsNoError_WhenNotInTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err := sut.Rollback()

	// Verify
	ExpectNoError(err, t)
}

func TestThat_Connection_Rollback_ReturnsError_WhenErrorOnTransactionRollback(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedMsg := "bogus failure"
	(*mockInfo.Mock).ExpectRollback().WillReturnError(fmt.Errorf(expectedMsg))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	err := sut.Rollback()

	// Verify
	ExpectError(err, t)
	ExpectString(expectedMsg, err.Error(), t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Rollback_ReturnsNoError_WhenTransactionRollsBack_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectRollback()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err := sut.Rollback()

	// Verify
	ExpectNoError(err, t)
	// Expect this expectation to fail: Rollback() should NOT be called when no transaction is active
	ExpectError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithoutArgs_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualResult, err := sut.Exec(expectedQuery)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualResult, t)
	actualInsertId, _ := actualResult.LastInsertId()
	ExpectInt64(expectedInsertId, actualInsertId, t)
	actualAffectedRows, _ := actualResult.RowsAffected()
	ExpectInt64(expectedAffectedRows, actualAffectedRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualResult, err := sut.Exec(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithoutArgs_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualResult, err := sut.Exec(expectedQuery)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualResult, t)
	actualInsertId, _ := actualResult.LastInsertId()
	ExpectInt64(expectedInsertId, actualInsertId, t)
	actualAffectedRows, _ := actualResult.RowsAffected()
	ExpectInt64(expectedAffectedRows, actualAffectedRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualResult, err := sut.Exec(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithArgs_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs(expectedArgs).
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualResult, err := sut.Exec(expectedQuery, expectedArgs)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualResult, t)
	actualInsertId, _ := actualResult.LastInsertId()
	ExpectInt64(expectedInsertId, actualInsertId, t)
	actualAffectedRows, _ := actualResult.RowsAffected()
	ExpectInt64(expectedAffectedRows, actualAffectedRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs(expectedArgs).
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualResult, err := sut.Exec(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithArgs_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualResult, err := sut.Exec(expectedQuery, expectedArgs)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualResult, t)
	actualInsertId, _ := actualResult.LastInsertId()
	ExpectInt64(expectedInsertId, actualInsertId, t)
	actualAffectedRows, _ := actualResult.RowsAffected()
	ExpectInt64(expectedAffectedRows, actualAffectedRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Exec_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualResult, err := sut.Exec(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithoutArgs_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithoutArgs_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithArgs_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithArgs_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Query_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithoutArgs_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualRow := sut.QueryRow(expectedQuery)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithoutArgs_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualRow := sut.QueryRow(expectedQuery)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithArgs_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	actualRow := sut.QueryRow(expectedQuery, expectedArgs)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithArgs_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockInfo.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	actualRow := sut.QueryRow(expectedQuery, expectedArgs)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mockInfo.Mock).ExpectationsWereMet(), t)
}
