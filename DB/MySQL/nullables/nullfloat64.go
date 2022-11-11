package nullables

import (
	"fmt"
	"time"
	"strconv"
	"encoding/json"
	"database/sql"
)

// NullFloat has non-exported sql.NullFloat, requires use of exported receiver functions to access
type NullFloat64Ifc interface {
	GetValue() *float64
	SetValue(value *float64)
}

type NullFloat64 struct {
	n	sql.NullFloat64
}


// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewNullFloat64(value float64) *NullFloat64 {
	r := NullFloat64{}
	r.SetValue(&value)
	return &r
}

// -------------------------------------------------------------------------------------------------
// NullString Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullFloat64) GetValue() *float64 {
	if ! r.n.Valid { return nil }
	return  &r.n.Float64
}

func (r *NullFloat64) SetValue(value *float64) {
	if nil != value { r.n.Float64 = *value }
	r.n.Valid = (nil != value)
}

// -------------------------------------------------------------------------------------------------
// NullableValueIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullFloat64) GetType() NullableType {
	return NULLABLE_FLOAT64
}

func (r *NullFloat64) GetInt64() *int64 {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := int64(*rv)
	return &v
}

func (r *NullFloat64) GetBool() *bool {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := (int64(*rv) != 0)
	return &v
}

func (r *NullFloat64) GetFloat64() *float64 {
	return r.GetValue()
}

func (r *NullFloat64) GetString() *string {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := strconv.FormatFloat(*rv, 'E', -1, 64)
	return &v
}

func (r *NullFloat64) GetTime() *time.Time {
	rv := r.GetValue()
	if nil == rv { return nil }
	v := time.Unix(int64(*rv), 0)
	return &v
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

