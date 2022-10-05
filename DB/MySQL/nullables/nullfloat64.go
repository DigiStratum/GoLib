package nullables

import (
	"fmt"
	"encoding/json"
	"database/sql"
)

// NullFloat has non-exported sql.NullFloat, requires use of exported receiver functions to access
type NullFloatIfc interface {
	GetValue() *float64
	SetValue(value *float64)
}

type NullFloat struct {
	n	sql.NullFloat
}

// -------------------------------------------------------------------------------------------------
// NullString Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullFloat) GetValue() *float64 {
	if ! r.n.Valid { return nil }
	return  &r.n.Float64
}

func (r *NullFloat) SetValue(value *float64) {
	if nil != value { r.n.Float64 = *value }
	r.n.Valid = (nil != value)
}

// -------------------------------------------------------------------------------------------------
// NullableValueIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullFloat64) GetType() NullableType {
	return NULLABLE_FLOAT64
}

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

