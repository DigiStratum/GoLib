package mysql

import(
	"github.com/DigiStratum/GoLib/DB"
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMySQLConnectionFactory() *db.DBConnectionFactory {
	return db.NewDBConnectionFactory("mysql")
}
