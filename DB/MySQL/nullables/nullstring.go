package nullables

/*
NullString is an alias for sql.NullString data type extended for JSON Un|Marshaling support
ref: https://golang.org/src/database/sql/sql.go?s=4943:5036#L177
*/

import (
	"reflect"
	"encoding/json"
	"database/sql"
)

type NullString sql.NullString

func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil { return err }

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ! ns.Valid { return []byte("null"), nil }
	return json.Marshal(ns.String)
}

func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}