package nullables

/*
Nullable primitive data types extended to work for JSON Marshaling

Nullable is a compound structure that supports all of the nullable types with additional support methods. The idea
is to be able to support loose typing of sorts from MySQL data. This sounds easier than it is with the way Scan()
gets hints from interface{}, etc. But this is a better start in this direction than we had with our earlier prototype
based model which began to reveal over-complex interface compliance as we begin to shift away from DTO type structures
with all the database record fields exported and towards interface driven models that lend themselves well to a more
generalized approach - this enables us to move more of our boilerplate implementation to this shared library level to
reduce requirements at the application/service model layer for faster, easier develpoment of mysql-backed models.

Even though we want to read records from the database into simple string, int, etc. the reality is
that these values could be null in the database... and where that is the case, they must be
nullable in our Result object as well - otherwise we'll get an error from the query Result Scan()
when attempting to write a nul into a non-nullable field.

ref: https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267
ref: https://kylewbanks.com/blog/query-result-to-map-in-golang

We define the following nullable data types as extensions of the same-named types from the sql package:

* NullInt64
* NullBool
* NullFloat64
* NullString
* NullTime

What these allow for is a query to return null for one of the values and store it into the nullable. If a value were, say,
a straight string or int, Go does not allow this to be nil, so things get difficult.

TODO:
 * Add mapping for unsigned ints in addition to the signed ones - MySQL supports these natively, but database/sql does not!

*/

import (
	"fmt"
	"time"
	"strconv"
	"strings"
)

type NullableType int8

const (
	NULLABLE_NIL NullableType = iota
	NULLABLE_INT64
	NULLABLE_BOOL
	NULLABLE_FLOAT64
	NULLABLE_STRING
	NULLABLE_TIME
)

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
        MarshalJSON() ([]byte, error)
        UnmarshalJSON(b []byte) error
        Scan(value interface{}) error
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

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewNullable(value interface{}) *Nullable {
	n := Nullable{}
	if res := n.SetValue(value); res { return &n }
	return nil
}

// -------------------------------------------------------------------------------------------------
// NullableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r Nullable) IsNil() bool { return r.isNil }

// Convert value to appropriate Nullable; return true on success, else false
func (r *Nullable) SetValue(value interface{}) bool {
	// Set Valid=false property for each nullable; this indicates to base object that it is nil
	r.ni = NullInt64{ Valid: false }
	r.nb = NullBool{ Valid: false }
	r.nf = NullFloat64{ Valid: false }
	//r.ns = NullString{ Valid: false }
	r.ns = NullString{ }
	r.nt = NullTime{ Valid: false }
	if nil == value { return r.setNil() }
	if v, ok := value.(int); ok { return r.setInt64(int64(v)) }
	if v, ok := value.(int8); ok { return r.setInt64(int64(v)) }
	if v, ok := value.(int16); ok { return r.setInt64(int64(v)) }
	if v, ok := value.(int32); ok { return r.setInt64(int64(v)) }
	if v, ok := value.(int64); ok {	return r.setInt64(v) }
	if v, ok := value.(bool); ok { return r.setBool(v) }
	if v, ok := value.(float32); ok { return r.setFloat64(float64(v)) }
	if v, ok := value.(float64); ok { return r.setFloat64(v) }
	if v, ok := value.(string); ok { return r.setString(v) }
	if v, ok := value.(time.Time); ok { return r.setTime(v) }
	return false
}

func (r *Nullable) GetType() NullableType { return r.nullableType }

func (r *Nullable) IsInt64() bool { return r.nullableType == NULLABLE_INT64 }
func (r *Nullable) IsBool() bool { return r.nullableType == NULLABLE_BOOL }
func (r *Nullable) IsFloat64() bool { return r.nullableType == NULLABLE_FLOAT64 }
func (r *Nullable) IsString() bool { return r.nullableType == NULLABLE_STRING }
func (r *Nullable) IsTime() bool { return r.nullableType == NULLABLE_TIME }

// Return the value as an Int64, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetInt64() *int64 {
	switch {
		case NULLABLE_INT64==r.nullableType && r.ni.Valid:
			// NullInt64 passes through unmodified
			return &r.ni.Int64
		case NULLABLE_BOOL==r.nullableType && r.nb.Valid:
			// NullBool converts to a int64
			var v int64 = 0
			if r.nb.Bool { v = 1 }
			return &v
		case NULLABLE_FLOAT64==r.nullableType && r.nf.Valid:
			// NullFloat64 converts to an int64
			v := int64(r.nf.Float64)
			return &v
		//case NULLABLE_STRING==r.nullableType && r.ns.Valid:
		case NULLABLE_STRING==r.nullableType && r.ns.IsValid():
			// NullString converts to an int64
			//if vc, err := strconv.ParseInt(r.ns.String, 0, 64); nil == err {
			s := r.ns.GetValue()
			if nil == s { return nil }
			if vc, err := strconv.ParseInt(*s, 0, 64); nil == err {
				return &vc
			}
		case NULLABLE_TIME==r.nullableType && r.nt.Valid:
			// NullTime converts to an int64 (timestamp)
			vc := r.nt.Time.Unix()
			return &vc
	}
	return nil
}

// Return the value as a bool, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetBool() *bool {
	switch {
		case NULLABLE_INT64==r.nullableType && r.ni.Valid:
			// NullInt64 converts to a bool
			v := (r.ni.Int64 == 1)
			return &v
		case NULLABLE_BOOL==r.nullableType && r.nb.Valid:
			// NullBool passes through unmodified
			return &r.nb.Bool
		case NULLABLE_FLOAT64==r.nullableType && r.nf.Valid:
			// NullFloat64 converts to a bool (true if we drop the decimal and the remaining int != 0)
			v := (int64(r.nf.Float64) != 0)
			return &v
		//case NULLABLE_STRING==r.nullableType && r.ns.Valid:
		case NULLABLE_STRING==r.nullableType && r.ns.IsValid():
			// NullString converts to a bool (true if "true" or stringified int and != 0 )
			//lcv := strings.ToLower(r.ns.String)
			s := r.ns.GetValue()
			if nil == s { return nil }
			lcv := strings.ToLower(*s)
			if lcv == "true" { v := true; return &v }
			//if vc, err := strconv.ParseInt(r.ns.String, 0, 64); nil != err {
			if vc, err := strconv.ParseInt(*s, 0, 64); nil != err {
				v := (vc != 0)
				return &v
			}
		case NULLABLE_TIME==r.nullableType && r.nt.Valid:
			// Any non-nil NullTime converts to a bool=true
			v := true
			return &v
	}
	return nil
}

// Return the value as a Float64, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetFloat64() *float64 {
	switch {
		case NULLABLE_INT64==r.nullableType && r.ni.Valid:
			// NullInt64 converts to a Float64
			vc := float64(r.ni.Int64)
			return &vc
		case NULLABLE_BOOL==r.nullableType && r.nb.Valid:
			// NullBool converts to a Float64
			// we use 2.0|0.0 for true|false, respectively so that inverse conversion works.
			// Precision rounding reduces 1.0 to < 1 (0.999) which when converted back would yield 0 decimal value (false)
			var v float64 = 0.0
			if r.nb.Bool { v = 2.0 }
			return &v
		case NULLABLE_FLOAT64==r.nullableType && r.nf.Valid:
			// NullFloat64 passes through unmodified
			return &r.nf.Float64
		//case NULLABLE_STRING==r.nullableType && r.ns.Valid:
		case NULLABLE_STRING==r.nullableType && r.ns.IsValid():
			// NullString converts to a Float64
			//if vc, err := strconv.ParseFloat(r.ns.String, 64); nil == err {
			s := r.ns.GetValue()
			if nil == s { return nil }
			if vc, err := strconv.ParseFloat(*s, 64); nil == err {
				return &vc
			}
		// NullTime conversion to a Float64 not supported
		// (timestamp) would lose precision, so we will 0 it out, no value
	}
	return nil
}

// Return the value as a *string, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetString() *string {
	switch {
		case NULLABLE_INT64==r.nullableType && r.ni.Valid:
			// NullInt64 converts to a string
			v := fmt.Sprintf("%d", r.ni.Int64)
			return &v
		case NULLABLE_BOOL==r.nullableType && r.nb.Valid:
			// NullBool converts to a string
			v := "false"
			if r.nb.Bool { v = "true" }
			return &v
		case NULLABLE_FLOAT64==r.nullableType && r.nf.Valid:
			// NullFloat64 converts to a string
			v := strconv.FormatFloat(r.nf.Float64, 'E', -1, 64)
			return &v
		//case NULLABLE_STRING==r.nullableType && r.ns.Valid:
		case NULLABLE_STRING==r.nullableType && r.ns.IsValid():
			// NullString passes through unmodified
			//return &r.ns.String
			return r.ns.GetValue()
		case NULLABLE_TIME==r.nullableType && r.nt.Valid:
			// NullTime converts to a string
			// ref: https://stackoverflow.com/questions/33119748/convert-time-time-to-string
			// ref: (so annoying...) https://pkg.go.dev/time#Time.Format
			v := r.nt.Time.Format("2006-01-02T15:04:05Z")
			return &v
	}
	return nil
}

// Return the value as a *time.Time, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetTime() *time.Time {
	switch {
		case NULLABLE_INT64==r.nullableType && r.ni.Valid:
			// NullInt64 converts to a time.Time (unix timestamp)
			v := time.Unix(r.ni.Int64, 0)
			return &v
		case NULLABLE_FLOAT64==r.nullableType && r.nf.Valid:
			// NullFloat64 converts to an int64, then to a time
			v := time.Unix(int64(r.nf.Float64), 0)
			return &v
		//case NULLABLE_STRING==r.nullableType && r.ns.Valid:
		case NULLABLE_STRING==r.nullableType && r.ns.IsValid():
			// NullString parses as a datetime (MySQL style)
			s := r.ns.GetValue()
			if nil == s { return nil }
			//v, err := time.Parse("2006-01-02T15:04:05Z", r.ns.String)
			v, err := time.Parse("2006-01-02T15:04:05Z", *s)
			if nil != err { return nil }
			return &v
		case NULLABLE_TIME==r.nullableType && r.nt.Valid:
			// NullTime passes through unmodified
			return &r.nt.Time
		// NullBool does not convert...
	}
	return nil
}

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

// Scan for Nullable - we just sub it out to the underlying Nullable type
func (r *Nullable) Scan(value interface{}) error {
	switch r.nullableType {
		case NULLABLE_NIL: return nil
		case NULLABLE_INT64: return r.ni.Scan(value)
		case NULLABLE_BOOL: return r.nb.Scan(value)
		case NULLABLE_FLOAT64: return r.nf.Scan(value)
		case NULLABLE_STRING: return r.ns.Scan(value)
		case NULLABLE_TIME: return r.nt.Scan(value)
	}
	return fmt.Errorf("Nullable.Scan - Unsupported Nullable Type (oversight in implementation for type=%d!)", r.nullableType)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r Nullable) MarshalJSON() ([]byte, error) {
	// Sub it out to the underlying Nullable type
	switch r.nullableType {
		case NULLABLE_INT64: return r.ni.MarshalJSON()
		case NULLABLE_BOOL: return r.nb.MarshalJSON()
		case NULLABLE_FLOAT64: return r.nf.MarshalJSON()
		case NULLABLE_STRING: return r.ns.MarshalJSON()
		case NULLABLE_TIME: return r.nt.MarshalJSON()
		case NULLABLE_NIL: return []byte("null"), nil
	}
	return make([]byte, 0), fmt.Errorf("Nullable.MarshalJSON - Unsupported Nullable Type (oversight in implementation for type=%d!)", r.nullableType)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Nullable) UnmarshalJSON(b []byte) error {
	// Sub it out to the underlying Nullable type
	switch r.nullableType {
		case NULLABLE_NIL: return nil
		case NULLABLE_INT64: return r.ni.UnmarshalJSON(b)
		case NULLABLE_BOOL: return r.nb.UnmarshalJSON(b)
		case NULLABLE_FLOAT64: return r.nf.UnmarshalJSON(b)
		case NULLABLE_STRING: return r.ns.UnmarshalJSON(b)
		case NULLABLE_TIME: return r.nt.UnmarshalJSON(b)
	}
	return fmt.Errorf("Nullable.UnmarshalJSON - Unsupported Nullable Type (oversight in implementation for type=%d!)", r.nullableType)
}

// -------------------------------------------------------------------------------------------------
// Nullable private supporting functions
// -------------------------------------------------------------------------------------------------

func (r *Nullable) setNil() bool {
	r.nullableType = NULLABLE_NIL
	r.isNil = true
	return true
}

func (r *Nullable) setInt64(value int64) bool {
	r.nullableType = NULLABLE_INT64
	r.ni.Int64 = value
	r.ni.Valid = true
	r.isNil = false
	return true
}

func (r *Nullable) setBool(value bool) bool {
	r.nullableType = NULLABLE_BOOL
	r.nb.Bool = value
	r.nb.Valid = true
	r.isNil = false
	return true
}

func (r *Nullable) setFloat64(value float64) bool {
	r.nullableType = NULLABLE_FLOAT64
	r.nf.Float64 = value
	r.nf.Valid = true
	r.isNil = false
	return true
}

func (r *Nullable) setString(value string) bool {
	r.nullableType = NULLABLE_STRING
	//r.ns.String = value
	r.ns.SetValue(&value)
	//r.ns.Valid = true
	r.isNil = false
	return true
}

func (r *Nullable) setTime(value time.Time) bool {
	r.nullableType = NULLABLE_TIME
	r.nt.Time = value
	r.nt.Valid = true
	r.isNil = false
	return true
}

