package data

/*

Represent a data structure as a loosely typed object tree with JavaScript-like selectors and other conveniences.

TODO:
 * Add support for [de]serialization
 * Add support for [de]referencing; make references an embeddable string (like mustache), use
   configurable start/stop delimiters with default; for whole-string references like "{{sel.ect.or}}"
   convert the value type to that of the selected reference, null if it doesn't exist. For partial
   references like "See also: {{sel.ect.or}}", convert the value type to string and perform string
   replacement, empty string if it doesn't exist. Introduce "RDATA" envelope to encode metadata to
   describe the encoding within, versioning, etc. to help with future-proofing, versioning, etc.
 * Add support for immutability; once the immutable flag is set, no more changes allowed, read-only?
 * Add support for start/stop DataValue callback events
 * Add support for chunked document loading for streaming data sources (avoid loading entire
   document into memory before lexing into structured data)
 * Add YAML loader/lexer like json
 * Add INI loader/lexer like json
 * Add XML loader/lexer like json
 * Add support for conveniences of Hashmap, Config, and other popular libraries like underscore.js
   with "pluck", etc
*/

import (
	"fmt"
	"strconv"
	"unicode"

	"GoLib/Data/iterable"
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

type DataValueIfc interface {
	iterable.IterableIfc

	// Validity
	IsValid() bool
	GetType() ValueType
	GetError() error

	// Nulls
	IsNull() bool
	SetNull() *DataValue

	// Strings
	IsString() bool
	GetString() string
	SetString(value string) *DataValue

	// Objects
	IsObject() bool
	PrepareObject() *DataValue
	SetObjectProperty(name string, dataValue *DataValue) *DataValue
	DropObjectProperty(name string) *DataValue
	HasObjectProperty(name string) bool
	GetObjectPropertyNames() []string
	GetObjectProperty(name string) *DataValue

	// Booleans
	IsBoolean() bool
	GetBoolean() bool
	SetBoolean(value bool) *DataValue

	// Arrays
	IsArray() bool
	PrepareArray() *DataValue
	GetArraySize() int
	GetArrayValue(index int) *DataValue
	AppendArrayValue(dataValue *DataValue) *DataValue

	// Floats
	IsFloat() bool
	GetFloat() float64
	SetFloat(value float64) *DataValue

	// Integers
	IsInteger() bool
	GetInteger() int64
	SetInteger(value int64) *DataValue

	// Modern amenities ;^)
	Select(selector string) (*DataValue, error)
}

type DataValue struct {
	err			error
	valueType		ValueType
	valueBoolean		bool
	valueInteger		int64
	valueFloat		float64
	valueString		string
	valueArr		[]*DataValue
	valueObject		map[string]*DataValue
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these for programmatic construction
func NewDataValue() *DataValue {
	r := DataValue{
		valueType:	VALUE_TYPE_INVALID,
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// DataValueIfc
// -------------------------------------------------------------------------------------------------

// -----------------------------------------------
// Validity

func (r *DataValue) IsValid() bool {
	r.err = nil
	return r.valueType > VALUE_TYPE_INVALID
}

func (r *DataValue) GetType() ValueType {
	r.err = nil
	return r.valueType
}

func (r *DataValue) GetError() error {
	return r.err
}

// -----------------------------------------------
// Nulls

func (r *DataValue) IsNull() bool {
	r.err = nil
	return r.valueType == VALUE_TYPE_NULL
}

func (r *DataValue) SetNull() *DataValue {
	r.err = nil
	r.valueType = VALUE_TYPE_NULL
	return r
}

// -----------------------------------------------
// Strings

func (r *DataValue) IsString() bool {
	r.err = nil
	return r.valueType == VALUE_TYPE_STRING
}

func (r *DataValue) GetString() string {
	r.err = nil
	if ! r.IsString() { return "" }
	return r.valueString
}

func (r *DataValue) SetString(value string) *DataValue {
	r.err = nil
	r.valueType = VALUE_TYPE_STRING
	r.valueString = value
	return r
}

// -----------------------------------------------
// Objects

func (r *DataValue) IsObject() bool {
	r.err = nil
	return r.valueType == VALUE_TYPE_OBJECT
}

func (r *DataValue) PrepareObject() *DataValue {
	r.err = nil
	r.valueType = VALUE_TYPE_OBJECT
	r.valueObject = make(map[string]*DataValue)
	return r
}

func (r *DataValue) SetObjectProperty(name string, dataValue *DataValue) *DataValue {
	r.err = nil
	if VALUE_TYPE_OBJECT != r.valueType {
		r.err = fmt.Errorf("Not an object type, cannot set object property; use PrepareObject() first!")
	} else {
		// Don't add nil DataValue into map; Use VALUE_TYPE_NULL DataValue for JSON NULL value
		if nil != dataValue { r.valueObject[name] = dataValue }
	}
	return r
}

func (r *DataValue) DropObjectProperty(name string) *DataValue {
	r.err = nil
	if VALUE_TYPE_OBJECT != r.valueType {
		r.err = fmt.Errorf("Not an object type, cannot drop object property; use PrepareObject() first!")
	} else {

		// Delete property if exists; non-existent is non-error: caller already has desired result
		if _, ok := r.valueObject[name]; ok { delete(r.valueObject, name) }
	}
	return r
}

func (r *DataValue) HasObjectProperty(name string) bool {
	r.err = nil
	if ! r.IsObject() { return false }
	_, ok := r.valueObject[name]
	return ok
}

func (r *DataValue) GetObjectPropertyNames() []string {
	r.err = nil
	// TODO: Cache this internally so that it doesn't need to be done on-the-fly for subsequent requests
	names := make([]string, 0)
	if r.IsObject() {
		for name, _ := range r.valueObject { names = append(names, name) }
	}
	return names
}

func (r *DataValue) GetObjectProperty(name string) *DataValue {
	r.err = nil
	if ! r.IsObject() { return nil }
	value, _ := r.valueObject[name]
	return value
}

// -----------------------------------------------
// Booleans

func (r *DataValue) IsBoolean() bool {
	r.err = nil
	return r.valueType == VALUE_TYPE_BOOLEAN
}

func (r *DataValue) GetBoolean() bool {
	r.err = nil
	return r.IsBoolean() && r.valueBoolean
}

func (r *DataValue) SetBoolean(value bool) *DataValue {
	r.err = nil
	r.valueType = VALUE_TYPE_BOOLEAN
	r.valueBoolean = value
	return r
}

// -----------------------------------------------
// Arrays

func (r *DataValue) IsArray() bool {
	r.err = nil
	return r.valueType == VALUE_TYPE_ARRAY
}

func (r *DataValue) PrepareArray() *DataValue {
	r.err = nil
	r.valueType = VALUE_TYPE_ARRAY
	r.valueArr = make([]*DataValue, 0)
	return r
}

func (r *DataValue) GetArraySize() int {
	r.err = nil
	if ! r.IsArray() { return 0 }
	return len(r.valueArr)
}

func (r *DataValue) GetArrayValue(index int) *DataValue {
	r.err = nil
	if ! r.IsArray() { return nil }
	if (index < 0) || (index >= len(r.valueArr)) { return nil }
	return r.valueArr[index]
}

func (r *DataValue) AppendArrayValue(dataValue *DataValue) *DataValue {
	r.err = nil
	if nil == dataValue {
		r.err = fmt.Errorf("nil DataValue cannot be appended to Array value")
	} else {
		r.valueArr = append(r.valueArr, dataValue)
	}
	return r
}

// -----------------------------------------------
// Floats

func (r *DataValue) IsFloat() bool {
	r.err = nil
	return r.valueType == VALUE_TYPE_FLOAT
}

func (r *DataValue) GetFloat() float64 {
	r.err = nil
	if ! r.IsFloat() { return float64(0.0) }
	return r.valueFloat
}

func (r *DataValue) SetFloat(value float64) *DataValue {
	r.err = nil
	r.valueType = VALUE_TYPE_FLOAT
	r.valueFloat = value
	return r
}

// -----------------------------------------------
// Integers

func (r *DataValue) IsInteger() bool {
	r.err = nil
	return r.valueType == VALUE_TYPE_INTEGER
}

func (r *DataValue) GetInteger() int64 {
	r.err = nil
	if ! r.IsInteger() { return int64(0) }
	return r.valueInteger
}

func (r *DataValue) SetInteger(value int64) *DataValue {
	r.err = nil
	r.valueType = VALUE_TYPE_INTEGER
	r.valueInteger = value
	return r
}

// -----------------------------------------------
// Conveniences

func (r *DataValue) Select(selector string) (*DataValue, error) {
	r.err = nil
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

func (r *DataValue) selectNextElement(selector string) (objectProperty *string, arrayIndex *int, newSelector string, err error) {
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

func (r *DataValue) selectArrayIndexElement(selector string) (arrayIndex *int, newSelector string, err error) {
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

func (r *DataValue) selectObjectPropertyElement(selector string) (objectProperty *string, newSelector string, err error) {
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

	return
}

// -------------------------------------------------------------------------------------------------
// IterableIfc
// -------------------------------------------------------------------------------------------------

type KeyValuePair struct {
        Key     string
        Value   *DataValue
}

// Returns iterator func of []KeyValuePair for Objects, []*DataValue for Arrays, nil for other types
func (r *DataValue) GetIterator() func () interface{} {
	r.err = nil
	// Return object KeyValuePairs
	if r.IsObject() {
		kvps := make([]KeyValuePair, 0)
		var idx int = 0
		for k, v := range r.valueObject {
			kvps = append(kvps, KeyValuePair{ Key: k, Value: v })
			idx++
		}
		idx = 0
		return func () interface{} {
			// If we're done iterating, return nothing
			if idx >= len(kvps) { return nil }
			prev_idx := idx
			idx++
			return kvps[prev_idx]
		}
	}
	// Return Array values
	if r.IsArray() {
		idx := 0
		var data_len = r.GetArraySize()
		return func () interface{} {
			// If we're done iterating, return nothing
			if idx >= data_len { return nil }
			prev_idx := idx
			idx++
			return r.GetArrayValue(prev_idx)
		}
	}
	r.err = fmt.Errorf("DataValue is neither an Object, nor Array, so cannot iterate!")
	return nil
}

