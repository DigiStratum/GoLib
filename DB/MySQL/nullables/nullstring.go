package nullables

import (
	"fmt"
	"encoding/json"
	"database/sql"
)

type NullStringIfc interface {
	GetValue() *string
	SetValue(value *string)
}

// NullString has non-exported sql.NullString, requires use of exported receiver functions to access
type NullString struct {
	n	sql.NullString
}


// -------------------------------------------------------------------------------------------------
// NullString Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) GetValue() *string {
	if ! r.IsValid() { return nil }
	return  &r.n.String
}

func (r *NullString) SetValue(value *string) {
	if nil != value { r.n.String = *value }
	r.n.Valid = (nil != value)
}

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) Scan(value interface{}) error {
	if nil == r { return fmt.Errorf("NullString.Scan() - cannot scan into nil receiver") }
	var s sql.NullString
	err := s.Scan(value)
	r.n.String = s.String
	r.n.Valid = s.Valid
	if r.n.Valid { return nil }
	if nil != err { return err }
	return fmt.Errorf("NullString.Scan() - Invalid result without error")
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) MarshalJSON() ([]byte, error) {
	// Nil reciever? Bogus request!
	if nil == r { return make([]byte, 0), fmt.Errorf("NullString.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.n.Valid { return []byte("null"), nil }
	return json.Marshal(r.n.String)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) UnmarshalJSON(b []byte) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullString.UnmarshalJSON() - cannot decode JSON into nil receiver") }
	err := json.Unmarshal(b, &r.n.String)
	r.n.Valid = (nil == err)
	return err
}

