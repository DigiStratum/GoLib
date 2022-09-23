package prototype

/*
An ObjectPrototype is an established collection of named and typed fields from which to derive new Objects
*/

import (
	obj "github.com/DigiStratum/GoLib/Object"
	objf "github.com/DigiStratum/GoLib/Object/field"
)

type ObjectPrototypeIfc interface {
	SetFieldType(fieldName string, ofType objf.ObjectFieldType)
	NewObject() *obj.Object		// TODO: change this to interface
}

type ObjectPrototype struct {
	fields		map[string]objf.ObjectFieldType
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObjectPrototype() *ObjectPrototype {
	return &ObjectPrototype{
		fields:		make(map[string]objf.ObjectFieldType),
	}
}

// Make a new one of these and initialize it from Json
func NewObjectPrototypeFromJson(json *string) *ObjectPrototype {
	objectPrototype := NewObjectPrototype()
	// TODO: Read JSON map of field name=type collection and create a new ObjectField for each
	return objectPrototype
}

// -------------------------------------------------------------------------------------------------
// ObjectPrototypeIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectPrototype) SetFieldType(fieldName string, ofType objf.ObjectFieldType) {
	oft := objf.NewObjectFieldType()
	oft.SetType(ofType)
	r.fields[fieldName] = oft
}

func (r ObjectPrototype) NewObject() *obj.Object {
	object := obj.NewObject()
	for fieldName, ofType := range r.fields {
		objectField := objf.NewObjectField()
		objectField.Type = ofType
		object.SetField(fieldName, objectField)
	}
	return &object
}
