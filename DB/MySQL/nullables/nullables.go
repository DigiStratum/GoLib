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
	IsNil() bool
	Scan(value interface{}) error
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(b []byte) error
}

type NullableIfc interface {
	NullableValueIfc
	GetInt64Default(d int64) int64
	GetBoolDefault(d bool) bool
	GetFloat64Default(d float64) float64
	GetStringDefault(d string) string
	GetTimeDefault(d time.Time) time.Time
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
		case int: r.value = NewNullInt64(int64(v.(int)))
		case int8: r.value = NewNullInt64(int64(v.(int8)))
		case int16: r.value = NewNullInt64(int64(v.(int16)))
		case int32: r.value = NewNullInt64(int64(v.(int32)))
		case int64: r.value = NewNullInt64(v.(int64))
		case float32: r.value = NewNullFloat64(float64(v.(float32)))
		case float64: r.value = NewNullFloat64(v.(float64))
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

func (r *Nullable) GetInt64() *int64 {
	if nil == r.value { return nil }
	return r.value.GetInt64()
}

func (r *Nullable) GetInt64Default(d int64) int64 {
	if (nil == r) { return d }
	if v := r.GetInt64(); (nil != v) { return *v }
	return d;
}

func (r *Nullable) GetBool() *bool {
	if nil == r.value { return nil }
	return r.value.GetBool()
}

func (r *Nullable) GetBoolDefault(d bool) bool {
	if (nil == r) { return d }
	if v := r.GetBool(); (nil != v) { return *v }
	return d;
}

func (r *Nullable) GetFloat64() *float64 {
	if nil == r.value { return nil }
	return r.value.GetFloat64()
}

func (r *Nullable) GetFloat64Default(d float64) float64 {
	if (nil == r) { return d }
	if v := r.GetFloat64(); (nil != v) { return *v }
	return d;
}

func (r *Nullable) GetString() *string {
	if nil == r.value { return nil }
	return r.value.GetString()
}

func (r *Nullable) GetStringDefault(d string) string {
	if (nil == r) { return d }
	if v := r.GetString(); (nil != v) { return *v }
	return d;
}

func (r *Nullable) GetTime() *time.Time {
	if nil == r.value { return nil }
	return r.value.GetTime()
}

func (r *Nullable) GetTimeDefault(d time.Time) time.Time {
	if (nil == r) { return d }
	if v := r.GetTime(); (nil != v) { return *v }
	return d;
}

func (r *Nullable) IsNil() bool {
	if nil == r.value { return true }
	return r.value.IsNil()
}

// -------------------------------------------------------------------------------------------------
// database/sql.Scanner Public Interface
// -------------------------------------------------------------------------------------------------

// Scan for Nullable - we just sub it out to the underlying Nullable type
func (r *Nullable) Scan(value interface{}) error {
	if nil == r.value { return nil }
	return r.value.Scan(value)
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r *Nullable) MarshalJSON() ([]byte, error) {
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

