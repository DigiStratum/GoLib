package mysql

import db "github.com/DigiStratum/GoLib/DB"

// DB Key public Interface
type DBKeyIfc interface {
	GetKey() string
}

// DB Key (Connection identifier for Manager)
type dBKey struct {
	key	string
}

// Make a new one of these from an existing key
func NewDBKey(key string) DBKeyIfc {
	return &dBKey{ key: key }
}

// Make a new one of these from DSN
func NewDBKeyFromDSN(dsn string) DBKeyIfc {
	return &dBKey{ key: db.GetDSNHash(dsn) }
}

// Get the key
func (dbk *dBKey) GetKey() string {
	return dbk.key
}

