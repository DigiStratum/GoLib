package mockdb

import(
	"fmt"
	"database/sql"

        "github.com/DATA-DOG/go-sqlmock"

	"github.com/DigiStratum/GoLib/DB"
)

// ref: https://medium.com/easyread/unit-test-sql-in-golang-5af19075e68e
// ref: https://dev.to/techschoolguru/mock-db-for-testing-http-api-in-go-and-achieve-100-coverage-4pa9

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

func NewMockDBConnection(driverName string, dsn db.DSNIfc) (*sql.DB, error) {
	conn, mock, err := sqlmock.New()
	if nil != err { return nil, err }
	key := fmt.Sprintf("%s:%s", driverName, dsn.ToHash())
	i := getInstance()
	i.mocks[key] = mockInfo{
		Conn:	conn,
		Mock:	&mock,
	}

	return conn, nil
}

func GetDBConnectionMockInfo(driverName string, dsn db.DSNIfc) *mockInfo {
	key := fmt.Sprintf("%s:%s", driverName, dsn.ToHash())
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
