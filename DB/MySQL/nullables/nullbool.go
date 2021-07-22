package nullables

/*
NullBool is an alias for sql.NullBool data type extended for JSON Un|Marshaling support
*/

import (
	"reflect"
	"encoding/json"
	"database/sql"
)

type NullBool sql.NullBool

func (nb *NullBool) Scan(value interface{}) error {
	var b sql.NullBool
	if err := b.Scan(value); err != nil { return err }

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nb = NullBool{b.Bool, false}
	} else {
		*nb = NullBool{b.Bool, true}
	}

	return nil
}

func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if ! nb.Valid { return []byte("null"), nil }
	return json.Marshal(nb.Bool)
}

func (nb *NullBool) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nb.Bool)
	nb.Valid = (err == nil)
	return err
}