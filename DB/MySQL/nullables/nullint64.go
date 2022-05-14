package nullables

import (
	"fmt"
	"encoding/json"
	"database/sql"
)

// NullInt64 is an alias for sql.NullInt64 data type which we extend
type NullInt64 sql.NullInt64

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) Scan(value interface{}) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullInt64.Scan() - cannot scan into nil receiver") }
	var i sql.NullInt64
	err := i.Scan(value)
	r.Int64 = i.Int64
	r.Valid = i.Valid
	if r.Valid { return nil }
	if nil != err { return err }
	return fmt.Errorf("NullInt64.Scan() - Invalid result without error")
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) MarshalJSON() ([]byte, error) {
	// Nil reciever? Bogus request!
	if nil == r { return make([]byte, 0), fmt.Errorf("NullInt64.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.Valid { return []byte("null"), nil }
	return json.Marshal(r.Int64)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) UnmarshalJSON(b []byte) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullInt64.UnmarshalJSON() - cannot decode JSON into nil receiver") }
	err := json.Unmarshal(b, &r.Int64)
	r.Valid = (nil == err)
	return err
}
