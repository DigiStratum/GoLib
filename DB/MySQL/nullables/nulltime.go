package nullables

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/go-sql-driver/mysql"
)

// NullTime has non-exported sql.NullTime, requires use of exported receiver functions to access
type NullTimeIfc interface {
	GetValue() *time
	SetValue(value *time)
}

// NullTime is an alias for sql.NullTime data type which we extend
type NullTime struct {
	n	mysql.NullTime
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewNullTime(v *time.Time) *NullTime {
	n := NullTime{}
	n.SetValue(v)
	return &n
}

// -------------------------------------------------------------------------------------------------
// NullString Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullTime) GetValue() *time.Time {
	if ! r.n.Valid { return nil }
	return &r.n.Time
}

func (r *NullTime) SetValue(v *time.Time) {
	if nil != v { r.n.Time = *value }
	r.n.Valid = (nil != value)
}

// -------------------------------------------------------------------------------------------------
// NullableValueIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullTime) GetType() NullableType {
	return NULLABLE_TIME
}

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullTime) Scan(value interface{}) error {
        // Nil reciever? Bogus request!
        if nil == r { return fmt.Errorf("NullTime.Scan() - cannot scan into nil receiver") }
	var t mysql.NullTime
	err := t.Scan(value)
        r.Time = t.Time
        r.Valid = t.Valid
        if r.Valid { return nil }
        if nil != err { return err }
        return fmt.Errorf("NullTime.Scan() - Invalid result without error")

}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullTime) MarshalJSON() ([]byte, error) {
        // Nil reciever? Bogus request!
        if nil == r { return make([]byte, 0), fmt.Errorf("NullTime.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.Valid { return []byte("null"), nil }
	val := fmt.Sprintf("\"%s\"", r.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullTime) UnmarshalJSON(b []byte) error {
        // Nil reciever? Bogus request!
        if nil == r { return fmt.Errorf("NullTime.UnmarshalJSON() - cannot decode JSON into nil receiver") }
	// Unmarshal the JSON to a string first...
	var s string
	if err := json.Unmarshal(b, &s); nil != err { return err }
	// Then parse the string as a datetime per RFC3339 formatting
	x, err := time.Parse(time.RFC3339, s)
	r.Time = x
        r.Valid = (nil == err)
	return err
}
