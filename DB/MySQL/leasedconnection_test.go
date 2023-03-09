package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB"
	dep "github.com/DigiStratum/GoLib/Dependencies"
	cfg "github.com/DigiStratum/GoLib/Config"

	. "github.com/DigiStratum/GoLib/Testing"
	. "github.com/DigiStratum/GoLib/Testing/mocks"
)

func TestThat_NewLeasedConnection_ReturnsSomething(t *testing.T) {
	// Setup
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	connectionPool := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	connectionPool.InjectDependencies(
		dep.NewDependencyInstance("connectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	connectionPool.Configure(config)

	conn, _ := connectionFactory.NewConnection(dsn)
	newConnection, _ := NewConnection(conn)
	newPooledConnection, _ := NewPooledConnection(newConnection, connectionPool)

	leaseKey := int64(333)

	// Test
	var actual *LeasedConnection = NewLeasedConnection(newPooledConnection, leaseKey)

	// Verify
	ExpectNonNil(actual, t)
}

func TestThat_LeasedConnection_Release_Returns_WithoutError(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
}

func TestThat_LeasedConnection_Release_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()
	err2 := sut.Release()

	// Test
	err3 := sut.Release()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectError(err3, t)
}

func TestThat_LeasedConnection_ConnectionIfc_IsConnected_ReturnsTrue_ForGoodLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	res := sut.IsConnected()

	// Verify
	ExpectNoError(err1, t)
	ExpectTrue(res, t)
}

func TestThat_LeasedConnection_ConnectionIfc_IsConnected_ReturnsFalse_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	res := sut.IsConnected()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectFalse(res, t)
}

func TestThat_LeasedConnection_ConnectionIfc_Disconnect_LeavesLeasedConnectionIntact(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	sut.Disconnect()
	res := sut.IsConnected()

	// Verify
	ExpectNoError(err1, t)
	ExpectTrue(res, t)
}

func TestThat_LeasedConnection_ConnectionIfc_Connect_ReturnsError(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Connect()
	res := sut.IsConnected()

	// Verify
	ExpectNoError(err1, t)
	ExpectError(err2, t)
	ExpectTrue(res, t)
}

func TestThat_LeasedConnection_ConnectionIfc_InTransaction_ReturnsFalse_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	res := sut.InTransaction()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectFalse(res, t)
}

func TestThat_LeasedConnection_ConnectionIfc_Begin_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	err3 := sut.Begin()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectError(err3, t)
}

func TestThat_LeasedConnection_ConnectionIfc_NewQuery_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	res, err3 := sut.NewQuery(NewSQLQuery("bogus sql"))

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectNil(res, t)
	ExpectError(err3, t)
}

func TestThat_LeasedConnection_ConnectionIfc_Commit_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	err3 := sut.Commit()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectError(err3, t)
}

func TestThat_LeasedConnection_ConnectionIfc_Rollback_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	err3 := sut.Rollback()

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectError(err3, t)
}

func TestThat_LeasedConnection_ConnectionIfc_Exec_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	res, err3 := sut.Exec(NewSQLQuery("bogus sql"))

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectNil(res, t)
	ExpectError(err3, t)
}

func TestThat_LeasedConnection_ConnectionIfc_Query_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	res, err3 := sut.Query(NewSQLQuery("bogus sql"))

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectNil(res, t)
	ExpectError(err3, t)
}

func TestThat_LeasedConnection_ConnectionIfc_QueryRow_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection()

	// Test
	err2 := sut.Release()
	res := sut.QueryRow(NewSQLQuery("bogus sql"))

	// Verify
	ExpectNoError(err1, t)
	ExpectNoError(err2, t)
	ExpectNil(res, t)
}

func getGoodLeasedConnection() (*LeasedConnection, error) {
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	connectionPool := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	connectionPool.InjectDependencies(
		dep.NewDependencyInstance("connectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	connectionPool.Configure(config)

	return connectionPool.GetConnection()
}

