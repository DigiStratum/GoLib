package mysql

/*
Nullable primitive data types extended to work for JSON Marshaling

Even though we want to read records from the database into simple string, int, etc. the reality is
that these values could be null in the database... and where that is the case, they must be
nullable in our Result object as well - otherwise we'll get an error from the query Result Scan()
when attempting to write a nul into a non-nullable field.

ref: https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267

We define the following nullable data types:

* NullInt64
* NullBool
* NullFloat64
* NullString
* NullTime

*/

import (
	"fmt"
	"time"
	"reflect"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
)


// -------------------------------------------------------------------------------------------------
// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 sql.NullInt64

// Scan implements the Scanner interface for NullInt64
func (ni *NullInt64) Scan(value interface{}) error {
	var i sql.NullInt64
	if err := i.Scan(value); err != nil { return err }

	// if nil the make Valid false
	if reflect.TypeOf(value) == nil {
		*ni = NullInt64{i.Int64, false}
	} else {
		*ni = NullInt64{i.Int64, true}
	}
	return nil
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if ! ni.Valid { return []byte("null"), nil }
	return json.Marshal(ni.Int64)
}

// UnmarshalJSON for NullInt64
func (ni *NullInt64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ni.Int64)
	ni.Valid = (err == nil)
	return err
}

// -------------------------------------------------------------------------------------------------
// NullBool is an alias for sql.NullBool data type
type NullBool sql.NullBool

// Scan implements the Scanner interface for NullBool
func (nb *NullBool) Scan(value interface{}) error {
	var b sql.NullBool
	if err := b.Scan(value); err != nil { return err }

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nb = NullBool{b.Bool, false}
	} else {
		*nb = NullBool{b.Bool, true}
	}

	return nil
}

// MarshalJSON for NullBool
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if ! nb.Valid { return []byte("null"), nil }
	return json.Marshal(nb.Bool)
}

// UnmarshalJSON for NullBool
func (nb *NullBool) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nb.Bool)
	nb.Valid = (err == nil)
	return err
}

// -------------------------------------------------------------------------------------------------
// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 sql.NullFloat64

// Scan implements the Scanner interface for NullFloat64
func (nf *NullFloat64) Scan(value interface{}) error {
	var f sql.NullFloat64
	if err := f.Scan(value); err != nil { return err }

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*nf = NullFloat64{f.Float64, false}
	} else {
		*nf = NullFloat64{f.Float64, true}
	}

	return nil
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if ! nf.Valid { return []byte("null"), nil }
	return json.Marshal(nf.Float64)
}

// UnmarshalJSON for NullFloat64
func (nf *NullFloat64) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &nf.Float64)
	nf.Valid = (err == nil)
	return err
}

// -------------------------------------------------------------------------------------------------
// NullString is an alias for sql.NullString data type
// ref: https://golang.org/src/database/sql/sql.go?s=4943:5036#L177
type NullString sql.NullString

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil { return err }

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if ! ns.Valid { return []byte("null"), nil }
	return json.Marshal(ns.String)
}

// UnmarshalJSON for NullString
func (ns *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &ns.String)
	ns.Valid = (err == nil)
	return err
}

// -------------------------------------------------------------------------------------------------
// NullTime is an alias for mysql.NullTime data type
type NullTime mysql.NullTime

// Scan implements the Scanner interface for NullTime
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

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if ! nt.Valid { return []byte("null"), nil }
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

// UnmarshalJSON for NullTime
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

// -------------------------------------------------------------------------------------------------
// Nullable - a compound structure that supports all of the nullable types with additional support methods
type NullableType int8

const (
	NULLABLE_NIL NullableType = iota
	NULLABLE_INT64
	NULLABLE_BOOL
	NULLABLE_FLOAT64
	NULLABLE_STRING
	NULLABLE_TIME
	NULLABLE_UNKNOWN
)

func  GetNullableTypeString(nullableType NullableType) string {
	switch nullableType {
		case NULLABLE_NIL:		return "nil"
		case NULLABLE_INT64:		return "int64"
		case NULLABLE_BOOL:		return "bool"
		case NULLABLE_FLOAT64:		return "float64"
		case NULLABLE_STRING:		return "string"
		case NULLABLE_TIME:		return "time"
	}
	return "unknown"
}

type NullableIfc interface {
	IsNil() bool
	SetValue(value interface{}) bool
	GetType() NullableType
	IsInt64() bool
        IsBool() bool
        IsFloat64() bool
        IsString() bool
        IsTime() bool
        GetInt64() *int64
        GetBool() *bool
        GetFloat64() *float64
        GetString() *string
}

type Nullable struct {
	isNil		bool
	nullableType	NullableType
	ni		NullInt64
	nb		NullBool
	nf		NullFloat64
	ns		NullString
	nt		NullTime
}

// Make a new one of these!
func NewNullable(value interface{}) NullableIfc {
	n := Nullable{
		isNil:	true,
		ni:	NullInt64{ Valid: false },
		nb:	NullBool{ Valid: false },
		nf:	NullFloat64{ Valid: false },
		ns:	NullString{ Valid: false },
		nt:	NullTime{ Valid: false },
	}
	n.SetValue(value)
	return &n
}

func (n *Nullable) IsNil() bool { return (*n).isNil }

// Convert value to appropriate Nullable; return true on success, else false
func (n *Nullable) SetValue(value interface{}) bool {
	if v, ok := value.(int64); ok {
		n.setInt64(v)
	} else if v, ok := value.(bool); ok {
		n.setBool(v)
	} else if v, ok := value.(float64); ok {
		n.setFloat64(v)
	} else if v, ok := value.(string); ok {
		n.setString(v)
	} else if v, ok := value.(time.Time); ok {
		n.setTime(v)
	} else { return false }
	return true
}

func (n *Nullable) setInt64(value int64) {
	(*n).nullableType = NULLABLE_INT64
	(*n).ni.Int64 = value
	(*n).ni.Valid = true
	(*n).isNil = false
}

func (n *Nullable) setBool(value bool) {
	(*n).nullableType = NULLABLE_BOOL
	(*n).nb.Bool = value
	(*n).nb.Valid = true
	(*n).isNil = false
}

func (n *Nullable) setFloat64(value float64) {
	(*n).nullableType = NULLABLE_FLOAT64
	(*n).nf.Float64 = value
	(*n).nf.Valid = true
	(*n).isNil = false
}

func (n *Nullable) setString(value string) {
	(*n).nullableType = NULLABLE_STRING
	(*n).ns.String = value
	(*n).ns.Valid = true
	(*n).isNil = false
}

func (n *Nullable) setTime(value time.Time) {
	(*n).nullableType = NULLABLE_TIME
	(*n).nt.Time = value
	(*n).nt.Valid = true
	(*n).isNil = false
}

func (n *Nullable) GetType() NullableType {
	if n.IsNil() { return NULLABLE_NIL }
	if n.IsInt64() { return NULLABLE_INT64 }
	if n.IsBool() { return NULLABLE_BOOL }
	if n.IsFloat64() { return NULLABLE_FLOAT64 }
	if n.IsString() { return NULLABLE_STRING }
	if n.IsTime() { return NULLABLE_TIME }
	return NULLABLE_UNKNOWN
}

func (n *Nullable) IsInt64() bool { return (*n).nullableType == NULLABLE_NIL }
func (n *Nullable) IsBool() bool { return (*n).nullableType == NULLABLE_BOOL }
func (n *Nullable) IsFloat64() bool { return (*n).nullableType == NULLABLE_FLOAT64 }
func (n *Nullable) IsString() bool { return (*n).nullableType == NULLABLE_STRING }
func (n *Nullable) IsTime() bool { return (*n).nullableType == NULLABLE_TIME }

// Return the value as an Int64, complete with data conversions, or nil if nil or conversion problem
func (n *Nullable) GetInt64() *int64 {
	if n.IsNil() { return nil }

	// NullInt64 passes through unmodified
	if (*n).ni.Valid { return &(*n).ni.Int64 }

	// NullBool converts to a int64
	if (*n).nb.Valid {
		var v int64 = 0
		if (*n).nb.Bool { v = 1 }
		return &v
	}

	// NullFloat64 converts to an int64
	if (*n).nf.Valid { v := int64((*n).nf.Float64); return &v }

	// NullString converts to an int64
	if (*n).ns.Valid {
		if vc, err := strconv.ParseInt((*n).ns.String, 0, 64); nil == err { return &vc }
		return nil
	}

	// NullTime converts to an int64 (timestamp)
	if (*n).nt.Valid { vc := (*n).nt.Time.Unix(); return &vc }

	return nil
}

// Return the value as a bool, complete with data conversions, or nil if nil or conversion problem
func (n *Nullable) GetBool() *bool {
	if n.IsNil() { return nil }

	// NullInt64 converts to a bool
	if (*n).ni.Valid { v := ((*n).ni.Int64 == 0); return &v }

	// NullBool passes through unmodified
	if (*n).nb.Valid { return &(*n).nb.Bool }

	// NullFloat64 converts to a bool (true if we drop the decimal and the remaining int != 0)
	if (*n).nf.Valid { v := (int64((*n).nf.Float64) != 0); return &v }

	// NullString converts to a bool (true if "true" or stringified int and != 0 )
	if (*n).ns.Valid {
		lcv := strings.ToLower((*n).ns.String)
		if lcv == "true" { v := true; return &v }
		if vc, err := strconv.ParseInt((*n).ns.String, 0, 64); nil != err {
			v := (vc != 0);
			return &v
		}
		return nil
	}

	// NullTime converts to a bool (true if non-null)
	if (*n).nt.Valid { v := true; return &v }

	return nil
}

// Return the value as a Float64, complete with data conversions, or nil if nil or conversion problem
func (n *Nullable) GetFloat64() *float64 {
	if n.IsNil() { return nil }

	// NullInt64 converts to a Float64
	if (*n).ni.Valid { vc := float64((*n).ni.Int64); return &vc }

	// NullBool converts to a Float64
	// we use 2.0|0.0 for true|false, respectively so that inverse conversion works.
	// Precision rounding reduces 1.0 to < 1 (0.999) which when converted back would yield 0 decimal value (false)
	if (*n).nb.Valid {
		var v float64 = 0.0
		if (*n).nb.Bool { v = 2.0 }
		return &v
	}

	// NullFloat64 passes through unmodified
	if (*n).nf.Valid { return &(*n).nf.Float64 }

	// NullString converts to a Float64
	if (*n).ns.Valid {
		if vc, err := strconv.ParseFloat((*n).ns.String, 64); nil != err { return &vc }
		return nil
	}

	// NullTime conversion to a Float64 not supported
	// (timestamp) would lose precision, so we will 0 it out, no value
	if (*n).nt.Valid { return nil }

	return nil
}

// Return the value as a *string, complete with data conversions, or nil if nil or conversion problem
func (n *Nullable) GetString() *string {
	if n.IsNil() { return nil }

	// NullInt64 converts to a string
	if (*n).ni.Valid {
		v := fmt.Sprintf("%d", (*n).ni.Int64)
		return &v
	}

	// NullBool converts to a string
	if (*n).nb.Valid {
		v := "false"
		if (*n).nb.Bool { v = "true" }
		return &v
	}

	// NullFloat64 converts to a string
	if (*n).nf.Valid {
		v := strconv.FormatFloat((*n).nf.Float64, 'E', -1, 64)
		return &v
	}

	// NullString passes through unmodified
	if (*n).ns.Valid { return &(*n).ns.String }

	// NullTime converts to a string
	// ref: https://stackoverflow.com/questions/33119748/convert-time-time-to-string
	// ref: (so annoying...) https://pkg.go.dev/time#Time.Format
	if (*n).nt.Valid {
		v := (*n).nt.Time.Format("2006-01-02 15:04:05")
		return &v
	}

	return nil
}

// Return the value as a *time.Time, complete with data conversions, or nil if nil or conversion problem
func (n *Nullable) GetTime() *time.Time {
	if n.IsNil() { return nil }

	// NullInt64 converts to a time.Time (unix timestamp)
	if (*n).ni.Valid {
		v := time.Unix((*n).ni.Int64, 0)
		return &v
	}

	// NullBool does not convert...
	if (*n).nb.Valid { return nil }

	// NullFloat64 converts to an int64, then to a time
	if (*n).nf.Valid {
		v := time.Unix(int64((*n).nf.Float64), 0)
		return &v
	}

	// NullString parses as a datetime (MySQL style)
	if (*n).ns.Valid {
		v, err := time.Parse("2006-01-02 15:04:05", (*n).ns.String)
		if nil != err { return nil }
		return &v
	}

	// NullTime passes through unmodified
	if (*n).nt.Valid { return &(*n).nt.Time	}

	return nil
}
