package objectfield

/*
An ObjectFieldType allows us to emulate loose typing for ObjectField Values
*/

type OFType int

type ObjectFieldTypeIfc interface {
	SetType(ofType OFType)
	GetType() OFType
	ToString() string
	IsValid(value *string) bool
}

type ObjectFieldType struct {
	ofType		OFType
}

const (
	OFT_UNKNOWN OFType = iota
	OFT_NUMERIC	// Any base 10 numeric form
	OFT_TEXTUAL	// Any string/text form
	OFT_DATETIME	// Any valid date and/or time form
	OFT_BOOLEAN	// Any boolean form
	OFT_BYTE	// any 8 bit form
	OFT_SHORT	// any 16 bit form
	OFT_INT		// any 32 bit form
	OFT_LONG	// any 64 but form
	OFT_FLOAT	// any floating point "real" value
	OFT_DOUBLE	// any double-precision "real" value
	OFT_FIXED	// any fixed point "real" value
	OFT_STRING	// any ASCII string
	OFT_CHAR	// any ASCII single character
	OFT_MBSTRING	// any multibyte string
	OFT_MBCHAR	// any multibyte single character
)

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewObjectFieldType() *ObjectFieldType {
	return &ObjectFieldType{}
}

func NewObjectFieldTypeFromString(strType string) *ObjectFieldType {
	oft := NewObjectFieldType()
	oft.ofType := oft.getOFType(strType)
	return &ObjectFieldType{}
}

// -------------------------------------------------------------------------------------------------
// ObjectFieldTypeIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectFieldType) SetType(ofType OFType) {
	r.ofType = ofType
}

func (r ObjectFieldType) GetType() OFType {
	return r.ofType
}

// Return a readable string for each possible type
func (r ObjectFieldType) ToString() string {
	switch (r.ofType) {
		case OFT_UNKNOWN: return "uninitialized"
		case OFT_NUMERIC: return "numeric"
		case OFT_TEXTUAL: return "textual"
		case OFT_DATETIME: return "datetime"
		case OFT_BOOLEAN: return "boolean"
		case OFT_BYTE: return "byte"
		case OFT_SHORT: return "short"
		case OFT_INT: return "int"
		case OFT_LONG: return "long"
		case OFT_FLOAT: return "float"
		case OFT_DOUBLE: return "double"
		case OFT_FIXED: return "fixed"
		case OFT_STRING: return "string"
		case OFT_CHAR: return "char"
		case OFT_MBSTRING: return "mbstring"
		case OFT_MBCHAR: return "mbchar"
		default: return "unknown type"
	}
}

// Validate the supplied value against our type
func (r ObjectFieldType) IsValid(value *string) bool {
	// TODO: Switch against r.ofType to check the value against the expected type
	switch (r.ofType) {
		case OFT_UNKNOWN: return false
	}
	return true
}

// -------------------------------------------------------------------------------------------------
// ObjectFieldType Private Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectFieldType) getOFType(strType string) OFType {
	switch (strType) {
		case "uninitialized": return OFT_UNKNOWN
		case "numeric": return OFT_NUMERIC
		case "textual": return OFT_TEXTUAL
		case "datetime": return OFT_DATETIME
		case "boolean": return OFT_BOOLEAN
		case "byte": return OFT_BYTE
		case "short": return OFT_SHORT
		case "int": return OFT_INT
		case "long": return OFT_LONG
		case "float": return OFT_FLOAT
		case "double": return OFT_DOUBLE
		case "fixed": return OFT_FIXED
		case "string": return OFT_STRING
		case "char": return OFT_CHAR
		case "mbstring": return OFT_MBSTRING
		case "mbchar": return OFT_MBCHAR
		default: return OFT_UNKNOWN
	}
}
