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
	"unicode/utf8"
)

type JsonLexerIfc interface {
	LexJsonValue(json *[]rune) (*JsonValue, error)
}

type JsonLexer struct {
	lexerJson		*[]rune
	lexerPosition		int
	lexerErr		error
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

// Lexically parse one of these out of existing JSON
func (r *JsonLexer) LexJsonValue(json *[]rune) (*JsonValue, error) {
	r.lexerJson = json
	r.lexerPosition = 0
	return r.lexNextValue()
}

const (
	_LEXER_STATE_SEEK_NEXT_VALUE int = iota
	_LEXER_STATE_DONE
)


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

func (r *JsonLexer) lexNextValue(json *[]rune) (*JsonValue, error) {
	// No JSON is a coding mistake!
	if nil == json {
		return nil, fmt.Errorf("JSON string was nil, nothing to unmarshal")
	}

	// Require UTF-8 
	if ! utf8.ValidString(string(*json)) {
		return nil, fmt.Errorf("JSON has invalid UTF-8 multibyte sequences")
	}

	// Scaffold a JsonValue to return
	jsonValue := JsonValue{
		valueType:	VALUE_TYPE_INVALID,
		startPos:	position,
	}

	// Empty JSON is invalid
	jsonLen := len(*json)
	if 0 == jsonLen { return &jsonValue, nil }

	// Time for some lexing!
	for state := _LEXER_STATE_SEEK_NEXT_VALUE; _LEXER_STATE_DONE != state; {
		switch state {
			// Look for a JSON value
			case _LEXER_STATE_SEEK_NEXT_VALUE:
				// 1) Consume any white-space until we get to something juicy
				position = lexSkipJsonWhitespace(position, json)
				jsonValue.startPos = position

				// 2) TODO: Use the next character to determine data type for the value
				switch json[position] {
					case '"':
						// String value
						value, err := lexConsumeValueString(position, json)
						if nil != err { return nil, err }
						jsonValue.valueType = VALUE_TYPE_STRING
						jsonValue.

					case '{':
						// Object value

					case '[':
						// Array value

					case 't':
						fallthrough
					case 'T':
						fallthrough
					case 'f':
						fallthrough
					case 'F':
						// Boolean value

					case 'n':
						// Null value
				}
				state = _LEXER_STATE_DONE
		}
	}
	jsonValue.stopPos = position - 1
	return &jsonValue, nil
}

func lexConsumeValueString(position int, json*[]rune) (*string, error) {
	value := ""

	return value, nil
}

func lexSkipJsonWhitespace(position int, json *[]rune) int {
	for ; (position < jsonLen) && isWhiteSpace((*json)[position]); position++ {
		// TODO: Validate json as UTF-8 with unicode/utf8.ValidRune() 
		// ref: https://stackoverflow.com/questions/18130859/how-can-i-iterate-over-a-string-by-runes-in-go
	}
}

func isWhiteSpace(r rune) bool {
	// FIXME: use regex to detect whitespace match
	// ref: https://go.dev/blog/strings
	// ref: https://pkg.go.dev/unicode/utf8
	// ref: https://www.reddit.com/r/learnprogramming/comments/sqa4p5/why_does_multibyte_utf8_characters_start_with_a/
	return false
}

