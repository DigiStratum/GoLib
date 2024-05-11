package json

/*

Represent a JSON structure as an object tree with JavaScript-like selectors and other conveniences.

Structuress under consideration:

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

TODO:
 * If we make our own tokenizer/parser, then we can unmarshall into a custom tree structure
   that separates JSON types from values to make it easier to traverse, access, and assert

*/

import(
	"fmt"
	"strings"
	gojson "encoding/json"
)

type JsonTreeIfc interface {

}

type JsonTree struct {
	data		map[string]interface{}
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewJsonTree(jsonString *string) *JsonTree {
	if nil == jsonString { return nil }
	var data map[string]interface{}
	if err := gojson.Unmarshal([]byte(*jsonString), &data); err != nil { return nil }
	r := JsonTree{
		data:	data,
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// JsonTreeIfc
// -------------------------------------------------------------------------------------------------

func (r *JsonTree) GetString(selector string) (*string, error) {
	node, err := r.getNode(selector)
	if nil != err { return nil, err }
	if nil == node { return nil, nil }
	if v, ok := node.(string); ok { return &v, nil }
	return nil, fmt.Errorf("Failed type assertion to string");
}

// -------------------------------------------------------------------------------------------------
// JsonTree
// -------------------------------------------------------------------------------------------------

func (r *JsonTree) getNode(selector string) (interface{}, error) {

	// We start walking nodes from the base
	node := r.data

	// Tokenize the JavaTree selector and validate the base before recursively walking the tree
	// "$"
	// "$.prop1"
	// "$.prop1.prop2"
	// "$[index]"
	// "$[index].prop1"
	// "$.prop1[index1][index2].prop2"
	tokens := strings.Split(selector, ".")

	// A blank selector is unacceptable
	if len(tokens) == 0 { return nil, fmt.Errorf("Invalid selector (blank)") }

	// The first token must begin with '$' for the base
	if tokens[0][0] != '$' { return nil, fmt.Errorf("Invalid selector (base)") }

	// Consume the the base '$'
	basetoken := tokens[0][1:len(tokens[0])-1]
	newselector := strings.Join(tokens[1:], ",")
	if len(basetoken) > 0 { newselector = basetoken + "." + newselector; }

	return r.walkNode(newselector, node)
}

// Recursive function to walk the selector tree until the end or error
func (r *JsonTree) walkNode(selector string, node interface{}) (interface{}, error) {
	// Tokenize the JavaTree selector and walk the tree of indexes and properties until we reach the final or fail
	// ""
	// "prop1"
	// "prop1.prop2"
	// "prop1[index1].prop2"
	// "prop1[index1][index2].prop2"
	// "[index]"
	// "[index].prop1"

	// If the selector is empty, then the node is the final
	if selector == "" { return node, nil }

	return nil, fmt.Errorf("Not Implemented")
}

