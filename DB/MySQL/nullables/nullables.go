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
	"errors"
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
	n := Nullable{
		isNil:		true,
		nullableType:	NULLABLE_NIL,
		ni:		NullInt64{ Valid: false },
		nb:		NullBool{ Valid: false },
		nf:		NullFloat64{ Valid: false },
		ns:		NullString{ Valid: false },
		nt:		NullTime{ Valid: false },
	}
	res := n.SetValue(value)
	if res { return &n }
	return nil
}

// -------------------------------------------------------------------------------------------------
// NullableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r Nullable) IsNil() bool { return r.isNil }

// Convert value to appropriate Nullable; return true on success, else false
func (r *Nullable) SetValue(value interface{}) bool {
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
	r.ns.String = value
	r.ns.Valid = true
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

func (r *Nullable) GetType() NullableType { return r.nullableType }

func (r *Nullable) IsInt64() bool { return r.nullableType == NULLABLE_INT64 }
func (r *Nullable) IsBool() bool { return r.nullableType == NULLABLE_BOOL }
func (r *Nullable) IsFloat64() bool { return r.nullableType == NULLABLE_FLOAT64 }
func (r *Nullable) IsString() bool { return r.nullableType == NULLABLE_STRING }
func (r *Nullable) IsTime() bool { return r.nullableType == NULLABLE_TIME }

// Return the value as an Int64, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetInt64() *int64 {
	switch r.nullableType {
		case NULLABLE_INT64:	// NullInt64 passes through unmodified
			if r.ni.Valid {
				return &r.ni.Int64
			}
		case NULLABLE_BOOL:	// NullBool converts to a int64
			if r.nb.Valid {
				var v int64 = 0
				if r.nb.Bool { v = 1 }
				return &v
			}
		case NULLABLE_FLOAT64:	// NullFloat64 converts to an int64
			if r.nf.Valid {
				v := int64(r.nf.Float64)
				return &v
			}
		case NULLABLE_STRING:	// NullString converts to an int64
			if vc, err := strconv.ParseInt(r.ns.String, 0, 64); nil == err {
				return &vc
			}
		case NULLABLE_TIME:	// NullTime converts to an int64 (timestamp)
			if r.nt.Valid { vc := r.nt.Time.Unix(); return &vc }
	}
	return nil
}

// Return the value as a bool, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetBool() *bool {
	if r.IsNil() { return nil }

	// NullInt64 converts to a bool
	if r.ni.Valid { v := (r.ni.Int64 == 0); return &v }

	// NullBool passes through unmodified
	if r.nb.Valid { return &r.nb.Bool }

	// NullFloat64 converts to a bool (true if we drop the decimal and the remaining int != 0)
	if r.nf.Valid { v := (int64(r.nf.Float64) != 0); return &v }

	// NullString converts to a bool (true if "true" or stringified int and != 0 )
	if r.ns.Valid {
		lcv := strings.ToLower(r.ns.String)
		if lcv == "true" { v := true; return &v }
		if vc, err := strconv.ParseInt(r.ns.String, 0, 64); nil != err {
			v := (vc != 0);
			return &v
		}
		return nil
	}

	// NullTime converts to a bool (true if non-null)
	if r.nt.Valid { v := true; return &v }

	return nil
}

// Return the value as a Float64, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetFloat64() *float64 {
	if r.IsNil() { return nil }

	// NullInt64 converts to a Float64
	if r.ni.Valid { vc := float64(r.ni.Int64); return &vc }

	// NullBool converts to a Float64
	// we use 2.0|0.0 for true|false, respectively so that inverse conversion works.
	// Precision rounding reduces 1.0 to < 1 (0.999) which when converted back would yield 0 decimal value (false)
	if r.nb.Valid {
		var v float64 = 0.0
		if r.nb.Bool { v = 2.0 }
		return &v
	}

	// NullFloat64 passes through unmodified
	if r.nf.Valid { return &r.nf.Float64 }

	// NullString converts to a Float64
	if r.ns.Valid {
		if vc, err := strconv.ParseFloat(r.ns.String, 64); nil != err { return &vc }
		return nil
	}

	// NullTime conversion to a Float64 not supported
	// (timestamp) would lose precision, so we will 0 it out, no value
	if r.nt.Valid { return nil }

	return nil
}

// Return the value as a *string, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetString() *string {
	if r.IsNil() { return nil }

	// NullInt64 converts to a string
	if r.ni.Valid {
		v := fmt.Sprintf("%d", r.ni.Int64)
		return &v
	}

	// NullBool converts to a string
	if r.nb.Valid {
		v := "false"
		if r.nb.Bool { v = "true" }
		return &v
	}

	// NullFloat64 converts to a string
	if r.nf.Valid {
		v := strconv.FormatFloat(r.nf.Float64, 'E', -1, 64)
		return &v
	}

	// NullString passes through unmodified
	if r.ns.Valid { return &r.ns.String }

	// NullTime converts to a string
	// ref: https://stackoverflow.com/questions/33119748/convert-time-time-to-string
	// ref: (so annoying...) https://pkg.go.dev/time#Time.Format
	if r.nt.Valid {
		v := r.nt.Time.Format("2006-01-02T15:04:05Z")
		return &v
	}

	return nil
}

// Return the value as a *time.Time, complete with data conversions, or nil if nil or conversion problem
func (r Nullable) GetTime() *time.Time {
	if r.IsNil() { return nil }

	// NullInt64 converts to a time.Time (unix timestamp)
	if r.ni.Valid {
		v := time.Unix(r.ni.Int64, 0)
		return &v
	}

	// NullBool does not convert...
	if r.nb.Valid { return nil }

	// NullFloat64 converts to an int64, then to a time
	if r.nf.Valid {
		v := time.Unix(int64(r.nf.Float64), 0)
		return &v
	}

	// NullString parses as a datetime (MySQL style)
	if r.ns.Valid {
		v, err := time.Parse("2006-01-02T15:04:05Z", r.ns.String)
		if nil != err { return nil }
		return &v
	}

	// NullTime passes through unmodified
	if r.nt.Valid { return &r.nt.Time	}

	return nil
}

// MarshalJSON for Nullable - we just sub it out to the underlying Nullable type
func (r Nullable) MarshalJSON() ([]byte, error) {
	switch r.nullableType {
		case NULLABLE_INT64: return r.ni.MarshalJSON()
		case NULLABLE_BOOL: return r.nb.MarshalJSON()
		case NULLABLE_FLOAT64: return r.nf.MarshalJSON()
		case NULLABLE_STRING: return r.ns.MarshalJSON()
		case NULLABLE_TIME: return r.nt.MarshalJSON()
		default: return []byte("[\"nullnullable\"]"), nil
	}
}

// UnmarshalJSON for Nullable - we just sub it out to the underlying Nullable type
func (r *Nullable) UnmarshalJSON(b []byte) error {
	switch r.nullableType {
		case NULLABLE_INT64: return r.ni.UnmarshalJSON(b)
		case NULLABLE_BOOL: return r.nb.UnmarshalJSON(b)
		case NULLABLE_FLOAT64: return r.nf.UnmarshalJSON(b)
		case NULLABLE_STRING: return r.ns.UnmarshalJSON(b)
		case NULLABLE_TIME: return r.nt.UnmarshalJSON(b)
		default: return nil
	}
}

// Scan for Nullable - we just sub it out to the underlying Nullable type
func (r *Nullable) Scan(value interface{}) error {
	switch r.nullableType {
		case NULLABLE_INT64: return r.ni.Scan(value)
		case NULLABLE_BOOL: return r.nb.Scan(value)
		case NULLABLE_FLOAT64: return r.nf.Scan(value)
		case NULLABLE_STRING: return r.ns.Scan(value)
		case NULLABLE_TIME: return r.nt.Scan(value)
		default: return errors.New("Unsupported Nullable Type (oversight in implementation!)")
	}
}

// -------------------------------------------------------------------------------------------------
// Nullable supporting functions
// -------------------------------------------------------------------------------------------------

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

