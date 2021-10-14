package mocks

import(
	"time"
	"context"
	"database/sql"
        "database/sql/driver"
)

type MockDBConnection struct {}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMockDBConnection(driverName, dataSourceName string) (*MockDBConnection, error) {
	return &MockDBConnection{}, nil
}

// -------------------------------------------------------------------------------------------------
// db.DBConnectionIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *MockDBConnection) Begin() (*sql.Tx, error) {
	return nil, nil
}

func (r *MockDBConnection) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return nil, nil
}

func (r *MockDBConnection) Close() error {
	return nil
}

func (r *MockDBConnection) Conn(ctx context.Context) (*sql.Conn, error) {
	return nil, nil
}

func (r *MockDBConnection) Driver() driver.Driver {
	return nil
}

func (r *MockDBConnection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (r *MockDBConnection) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (r *MockDBConnection) Ping() error {
	return nil
}

func (r *MockDBConnection) PingContext(ctx context.Context) error {
	return nil
}

func (r *MockDBConnection) Prepare(query string) (*sql.Stmt, error) {
	return nil, nil
}

func (r *MockDBConnection) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return nil, nil
}

func (r *MockDBConnection) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (r *MockDBConnection) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (r *MockDBConnection) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}

func (r *MockDBConnection) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

func (r *MockDBConnection) SetConnMaxIdleTime(d time.Duration) {
}

func (r *MockDBConnection) SetConnMaxLifetime(d time.Duration) {
}

func (r *MockDBConnection) SetMaxIdleConns(n int) {
}

func (r *MockDBConnection) SetMaxOpenConns(n int) {
}

func (r *MockDBConnection) Stats() sql.DBStats {
	return sql.DBStats{}
}
