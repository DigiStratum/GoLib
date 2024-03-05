package json

/*

Represent a JSON structure as an object tree with JavaScript-like selectors and other conveniences.

Structuress under consideration:

NIL:
  ''
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

TODO:
 * If we make our own tokenizer/parser, then we can unmarshall into a custom tree structure
   that separates JSON types from values to make it easier to traverse, access, and assert

*/

import(
	"fmt"
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

func (r *JsonTree) GetString(jtpath string) (*string, error) {
	node, err := r.getNode(jtpath)
	if nil != err { return nil, err }
	if nil == node { return nil, nil }
	if v, ok := node.(string); ok { return &v, nil }
	return nil, fmt.Errorf("Failed type assertiong to string");
}

// -------------------------------------------------------------------------------------------------
// JsonTree
// -------------------------------------------------------------------------------------------------

func (r *JsonTree) getNode(jtpath string) (interface{}, error) {
	// TODO: Tokenize the JavaTree Path and traverse the tree until the node is found or error
	return nil, fmt.Errorf("Not Implemented Yet!")
}

