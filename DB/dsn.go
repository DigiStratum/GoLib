package dsn

/*
A library of functions to deal with Data Source Name specifiers

ref: https://en.wikipedia.org/wiki/Data_source_name

Stateless, no point in making this an object/class
*/

import (
	"fmt"
	"crypto/md5"
)

func MakeDSN (user, pass, host, port, name string) string {
	// ref: https://stackoverflow.com/questions/23550453/golang-how-to-open-a-remote-mysql-connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)
        return dsn
}

func GetDSNHash(dsn string) string {
	// ref: https://golang.org/pkg/crypto/md5/
	data := []byte(dsn)
	return fmt.Sprintf("%x", md5.Sum(data))
}

