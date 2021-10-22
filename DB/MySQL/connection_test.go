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
	ExpectNonNil(err, t)
}

func TestThat_NewConnection_ReturnsConnection_WhenGivenDBConnection(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)

	// Test
	sut, err := NewConnection(mockDBConnection)

	// Verify
	ExpectNil(err, t)
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
	ExpectNil(err, t)
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
	ExpectNil(err1, t)
	ExpectNil(err2, t)
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
	ExpectNil(err1, t)
	ExpectNil(err2, t)
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
	ExpectNil(err1, t)
	ExpectNil(err2, t)
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
	ExpectNonNil(err, t)
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
	ExpectNil(err, t)
}

func TestThat_Connection_NewQuery_ReturnsError(t *testing.T) {
	// Setup
	query := "bogus query"
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectPrepare(query).WillReturnError(fmt.Errorf("bogus error"))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	res, err := sut.NewQuery("bogus query")

	// Verify
	ExpectNil(res, t)
	ExpectNonNil(err, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Commit_ReturnsError_WhenNotInTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err := sut.Commit()

	// Verify
	ExpectNonNil(err, t)
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
	ExpectNonNil(err, t)
	ExpectString(expectedMsg, err.Error(), t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
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
	ExpectNil(err, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Rollback_ReturnsNoError_WhenNotInTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	err := sut.Rollback()

	// Verify
	ExpectNil(err, t)
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
	ExpectNonNil(err, t)
	ExpectString(expectedMsg, err.Error(), t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
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
	ExpectNil(err, t)
	// Expect this expectation to fail: Rollback() should NOT be called when no transaction is active
	ExpectNonNil((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Prepare_ReturnsStatementNoError_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	query := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(query)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	stmt, err := sut.Prepare(query)

	// Verify
	ExpectNil(err, t)
	ExpectNonNil(stmt, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Prepare_ReturnsStatementNoError(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	query := "bogus query"
	(*mockInfo.Mock).ExpectPrepare(query)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	stmt, err := sut.Prepare(query)

	// Verify
	ExpectNil(err, t)
	ExpectNonNil(stmt, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_Prepare_ReturnsError(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	query := "bogus query"
	expectedMsg := "bogus failure"
	(*mockInfo.Mock).ExpectPrepare(query).WillReturnError(fmt.Errorf(expectedMsg))
	sut, _ := NewConnection(mockDBConnection)

	// Test
	stmt, err := sut.Prepare(query)

	// Verify
	ExpectNonNil(err, t)
	ExpectString(expectedMsg, err.Error(), t)
	ExpectNil(stmt, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_StmtExec_ReturnsResultNoError_OutsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	query := "bogus query"
	mockPrepare := (*mockInfo.Mock).ExpectPrepare(query)
	var expectedInsertId int64 = 10
	var expectedAffectedRows int64 = 20
	result := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	mockPrepare.ExpectExec().WillReturnResult(result)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	stmt, err1 := sut.Prepare(query)
	ExpectNoError(err1, t)
	ExpectNonNil(stmt, t)
	res, err2 := sut.StmtExec(stmt)
	actualInsertId, err3 := res.LastInsertId()
	ExpectNoError(err3, t)
	actualRowsAffected, err4 := res.RowsAffected()
	ExpectNoError(err4, t)

	// Verify
	ExpectNoError(err2, t)
	ExpectNonNil(res, t)
	ExpectInt64(expectedInsertId, actualInsertId, t)
	ExpectInt64(expectedAffectedRows, actualRowsAffected, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
}

func TestThat_Connection_StmtExec_ReturnsResultNoError_InsideTransaction(t *testing.T) {
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	query := "bogus query"
	mockPrepare := (*mockInfo.Mock).ExpectPrepare(query)
	result := sqlmock.NewResult(1, 1)
	mockPrepare.ExpectExec().WillReturnResult(result)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	stmt, err1 := sut.Prepare(query)
	ExpectNoError(err1, t)
	ExpectNonNil(stmt, t)
	res, err2 := sut.StmtExec(stmt)

	// Verify
	ExpectNoError(err2, t)
	ExpectNonNil(res, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)

/*
	// Setup
	mockDBConnection, _ := NewMockDBConnection(driverName, dataSourceName)
	mockInfo := GetDBConnectionMockInfo(driverName, dataSourceName)
	(*mockInfo.Mock).ExpectBegin()
	query := "bogus query"
	mockPrepare := (*mockInfo.Mock).ExpectPrepare(query)
	var expectedInsertId int64 = 10
	var expectedAffectedRows int64 = 20
	result := sqlmock.NewResult(expectedInsertId, expectedAffectedRows)
	mockPrepare.ExpectExec().WillReturnResult(result)
	sut, _ := NewConnection(mockDBConnection)

	// Test
	sut.Begin()
	stmt, err1 := sut.Prepare(query)
	ExpectNoError(err1, t)
	ExpectNonNil(stmt, t)
	res, err2 := sut.StmtExec(stmt)

	// Verify
	ExpectNoError(err2, t)
	ExpectNonNil(res, t)
	ExpectNil((*mockInfo.Mock).ExpectationsWereMet(), t)
*/
}
