package nullables

/*
NullTime is an alias for sql.NullTime data type extended for JSON Un|Marshaling support
*/

import (
	"fmt"
	"time"
	"reflect"

	"github.com/go-sql-driver/mysql"
)

type NullTime mysql.NullTime

func (nt *NullTime) Scan(value interface{}) error {
	var t mysql.NullTime
	if err := t.Scan(value); err != nil { return err }

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nt = NullTime{t.Time, false}
	} else {
		*nt = NullTime{t.Time, true}
	}

	return nil
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if ! nt.Valid { return []byte("null"), nil }
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func (nt *NullTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	x, err := time.Parse(time.RFC3339, s)
	if err != nil {
		nt.Valid = false
		return err
	}

	nt.Time = x
	nt.Valid = true
	return nil
}
