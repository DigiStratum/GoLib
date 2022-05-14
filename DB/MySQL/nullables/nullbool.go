package nullables

import (
	"fmt"
	"encoding/json"
	"database/sql"
)

// NullBool is an alias for sql.NullBool data type which we extend
type NullBool sql.NullBool

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) Scan(value interface{}) error {
        // Nil reciever? Bogus request!
        if nil == r { return fmt.Errorf("NullBool.Scan() - cannot scan into nil receiver") }
	var b sql.NullBool
	err := b.Scan(value)
        r.Bool = b.Bool
        r.Valid = b.Valid
        if r.Valid { return nil }
        if nil != err { return err }
        return fmt.Errorf("NullBool.Scan() - Invalid result without error")

}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) MarshalJSON() ([]byte, error) {
	// Nil reciever? Bogus request!
	if nil == r { return make([]byte, 0), fmt.Errorf("NullBool.MarshalJSON() - cannot make nothing into JSON") }
	if ! r.Valid { return []byte("null"), nil }
	return json.Marshal(r.Bool)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *NullBool) UnmarshalJSON(b []byte) error {
        // Nil reciever? Bogus request!
        if nil == r { return fmt.Errorf("NullBool.UnmarshalJSON() - cannot decode JSON into nil receiver") }
        err := json.Unmarshal(b, &r.Bool)
        r.Valid = (nil == err)
        return err
}
