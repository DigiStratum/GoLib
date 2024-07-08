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
	var actual LeasedConnectionIfc = NewLeasedConnection(newPooledConnection, leaseKey)	// <- Ensures that we satisfy our interface

	// Verify
	if ! ExpectNonNil(actual, t) { return }
}

func TestThat_LeasedConnection_Release_Returns_WithoutError(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
}

func TestThat_LeasedConnection_Release_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)
	err2 := sut.Release()

	// Test
	err3 := sut.Release()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectError(err3, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_IsConnected_ReturnsTrue_ForGoodLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	res := sut.IsConnected()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectTrue(res, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_IsConnected_ReturnsFalse_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	res := sut.IsConnected()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectFalse(res, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_Disconnect_LeavesLeasedConnectionIntact(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	sut.Disconnect()
	res := sut.IsConnected()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectTrue(res, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_Connect_ReturnsError(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Connect()
	res := sut.IsConnected()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectError(err2, t) { return }
	if ! ExpectTrue(res, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_InTransaction_ReturnsFalse_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	res := sut.InTransaction()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectFalse(res, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_Begin_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	err3 := sut.Begin()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectError(err3, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_NewQuery_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	res, err3 := sut.NewQuery(NewSQLQuery("bogus sql"))

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectNil(res, t) { return }
	if ! ExpectError(err3, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_Commit_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	err3 := sut.Commit()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectError(err3, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_Rollback_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	err3 := sut.Rollback()

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectError(err3, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_Exec_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	res, err3 := sut.Exec(NewSQLQuery("bogus sql"))

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectNil(res, t) { return }
	if ! ExpectError(err3, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_Query_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	res, err3 := sut.Query(NewSQLQuery("bogus sql"))

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectNil(res, t) { return }
	if ! ExpectError(err3, t) { return }
}

func TestThat_LeasedConnection_ConnectionIfc_QueryRow_ReturnsError_ForInvalidLease(t *testing.T) {
	// Setup
	sut, err1 := getGoodLeasedConnection(t)

	// Test
	err2 := sut.Release()
	res := sut.QueryRow(NewSQLQuery("bogus sql"))

	// Verify
	if ! ExpectNoError(err1, t) { return }
	if ! ExpectNoError(err2, t) { return }
	if ! ExpectNil(res, t) { return }
}

func getGoodLeasedConnection(t *testing.T) (*leasedConnection, error) {
	dsn, _ := db.NewDSN("user:pass@tcp(host:333)/name")
	connectionPool := NewConnectionPool(*dsn)
	connectionFactory := NewMockDBConnectionFactory()
	connectionPool.InjectDependencies(
		dep.NewDependencyInstance("ConnectionFactory", connectionFactory),
	)

	config := cfg.NewConfig()
	config.Set("min_connections", "1")
	config.Set("max_connections", "1")
	config.Set("max_idle", "1")
	connectionPool.Configure(config)
	err := connectionPool.Start()
	ExpectNoError(err, t)

	return connectionPool.GetConnection()
}

