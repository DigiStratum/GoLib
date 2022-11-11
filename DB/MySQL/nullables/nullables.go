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
 * Convert Nil to a separate nullable type that implements the NullableValueIfc - no need for special treatment

*/

import (
	"fmt"
	"time"
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

type NullableValueIfc interface {
	GetType() NullableType
	GetInt64() *int64
	GetBool() *bool
	GetFloat64() *float64
	GetString() *string
	GetTime() *time.Time
	Scan(value interface{}) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

type NullableIfc interface {
	NullableValueIfc
	//SetValue(value interface{}) error
	//MarshalJSON() ([]byte, error)
	//UnmarshalJSON(b []byte) error
	//Scan(value interface{}) error
}

type Nullable struct {
	value	NullableValueIfc
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewNullable(value interface{}) *Nullable {
	n := Nullable{}
	err := n.SetValue(value)
	if nil == err { return &n }
	return nil
}

// -------------------------------------------------------------------------------------------------
// NullableIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Convert value to appropriate Nullable; return true on success, else false
func (r *Nullable) SetValue(v interface{}) error {
	switch v.(type) {
		case nil: r.value = nil
		case int, int8, int16, int32, int64: r.value = NewNullInt64(v.(int64))
		case float32, float64: r.value = NewNullFloat64(v.(float64))
		case bool: r.value = NewNullBool(v.(bool))
		case string: r.value = NewNullString(v.(string))
		case time.Time: r.value = NewNullTime(v.(time.Time))
		default: return fmt.Errorf("Supplied value did not match a supported type")
	}
	return nil
}

func (r *Nullable) GetType() NullableType {
	if nil == r.value { return NULLABLE_NIL }
	return r.value.GetType()
}

func (r Nullable) GetInt64() *int64 {
	if nil == r.value { return nil }
	return r.value.GetInt64()
}

func (r Nullable) GetBool() *bool {
	if nil == r.value { return nil }
	return r.value.GetBool()
}

func (r Nullable) GetFloat64() *float64 {
	if nil == r.value { return nil }
	return r.value.GetFloat64()
}

func (r Nullable) GetString() *string {
	if nil == r.value { return nil }
	return r.value.GetString()
}

func (r Nullable) GetTime() *time.Time {
	if nil == r.value { return nil }
	return r.value.GetTime()
}

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

// Scan for Nullable - we just sub it out to the underlying Nullable type
func (r *Nullable) Scan(value interface{}) error {
	if nil == r.value { return nil }
	return r.Scan(value)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r Nullable) MarshalJSON() ([]byte, error) {
	if (r.value == nil) || (r.value.GetType() == NULLABLE_NIL) { return []byte("null"), nil }
	return r.value.MarshalJSON()
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Unmarshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Nullable) UnmarshalJSON(b []byte) error {
	if (r.value == nil) || (r.value.GetType() == NULLABLE_NIL) { return nil }
	return r.value.UnmarshalJSON(b)
}

/*
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
//	r.nullableType = NULLABLE_STRING
//	//r.ns.String = value
//	r.ns.SetValue(&value)
//	//r.ns.Valid = true
//	r.isNil = false
//	return true
}

func (r *Nullable) setTime(v *time.Time) error {
//	r.nullableType = NULLABLE_TIME
//	r.nt.Time = value
//	r.nt.Valid = true
//	r.isNil = false
//	return true
//	r.nullableType
	r.value = NewNullTime(v)
	return nil
}

func (r *Nullable) setNil() error {
	r.value = nil
	return nil
}
*/
