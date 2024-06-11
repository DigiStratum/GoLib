package json

/*

Represent a JSON structure as an object tree with JavaScript-like selectors and other conveniences.

TODO:
 * Add support for DOM-like mutability
 * Add support for Marshal() to spit it back out as a JSON string again!

*/

import (
	"fmt"
	//"unicode/utf8"
)

type ValueType int

const (
	VALUE_TYPE_INVALID ValueType = iota
	VALUE_TYPE_NULL
	VALUE_TYPE_BOOLEAN
	VALUE_TYPE_INTEGER
	VALUE_TYPE_FLOAT
	VALUE_TYPE_STRING
	VALUE_TYPE_OBJECT
	VALUE_TYPE_ARRAY
)

type JsonValueIfc interface {
	// Accessors
	IsValid() bool
	IsNull() bool
	IsBoolean() bool
	IsInteger() bool
	IsFloat() bool
	IsString() bool
	IsArray() bool
	IsObject() bool

	GetBoolean() bool

	GetInteger() int64

	GetFloat() float64

	GetString() string

	GetArraySize() int
	GetArrayElement(index int) *JsonValue

	HasObjectProperty(name string) bool
	GetObjectPropertyNames() []string
	GetObjectProperty(name string) *JsonValue

	// TODO: Implement this bad boy!
	//Select(selector string) (*JsonValue, error)

	// Mutators
	SetString(value string)
	PrepareObject()
	SetObjectProperty(name string, jsonValue *JsonValue) error
	DropObjectProperty(name string) error

	SetBoolean(value bool)
	SetNull()

	PrepareArray()
	AppendArrayValue(jsonValue *JsonValue) error
}

type JsonValue struct {
	valueType		ValueType
	valueBoolean		bool
	valueInteger		int64
	valueFloat		float64
	valueString		string
	valueArr		[]*JsonValue
	valueObject		map[string]*JsonValue
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these for programmatic construction
func NewJsonValue() *JsonValue {
	r := JsonValue{
		valueType:	VALUE_TYPE_INVALID,
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// JsonValueIfc
// -------------------------------------------------------------------------------------------------

func (r *JsonValue) IsValid() bool {
	return r.valueType > VALUE_TYPE_INVALID
}

func (r *JsonValue) IsNull() bool {
	return r.valueType == VALUE_TYPE_NULL
}

func (r *JsonValue) IsBoolean() bool {
	return r.valueType == VALUE_TYPE_BOOLEAN
}

func (r *JsonValue) IsInteger() bool {
	return r.valueType == VALUE_TYPE_INTEGER
}

func (r *JsonValue) IsFloat() bool {
	return r.valueType == VALUE_TYPE_FLOAT
}

func (r *JsonValue) IsString() bool {
	return r.valueType == VALUE_TYPE_STRING
}

func (r *JsonValue) IsArray() bool {
	return r.valueType == VALUE_TYPE_ARRAY
}

func (r *JsonValue) IsObject() bool {
	return r.valueType == VALUE_TYPE_OBJECT
}

func (r *JsonValue) GetBoolean() bool {
	return r.IsBoolean() && r.valueBoolean
}

func (r *JsonValue) GetInteger() int64 {
	if ! r.IsInteger() { return int64(0) }
	return r.valueInteger
}

func (r *JsonValue) GetFloat() float64 {
	if ! r.IsFloat() { return float64(0.0) }
	return r.valueFloat
}

func (r *JsonValue) GetString() string {
	if ! r.IsString() { return "" }
	return r.valueString
}

func (r *JsonValue) SetString(value string) {
	r.valueType = VALUE_TYPE_STRING
	r.valueString = value
}

func (r *JsonValue) GetArraySize() int {
	if ! r.IsArray() { return 0 }
	return len(r.valueArr)
}

func (r *JsonValue) GetArrayElement(index int) *JsonValue {
	if ! r.IsArray() { return nil }
	if (index < 0) || (index >= len(r.valueArr)) { return nil }
	return r.valueArr[index]
}

// Objects

func (r *JsonValue) PrepareObject() {
	r.valueType = VALUE_TYPE_OBJECT
	r.valueObject = make(map[string]*JsonValue)
}

func (r *JsonValue) SetObjectProperty(name string, jsonValue *JsonValue) error {
	if VALUE_TYPE_OBJECT != r.valueType {
		return fmt.Errorf("Not an object type, cannot set object property; use PrepareObject() first!")
	}

	// Don't add nil JsonValue into map; Use VALUE_TYPE_NULL JsonValue for JSON NULL value
	if nil != jsonValue { r.valueObject[name] = jsonValue }
	return nil
}

func (r *JsonValue) DropObjectProperty(name string) error {
	if VALUE_TYPE_OBJECT != r.valueType {
		return fmt.Errorf("Not an object type, cannot drop object property; use PrepareObject() first!")
	}

	// Delete property if exists; non-existent is non-error: caller already has desired result
	if _, ok := r.valueObject[name]; ok { delete(r.valueObject, name) }
	return nil
}

func (r *JsonValue) HasObjectProperty(name string) bool {
	if ! r.IsObject() { return false }
	_, ok := r.valueObject[name]
	return ok
}

func (r *JsonValue) GetObjectPropertyNames() []string {
	// TODO: Cache this internally so that it doesn't need to be done on-the-fly for subsequent requests
	names := make([]string, 0)
	if r.IsObject() {
		for name, _ := range r.valueObject { names = append(names, name) }
	}
	return names
}

func (r *JsonValue) GetObjectProperty(name string) *JsonValue {
	if ! r.IsObject() { return nil }
	value, _ := r.valueObject[name]
	return value
}

// Booleans

func (r *JsonValue) SetBoolean(value bool) {
	r.valueType = VALUE_TYPE_BOOLEAN
	r.valueBoolean = value
}

// Nulls

func (r *JsonValue) SetNull() {
	r.valueType = VALUE_TYPE_NULL
}

// Arrays

func (r *JsonValue) PrepareArray() {
	r.valueType = VALUE_TYPE_ARRAY
	r.valueArr = make([]*JsonValue, 0)
}

func (r *JsonValue) AppendArrayValue(jsonValue *JsonValue) error {
	if nil == jsonValue { return fmt.Errorf("nil JsonValue cannot be appended to Array value") }
	r.valueArr = append(r.valueArr, jsonValue)
	return nil
}

