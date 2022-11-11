package nullables

import (
	"fmt"
	"time"
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
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewNullInt64(value int64) *NullInt64 {
	r := NullInt64{}
	r.SetValue(&value)
	return &r
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
// NullableValueIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) GetType() NullableType {
	return NULLABLE_INT64
}

func (r *NullInt64) GetInt64() *int64 {
	return r.GetValue()
}

func (r *NullInt64) GetBool() *bool {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := (*rv == 1)
	return &v
}

func (r *NullInt64) GetFloat64() *float64 {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := float64(*rv)
	return &v
}

func (r *NullInt64) GetString() *string {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := fmt.Sprintf("%d", *rv)
	return &v
}

func (r *NullInt64) GetTime() *time.Time {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := time.Unix(*rv, 0)
	return &v
}

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) Scan(value interface{}) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullInt64.Scan() - cannot scan into nil receiver") }
	var i sql.NullInt64
	err := i.Scan(value)
	r.n.Int64 = i.Int64
	r.n.Valid = i.Valid
	if r.n.Valid { return nil }
	if nil != err { return err }
	return fmt.Errorf("NullInt64.Scan() - Invalid result without error")
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) MarshalJSON() ([]byte, error) {
	// Nil reciever? Bogus request!
	if nil == r { return make([]byte, 0), fmt.Errorf("NullInt64.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.n.Valid { return []byte("null"), nil }
	return json.Marshal(r.n.Int64)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullInt64) UnmarshalJSON(b []byte) error {
	// Nil reciever? Bogus request!
	if nil == r { return fmt.Errorf("NullInt64.UnmarshalJSON() - cannot decode JSON into nil receiver") }
	err := json.Unmarshal(b, &r.n.Int64)
	r.n.Valid = (nil == err)
	return err
}
