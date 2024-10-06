package data

type DataType int

const (
	DATA_TYPE_INVALID DataType = iota
	DATA_TYPE_NULL
	DATA_TYPE_BOOLEAN
	DATA_TYPE_INTEGER
	DATA_TYPE_FLOAT
	DATA_TYPE_STRING
	DATA_TYPE_OBJECT
	DATA_TYPE_ARRAY
)

func (r DataType) ToString() string {
	switch r {
		case DATA_TYPE_INVALID: return "invalid"
		case DATA_TYPE_NULL: return "null"
		case DATA_TYPE_BOOLEAN: return "boolean"
		case DATA_TYPE_INTEGER: return "integer"
		case DATA_TYPE_FLOAT: return "float"
		case DATA_TYPE_STRING: return "string"
		case DATA_TYPE_OBJECT: return "object"
		case DATA_TYPE_ARRAY: return "array"
	}
	return ""
}

