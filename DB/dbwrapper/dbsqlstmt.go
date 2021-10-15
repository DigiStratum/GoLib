package dbwrapper

import(
	"context"
	"database/sql"
)

type DBSqlStmtIfc interface {
	Close() error
	Exec(args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
	Query(args ...interface{}) (DBSQLRowsIfc, error)
	QueryContext(ctx context.Context, args ...interface{}) (DBSQLRowsIfc, error)
}