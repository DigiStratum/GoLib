package mysql

import(
	"github.com/DigiStratum/GoLib/DB"
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMySQLConnectionFactory(driver string) *db.DBConnectionFactory {
	return db.NewDBConnectionFactory("mysql")
}
