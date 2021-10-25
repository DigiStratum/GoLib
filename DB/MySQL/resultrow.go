package mysql

/*
A single result row from a MySQL query result set

TODO: See if there is a way to encode each Nullable value as it's native JSON data type instead of making them all strings

Interesting:
 * http://go-database-sql.org/varcols.html
 * http://jmoiron.github.io/sqlx/
*/

import (
	"encoding/json"

	nullables "github.com/DigiStratum/GoLib/DB/MySQL/nullables"
)

type ResultRowIfc interface {
	Get(field string) nullables.NullableIfc
	Set(field string, value nullables.Nullable)
	Fields() []string
	ToJson() (*string, error)
}

// Non-exported structure with exported properties that we can serialize
type resultRowSerializableProperties struct {
	values		map[string]nullables.Nullable
}

// Exported structure with non-exported properties to prevent consumer from accessing directly
type ResultRow struct {
	props		resultRowSerializableProperties
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewResultRow() *ResultRow {
	return &ResultRow{
		props:		resultRowSerializableProperties{
			values:		make(map[string]nullables.Nullable),
		},
	}
}

// -------------------------------------------------------------------------------------------------
// ResultRowIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r ResultRow) Get(field string) nullables.NullableIfc {
	if value, ok := r.props.values[field]; ok { return &value }
	return nil
}

func (r *ResultRow) Set(field string, value nullables.Nullable) {
	r.props.values[field] = value
}

// Pluck the fields out of the result and just return them so that caller can iterate with Get()
func (r ResultRow) Fields() []string {
	fields := make([]string, len(r.props.values))
	i := 0
	for field, _ := range r.props.values { fields[i] = field; i++ }
	return fields
}

func (r ResultRow) ToJson() (*string, error) {
	jsonBytes, err := r.MarshalJSON()
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

// ref: http://gregtrowbridge.com/golang-json-serialization-with-interfaces/

func (r ResultRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.props.values)
}
