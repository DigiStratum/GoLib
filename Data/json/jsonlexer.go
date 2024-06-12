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
 * Make DependencyInjectable to accept a logger for errors, debug, etc (?)

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
	"strconv"
)

type JsonLexerIfc interface {
	LexJsonValue(json string) (*JsonValue, error)
}

type JsonLexer struct {
	lexerJson		[]rune
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
func (r *JsonLexer) LexJsonValue(json string) (*JsonValue, error) {
	// Require UTF-8
	if ! utf8.ValidString(json) {
		return nil, fmt.Errorf("JSON has invalid UTF-8 multibyte sequences")
	}

	r.lexerJson = []rune(json)
	r.lexerPosition = 0
	r.lexerJsonLen = len(r.lexerJson)
	r.humanLine = 1
	r.humanPosition = 1
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
		return nil, fmt.Errorf("JSON string was nil, nothing to lex")
	}

	// 1) Consume any white-space, if any, between useful bits
	r.lexConsumeWhitespace()
	if r.lexAtEOF() { return NewJsonValue(), nil }

	// 2) Use the next character to determine data type for the value
	switch r.lexPeekCharacter() {
		// String value
		case '"': return r.lexNextValueString()

		// Object value
		case '{': return r.lexNextValueObject()

		// Array value
		case '[': return r.lexNextValueArray()

		// Boolean value
		case 't': fallthrough
		case 'T': fallthrough
		case 'f': fallthrough
		case 'F':
			return r.lexNextValueBool()

		// Null value
		case 'n':
			return r.lexNextValueNull()

		// Number value
		default: return r.lexNextValueNumber()
	}
	return NewJsonValue(), nil
}

// Is this the end?
func (r *JsonLexer) lexAtEOF() bool {
	return r.lexerPosition >= r.lexerJsonLen
}

// Peek at the character for lexer's current position without consuming it
func (r *JsonLexer) lexPeekCharacter() rune {
	char := r.lexerJson[r.lexerPosition]
	return char
}

// Every character must be consumed one at a time to track position
func (r *JsonLexer) lexConsumeCharacter() rune {
	char := r.lexPeekCharacter()
//fmt.Printf("[%c] %d", char, r.lexerPosition)
	r.lexerPosition++
	if '\n' == char {
		r.humanLine++
		r.humanPosition = 1
	}
	r.humanPosition++
	return char
}

// Consume sequential white space characters to get to the next useful thing
// Returns true if EOF reached, else false
func (r *JsonLexer) lexConsumeWhitespace() bool {
	for ; ; {
		if r.lexAtEOF() { return true }
		// ref: https://www.geeksforgeeks.org/check-if-the-rune-is-a-space-character-or-not-in-golang/
		if ! unicode.IsSpace(r.lexPeekCharacter()) { return false}
		r.lexConsumeCharacter()
	}
}

// Extract a quoted string JsonValue one character at a time
func (r *JsonLexer) lexNextValueString() (*JsonValue, error) {
	// Expect first character is double-quote string opener
	char := r.lexPeekCharacter()
	if '"' != char {
		return nil, r.lexError(fmt.Sprintf(
			"Expected string start with '\"' but got '%c' instead", char,
		))
	}
	r.lexConsumeCharacter()

	// We opened a string value! Scaffold a JsonValue to return
	jsonValue := NewJsonValue()
	stringValue := make([]rune, 0)

	// Read characters into the string value until the terminating quote comes
	escaped := false
	for ; ; {
		char = r.lexConsumeCharacter()
		// if we're NOT escaped, then we care if this char is a '"' or an escape
		if ! escaped {
			// The closure! Return our value
			if char == '"' {
				jsonValue.SetString(string(stringValue))
				return jsonValue, nil
			}
			// New escape sequence!
			if char == '\\' { escaped = true }
		} else {
			// Otherwise cancel the escape sequence and continue
			escaped = false
		}
		// Add the character to the string value
		stringValue = append(stringValue, char)

		if r.lexAtEOF() { break }
	}

	// If we got here then it's because we got to EOF before string closure
	return nil, r.lexError("String runs past EOF without closing")
}

// Extract an object JsonValue one name-value pair at a time
func (r *JsonLexer) lexNextValueObject() (*JsonValue, error) {
	// Expect first character is curly brace opener
	if char := r.lexConsumeCharacter(); char != '{' {
		return nil, r.lexError(fmt.Sprintf(
			"Expected object start with '{' but got '%c' instead", char,
		))
	}

	// We opened an Object value! Scaffold a JsonValue to return
	jsonValue := NewJsonValue()
	jsonValue.PrepareObject()

	// Read comma-separated name:value pairs until '}' token
	for ; ; {
		if r.lexConsumeWhitespace() { break }

		// If the next character closes the object, then we're done!
		char := r.lexPeekCharacter()
		if '}' == char {
			r.lexConsumeCharacter()
			return jsonValue, nil
		}

		// Expect a non-empty, quoted name string value for a property name, then...
		if '"' != char {
			return nil, r.lexError("Expected quoted object property name, but got something else instead")
		}
		nameValue, err := r.lexNextValueString()
		if nil != err { return nil, err }
		propertyName := nameValue.GetString()
		if 0 == len(propertyName) {
			return nil, r.lexError("Expected non-empty object property name, but got empty string instead")
		}

		// After the name may be whitespace
		if r.lexConsumeWhitespace() { break }

		// Expect a ':' separator between the name and value
		if ':' != r.lexPeekCharacter() {
			return nil, r.lexError(fmt.Sprintf(
				"Expected ':' object property name separator, but got '%c' instead", char,
			))
		}

		// Consume the separator; After may be whitespace
		r.lexConsumeCharacter()
		if r.lexConsumeWhitespace() { break }

		// Receive any possible valid value that follows
		propertyValue, err := r.lexNextValue() // <- BEWARE: Recursion!
		if nil != err { return nil, err }
		if ! propertyValue.IsValid() {
			return nil, r.lexError(fmt.Sprintf(
				"Expected valid value for object property '%s', but got something else instead", propertyName,
			))
		}
		if err = jsonValue.SetObjectProperty(propertyName, propertyValue); nil != err { return nil, err }

		// After the value may be whitespace
		if r.lexConsumeWhitespace() { break }

		// Expect a ',' separator between the name:value pairs or closing '}'
		char = r.lexPeekCharacter()
		if ',' == char {
			r.lexConsumeCharacter()
		} else if '}' != char { break }

		if r.lexAtEOF() { break }
	}

	// If we got here then it's because we got to EOF before object closure
	return nil, r.lexError("Object runs past EOF without closing")
}

// Extract a boolean JsonValue one character at a time
func (r *JsonLexer) lexNextValueBool() (*JsonValue, error) {
	truthy := "TRUE"
	falsey := "FALSE"
	value := ""
	rawValue := ""
	for ; ; {
		char := r.lexConsumeCharacter()
		rawValue = rawValue + string(char)
		value = value + string(unicode.ToUpper(char))
		if value == truthy {
			jsonValue := NewJsonValue()
			jsonValue.SetBoolean(true)
			return jsonValue, nil
		} else if value == falsey {
			jsonValue := NewJsonValue()
			jsonValue.SetBoolean(false)
			return jsonValue, nil
		} else if (len(value) > 5) { break }
		if r.lexAtEOF() { break }
	}
	return nil, r.lexError(fmt.Sprintf(
		"Expected valid value for boolean, but got '%s' instead", rawValue,
	))
}

// Extract a null JsonValue one character at a time
func (r *JsonLexer) lexNextValueNull() (*JsonValue, error) {
	null := "NULL"
	value := ""
	rawValue := ""
	for ; ; {
		char := r.lexConsumeCharacter()
		rawValue = rawValue + string(char)
		value = value + string(unicode.ToUpper(char))
		if value == null {
			jsonValue := NewJsonValue()
			jsonValue.SetNull()
			return jsonValue, nil
		} else if (len(value) > 4) { break }
		if r.lexAtEOF() { break }
	}
	return nil, r.lexError(fmt.Sprintf(
		"Expected valid value for null, but got '%s' instead", rawValue,
	))
}

func (r *JsonLexer) lexNextValueArray() (*JsonValue, error) {
	// Expect first character is square bracket opener
	if char := r.lexConsumeCharacter(); char != '[' {
		return nil, r.lexError(fmt.Sprintf(
			"Expected array start with '[' but got '%c' instead", char,
		))
	}

	// We opened an Array value! Scaffold a JsonValue to return
	jsonValue := NewJsonValue()
	jsonValue.PrepareArray()

	// Read comma-separated values until ']' token
	for ; ; {
		if r.lexConsumeWhitespace() { break }

		// If the next character closes the array, then we're done!
		char := r.lexPeekCharacter()
		if ']' == char {
			r.lexConsumeCharacter()
			return jsonValue, nil
		}

		// Receive any possible valid value that follows
		value, err := r.lexNextValue() // <- BEWARE: Recursion!
		if nil != err { return nil, err }
		if ! value.IsValid() {
			return nil, r.lexError(fmt.Sprintf(
				"Expected valid value for array entry, but got something else instead",
			))
		}
		if err = jsonValue.AppendArrayValue(value); nil != err { return nil, err }

		// After the value may be whitespace
		if r.lexConsumeWhitespace() { break }

		// Expect a ',' separator between values or closing ']'
		char = r.lexPeekCharacter()
		if ',' == char {
			r.lexConsumeCharacter()
		} else if ']' != char { break }

		if r.lexAtEOF() { break }
	}

	// If we got here then it's because we got to EOF before array closure
	return nil, r.lexError("Array runs past EOF without closing")
}

/*
Note:
 * int64 range is -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807; that's 20 chars, less the
   commas, max readalbe

 * float64 numeric grammar is more elaborate; Ref: https://www.rfc-editor.org/rfc/rfc7159.html#section-6

 	number = [ minus ] int [ frac ] [ exp ]
	decimal-point = %x2E       ; .
	digit1-9 = %x31-39         ; 1-9
	e = %x65 / %x45            ; e E
	exp = e [ minus / plus ] 1*DIGIT
	frac = decimal-point 1*DIGIT
	int = zero / ( digit1-9 *DIGIT )
	minus = %x2D               ; -
	plus = %x2B                ; +
	zero = %x30                ; 0

* Some string representations of float appear to support up to 26 characters in various software
   implementations. From RFC7159 above:

	Note that when such software is used, numbers that are integers and are in the range
	[-(2**53)+1, (2**53)-1] are interoperable in the sense that implementations will agree
	exactly on their numeric values.

 * Therefore, TODO: Implement a configurable precision/range limit that defaults to the above
   interoperability standards, with optional override, and throw error as we parse and discover
   values that are out of range.
*/
func (r *JsonLexer) lexNextValueNumber() (*JsonValue, error) {
	isFloat := false
	valueStr := ""

	// TODO: Consume numeric value [-]*([0]|[1-9][0-9]*)(\.[0-9]+)([eE](+-)[1-9][0-9]+)*

	// 1) Optional negative sign
	if r.lexAtEOF() { return nil, r.lexError("Number value runs past EOF without closing [1]") }
	if '-' == r.lexPeekCharacter() {
		r.lexConsumeCharacter()
		valueStr = "-"
	}

	// 2) Integer (leading 0's not allowed, but we needn't enforce here)
	if r.lexAtEOF() { return nil, r.lexError("Number value runs past EOF without closing [2]") }
	char := r.lexPeekCharacter()
	// Get the integer digits!
	digits, err := r.lexNextDigits()
	if (nil != err) || (len(*digits) == 0) { return nil, r.lexError("Number value runs past EOF without closing [3]") }
	valueStr = valueStr + *digits

	// 3) Optional decimal (but only if we have an integer part, doh!)
	if ! r.lexAtEOF() {
		if '.' == char {
			r.lexConsumeCharacter()
			valueStr = valueStr + "."
			isFloat = true
			// Get the decimal digits!
			digits, err := r.lexNextDigits()
			if nil != err { return nil, err }
			valueStr = valueStr + *digits
		}

		// 4) Optional exponent
		if 'E' == unicode.ToUpper(char) {
			r.lexConsumeCharacter()
			valueStr = valueStr + "e"
			if r.lexAtEOF() { return nil, r.lexError("Number value runs past EOF without closing [4]") }
			char = r.lexPeekCharacter()
			if '-' == char {
				valueStr = valueStr + string('-')
				r.lexConsumeCharacter()
				if r.lexAtEOF() { return nil, r.lexError("Number value runs past EOF without closing [5]") }
				char = r.lexPeekCharacter()
			} else if '+' == char {
				r.lexConsumeCharacter()
				if r.lexAtEOF() { return nil, r.lexError("Number value runs past EOF without closing [6]") }
				char = r.lexPeekCharacter()
			}
			// Get the exponent digits!
			digits, err := r.lexNextDigits()
			if nil != err { return nil, err }
			valueStr = valueStr + *digits
		}
	}

	// Return the valueStr as an int64 or float64!
	jsonValue := NewJsonValue()
	if isFloat {
		valueFloat, err := strconv.ParseFloat(valueStr, 64)
		if nil != err {
			return nil, r.lexError(fmt.Sprintf("Error converting parsed number '%s' to float64: %s", valueStr, err.Error()))
		}
		jsonValue.SetFloat(valueFloat)
	} else {
		valueInteger, err := strconv.ParseInt(valueStr, 10, 64)
		if nil != err {
			return nil, r.lexError(fmt.Sprintf("Error converting parsed number '%s' to int64: %s", valueStr, err.Error()))
		}
		jsonValue.SetInteger(valueInteger)
	}

	return jsonValue, nil
}

func (r *JsonLexer) lexNextDigits() (*string, error) {
	// Gather digits!
	valueStr := ""
	if r.lexAtEOF() { return nil, r.lexError("Number value runs past EOF without closing [digits]") }
	char := r.lexPeekCharacter()
	for ; ('0' <= char) && ('9' >= char) ; {
		valueStr = valueStr + string(r.lexConsumeCharacter())
		if r.lexAtEOF() { break }
		char = r.lexPeekCharacter()
	}
//fmt.Printf("Got digits [%s] at line %d, pos %d\n", valueStr, r.humanLine, r.humanPosition)
	return &valueStr, nil
}

func (r *JsonLexer) lexError(msg string) error {
	return fmt.Errorf( "%s at line %d, pos %d", msg, r.humanLine, r.humanPosition)
}

