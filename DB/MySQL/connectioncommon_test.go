package mysql

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

var driverName string = "mockdriver"

func testThat_ConnectionCommon_Connection_InTransaction_ReturnsFalse_WhenNotInTransaction(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	res := sut.InTransaction()

	// Verify
	ExpectFalse(res, t)
}

func testThat_ConnectionCommon_Connection_InTransaction_ReturnsTrue_WhenInTransaction(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	err := sut.Begin()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err, t)
	ExpectTrue(res, t)
}

func testThat_ConnectionCommon_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenRollback(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	err1 := sut.Begin()
	err2 := sut.Rollback()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectFalse(res, t)
}

func testThat_ConnectionCommon_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenCommit(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	err1 := sut.Begin()
	err2 := sut.Commit()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectFalse(res, t)
}

func testThat_ConnectionCommon_Connection_Begin_ReturnsNoError_WhenCalledTwice(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Test
	err1 := sut.Begin()
	mock := mockDB.GetMock()
	(*mock).ExpectBegin()
	err2 := sut.Begin()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectTrue(res, t)
}

// TODO: Figure out whether this is passing because the mock connection is closed, or because
// of ExpectBegin(), or both or neither; We want the result of Begin() on a Closed connection
func testThat_ConnectionCommon_Connection_Begin_ReturnsError_WhenConnectionClosed(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	err := sut.Begin()
	res := sut.InTransaction()

	// Verify
	ExpectError(err, t)
	ExpectFalse(res, t)
}

func testThat_ConnectionCommon_Connection_NewQuery_ReturnsQueryNoError(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	res, err := sut.NewQuery("bogus query")

	// Verify
	ExpectNonNil(res, t)
	ExpectNoError(err, t)
}

func testThat_ConnectionCommon_Connection_Commit_ReturnsError_WhenNotInTransaction(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	err := sut.Commit()

	// Verify
	ExpectError(err, t)
}

func testThat_ConnectionCommon_Connection_Commit_ReturnsError_WhenErrorOnTransactionCommit(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedMsg := "bogus failure"
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	err := sut.Commit()

	// Verify
	ExpectError(err, t)
	ExpectString(expectedMsg, err.Error(), t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Commit_ReturnsNoError_WhenTransactionCommits(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	err := sut.Commit()

	// Verify
	ExpectNoError(err, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Rollback_ReturnsNoError_WhenNotInTransaction(sut ConnectionCommonIfc, t *testing.T) {
	// Test
	err := sut.Rollback()

	// Verify
	ExpectNoError(err, t)
}

func testThat_ConnectionCommon_Connection_Rollback_ReturnsError_WhenErrorOnTransactionRollback(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedMsg := "bogus failure"
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	err := sut.Rollback()

	// Verify
	ExpectError(err, t)
	ExpectString(expectedMsg, err.Error(), t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Rollback_ReturnsNoError_WhenTransactionRollsBack_OutsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Test
	err := sut.Rollback()
	mock := mockDB.GetMock()

	// Verify
	ExpectNoError(err, t)
	// Expect this expectation to fail: Rollback() should NOT be called when no transaction is active
	ExpectError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithoutArgs_OutsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	actualResult, err := sut.Exec(expectedQuery)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualResult, t)
	actualInsertId, _ := actualResult.LastInsertId()
	ExpectInt64(expectedInsertId, actualInsertId, t)
	actualAffectedRows, _ := actualResult.RowsAffected()
	ExpectInt64(expectedAffectedRows, actualAffectedRows, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	actualResult, err := sut.Exec(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithoutArgs_InsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

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
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualResult, err := sut.Exec(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithArgs_OutsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedArgs := 333
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	actualResult, err := sut.Exec(expectedQuery, expectedArgs)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualResult, t)
	actualInsertId, _ := actualResult.LastInsertId()
	ExpectInt64(expectedInsertId, actualInsertId, t)
	actualAffectedRows, _ := actualResult.RowsAffected()
	ExpectInt64(expectedAffectedRows, actualAffectedRows, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedArgs := 333
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	actualResult, err := sut.Exec(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithArgs_InsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

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
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualResult, err := sut.Exec(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualResult, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithoutArgs_OutsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithoutArgs_InsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithArgs_OutsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

	// Test
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

	// Test
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithArgs_InsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRows, t)
	ExpectTrue(actualRows.Next(), t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_Query_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualRows, err := sut.Query(expectedQuery, expectedArgs)

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRows, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithoutArgs_OutsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	actualRow := sut.QueryRow(expectedQuery)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithoutArgs_InsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualRow := sut.QueryRow(expectedQuery)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithArgs_OutsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

	// Test
	actualRow := sut.QueryRow(expectedQuery, expectedArgs)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}

func testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithArgs_InsideTransaction(sut ConnectionCommonIfc, mockDB MockDBConnectionIfc, t *testing.T) {
	// Setup
	expectedQuery := "bogus query"
	expectedArgs := 333
	mock := mockDB.GetMock()

	// Test
	sut.Begin()
	actualRow := sut.QueryRow(expectedQuery, expectedArgs)
	err := actualRow.Err()

	// Verify
	ExpectNonNil(actualRow, t)
	ExpectNoError(err, t)
	ExpectNoError((*mock).ExpectationsWereMet(), t)
}
