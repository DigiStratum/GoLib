package nullables

import (
	"fmt"
	"encoding/json"
	"database/sql"
)

// NullString is an alias for sql.NullString data type which we extend
type NullString sql.NullString

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) Scan(value interface{}) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullString.Scan() - cannot scan into nil receiver") }
	var s sql.NullString
	err := s.Scan(value)
	r.String = s.String
	r.Valid = s.Valid
	if r.Valid { return nil }
	if nil != err { return err }
	return fmt.Errorf("NullString.Scan() - Invalid result without error")
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) MarshalJSON() ([]byte, error) {
	// Nil reciever? Bogus request!
	if nil == r { return make([]byte, 0), fmt.Errorf("NullString.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.Valid { return []byte("null"), nil }
	return json.Marshal(r.String)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) UnmarshalJSON(b []byte) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullString.UnmarshalJSON() - cannot decode JSON into nil receiver") }
	err := json.Unmarshal(b, &r.String)
	r.Valid = (nil == err)
	return err
}
