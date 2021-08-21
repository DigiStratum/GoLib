package mysql

import db "github.com/DigiStratum/GoLib/DB"

// DB Key public Interface
type DBKeyIfc interface {
	GetKey() string
}

// DB Key (Connection identifier for Manager)
type DBKey struct {
	key	string
}

// Factory Functions

// Make a new one of these from an existing key
func NewDBKey(key string) *DBKey {
	return &DBKey{ key: key }
}

// Make a new one of these from DSN
func NewDBKeyFromDSN(dsn string) *DBKey {
	return NewDBKey(db.GetDSNHash(dsn))
}

// -------------------------------------------------------------------------------------------------
// DBKeyIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get the key
func (r DBKey) GetKey() string {
	return r.key
}

