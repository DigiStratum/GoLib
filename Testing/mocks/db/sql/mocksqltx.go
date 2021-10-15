package mockdbsql

import(
	"time"
	"context"
	"database/sql"
        "database/sql/driver"
)

func NewMockTx() *sql.Tx {
	mockTx := sql.Tx{}
	return &mockTx
}

