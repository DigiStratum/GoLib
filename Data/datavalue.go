package data

/*

Represent a data structure as a loosely typed object tree with JavaScript-like selectors and other conveniences.

On Error handling - while we would normally return multiple values along with an error as needed for
a given function result, because we had a goal of a chainable builder-style pattern here, and
because some of the builder mutators could result in an error, it was necessary to capture the error
and stash it instead of returning multiple value, since functions with multiple return values are
not chainable. So we stash the "last error" in a structure member property and make it available via
the GetError() method. In a chaining scenario, such an error would be lost if it were to occur mid-
stream in the chain after the next operation clears it. Either some sort of error history or, more
simply, a logger, would be necessary to preserve this through the chain to make it known to the
caller and/or developer. Although not all of the methods are intended for the purpose of chaining,
certainly only those that return the receiver pointer (a *DataValue) as the result, this same error
capture/get approach is used throughout for consistency. Thus any consumer/extension of this package
will be able to use a uniform method of error discovery and handling instead of varying by method.

TODO:
 * Add a generic selector Drop(selector string) method to Drop ANY matched selector from the Data?
 * Add a Pluck() method to pluck out one or more selectors as a new DataValue, exxectively a subset
   of the original, though capable of effectively replicating the entire original as with Copy()
 * Add support for chunked document loading for streaming data sources (avoid loading entire
   document into memory before lexing into structured data)
 * Add YAML loader/lexer like json
 * Add INI loader/lexer like json
 * Add XML loader/lexer like json
 * Add CSV loader/lexer like json
 * Add loader/lexers for Google Protocol Buffers (AKA protobuf), MessagePack, BSON (Binary JSON),
   and Avro (from Apache Hadoop) for faster/tighter data handling, application-to-application data
   exchange where human readability is less important
 * Add support for binary (bytearray) data type
 * Consider Iterating tree recursively for all data types, not just Object|Array; maybe some new type
   of iterator with an onMutation circuit breaker and callable (i.e. Iterator calls callable
*/

import (
	"fmt"
	"strings"
	"strconv"
	"unicode"

	"GoLib/Data/iterable"
)

type DataValueIfc interface {
	iterable.IterableIfc

	// State
	IsValid() bool
	GetType() DataType
	GetError() error
	IsImmutable() bool
	SetImmutable() *DataValue

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
	GetObjectProperties() []string
	GetObjectProperty(name string) *DataValue
	HasAllObjectProperties(names ...string) bool
        GetMissingObjectProperties(names ...string) *[]string
	DropObjectProperties(names ...string) *DataValue

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
	ReplaceArrayValue(index int, dataValue *DataValue) *DataValue

	// Floats
	IsFloat() bool
	GetFloat() float64
	SetFloat(value float64) *DataValue

	// Integers
	IsInteger() bool
	GetInteger() int64
	SetInteger(value int64) *DataValue

	// Modern amenities ;^)
	Select(selector string) *DataValue
	HasAll(selectors ...string) bool
	GetMissing(selectors ...string) []string
	Merge(dataValue DataValueIfc) *DataValue
	ToString() string
	ToJson() string
	Clone() *DataValue
}

type DataValue struct {
	err			error
	isImmutable		bool
	dataType		DataType
	valueBoolean		bool
	valueInteger		int64
	valueFloat		float64
	valueString		string
	valueArray		[]*DataValue
	valueObject		map[string]*DataValue
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewDataValue() *DataValue {
	r := DataValue{
		dataType:	DATA_TYPE_INVALID,
	}
	return &r
}

func NewNull() *DataValue { return NewDataValue().SetNull() }

func NewString(value string) *DataValue { return NewDataValue().SetString(value) }

func NewObject() *DataValue { return NewDataValue().PrepareObject() }

func NewBoolean(value bool) *DataValue { return NewDataValue().SetBoolean(value) }

func NewArray() *DataValue { return NewDataValue().PrepareArray() }

func NewFloat(value float64) *DataValue { return NewDataValue().SetFloat(value) }

func NewInteger(value int64) *DataValue { return NewDataValue().SetInteger(value) }

// -------------------------------------------------------------------------------------------------
// DataValueIfc
// -------------------------------------------------------------------------------------------------

// -----------------------------------------------
// State

func (r *DataValue) IsValid() bool {
	r.err = nil
	return r.dataType > DATA_TYPE_INVALID
}

func (r *DataValue) GetType() DataType {
	r.err = nil
	return r.dataType
}

func (r *DataValue) GetError() error {
	return r.err
}

func (r *DataValue) IsImmutable() bool {
	return r.isImmutable
}

func (r *DataValue) SetImmutable() *DataValue {
	r.isImmutable = true
	return r
}

// -----------------------------------------------
// Nulls

func (r *DataValue) IsNull() bool {
	r.err = nil
	return r.dataType == DATA_TYPE_NULL
}

func (r *DataValue) SetNull() *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	r.dataType = DATA_TYPE_NULL
	return r
}

// -----------------------------------------------
// Strings

func (r *DataValue) IsString() bool {
	r.err = nil
	return r.dataType == DATA_TYPE_STRING
}

func (r *DataValue) GetString() string {
	r.err = nil
	if ! r.IsString() { return "" }
	return r.valueString
}

func (r *DataValue) SetString(value string) *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	r.dataType = DATA_TYPE_STRING
	r.valueString = value
	return r
}

// -----------------------------------------------
// Objects

func (r *DataValue) IsObject() bool {
	r.err = nil
	return r.dataType == DATA_TYPE_OBJECT
}

func (r *DataValue) PrepareObject() *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	r.dataType = DATA_TYPE_OBJECT
	r.valueObject = make(map[string]*DataValue)
	return r
}

func (r *DataValue) SetObjectProperty(name string, dataValue *DataValue) *DataValue {
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, cannot set property; use PrepareObject() first!")
		return r
	}
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil

	// Don't add nil DataValue into map; Use DATA_TYPE_NULL DataValue for JSON NULL value
	if nil != dataValue { r.valueObject[name] = dataValue }
	return r
}

func (r *DataValue) DropObjectProperty(name string) *DataValue {
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, cannot drop property; use PrepareObject() first!")
		return r
	}
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil

	// Delete property if exists; non-existent is non-error: caller already has desired result
	if _, ok := r.valueObject[name]; ok { delete(r.valueObject, name) }
	return r
}

func (r *DataValue) HasObjectProperty(name string) bool {
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, cannot check property; use PrepareObject() first!")
		return false
	}
	r.err = nil
	_, ok := r.valueObject[name]
	return ok
}

func (r *DataValue) GetObjectProperties() []string {
	// TODO: Cache this internally so that it doesn't need to be done on-the-fly for subsequent requests
	names := make([]string, 0)
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, cannot set property; use PrepareObject() first!")
		return names
	}
	r.err = nil
	if r.IsObject() {
		for name, _ := range r.valueObject { names = append(names, name) }
	}
	return names
}

func (r *DataValue) GetObjectProperty(name string) *DataValue {
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, cannot get property; use PrepareObject() first!")
		return nil
	}
	r.err = nil
	value, _ := r.valueObject[name]
	return value
}

func (r *DataValue) HasAllObjectProperties(names ...string) bool {
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, it has no properties; use PrepareObject() first!")
		return false
	}
	r.err = nil
	for _, name := range names {
		if _, ok := r.valueObject[name]; ! ok { return false }
	}
	return true
}

func (r *DataValue) GetMissingObjectProperties(names ...string) *[]string {
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, it has no properties; use PrepareObject() first!")
		return nil
	}
	r.err = nil
	missing := make([]string, 0)
	for _, name := range names {
		if _, ok := r.valueObject[name]; ! ok { missing = append(missing, name) }
	}
	return &missing
}

func (r *DataValue) DropObjectProperties(names ...string) *DataValue {
	if DATA_TYPE_OBJECT != r.dataType {
		r.err = fmt.Errorf("Not an object type, it has no properties; use PrepareObject() first!")
		return r
	}
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	for _, name := range names {
		if _, ok := r.valueObject[name]; ok { delete(r.valueObject, name) }
	}
	return r
}

// -----------------------------------------------
// Booleans

func (r *DataValue) IsBoolean() bool {
	r.err = nil
	return r.dataType == DATA_TYPE_BOOLEAN
}

func (r *DataValue) GetBoolean() bool {
	r.err = nil
	return r.IsBoolean() && r.valueBoolean
}

func (r *DataValue) SetBoolean(value bool) *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	r.dataType = DATA_TYPE_BOOLEAN
	r.valueBoolean = value
	return r
}

// -----------------------------------------------
// Arrays

func (r *DataValue) IsArray() bool {
	r.err = nil
	return r.dataType == DATA_TYPE_ARRAY
}

func (r *DataValue) PrepareArray() *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	r.dataType = DATA_TYPE_ARRAY
	r.valueArray = make([]*DataValue, 0)
	return r
}

func (r *DataValue) GetArraySize() int {
	if DATA_TYPE_ARRAY != r.dataType {
		r.err = fmt.Errorf("Not an array type; use PrepareArray() first!")
		return 0
	}
	r.err = nil
	return len(r.valueArray)
}

func (r *DataValue) GetArrayValue(index int) *DataValue {
	if DATA_TYPE_ARRAY != r.dataType {
		r.err = fmt.Errorf("Not an array type; use PrepareArray() first!")
		return nil
	}
	r.err = nil
	if (index < 0) || (index >= len(r.valueArray)) {
		r.err = fmt.Errorf("Array index %d out of bounds; valid range is 0 to %d", index, (len(r.valueArray) - 1))
		return nil
	}
	return r.valueArray[index]
}

func (r *DataValue) AppendArrayValue(dataValue *DataValue) *DataValue {
	if DATA_TYPE_ARRAY != r.dataType {
		r.err = fmt.Errorf("Not an array type; use PrepareArray() first!")
		return r
	}
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	if nil == dataValue {
		r.err = fmt.Errorf("nil DataValue cannot be appended to Array value")
		return r
	}
	r.err = nil
	r.valueArray = append(r.valueArray, dataValue)
	return r
}
func (r *DataValue) ReplaceArrayValue(index int, dataValue *DataValue) *DataValue {
	if DATA_TYPE_ARRAY != r.dataType {
		r.err = fmt.Errorf("Not an array type; use PrepareArray() first!")
		return nil
	}
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	if (index < 0) || (index >= len(r.valueArray)) {
		r.err = fmt.Errorf("Array index %d out of bounds; valid range is 0 to %d", index, (len(r.valueArray) - 1))
		return nil
	}
	if nil == dataValue { dataValue = NewNull() }
	r.valueArray[index] = dataValue
	return r
}

// -----------------------------------------------
// Floats

func (r *DataValue) IsFloat() bool {
	r.err = nil
	return r.dataType == DATA_TYPE_FLOAT
}

func (r *DataValue) GetFloat() float64 {
	r.err = nil
	if ! r.IsFloat() { return float64(0.0) }
	return r.valueFloat
}

func (r *DataValue) SetFloat(value float64) *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	r.dataType = DATA_TYPE_FLOAT
	r.valueFloat = value
	return r
}

// -----------------------------------------------
// Integers

func (r *DataValue) IsInteger() bool {
	r.err = nil
	return r.dataType == DATA_TYPE_INTEGER
}

func (r *DataValue) GetInteger() int64 {
	r.err = nil
	if ! r.IsInteger() { return int64(0) }
	return r.valueInteger
}

func (r *DataValue) SetInteger(value int64) *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	r.dataType = DATA_TYPE_INTEGER
	r.valueInteger = value
	return r
}

// -----------------------------------------------
// Conveniences

func (r *DataValue) Select(selector string) *DataValue {
	r.err = nil
	// 1) An empty selector means we're already at the right place
	if 0 == len(selector) { return r }

	// 1) If this isn't an Array or Object value...
	if ! (r.IsArray() || r.IsObject()) {
		r.err = fmt.Errorf("Selectors are only valid for Object or Array values")
		return nil
	}

	// 2) Traverse the selector one element at a time
	objectProperty, arrayIndex, newSelector, err := r.selectNextElement(selector)
	if nil != err { r.err = err; return nil }
	if nil != objectProperty {
		if r.HasObjectProperty(*objectProperty) {
			// If the new selector starts with a '.' (object property separator) then chop it off
			return r.GetObjectProperty(*objectProperty).Select(newSelector) // <- BEWARE: recursion!
		}
		r.err = fmt.Errorf("Selected Object Property '%s' doesn't exist", *objectProperty)
		return nil
	}
	if nil != arrayIndex {
		if r.GetArraySize() > *arrayIndex {
			return r.GetArrayValue(*arrayIndex).Select(newSelector) // <- BEWARE: recursion!
		}
		r.err = fmt.Errorf("Selected Array Index '%d' is out of bounds; Array size is %d", *arrayIndex, r.GetArraySize())
		return nil
	}

	// selectNextElement() must return objectProperty, arrayIndex, or error and be handled above
	r.err = fmt.Errorf("Unexpected error for selector '%s'", selector)
	return nil
}

func (r *DataValue) HasAll(selectors ...string) bool {
	r.err = nil
	// For each selector in the variadic list...
	for _, selector := range selectors {
		// If we found no DataValue or hit an error, then we don't have it, therefore FALSE!
		if res := r.Select(selector); nil == res { return false }
	}
	return true
}

func (r *DataValue) GetMissing(selectors ...string) []string {
	r.err = nil
	missing := make([]string, 0)
	// For each selector in the variadic list...
	for _, selector := range selectors {
		// If we found no DataValue or hit an error, then we don't have it, therefore MISSING!
		if res := r.Select(selector); nil == res {
			missing = append(missing, selector)
		}
	}
	return missing
}

func (r *DataValue) Merge(dataValue DataValueIfc) *DataValue {
	if r.isImmutable {
		r.err = fmt.Errorf("Data is immutable, cannot modify!")
		return r
	}
	r.err = nil
	if nil == dataValue {
		r.err = fmt.Errorf("nil merge value, nothing possible!")
		return r
	}
	if r.IsObject() && dataValue.IsObject() {
		// Key used to deduplicate object properties; existing key values will be overwritten
		it := dataValue.GetIterator()
		for kvpi := it(); nil != kvpi; kvpi = it() {
			if kvp, ok := kvpi.(KeyValuePair); ok {
				r.SetObjectProperty(kvp.Key, kvp.Value)
			}
		}
	} else if r.IsArray() && dataValue.IsArray() {
		// Note: No deduplication for array values; pile dataValue's entries onto our tail
		it := dataValue.GetIterator()
		for ivpi := it(); nil != ivpi; ivpi = it() {
			if ivp, ok := ivpi.(IndexValuePair); ok {
				r.AppendArrayValue(ivp.Value)
			}
		}
	} else {
		r.err = fmt.Errorf("mismatch value type, nothing possible!")
		return r
	}
	return r
}

func (r *DataValue) ToString() string {
	r.err = nil
	return r.stringify(false)
}

func (r *DataValue) ToJson() string {
	r.err = nil
	return r.stringify(true)
}

// Clone a deep copy of this entire thing; all pointers dereferenced and copied by value
func (r *DataValue) Clone() *DataValue {
	dv := DataValue{
		err:			r.err,
		isImmutable:		r.isImmutable,
		dataType:		r.dataType,
		valueBoolean:		r.valueBoolean,
		valueInteger:		r.valueInteger,
		valueFloat:		r.valueFloat,
		valueString:		r.valueString,
	}
	switch r.dataType {
		case DATA_TYPE_ARRAY:
			dv.valueArray = make([]*DataValue, 0)
			for _, arrayValue := range r.valueArray {
				dv.valueArray = append(dv.valueArray, arrayValue.Clone()) // <- BEWARE: recursion!
			}

		case DATA_TYPE_OBJECT:
			dv.valueObject = make(map[string]*DataValue)
			for key, objectValue := range r.valueObject {
				dv.valueObject[key] = objectValue.Clone() // <- BEWARE: recursion!
			}
	}

	return &dv
}

// -------------------------------------------------------------------------------------------------
// IterableIfc
// -------------------------------------------------------------------------------------------------

type KeyValuePair struct {
        Key	string
        Value	*DataValue
}

type IndexValuePair struct {
        Index	int
        Value	*DataValue
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
			return IndexValuePair{Index: prev_idx, Value: r.GetArrayValue(prev_idx)}
		}
	}
	r.err = fmt.Errorf("DataValue is neither an Object, nor Array, so cannot iterate!")
	return nil
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
	if r.isValueArrayStartChar(selector[cursor]) {
		arrayIndex, newSelector, err = r.selectArrayIndexElement(selector)
	} else if r.isValueObjectStartChar(selector[cursor]) {
		objectProperty, newSelector, err = r.selectObjectPropertyElement(selector)
	}
	return
}

func (r *DataValue) isValueArrayStartChar(ch byte) bool {
	return '[' == ch
}

func (r *DataValue) isValueObjectStartChar(ch byte) bool {
	return (('a' <= ch) && ('z' >= ch) || ('A' <= ch) && ('Z' >= ch))
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
		if ('0' <= char) && ('9' >= char) {
			arrayIndexStr = arrayIndexStr + string(char)
			continue
		}
		// Something else unexpected!
		err = fmt.Errorf("Unsupported character '%c' reading numeric array index in selector", char)
		return
	}
	// If we extracted something...
	if len(arrayIndexStr) > 0 {
		// Parse it into an array index integer
		ai, _ := strconv.ParseInt(arrayIndexStr, 10, 32)
		aiInt := int(ai)
		arrayIndex = &aiInt
		// Chop it and the '[]' delimiters off the selector
		nextPos := len(arrayIndexStr) + 2
		if len(selector) > nextPos {
			if '.' == selector[nextPos] {
				newSelector = selector[nextPos+1:]
			} else if '[' != selector[nextPos] {
				err = fmt.Errorf("No valid separator found trailing this selector segment ")
			}
			//fmt.Printf("newSelector:[%s]\n", newSelector)
		}
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
	for cursor = 0; cursor < len(selector); cursor++ {
		char := selector[cursor]
		// Some other array index or property interupting?
		if ('[' == char) { break }
		if ('.' == char) { cursor++; break }
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
		//fmt.Printf("newSelector:[%s]\n", newSelector)
	} else {
		err = fmt.Errorf("Missing object property name in selector")
	}

	return
}

func (r *DataValue) stringify(quoteStrings bool) string {
	switch r.dataType {
		case DATA_TYPE_NULL: return "null"

		case DATA_TYPE_STRING:
			if quoteStrings { return strconv.Quote(r.valueString) }
			return r.valueString

		case DATA_TYPE_BOOLEAN:
			if r.valueBoolean { return "true" }
			return "false"

		case DATA_TYPE_INTEGER: return fmt.Sprint(r.valueInteger)

		case DATA_TYPE_FLOAT: return fmt.Sprint(r.valueFloat)

		case DATA_TYPE_ARRAY:
			var sb strings.Builder
			sb.WriteString("[")
			sep := ""
			for _, value := range r.valueArray {
				// Note: always quote strings in structured data, otherwise they can break the structure!
				strValue := value.stringify(true) // <- Recusrsion Alert!
				sb.WriteString(fmt.Sprintf("%s%s", sep, strValue))
				sep = ","
			}
			sb.WriteString("]")
			return sb.String()

		case DATA_TYPE_OBJECT:
			var sb strings.Builder
			sb.WriteString("{")
			sep := ""
			for key, value := range r.valueObject {
				strKey := strconv.Quote(key)
				// Note: always quote strings in structured data, otherwise they can break the structure!
				strValue := value.stringify(true) // <- Recusrsion Alert!
				sb.WriteString(fmt.Sprintf("%s%s:%s", sep, strKey, strValue))
				sep = ","
			}
			sb.WriteString("}")
			return sb.String()
	}
	return ""
}

