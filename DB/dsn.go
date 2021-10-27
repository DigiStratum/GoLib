package db

/*
Data Source Name

ref: https://stackoverflow.com/questions/23550453/golang-how-to-open-a-remote-mysql-connection
ref: https://en.wikipedia.org/wiki/Data_source_name

*/

import (
	"fmt"
	"crypto/md5"
	"github.com/go-sql-driver/mysql" // ref: https://github.com/go-sql-driver/mysql/blob/master/dsn.go
)

type DSNIfc interface {
	ToHash() string
}

type DSN struct {
	dsnString			string
	dsnConfig			mysql.Config
}

func NewDSN(dsn string) (*DSN, error) {
	dsnConfig, err := mysql.ParseDSN(dsn)
	if nil != err { return nil, err }
	n := DSN{
		dsnString:	dsn,
		dsnConfig:	dsnConfig,
	}
	n.explodeParts()
	return &n
}

func NewDSNFromParts(user, pass, host, port, name string) (*DSN, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
	return NewDSN(dsn)
}

// Create a unique hash of this DSN so that we can log/associate it without revealing secrets
func (r DSN) ToHash() string {
	// ref: https://golang.org/pkg/crypto/md5/
	data := []byte(dsn)
	return fmt.Sprintf("%x", md5.Sum(data))
}
