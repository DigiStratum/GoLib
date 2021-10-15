package dbwrapper

/*
Interface for DB Connections; without this, then all usages of a raw DB connection expect this
explicit struct instead of interface which prevents creating mocks/stubs for testing.
*/

import(
	"time"
	"context"
	"database/sql"
        "database/sql/driver"
)

// Clone of the set of member functions called out in: https://pkg.go.dev/database/sql#DB
type DBConnectionIfc interface {
	Begin() (DBSQLTxIfc, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (DBSQLTxIfc, error)
	Close() error
	Conn(ctx context.Context) (*sql.Conn, error)
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Ping() error
	PingContext(ctx context.Context) error
	Prepare(query string) (DBSQLStmtIfc, error)
	PrepareContext(ctx context.Context, query string) (DBSQLStmtIfc, error)
	Query(query string, args ...interface{}) (DBSQLRowsIfc, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (DBSQLRowsIfc, error)
	QueryRow(query string, args ...interface{}) DBSQLRowIfc
	QueryRowContext(ctx context.Context, query string, args ...interface{}) DBSQLRowIfc
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() sql.DBStats
}

type DBConnection struct {
	sql.DB
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Return interface instead of struct so that other implementations may satisfy
func NewDBConnection(driverName, dataSourceName string) (DBConnectionIfc, error) {
	dbconn, err := sql.Open(driverName, dataSourceName)
	if nil != err { return nil, err }
	return &DBConnection{
		DB: 	*dbconn,
	}, nil
}