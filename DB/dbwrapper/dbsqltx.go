package dbwrapper

import(
	"context"
	"database/sql"
)

type DBSQLTxIfc interface {
	Commit() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (DBSQLStmtIfc, error)
	PrepareContext(ctx context.Context, query string) (DBSQLStmtIfc, error)
	Query(query string, args ...interface{}) (DBSQLRowsIfc, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (DBSQLRowsIfc, error)
	QueryRow(query string, args ...interface{}) DBSQLRowIfc
	QueryRowContext(ctx context.Context, query string, args ...interface{}) DBSQLRowIfc
	Rollback() error
	Stmt(stmt DBSQLStmtIfc) DBSQLStmtIfc
	StmtContext(ctx context.Context, stmt DBSQLStmtIfc) DBSQLStmtIfc
}