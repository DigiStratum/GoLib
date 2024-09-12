package json

/*

Lexically parse a JSON string ([]rune, really) into a DataValue object tree

Go has built-in capability for JSON Un/Marshal, however there are some limitations that make it
less convenient for certain usages. While there is strong support for static, predefined
structures that have known properties at build time, JSON structures that are highly variable
can only be handled as generic {}interface which leaves the code to potentially a LOT of interface
assertion. This leaves the door open to substantially improved conveniences for handling JSON
structures that are highly variable.

This package is designed to produce a tree-like structure from JSON data and make it possible to
programmatically access nodes of the tree with more apparent/readable code and none of the interface
assertions. We do this by parsing the provided JSON ourselves, one character at a time, and forming
the tree structure as we progress.

Additionally, there is support for programmatically building a tree from scratch which lends to
uniformity of tooling for applications that need to either produce or consume variable JSON
structures.

Ref: https://www.rfc-editor.org/rfc/rfc7159.html

Structures under consideration, no limit on recursion depth:

NULL:
  'null'

STRING:
  '"value"'

NUMBER:
  '1'
  '3.14'
  '1E2'
  '-2.79e-4'

BOOLEAN:
  'true'|'false'

OBJECT:
  '{}'
  '{"n": {NULL}, "s": {STRING}, "v": {NUMBER}, "b": {BOOLEAN}, "o": {OBJECT}, "a": {ARRAY} }'

ARRAY:
  '[]'
  '[{NULL}, {STRING}, {NUMBER}, {BOOLEAN}, {OBJECT}, {ARRAY}]'

TODO:
 * Review ideas in https://github.com/valyala/fastjson/
 * Consider JIT lexing: only lex what's requested; leave remainder as raw JSON for later processing

*/

import (
	"fmt"
	"unicode"
	"unicode/utf8"
	"strconv"

	"GoLib/Data"
)

type jsonLexer struct {
	lexerJson		[]rune
	lexerJsonLen		int
	lexerPosition		int
	humanLine		int
	humanPosition		int
}

// -------------------------------------------------------------------------------------------------
// jsonLexerIfc
// -------------------------------------------------------------------------------------------------

// Lexically parse a DataValue out of existing JSON
func (r *jsonLexer) LexDataValue(json string) (*data.DataValue, error) {
	if ! utf8.ValidString(json) { return nil, fmt.Errorf("JSON has invalid UTF-8 multibyte sequences") }
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

func (r *jsonLexer) lexNextValue() (*data.DataValue, error) {
	// No JSON is a coding mistake!
	if nil == r.lexerJson { return nil, fmt.Errorf("JSON string was nil, nothing to lex") }

	// 1) Consume any white-space, if any, between useful bits
	r.lexConsumeWhitespace()
	if r.lexAtEOF() { return data.NewDataValue(), nil }

	// 2) Use the next character to determine data type for the value
	switch unicode.ToUpper(r.lexPeekCharacter()) {
		// String value
		case '"': return r.lexNextValueString()

		// Object value
		case '{': return r.lexNextValueObject()

		// Array value
		case '[': return r.lexNextValueArray()

		// Boolean value
		case 'T': fallthrough
		case 'F': return r.lexNextValueBool()

		// Null value
		case 'N': return r.lexNextValueNull()
	}

	// Number value
	return r.lexNextValueNumber()
}

// Is this the end?
func (r *jsonLexer) lexAtEOF() bool {
	return r.lexerPosition >= r.lexerJsonLen
}

// Peek at the character for lexer's current position without consuming it
func (r *jsonLexer) lexPeekCharacter() rune {
	return r.lexerJson[r.lexerPosition]
}

// Every character must be consumed one at a time to track position
func (r *jsonLexer) lexConsumeCharacter() rune {
	char := r.lexPeekCharacter()
//fmt.Printf("'%c'@%d\n", char, r.lexerPosition)
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
func (r *jsonLexer) lexConsumeWhitespace() bool {
	for ; (! r.lexAtEOF()) ; {
		if ! unicode.IsSpace(r.lexPeekCharacter()) { return false}
		r.lexConsumeCharacter()
	}
	return true
}

// Extract a quoted string DataValue one character at a time
func (r *jsonLexer) lexNextValueString() (*data.DataValue, error) {
	str, err := r.lexConsumeQuotedString()
	if nil != err { return nil, err }
	dataValue := data.NewDataValue().SetString(*str)
	return dataValue, nil
}

func (r *jsonLexer) lexConsumeQuotedString() (*string, error) {
	// Expect first character is double-quote string opener
	char := r.lexPeekCharacter()
	if '"' != char { return nil, r.lexError("Expected '\"' for string but got '%c' instead", char) }
	r.lexConsumeCharacter()

	// Read characters into the string value until the terminating quote comes
	escaped := false
	stringValue := make([]rune, 0)
	for ; (! r.lexAtEOF()) ; {
		char = r.lexConsumeCharacter()
		// if we're NOT escaped, then we care if this char is a '"' or an escape
		if ! escaped {
			// The closure! Return our value
			if char == '"' {
				str := string(stringValue)
				return &str, nil
			}
			// New escape sequence!
			if char == '\\' { escaped = true }
		} else {
			// Otherwise cancel the escape sequence and continue
			escaped = false
		}
		// Add the character to the string value
		stringValue = append(stringValue, char)
	}

	// If we got here then it's because we got to EOF before string closure
	return nil, r.lexError("String runs past EOF without closing")
}

// Extract an object DataValue one name-value pair at a time
func (r *jsonLexer) lexNextValueObject() (*data.DataValue, error) {
	// Expect first character is curly brace opener
	if char := r.lexConsumeCharacter(); char != '{' {
		return nil, r.lexError("Expected object start with '{' but got '%c' instead", char)
	}

	// We opened an Object value! Scaffold a DataValue to return
	dataValue := data.NewDataValue().PrepareObject()

	// Read comma-separated name:value pairs until '}' token
	for ; (!  r.lexConsumeWhitespace()) ; {

		// 1) If the next character closes the object, then we're done!
		if '}' == r.lexPeekCharacter() {
			r.lexConsumeCharacter()
			return dataValue, nil
		}

		// 2) Expect a non-empty, quoted name string value for a property name, then...
		propertyName, err := r.lexConsumeQuotedString()
		if nil != err { return nil, err }
		if 0 == len(*propertyName) { return nil, r.lexError(
			"Expected non-empty object property name, but got empty string instead",
		)}
		if r.lexConsumeWhitespace() { break }

		// 3) Expect a ':' separator between the name and value
		if ':' != r.lexPeekCharacter() {
			return nil, r.lexError("Expected ':' object property name separator, but got '%c' instead", r.lexPeekCharacter())
		}
		r.lexConsumeCharacter()
		if r.lexConsumeWhitespace() { break }

		// 4) Receive any possible valid value that follows
		propertyValue, err := r.lexNextValue() // <- BEWARE: Recursion!
		if nil != err { return nil, err }
		if ! propertyValue.IsValid() { return nil, r.lexError(
			"Expected value for object property '%s', but got something else instead", *propertyName,
		)}
		if err = dataValue.SetObjectProperty(*propertyName, propertyValue); nil != err { return nil, err }
		if r.lexConsumeWhitespace() { break }

		// 5) Expect a ',' separator between the name:value pairs or closing '}'
		char := r.lexPeekCharacter()
		if ',' == char {
			r.lexConsumeCharacter()
		} else if '}' != char { break }
	}

	// If we got here then it's because we got to EOF before object closure
	return nil, r.lexError("Object runs past EOF without closing")
}

// Extract a boolean DataValue one character at a time
func (r *jsonLexer) lexNextValueBool() (*data.DataValue, error) {
	truthy := "TRUE"
	falsey := "FALSE"
	value := ""
	for ; (! r.lexAtEOF()) && (len(value) <= 5) ; {
		value = value + string(unicode.ToUpper(r.lexConsumeCharacter()))
		if value == truthy {
			return data.NewDataValue().SetBoolean(true), nil
		} else if value == falsey {
			return data.NewDataValue().SetBoolean(false), nil
		}
	}
	return nil, r.lexError("Expected valid value for boolean, but got '%s' instead", value)
}

// Extract a null DataValue one character at a time
func (r *jsonLexer) lexNextValueNull() (*data.DataValue, error) {
	value := ""
	for ; (! r.lexAtEOF()) && (len(value) <= 4) ; {
		value = value + string(unicode.ToUpper(r.lexConsumeCharacter()))
		if "NULL" == value { return data.NewDataValue().SetNull(), nil }
	}
	return nil, r.lexError("Expected valid value for null, but got '%s' instead", value)
}

func (r *jsonLexer) lexNextValueArray() (*data.DataValue, error) {
	// Expect first character is square bracket opener
	if char := r.lexConsumeCharacter(); char != '[' {
		return nil, r.lexError("Expected array start with '[' but got '%c' instead", char)
	}
	// We opened an Array value! Scaffold a DataValue to return
	dataValue := data.NewDataValue()
	dataValue.PrepareArray()
	expectValue := false

	// Read comma-separated values until ']' token
	for ; ! r.lexConsumeWhitespace() ; {

		// If we're not expecting an element to follow, then it's OK to close
		if (! expectValue) {
			// If the next character closes the array, then we're done!
			if ']' == r.lexPeekCharacter() {
				r.lexConsumeCharacter()
				return dataValue, nil
			}
		}

		// Receive any possible valid value that follows
		value, err := r.lexNextValue() // <- BEWARE: Recursion!
		if nil != err { return nil, err }
		if ! value.IsValid() {
			return nil, r.lexError("Expected array entry value but got something else instead")
		}
		if err = dataValue.AppendArrayValue(value); nil != err { return nil, err }

		// After the value may be whitespace
		if r.lexConsumeWhitespace() { break }

		// Expect a ',' separator if another array element is coming at us...
		if ',' == r.lexPeekCharacter() {
			r.lexConsumeCharacter()
			expectValue = true
			// After the value may be whitespace
			if r.lexConsumeWhitespace() { break }
		} else { expectValue = false }
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
func (r *jsonLexer) lexNextValueNumber() (*data.DataValue, error) {
	isFloat := false
	valueStr := ""
	var err error

	// Consume numeric value [-]*([0]|[1-9][0-9]*)(\.[0-9]+)([eE](+-)[1-9][0-9]+)*

	// 1) Optional: negative sign
	valueStr, err = r.lexConsumeAppendCharacter(valueStr, '-');

	// 2) Required: Integer digits (Note: leading 0's not allowed, but we needn't enforce here)
	if valueStr, err = r.lexExpectConsumeAppendDigits(valueStr); nil != err { return nil, err }

	// 3) Optional: Pretend we expect a decimal point for a float; if we got one...
	if valueStr, err = r.lexConsumeAppendCharacter(valueStr, '.'); nil == err {
		// ... then require digits to follow!
		isFloat = true
		if valueStr, err = r.lexExpectConsumeAppendDigits(valueStr); nil != err { return nil, err }
	}

	// 4) Optional: Pretend we expect an e|E for an exponent specifier; if we got one...
	if valueStr, err = r.lexConsumeAppendCharacter(valueStr, 'e', 'E'); nil == err {
		// ... then allow an optional sign to follow...
		valueStr, _ = r.lexConsumeAppendCharacter(valueStr, '+', '-')
		// ... then require digits to follow!
		if valueStr, err = r.lexExpectConsumeAppendDigits(valueStr); nil != err { return nil, err }
	}

	// 5) Return the valueStr as an int64 or float64!
	if isFloat {
		// Floats
		valueFloat, err := strconv.ParseFloat(valueStr, 64)
		if nil != err {
			return nil, r.lexError("Error converting number '%s' to float64: %s", valueStr, err.Error())
		}
		return data.NewDataValue().SetFloat(valueFloat), nil
	}
	// Integers
	valueInteger, err := strconv.ParseInt(valueStr, 10, 64)
	if nil != err {
		return nil, r.lexError("Error converting number '%s' to int64: %s", valueStr, err.Error())
	}
	return data.NewDataValue().SetInteger(valueInteger), nil
}

func (r *jsonLexer) lexConsumeAppendCharacter(base string, acceptedChars ...rune) (string, error) {
	if ! r.lexAtEOF() {
		char := r.lexPeekCharacter()
		for _, acceptedChar := range acceptedChars {
			if char == acceptedChar {
				r.lexConsumeCharacter()
				return base + string(char), nil
			}
		}
	}
	ac := ""
	for _, acceptedChar := range acceptedChars { ac = ac + string(acceptedChar) }
	return base, r.lexError("Expected acceptable character [%s] but got something else or EOF", ac)
}

func (r *jsonLexer) lexExpectConsumeAppendDigits(base string) (string, error) {
	if r.lexAtEOF() { return base, r.lexError("Expected digits but got EOF") }
	value := ""
	for char := r.lexPeekCharacter() ; ('0' <= char) && ('9' >= char) ; char = r.lexPeekCharacter() {
		value = value + string(r.lexConsumeCharacter())
		if r.lexAtEOF() { break }
	}
	return base + value, nil
}

func (r *jsonLexer) lexError(msg string, args ...interface{}) error {
	m := fmt.Sprintf(msg, args...)
	return fmt.Errorf("%s at line %d, pos %d", m, r.humanLine, r.humanPosition)
}

