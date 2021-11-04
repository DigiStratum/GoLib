package db

/*
Data Source Name

ref: https://stackoverflow.com/questions/23550453/golang-how-to-open-a-remote-mysql-connection
ref: https://en.wikipedia.org/wiki/Data_source_name
ref: https://github.com/go-sql-driver/mysql/blob/master/dsn.go

old DNS function: dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)

*/

import (
	"fmt"
	"crypto/md5"
	"github.com/go-sql-driver/mysql"
)

type DSNIfc interface {
	ToHash() string
	ToString() string
}

type DSN struct {
	dsnString			string
	dsnConfig			*mysql.Config
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewDSN(dsn string) (*DSN, error) {
	dsnConfig, err := mysql.ParseDSN(dsn)
	if nil != err { return nil, err }
	return &DSN{
		dsnString:	dsn,
		dsnConfig:	dsnConfig,
	}, nil
}

// -------------------------------------------------------------------------------------------------
// DSNIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Create a unique hash of this DSN so that we can log/associate it without revealing secrets
func (r DSN) ToHash() string {
	// ref: https://golang.org/pkg/crypto/md5/
	data := []byte(r.dsnString)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func (r DSN) ToString() string {
	return r.dsnString
}
