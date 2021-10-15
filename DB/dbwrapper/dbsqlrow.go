package dbwrapper

type DBSqlRowIfc interface {
	Err() error
	Scan(dest ...interface{}) error
}