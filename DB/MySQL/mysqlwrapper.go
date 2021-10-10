package mysql

/*
Abstraction of MySQL db driver types so that we can use Dependency Injection to enable unit testing

*/
import(
	db "database/sql"
)

type MySQLStatement db.Stmt
type MySQLResult db.Result
type MySQLRows db.Rows
type MySQLRow db.Row

type MySQLConnectionIfc interface {
	IsConnected() bool
	Connect(dsn string) error
	Disconnect() error
	Begin() error
	Prepare(query string) (*MySQLStatement, error)
	Exec(query string, args ...interface{}) (MySQLResult, error)
	Query(query string, args ...interface{}) (*MySQLRows, error)
	QueryRow(query string, args ...interface{}) *MySQLRow
}

type MySQLConnection struct {
}

type MySQLConnectionFactoryIfc interface {
	GetConnection() MySQLConnection
}

type MySQLConnectionFactory struct {
}

func NewMySQLConnectionFactory() *MySQLConnectionFactory {
	return &MySQLConnectionFactory{}
}

func (r *MySQLConnectionFactory) GetConnection() MySQLConnection {
}