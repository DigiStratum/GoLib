package json

/*

Represent a JSON structure as an object tree with JavaScript-like selectors and other conveniences.

TODO:
 * Add support for ToJson() to spit it back out as a JSON string again!
 * Add support for array element and object property iterator
 * Add support for de|referencing; make references an embeddable json string (like mustache), use
   configurable start/stop delimiters with default; for whole-string references like "{{sel.ect.or}}"
   convert the value type to that of the selected reference, null if it doesn't exist. For partial
   references like "See also: {{sel.ect.or}}", convert the value type to string and perform string
   replacement, empty string if it doesn't exist. Introduce "RJSON" envelope to encode json metadata
   to describe the JSON encoding within, versioning, etc. to help with future-proofing, versioning,
   etc.

*/

import (
	"fmt"
	"strconv"
	"unicode"
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

	// Validity
	IsValid() bool
	GetType() ValueType

	// Nulls
	IsNull() bool
	SetNull() *JsonValue

	// Strings
	IsString() bool
	GetString() string
	SetString(value string) *JsonValue

	// Objects
	IsObject() bool
	PrepareObject() *JsonValue
	SetObjectProperty(name string, jsonValue *JsonValue) error
	DropObjectProperty(name string) error
	HasObjectProperty(name string) bool
	GetObjectPropertyNames() []string
	GetObjectProperty(name string) *JsonValue

	// Booleans
	IsBoolean() bool
	GetBoolean() bool
	SetBoolean(value bool) *JsonValue

	// Arrays
	IsArray() bool
	PrepareArray() *JsonValue
	GetArraySize() int
	GetArrayValue(index int) *JsonValue
	AppendArrayValue(jsonValue *JsonValue) error

	// Floats
	IsFloat() bool
	GetFloat() float64
	SetFloat(value float64) *JsonValue

	// Integers
	IsInteger() bool
	GetInteger() int64
	SetInteger(value int64) *JsonValue

	// Modern amenities ;^)
	Select(selector string) (*JsonValue, error)
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

// -----------------------------------------------
// Validity

func (r *JsonValue) IsValid() bool {
	return r.valueType > VALUE_TYPE_INVALID
}

func (r *JsonValue) GetType() ValueType {
	return r.valueType
}

// -----------------------------------------------
// Nulls

func (r *JsonValue) IsNull() bool {
	return r.valueType == VALUE_TYPE_NULL
}

func (r *JsonValue) SetNull() *JsonValue {
	r.valueType = VALUE_TYPE_NULL
	return r
}

// -----------------------------------------------
// Strings

func (r *JsonValue) IsString() bool {
	return r.valueType == VALUE_TYPE_STRING
}

func (r *JsonValue) GetString() string {
	if ! r.IsString() { return "" }
	return r.valueString
}

func (r *JsonValue) SetString(value string) *JsonValue {
	r.valueType = VALUE_TYPE_STRING
	r.valueString = value
	return r
}

// -----------------------------------------------
// Objects

func (r *JsonValue) IsObject() bool {
	return r.valueType == VALUE_TYPE_OBJECT
}

func (r *JsonValue) PrepareObject() *JsonValue {
	r.valueType = VALUE_TYPE_OBJECT
	r.valueObject = make(map[string]*JsonValue)
	return r
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

// -----------------------------------------------
// Booleans

func (r *JsonValue) IsBoolean() bool {
	return r.valueType == VALUE_TYPE_BOOLEAN
}

func (r *JsonValue) GetBoolean() bool {
	return r.IsBoolean() && r.valueBoolean
}

func (r *JsonValue) SetBoolean(value bool) *JsonValue {
	r.valueType = VALUE_TYPE_BOOLEAN
	r.valueBoolean = value
	return r
}

// -----------------------------------------------
// Arrays

func (r *JsonValue) IsArray() bool {
	return r.valueType == VALUE_TYPE_ARRAY
}

func (r *JsonValue) PrepareArray() *JsonValue {
	r.valueType = VALUE_TYPE_ARRAY
	r.valueArr = make([]*JsonValue, 0)
	return r
}

func (r *JsonValue) GetArraySize() int {
	if ! r.IsArray() { return 0 }
	return len(r.valueArr)
}

func (r *JsonValue) GetArrayValue(index int) *JsonValue {
	if ! r.IsArray() { return nil }
	if (index < 0) || (index >= len(r.valueArr)) { return nil }
	return r.valueArr[index]
}

func (r *JsonValue) AppendArrayValue(jsonValue *JsonValue) error {
	if nil == jsonValue { return fmt.Errorf("nil JsonValue cannot be appended to Array value") }
	r.valueArr = append(r.valueArr, jsonValue)
	return nil
}

// -----------------------------------------------
// Floats

func (r *JsonValue) IsFloat() bool {
	return r.valueType == VALUE_TYPE_FLOAT
}

func (r *JsonValue) GetFloat() float64 {
	if ! r.IsFloat() { return float64(0.0) }
	return r.valueFloat
}

func (r *JsonValue) SetFloat(value float64) *JsonValue {
	r.valueType = VALUE_TYPE_FLOAT
	r.valueFloat = value
	return r
}

// -----------------------------------------------
// Integers

func (r *JsonValue) IsInteger() bool {
	return r.valueType == VALUE_TYPE_INTEGER
}

func (r *JsonValue) GetInteger() int64 {
	if ! r.IsInteger() { return int64(0) }
	return r.valueInteger
}

func (r *JsonValue) SetInteger(value int64) *JsonValue {
	r.valueType = VALUE_TYPE_INTEGER
	r.valueInteger = value
	return r
}

// -----------------------------------------------
// Conveniences

func (r *JsonValue) Select(selector string) (*JsonValue, error) {
	// 1) An empty selector means we're already at the right place
	if 0 == len(selector) { return r, nil }

	// 1) If this isn't an Array or Object value...
	if ! (r.IsArray() || r.IsObject()) {
		return nil, fmt.Errorf("Selectors are only valid of Object or Array values")
	}

	// 2) Traverse the selector one element at a time
	objectProperty, arrayIndex, newSelector, err := r.selectNextElement(selector)
	if nil != err { return nil, err }
	if nil != objectProperty {
		if r.HasObjectProperty(*objectProperty) {
			return r.GetObjectProperty(*objectProperty).Select(newSelector) // <- BEWARE: recursion!
		}
		return nil, fmt.Errorf("Selected Object Property '%s' doesn't exist", *objectProperty)
	}
	if nil != arrayIndex {
		if r.GetArraySize() > *arrayIndex {
			return r.GetArrayValue(*arrayIndex).Select(newSelector) // <- BEWARE: recursion!
		}
		return nil, fmt.Errorf("Selected Array Index '%d' is out of bounds; Array size is %d", *arrayIndex, r.GetArraySize())
	}

	// selectNextElement() must return objectProperty, arrayIndex, or error and be handled above
	return nil, fmt.Errorf("Unexpected error for selector '%s'", selector)
}

// -----------------------------------------------
// Internal implementation

func (r *JsonValue) selectNextElement(selector string) (objectProperty *string, arrayIndex *int, newSelector string, err error) {
	// Return value defaults
	objectProperty = nil
	arrayIndex = nil
	newSelector = ""
	err = nil
	if len(selector) == 0 { return }
	cursor := 0
	if '[' == selector[cursor] {
		arrayIndex, newSelector, err = r.selectArrayIndexElement(selector)
	} else if '.' == selector[cursor] {
		objectProperty, newSelector, err = r.selectObjectPropertyElement(selector)
	}
	return
}

func (r *JsonValue) selectArrayIndexElement(selector string) (arrayIndex *int, newSelector string, err error) {
	// Return value defaults
	arrayIndex = nil
	newSelector = ""
	err = nil

	// Expect an array index terminated by ']'
	cursor := 0
	arrayIndexStr := ""
	for cursor = 1; cursor < len(selector); cursor++ {
		char := selector[cursor]
		// End of an array index?
		if ']' == char { break }
		// Digits?
		if ('0' <= char) || ('9' >= char) {
			arrayIndexStr = arrayIndexStr + string(char)
			continue
		}
		// Something else unexpected!
		err = fmt.Errorf("Unsupported character '%c' reading numeric array index in selector", char)
		return
	}
	// If we extracted something...
	if len(arrayIndexStr) > 0 {
		// Parse it into an array inted integer
		ai, _ := strconv.ParseInt(arrayIndexStr, 10, 32)
		aiInt := int(ai)
		arrayIndex = &aiInt
		// Chop it and the '[]' delimiters off the selector
		newSelector = selector[len(arrayIndexStr) + 2:]
	} else {
		err = fmt.Errorf("Missing numeric array index in selector")
	}
	return
}

func (r *JsonValue) selectObjectPropertyElement(selector string) (objectProperty *string, newSelector string, err error) {
//fmt.Printf("selectObjectPropertyElement() 1")
	// Return value defaults
	objectProperty = nil
	newSelector = ""
	err = nil

	// Expect an array index terminated by ']'
	cursor := 0
	objectPropertyStr := ""
	for cursor = 1; cursor < len(selector); cursor++ {
		char := selector[cursor]
		// Some other array index or property interupting?
		if ('[' == char) || ('.' == char) { break }
		// If it's not some odd white-space character...
		if (! unicode.IsSpace(rune(char))) {
			objectPropertyStr = objectPropertyStr + string(char)
			continue
		}
		err = fmt.Errorf("Unsupported character '%c' reading object property name in selector", char)
		return
	}
	// If we extracted something...
	if len(objectPropertyStr) > 0 {
		objectProperty = &objectPropertyStr
		// If the thing that stopped us was a '.' separator, chop it off along with what we found
		if (cursor < len(selector)) && ('.' == selector[cursor]) { cursor++ }
		newSelector = selector[cursor:]
	} else {
		err = fmt.Errorf("Missing object property name in selector")
	}
//if nil != objectProperty {
//	fmt.Printf("Found object property element '%s'\n", *objectProperty)
//} else {
//	fmt.Printf("No object property starting selector '%s'\n", selector)
//}
	return
}

