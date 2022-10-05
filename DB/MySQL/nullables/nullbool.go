package nullables

import (
	"fmt"
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

