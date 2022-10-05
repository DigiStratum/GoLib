package nullables

import (
	"fmt"
	"encoding/json"
	"database/sql"
)

// NullInt64 has non-exported sql.NullInt64, requires use of exported receiver functions to access
type NullInt64Ifc interface {
	GetValue() *int64
	SetValue(value *int64)
}

type NullInt64 struct {
	n	sql.NullInt64
}

// -------------------------------------------------------------------------------------------------
// NullString Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) GetValue() *int64 {
	if ! r.n.Valid { return nil }
	return &r.n.Int64
}

func (r *NullInt64) SetValue(value *int64) {
	if nil != value { r.n.Int64 = *value }
	r.n.Valid = (nil != value)
}

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
