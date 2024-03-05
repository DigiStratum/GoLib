package json

import(
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


