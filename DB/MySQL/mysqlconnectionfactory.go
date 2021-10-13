package mysql

import(
	"github.com/DigiStratum/GoLib/DB"
)

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewMySQLConnectionFactory(driver string) *DBConnectionFactory {
	return NewDBConnectionFactory("mysql")
}
