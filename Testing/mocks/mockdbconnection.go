package mocks

import(
	"fmt"
	"database/sql"

        "github.com/DATA-DOG/go-sqlmock"

	"github.com/DigiStratum/GoLib/DB"
)

// ref: https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e
// ref: https://dev.to/techschoolguru/mock-db-for-testing-http-api-in-go-and-achieve-100-coverage-4pa9

type MockDBConnectionIfc interface {
	GetConn() *sql.DB
	GetMock() *sqlmock.Sqlmock
}

type mockDBConnection struct {
	Conn		*sql.DB
	Mock		*sqlmock.Sqlmock
}

type mockDBConnections struct {
	mocks		map[string]mockDBConnection
}

var instance *mockDBConnections

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMockDBConnection(driverName string, dsn db.DSNIfc) (*sql.DB, error) {
	conn, mock, err := sqlmock.New()
	if nil != err { return nil, err }
	key := fmt.Sprintf("%s:%s", driverName, dsn.ToHash())
	i := getInstance()
	i.mocks[key] = mockDBConnection{
		Conn:	conn,
		Mock:	&mock,
	}

	return conn, nil
}

func GetDBConnectionMockInfo(driverName string, dsn db.DSNIfc) *mockDBConnection {
	key := fmt.Sprintf("%s:%s", driverName, dsn.ToHash())
	i := getInstance()
	mocks := i.mocks
	if mi, ok := mocks[key]; ok {
		return &mi
	}
	return nil
}

// Singleton pattern for our mockDbConnection state into
func getInstance() *mockDBConnections {
	if nil == instance {
		instance = &mockDBConnections{
			mocks:		make(map[string]mockDBConnection),
		}
	}
	return instance
}

// -------------------------------------------------------------------------------------------------
// MockDBConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r mockDBConnection) GetConn() *sql.DB {
	return r.Conn
}

func (r mockDBConnection) GetMock() *sqlmock.Sqlmock {
	return r.Mock
}
