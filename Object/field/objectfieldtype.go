package field

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

var string_ofts map[string]OFType
var oft_strings map[OFType]string

func init() {
	oft_strings = make(map[OFType]string)
	oft_strings[OFT_UNKNOWN] = "unnkown"
	oft_strings[OFT_NUMERIC] = "numeric"
	oft_strings[OFT_TEXTUAL] = "textual"
	oft_strings[OFT_DATETIME] = "datetime"
	oft_strings[OFT_BOOLEAN] = "boolean"
	oft_strings[OFT_BYTE] = "byte"
	oft_strings[OFT_SHORT] = "short"
	oft_strings[OFT_INT] = "int"
	oft_strings[OFT_LONG] = "long"
	oft_strings[OFT_FLOAT] = "float"
	oft_strings[OFT_DOUBLE] = "double"
	oft_strings[OFT_FIXED] = "fixed"
	oft_strings[OFT_STRING] = "string"
	oft_strings[OFT_CHAR] = "char"
	oft_strings[OFT_MBSTRING] = "mbstring"
	oft_strings[OFT_MBCHAR] = "mbchar"
	string_ofts = make(map[string]OFType)
	for k, v := range oft_strings {
		string_ofts[v] = k
	}
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewObjectFieldType() *ObjectFieldType {
	return &ObjectFieldType{}
}

func NewObjectFieldTypeFromString(strType string) *ObjectFieldType {
	oft := NewObjectFieldType()
	oft.SetType(oft.getOFType(strType))
	return oft
}

func NewObjectFieldTypeFromOFType(ofType OFType) *ObjectFieldType {
	oft := NewObjectFieldType()
	oft.SetType(ofType)
	return oft
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
	if s, ok := oft_strings[r.ofType]; ok { return s }
	return "unknown type"
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
	if oft, ok := string_ofts[strType]; ok { return oft }
	return OFT_UNKNOWN
}

