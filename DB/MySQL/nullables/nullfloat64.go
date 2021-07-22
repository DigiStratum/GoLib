package nullables

/*
NullFloat64 is an alias for sql.NullFloat64 data type extended for JSON Un|Marshaling support
*/

import (
	"reflect"
	"encoding/json"
	"database/sql"
)

type NullFloat64 sql.NullFloat64

func (nf *NullFloat64) Scan(value interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(value); err != nil { return err }

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nf = NullFloat64{f.Float64, false}
	} else {
		*nf = NullFloat64{f.Float64, true}
	}

	return nil
}

func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if ! nf.Valid { return []byte("null"), nil }
	return json.Marshal(nf.Float64)
}

func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}
