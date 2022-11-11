package nullables

import (
	"fmt"
	"time"
	"strings"
	"strconv"
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
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewNullString(value string) *NullString {
	r := NullString{}
	r.SetValue(&value)
	return &r
}

// -------------------------------------------------------------------------------------------------
// NullString Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) GetValue() *string {
	if ! r.n.Valid { return nil }
	return  &r.n.String
}

func (r *NullString) SetValue(value *string) {
	if nil != value { r.n.String = *value }
	r.n.Valid = (nil != value)
}

// -------------------------------------------------------------------------------------------------
// NullableValueIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullString) GetType() NullableType {
	return NULLABLE_STRING
}

func (r *NullString) GetInt64() *int64 {
	rv := r.GetValue()
	if nil == rv { return nil }
	if vc, err := strconv.ParseInt(*rv, 0, 64); nil == err { return &vc }
	return nil
}

func (r *NullString) GetBool() *bool {
	rv := r.GetValue()
	if nil == rv { return nil }
	// NullString converts to a bool (true if "true" or stringified int and != 0 )
	var v bool;
	lcv := strings.ToLower(*rv)
	if lcv == "true" { v = true; return &v }
	if vc, err := strconv.ParseInt(*rv, 0, 64); nil != err {
		v = (vc != 0)
		return &v
	}
	return &v
}

func (r *NullString) GetFloat64() *float64 {
	rv := r.GetValue()
	if nil == rv { return nil }
	if vc, err := strconv.ParseFloat(*rv, 64); nil == err { return &vc }
	return nil
}

func (r *NullString) GetString() *string {
	return r.GetValue()
}

func (r *NullString) GetTime() *time.Time {
	rv := r.GetValue()
	if nil == rv { return nil }
	if v, err := time.Parse("2006-01-02T15:04:05Z", *rv); nil == err { return &v }
	return nil
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

