package db

import(
	"database/sql"
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDBConnection(driverName string, dsn DSNIfc) (*sql.DB, error) {
	dbconn, err := sql.Open(driverName, dsn.ToString())
	if nil != err { return nil, err }
	return dbconn, nil
}
