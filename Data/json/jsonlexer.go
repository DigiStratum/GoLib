package json

/*

Lexically parse a JSON string ([]rune, really) into a JsonValue object tree

TODO:
 * If we make our own tokenizer/parser, then we can unmarshall into a custom tree structure
   that separates JSON types from values to make it easier to traverse, access, and assert
 * Add support for de|referencing; make references an embeddable json string (like mustache), use
   configurable start/stop delimiters with default; for whole-string references like "{{sel.ect.or}}"
   convert the value type to that of the selected reference, null if it doesn't exist. For partial
   references like "See also: {{sel.ect.or}}", convert the value type to string and perform string
   replacement, empty string if it doesn't exist. Introduce "RJSON" envelope to encode json metadata
   to describe the JSON encoding within, versioning, etc. to help with future-proofing, versioning,
   etc.
 * Consider JIT lexing: only lex what's requested; leave remainder as raw JSON for later processing

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

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type JsonLexerIfc interface {
	LexJsonValue(json *[]rune) (*JsonValue, error)
}

type JsonLexer struct {
	lexerJson		*[]rune
	lexerJsonLen		int
	lexerPosition		int
	lexerErr		error

	// Human-readable position within the JSON for error messaging
	humanLine		int
	humanPosition		int
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewJsonLexer() *JsonLexer {
	r := JsonLexer{ }
	return &r
}

// -------------------------------------------------------------------------------------------------
// JsonLexerIfc
// -------------------------------------------------------------------------------------------------

// Lexically parse a JsonValue out of existing JSON
func (r *JsonLexer) LexJsonValue(json *[]rune) (*JsonValue, error) {
	// Require UTF-8 
	if ! utf8.ValidString(string(r.json)) {
		return nil, fmt.Errorf("JSON has invalid UTF-8 multibyte sequences")
	}

	r.lexerJson = json
	r.lexerPosition = 0
	r.lexerJsonLen = len(*json)
	r.humanLine = r.humanPostion = 1
	return r.lexNextValue()
}

// -------------------------------------------------------------------------------------------------
// Private implementation
// -------------------------------------------------------------------------------------------------

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

func (r *JsonLexer) lexNextValue() (*JsonValue, error) {
	// No JSON is a coding mistake!
	if nil == r.lexerJson {
		return nil, fmt.Errorf("JSON string was nil, nothing to unmarshal")
	}

	// Scaffold a JsonValue to return
	jsonValue := NewJsonValue()

	// 1) Consume any white-space, if any, between useful bits
	r.lexConsumeWhitespace()
	if r.lexAtEOF() { return jsonValue, nil }

	// 2) Use the next character to determine data type for the value
	switch r.(* lexerJson)[r.lexerPosition] {
		// String value
		case '"': return lexNextValueString()

		// Object value
		case '{': return lexNextValueObject()

		// Array value
		case '[':

		// Boolean value
		case 't': fallthrough
		case 'T': fallthrough
		case 'f': fallthrough
		case 'F':
			// TODO: Consume true|false value

		// Null value
		case 'n':
			// TODO: Consume null value

		// Number value
		default:
			// TODO: Consume numeric value [-]*[0-9+](\.[0-9+])*
	}
	jsonValue.stopPos = position - 1
	return &jsonValue, nil
}

// Is this the end?
func (r *JsonLexer) lexAtEOF() bool {
	return r.lexerPosition >= (r.lexerJsonLen - 1)
}

// Peek at the character for lexer's current position without consuming it
func (r *JsonLexer) lexPeekCharacter() rune {
	return (*r).(*lexerJson)[r.lexerPosition]
}

// Every character must be consumed one at a time to track position
func (r *JsonLexer) lexConsumeCharacter() rune {
	char := r.lexPeekCharacter(()
	r.lexerPosition++
	if '\n' == char {
		r.humanLine++
		r.humanPosition = 1
	}
	r.humanPosition++
	return char
}

// Consume sequential white space characters to get to the next useful thing
func (r *JsonLexer) lexConsumeWhitespace() {
	for ; ! lexAtEOF(); char := r.lexPeekCharacter() {
		// ref: https://www.geeksforgeeks.org/check-if-the-rune-is-a-space-character-or-not-in-golang/
		if ! unicode.IsSpace(char) { break }
		r.lexConsumeCharacter()
	}
}

// Extract a quoted string JsonValue one character at a time
func lexNextValueString(position int, json*[]rune) (*JsonValue, error) {
	// Expect first character is double-quote string opener
	if char := r.lexConsumeCharacter(); char != '"' {
		return nil, fmt.Errorf(
			"Expected opening \" to start string at line %d, pos %d, but got %c instead",
			r.humanLine,
			r.humanPosition,
			char,
		)
	}

	// We opened a string value! Scaffold a JsonValue to return
	jsonValue := NewJsonValue()
	value := make([]rune, 0)

	// Read characters into the string value until the terminating quote comes
	escaped := false
	for ; ! lexAtEOF(); char = r.lexConsumeCharacter() {
		// if we're NOT escaped, then we care if this char is a '"' or an escape
		if ! escaped {
			// The closure! Return our value
			if char == '"' {
				jsonValue.SetString(value)
				return jsonValue, nil
			}
			// New escape sequence!
			if char == '\\' { escaped = true }
		} else {
			// Otherwise cancel the escape sequence and continue
			escaped = false
		}
		// Add the character to the string value
		value = append(value, char)
	}

	// If we got here then it's because we got to EOF before string closure
	return nil, fmt.Errorf(
		"String runs past EOF without closing quote at line %d, pos %d",
		r.humanLine,
		r.humanPosition,
	)
}

// Extract an object JsonValue one name-value pair at a time
func lexNextValueObject(position int, json*[]rune) (*JsonValue, error) {

	// TODO: Read comma-separated name:value pairs until '}' token
	// i. Get a string for object name
}

