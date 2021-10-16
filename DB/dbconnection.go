package db

import(
	"database/sql"
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDBConnection(driverName, dataSourceName string) (*sql.DB, error) {
	dbconn, err := sql.Open(driverName, dataSourceName)
	if nil != err { return nil, err }
	return dbconn, nil
}