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

*/

import (
	"fmt"
	"time"
	"strconv"
	"strings"
	"errors"
)

// -------------------------------------------------------------------------------------------------
// Nullable supporting types, constants, and functions
// -------------------------------------------------------------------------------------------------

type NullableType int8

const (
	NULLABLE_NIL NullableType = iota
	NULLABLE_INT64
	NULLABLE_BOOL
	NULLABLE_FLOAT64
	NULLABLE_STRING
	NULLABLE_TIME
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
        MarshalJSON() ([]byte, error)
        UnmarshalJSON(b []byte) error
        Scan(value interface{}) error
}

type nullable struct {
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
	n := nullable{
		isNil:		true,
		nullableType:	NULLABLE_NIL,
		ni:		NullInt64{ Valid: false },
		nb:		NullBool{ Valid: false },
		nf:		NullFloat64{ Valid: false },
		ns:		NullString{ Valid: false },
		nt:		NullTime{ Valid: false },
	}
	n.SetValue(value)
	return &n
}

// -------------------------------------------------------------------------------------------------
// NullableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (n *nullable) IsNil() bool { return (*n).isNil }

// Convert value to appropriate Nullable; return true on success, else false
func (n *nullable) SetValue(value interface{}) bool {
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

func (n *nullable) setInt64(value int64) {
	(*n).nullableType = NULLABLE_INT64
	(*n).ni.Int64 = value
	(*n).ni.Valid = true
	(*n).isNil = false
}

func (n *nullable) setBool(value bool) {
	(*n).nullableType = NULLABLE_BOOL
	(*n).nb.Bool = value
	(*n).nb.Valid = true
	(*n).isNil = false
}

func (n *nullable) setFloat64(value float64) {
	(*n).nullableType = NULLABLE_FLOAT64
	(*n).nf.Float64 = value
	(*n).nf.Valid = true
	(*n).isNil = false
}

func (n *nullable) setString(value string) {
	(*n).nullableType = NULLABLE_STRING
	(*n).ns.String = value
	(*n).ns.Valid = true
	(*n).isNil = false
}

func (n *nullable) setTime(value time.Time) {
	(*n).nullableType = NULLABLE_TIME
	(*n).nt.Time = value
	(*n).nt.Valid = true
	(*n).isNil = false
}

func (n *nullable) GetType() NullableType { return (*n).nullableType }

func (n *nullable) IsInt64() bool { return (*n).nullableType == NULLABLE_NIL }
func (n *nullable) IsBool() bool { return (*n).nullableType == NULLABLE_BOOL }
func (n *nullable) IsFloat64() bool { return (*n).nullableType == NULLABLE_FLOAT64 }
func (n *nullable) IsString() bool { return (*n).nullableType == NULLABLE_STRING }
func (n *nullable) IsTime() bool { return (*n).nullableType == NULLABLE_TIME }

// Return the value as an Int64, complete with data conversions, or nil if nil or conversion problem
func (n *nullable) GetInt64() *int64 {
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
func (n *nullable) GetBool() *bool {
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
func (n *nullable) GetFloat64() *float64 {
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
func (n *nullable) GetString() *string {
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
		v := (*n).nt.Time.Format("2006-01-02T15:04:05Z")
		return &v
	}

	return nil
}

// Return the value as a *time.Time, complete with data conversions, or nil if nil or conversion problem
func (n *nullable) GetTime() *time.Time {
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
		v, err := time.Parse("2006-01-02T15:04:05Z", (*n).ns.String)
		if nil != err { return nil }
		return &v
	}

	// NullTime passes through unmodified
	if (*n).nt.Valid { return &(*n).nt.Time	}

	return nil
}

// MarshalJSON for Nullable - we just sub it out to the underlying Nullable type
func (n *nullable) MarshalJSON() ([]byte, error) {
	switch (*n).nullableType {
		case NULLABLE_INT64: return (*n).ni.MarshalJSON()
		case NULLABLE_BOOL: return (*n).nb.MarshalJSON()
		case NULLABLE_FLOAT64: return (*n).nf.MarshalJSON()
		case NULLABLE_STRING: return (*n).ns.MarshalJSON()
		case NULLABLE_TIME: return (*n).nt.MarshalJSON()
		default: return []byte("null"), nil
	}
}

// UnmarshalJSON for Nullable - we just sub it out to the underlying Nullable type
func (n *nullable) UnmarshalJSON(b []byte) error {
	switch (*n).nullableType {
		case NULLABLE_INT64: return (*n).ni.UnmarshalJSON(b)
		case NULLABLE_BOOL: return (*n).nb.UnmarshalJSON(b)
		case NULLABLE_FLOAT64: return (*n).nf.UnmarshalJSON(b)
		case NULLABLE_STRING: return (*n).ns.UnmarshalJSON(b)
		case NULLABLE_TIME: return (*n).nt.UnmarshalJSON(b)
		default: return nil
	}
}

// Scan for Nullable - we just sub it out to the underlying Nullable type
func (n *nullable) Scan(value interface{}) error {
	switch (*n).nullableType {
		case NULLABLE_INT64: return (*n).ni.Scan(value)
		case NULLABLE_BOOL: return (*n).nb.Scan(value)
		case NULLABLE_FLOAT64: return (*n).nf.Scan(value)
		case NULLABLE_STRING: return (*n).ns.Scan(value)
		case NULLABLE_TIME: return (*n).nt.Scan(value)
		default: return errors.New("Unsupported Nullable Type (oversight in implementation!)")
	}
}