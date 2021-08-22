package objects

/*
An ObjectField encapsulates the properties of a single field for a given Object
*/

// Association of Type and value for a single Object Field
type ObjectField struct {
	Type		ObjectFieldType
	Value		*string			// Significance varies with Type
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

func NewObjectField() *ObjectField {
	return &ObjectField{
		Type:	NewObjectFieldType(),
	}
}
