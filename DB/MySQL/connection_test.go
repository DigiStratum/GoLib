package mysql

import(
	"fmt"
	"testing"

        "github.com/DATA-DOG/go-sqlmock"

	"github.com/DigiStratum/GoLib/DB"
	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

func TestThat_NewConnection_ReturnsError_WhenGivenNilConnection(t *testing.T) {
	// Test
	sut, err := NewConnection(nil)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

func TestThat_NewConnection_ReturnsConnection_WhenGivenDBConnection(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)

	// Test
	sut, err := NewConnection(mockDBConnection)

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(sut, t)
}

func TestThat_Connection_IsConnected_ReturnsTrue_WhenConnected(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res := sut.IsConnected()

	// Verify
	ExpectTrue(res, t)
}

func TestThat_Connection_IsConnected_ReturnsFalse_WhenNotConnected(t *testing.T) {
	// Setup
	sut := connection{}

	// Test
	res := sut.IsConnected()

	// Verify
	ExpectFalse(res, t)
}

func TestThat_Connection_IsConnected_ReturnsFalse_WhenConnectedThenClosed(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res1 := sut.IsConnected()
	sut.Close()
	res2 := sut.IsConnected()

	// Verify
	ExpectTrue(res1, t)
	ExpectFalse(res2, t)
}

// -------------------------------------------------------------------------------------------------
// ConnectionCommonIfc Public Interface
//
// We offload the ConnectionCommonIfc compliance aspects so that we can reuse them to check other
// implementations of the same interface for compliance as well.
// -------------------------------------------------------------------------------------------------

func TestThat_Connection_InTransaction_ReturnsFalse_WhenNotInTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_InTransaction_ReturnsFalse_WhenNotInTransaction(sut, t)
}

func TestThat_Connection_InTransaction_ReturnsTrue_WhenInTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_InTransaction_ReturnsTrue_WhenInTransaction(sut, t)
}

func TestThat_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenRollback(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	(*mockDB.Mock).ExpectRollback()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenRollback(sut, t)
}

func TestThat_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenCommit(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	(*mockDB.Mock).ExpectCommit()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_InTransaction_ReturnsFalse_WhenInTransactionThenCommit(sut, t)
}

func TestThat_Connection_Begin_ReturnsNoError_WhenCalledTwice(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	(*mockDB.Mock).ExpectRollback()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Begin_ReturnsNoError_WhenCalledTwice(sut, mockDB, t)
}

// TODO: Figure out if this is passing because the mock connection is closed, or because of
// ExpectBegin(), or both or neither; We want the result of Begin() on a Closed connection
func TestThat_Connection_Begin_ReturnsError_WhenConnectionClosed(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	mockDBConnection.Close()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Begin_ReturnsError_WhenConnectionClosed(sut, t)
}

func TestThat_Connection_NewQuery_ReturnsQueryNoError(t *testing.T) {
	// Setup
	query := "bogus query"
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectPrepare(query)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_NewQuery_ReturnsQueryNoError(sut, t)
}

func TestThat_Connection_Commit_ReturnsError_WhenNotInTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Commit_ReturnsError_WhenNotInTransaction(sut, t)
}

func TestThat_Connection_Commit_ReturnsError_WhenErrorOnTransactionCommit(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedMsg := "bogus failure"
	(*mockDB.Mock).ExpectCommit().WillReturnError(fmt.Errorf(expectedMsg))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Commit_ReturnsError_WhenErrorOnTransactionCommit(sut, mockDB, t)
}

func TestThat_Connection_Commit_ReturnsNoError_WhenTransactionCommits(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	(*mockDB.Mock).ExpectCommit()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Commit_ReturnsNoError_WhenTransactionCommits(sut, mockDB, t)
}

func TestThat_Connection_Rollback_ReturnsNoError_WhenNotInTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Rollback_ReturnsNoError_WhenNotInTransaction(sut, t)
}

func TestThat_Connection_Rollback_ReturnsError_WhenErrorOnTransactionRollback(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedMsg := "bogus failure"
	(*mockDB.Mock).ExpectRollback().WillReturnError(fmt.Errorf(expectedMsg))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Rollback_ReturnsError_WhenErrorOnTransactionRollback(sut, mockDB, t)
}

func TestThat_Connection_Rollback_ReturnsNoError_WhenTransactionRollsBack_OutsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectRollback()
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Rollback_ReturnsNoError_WhenTransactionRollsBack_OutsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithoutArgs_OutsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithoutArgs_OutsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithoutArgs_InsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithoutArgs_InsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithArgs_OutsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs(expectedArgs).
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithArgs_OutsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs(expectedArgs).
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsResultNoError_WithArgs_InsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	var expectedInsertId int64 = 22
	var expectedAffectedRows int64 = 33
	expectedResult := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnResult(expectedResult)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsResultNoError_WithArgs_InsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Exec_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectExec().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Exec_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithoutArgs_OutsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithoutArgs_OutsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsError_WithoutArgs_OutsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithoutArgs_InsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithoutArgs_InsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsError_WithoutArgs_InsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithArgs_OutsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithArgs_OutsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsError_WithArgs_OutsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsRowsNoError_WithArgs_InsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsRowsNoError_WithArgs_InsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_Query_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_Query_ReturnsError_WithArgs_InsideTransaction_WhenPrepareFails(sut, mockDB, t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithoutArgs_OutsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithoutArgs_OutsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithoutArgs_InsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithoutArgs_InsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithArgs_OutsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	expectedArgs := 333
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs(expectedArgs).
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithArgs_OutsideTransaction(sut, mockDB, t)
}

func TestThat_Connection_QueryRow_ReturnsRow_WithArgs_InsideTransaction(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	mockDBConnection, _ := NewMockDBConnection(driverName, dsn)
	mockDB := GetDBConnectionMockInfo(driverName, dsn)
	(*mockDB.Mock).ExpectBegin()
	expectedCols := []string{ "boguscol" }
	expectedRows := sqlmock.NewRows(expectedCols).AddRow("bogusval")
	expectedQuery := "bogus query"
	(*mockDB.Mock).ExpectPrepare(expectedQuery).ExpectQuery().
		WithArgs().
		WillReturnRows(expectedRows)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	testThat_ConnectionCommon_Connection_QueryRow_ReturnsRow_WithArgs_InsideTransaction(sut, mockDB, t)
}
