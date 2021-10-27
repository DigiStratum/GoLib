package db

import(
	"database/sql"
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDBConnection(driverName, dsn string) (*sql.DB, error) {
	dbconn, err := sql.Open(driverName, dsn)
	if nil != err { return nil, err }
	return dbconn, nil
}