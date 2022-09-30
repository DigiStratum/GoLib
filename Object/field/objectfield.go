package field

/*
An ObjectField encapsulates the properties of a single field for a given Object
*/

type ObjectFieldIfc interface {
	GetName() string
	GetType() *ObjectFieldType
	SetType(ofType ObjectFieldTypeIfc)
	GetValue() *string
	SetValue(value *string)
}

// Association of Type and value for a single Object Field
type ObjectField struct {
	name	string
	ofType	ObjectFieldTypeIfc
	value	*string			// Significance varies with Type
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewObjectField(name string) *ObjectField {
	return &ObjectField{
		name:	name,
		ofType:	NewObjectFieldType(),
	}
}

// -------------------------------------------------------------------------------------------------
// ObjectFieldIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r ObjectField) GetName() string {
	return r.name
}

func (r ObjectField) GetType() *ObjectFieldType {
	return NewObjectFieldTypeFromOFType(r.ofType.GetType())
}

func (r *ObjectField) SetType(ofType ObjectFieldTypeIfc) {
	r.ofType = ofType
}

func (r ObjectField) GetValue() *string {
	return r.value
}

func (r *ObjectField) SetValue(value *string) {
	r.value = value
}

