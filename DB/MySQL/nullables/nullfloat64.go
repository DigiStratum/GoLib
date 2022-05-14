package nullables

import (
	"fmt"
	"encoding/json"
	"database/sql"
)

// NullFloat64 is an alias for sql.NullFloat64 data type which we extend
type NullFloat64 sql.NullFloat64

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullFloat64) Scan(value interface{}) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullFloat64.Scan() - cannot scan into nil receiver") }
	var f sql.NullFloat64
	err := f.Scan(value)
	r.Float64 = f.Float64
	r.Valid = f.Valid
	if r.Valid { return nil }
	if nil != err { return err }
	return fmt.Errorf("NullFloat64.Scan() - Invalid result without error")
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullFloat64) MarshalJSON() ([]byte, error) {
	// Nil reciever? Bogus request!
	if nil == r { return make([]byte, 0), fmt.Errorf("NullFloat64.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.Valid { return []byte("null"), nil }
	return json.Marshal(r.Float64)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullFloat64) UnmarshalJSON(b []byte) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullFloat64.UnmarshalJSON() - cannot decode JSON into nil receiver") }
	err := json.Unmarshal(b, &r.Float64)
	r.Valid = (nil == err)
	return err
}

