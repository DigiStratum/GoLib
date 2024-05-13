package json

/*

Represent a JSON structure as an object tree with JavaScript-like selectors and other conveniences.

TODO:
 * If we make our own tokenizer/parser, then we can unmarshall into a custom tree structure
   that separates JSON types from values to make it easier to traverse, access, and assert
 * Add support for mutability and Marshal() to spit it back out as a JSON string again!
 * Make DependencyInjectable to accept a logger for errors, debug, etc
*/

import (
	"fmt"
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
	GetString() []rune

	GetArraySize() int
	GetArrayElement(index int) *JsonValue

	HasObjectProperty(name string) bool
	GetObjectPropertyNames() []string
	GetObjectProperty(name string) *JsonValue

	// TODO: Implement this bad boy!
	//Select(selector string) (*JsonValue, error)
}

type JsonValue struct {
	startPos, stopPos	int
	valueType		ValueType
	valueBoolean		bool
	valueInteger		int64
	valueFloat		float64
	valueString		[]rune
	valueArr		[]*JsonValue
	valueMap		map[string]*JsonValue
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewJsonValue(json *[]rune) *JsonValue {
	r, _ := unmarshal(json)
	// TODO: pass error along to DI Logger
	return r
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

func (r *JsonValue) GetString() []rune {
	if ! r.IsString() { return []rune("") }
	return r.valueString
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

func (r *JsonValue) HasObjectProperty(name string) bool {
	if ! r.IsObject() { return false }
	_, ok := r.valueMap[name]
	return ok
}

func (r *JsonValue) GetObjectPropertyNames() []string {
	// TODO: Cache this internally so that it doesn't need to be done on-the-fly for subsequent requests
	names := make([]string, 0)
	if r.IsObject() {
		for name, _ := range r.valueMap { names = append(names, name) }
	}
	return names
}

func (r *JsonValue) GetObjectProperty(name string) *JsonValue {
	if ! r.IsObject() { return nil }
	value, _ := r.valueMap[name]
	return value
}

// -------------------------------------------------------------------------------------------------
// JsonValue
// -------------------------------------------------------------------------------------------------

/*
Structures under consideration:

NULL:
  'null'

STRING:
  '"value"'

NUMBER:
  '1'
  '3.14'

BOOL:
  'true'|'false'

OBJECT:
  '{}'
  '{"a": null, "b": "c", "d": 1, "e": 3.14, "f": true, "g": {}, "h": {"a": "b"}, "i": [], "j": ["apples"]}'

ARRAY:
  '[]'
  '[null, "a", 1, 3.14, false, {}, {"a": "b"}, [], ["apples"]]'

Tokenizer notes:
 * Parse JSON string character by character
   * Determine how to do this safely for multi-byte glyphs
   * Default to UTF-8, allow override to maximize utility value
 * Consume whitespace between tokens
 * Initial state: expect ValueToken
 * For each ValueToken:
   1) Assert ValueToken.Type based on first character only:
     NULL 	-> /[n]/i
     STRING 	-> /["]/
     NUMBER 	-> /[0-9]/
     BOOL 	-> /[tf]]i
     OBJECT 	-> /[\{]/
     ARRAY 	-> /[\[]/

   2) Append subsequent chars to ValueToken.RawValue until expected terminator char based on ValueToken.Type:
     NULL	-> whitespace|EOF
     STRING	-> /["]/ (handle escapes; beware of escaped escapes like '\\"'; the escape is escaped, not the quote!)
     NUMBER	-> whitespace|EOF (allow anything in, fail in validation for non-numeric garbage)
     BOOL	-> whitespace|EOF (allow anything in, fail in validation for non-boolean garbage)
     OBJECT	-> enter recursion to expect NameToken : ValueToken until /[\}]/ clean close
     ARRAY	-> enter recursion to expect ValueToken, repeat for each [,] until /[\]]/ clean close

   3) Validate ValueToken.RawValue based on ValueToken.Type (require ValueToken.CleanClose == true)
     NULL	-> /^null$/i
     STRING	-> /^".*"$/
     NUMBER 	-> /^[0-9]+(\.[0-9]+)*$/
     BOOL	-> /^(true|false)$/i
     OBJECT	-> [NameToken:ValueToken] collection, non-regex
     ARRAY	-> [ValueToken] collection, non-regex
*/

// -------------------------------------------------------------------------------------------------
// private implementation
// -------------------------------------------------------------------------------------------------

// JSON Lexer
const (
	_LEXER_STATE_SEEK_NEXT_VALUE int = iota
	_LEXER_STATE_DONE
)

func unmarshal(json *[]rune) (*JsonValue, error) {
	// Fetch the first (root) value token starting at position 0
	return unmarshalFromPosition(json, 0)
}

/*
	r := JsonValue{
		valueType:	VALUE_TYPE_INVALID,
	}
	// Make the root value token our own
	r.valueType = jv.valueType
	r.valueBoolean = jv.valueBoolean
	r.valueInteger = jv.valueInteger
	r.valueFloat = jv.valueFloat
	r.valueString = jv.valueString
	if VALUE_TYPE_ARRAY == jv.valueType {
		r.valueArr = make([]*JsonValue, 0)
		r.valueArr = append(r.valueArr, jv.valueArr...)
	}
	if VALUE_TYPE_OBJECT == jv.valueType {
		r.valueMap = make(map[string]*JsonValue)
	}
	return nil
*/


func unmarshalFromPosition(json *[]rune, position int) (*JsonValue, error) {
	if nil == json { return nil, fmt.Errorf("JSON string was nil, nothing to unmarshal") }
	jsonLen := len(*json)

	// Boilerplate JsonValue to set up
	jsonValue := JsonValue{
		valueType:	VALUE_TYPE_INVALID,
		startPos:	position,
	}

	// TODO: Time for some lexing!
	for state := _LEXER_STATE_SEEK_NEXT_VALUE; _LEXER_STATE_DONE != state; {
		switch state {
			case _LEXER_STATE_SEEK_NEXT_VALUE:
				// Consume any white-space until we get to something juicy
				for ; (position < jsonLen) && isWhiteSpace((*json)[position]); position++ {
					// ref: https://stackoverflow.com/questions/18130859/how-can-i-iterate-over-a-string-by-runes-in-go
				}
				state = _LEXER_STATE_DONE
		}
	}
	jsonValue.stopPos = position - 1
	return &jsonValue, nil
}

func isWhiteSpace(r rune) bool {
	// FIXME: use regex to detect whitespace match
	return false
}

