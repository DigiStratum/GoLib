package objects

/*
An ObjectPrototype is an established collection of named and typed fields from which to derive new Objects
*/

type ObjectPrototypeIfc interface {
}

type ObjectPrototype struct {
	fields		map[string]ObjectFieldType
}

// Factory Functions
func NewObjectPrototype() *ObjectPrototype {
	return &ObjectPrototype{
		fields:		make(map[string]ObjectFieldType),
	}
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