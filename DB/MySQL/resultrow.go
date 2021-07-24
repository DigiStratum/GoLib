package mysql

/*
A single result row from a MySQL query result set

TODO: See if there is a way to encode each Nullable value as it's native JSON data type instead of making them all strings

Interesting:
 * http://go-database-sql.org/varcols.html
 * http://jmoiron.github.io/sqlx/
*/

import (
	nullables "github.com/DigiStratum/GoLib/DB/MySQL/nullables"
)

type ResultRowIfc interface {
	Get(field string) nullables.NullableIfc
	Set(field string, nullables.NullableIfc)
	Fields() []string
}

type resultRow map[string]nullables.NullableIfc

func NewResultRow() ResultRowIfc {
	rr := make(map[string]nullables.NullableIfc)
	return &rr
}

// -------------------------------------------------------------------------------------------------
// ResultRowIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (rr *resultRow) Get(field string) nullables.NullableIfc {
	if value, ok := (*rr)[field]; ok { return value }
	return nil
}

func (rr *resultRow) Set(field string, value nullables.NullableIfc) {
	(*rr)[field] = value
}

// Pluck the fields out of the result and just return them so that caller can iterate with Get()
func (rr *resultRow) Fields() []string {
	fields := make([]string, len(*rr))
	i := 0
	for field, _ := range (*rr) { fields[i++] = field }
	return fields
}
