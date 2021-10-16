package mockdb

import(
	"fmt"
	"database/sql"

        "github.com/DATA-DOG/go-sqlmock"
)

// ref: https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e

type mockInfo struct {
	Conn		*sql.DB
	Mock		*sqlmock.Sqlmock
}

type mockDBConnection struct {
	mocks		map[string]mockInfo
}

var instance *mockDBConnection

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMockDBConnection(driverName, dataSourceName string) (*sql.DB, error) {
	conn, mock, err := sqlmock.New()
	if nil != err { return nil, err }
	key := fmt.Sprintf("%s:%s", driverName, dataSourceName)
	setDBConnectionMockInfo(key, conn, &mock)
	return conn, nil
}

func setDBConnectionMockInfo(key string, conn *sql.DB, mock *sqlmock.Sqlmock) {
	i := getInstance()
	i.mocks[key] = mockInfo{
		Conn:	conn,
		Mock:	mock,
	}
}

func GetDBConnectionMockInfo(driverName, dataSourceName string) *mockInfo {
	key := fmt.Sprintf("%s:%s", driverName, dataSourceName)
	i := getInstance()
	mocks := i.mocks
	if mi, ok := mocks[key]; ok {
		return &mi
	}
	return nil
}

// Singleton pattern for our mockDbConnection state into
func getInstance() *mockDBConnection {
	if nil == instance {
		instance = &mockDBConnection{
			mocks:		make(map[string]mockInfo),
		}
	}
	return instance
}
