package db

/*
Interface for DB Connections; without this, then all usages of a raw DB connection expect this
explicit struct instead of interface which prevents creating mocks/stubs for testing.
*/

import(
	"database/sql"
)

// Clone of the set of member functions called out in: https://pkg.go.dev/database/sql#DB
type DBConnectionIfc interface {
	Begin() (*Tx, error)
	BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error)
	Close() error
	Conn(ctx context.Context) (*Conn, error)
	Driver() driver.Driver
	Exec(query string, args ...interface{}) (Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Ping() error
	PingContext(ctx context.Context) error
	Prepare(query string) (*Stmt, error)
	PrepareContext(ctx context.Context, query string) (*Stmt, error)
	Query(query string, args ...interface{}) (*Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*Rows, error)
	QueryRow(query string, args ...interface{}) *Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *Row
	SetConnMaxIdleTime(d time.Duration)
	SetConnMaxLifetime(d time.Duration)
	SetMaxIdleConns(n int)
	SetMaxOpenConns(n int)
	Stats() DBStats
}

type DBConnection struct {
	sql.DB
}

func NewDBConnection(driverName, dataSourceName string) *DBConnection {
	return &DBConnection{
		sql.DB: 	sql.Open(driverName, dataSourceName),
	}
}