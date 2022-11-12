package nullables

import (
	"fmt"
	"time"
	"encoding/json"
	"database/sql"
)

type NullBoolIfc interface {
	GetValue() *bool
	SetValue(value *bool)
}

// NullBool has non-exported sql.NullBool, requires use of exported receiver functions to access
type NullBool struct {
	n	sql.NullBool
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewNullBool(value bool) *NullBool {
	r := NullBool{}
	r.SetValue(&value)
	return &r
}

// -------------------------------------------------------------------------------------------------
// NullBoolIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) GetValue() *bool {
	if ! r.n.Valid { return nil }
	return  &r.n.Bool
}

func (r *NullBool) SetValue(value *bool) {
	if nil != value { r.n.Bool = *value }
	r.n.Valid = (nil != value)
}

// -------------------------------------------------------------------------------------------------
// NullableValueIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) GetType() NullableType {
	return NULLABLE_BOOL
}

func (r *NullBool) GetInt64() *int64 {
	rv := r.GetValue()
	if nil == rv { return nil }
	var v int64 = 0
	if *rv { v = 1 }
	return &v
}

func (r *NullBool) GetBool() *bool {
	return r.GetValue()
}

func (r *NullBool) GetFloat64() *float64 {
	rv := r.GetValue()
	if nil == rv { return nil }
	// we use 2.0|0.0 for true|false, respectively so that inverse conversion works.
	// Precision rounding reduces 1.0 to < 1 (0.999) which when converted back would yield 0 decimal value (false)
	var v float64 = 0.0
	if *rv { v = 2.0 }
	return &v
}

func (r *NullBool) GetString() *string {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := "false"
	if *rv { v = "true" }
	return &v
}

func (r *NullBool) GetTime() *time.Time {
	// There's no sensible conversion from bool to Time
	return nil
}

func (r *NullBool) IsNil() bool {
	return (nil == r.GetValue())
}

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) Scan(value interface{}) error {
        // Nil reciever? Bogus request!
        if nil == r { return fmt.Errorf("NullBool.Scan() - cannot scan into nil receiver") }
	var b sql.NullBool
	err := b.Scan(value)
        r.n.Bool = b.Bool
        r.n.Valid = b.Valid
        if r.n.Valid { return nil }
        if nil != err { return err }
        return fmt.Errorf("NullBool.Scan() - Invalid result without error")

}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) MarshalJSON() ([]byte, error) {
	// Nil reciever? Bogus request!
	if nil == r { return make([]byte, 0), fmt.Errorf("NullBool.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.n.Valid { return []byte("null"), nil }
	return json.Marshal(r.n.Bool)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) UnmarshalJSON(b []byte) error {
        // Nil reciever? Bogus request!
        if nil == r { return fmt.Errorf("NullBool.UnmarshalJSON() - cannot decode JSON into nil receiver") }
        err := json.Unmarshal(b, &r.n.Bool)
        r.n.Valid = (nil == err)
        return err
}

