package objects

/*
An ObjectPrototype is an established collection of named and typed fields from which to derive new Objects
*/

type ObjectPrototypeIfc interface {
	SetFieldType(fieldName string, ofType objectFieldType)
	NewObject() *Object
}

type ObjectPrototype struct {
	fields		map[string]ObjectFieldType
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObjectPrototype() *ObjectPrototype {
	return &ObjectPrototype{
		fields:		make(map[string]ObjectFieldType),
	}
}

// Make a new one of these and initialize it from Json
func NewObjectPrototypeFromJson(json *string) *ObjectPrototype {
	objectPrototype := NewObjectPrototype()
	// TODO: Read JSON map of field name=type collection and create a new ObjectField for each
	return &objectPrototype
}

// -------------------------------------------------------------------------------------------------
// ObjectPrototypeIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectPrototype) SetFieldType(fieldName string, ofType objectFieldType) {
	oft := NewObjectFieldType()
	oft.SetType(ofType)
	r.fields[fieldName] = oft
}

func (r ObjectPrototype) NewObject() *Object {
	object := NewObject()
	for fieldName, ofType, := range r.fields {
		objectField := NewObjectField()
		objectField.Type = ofType
		object.SetField(fieldName, objectField)
	}
	return &object
}