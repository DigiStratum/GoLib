package mysql

/*
DB Results for MySQL - structure and support for reading MySQL query result data into structure
properties instead of a map.

This is important because an array of struct is more efficient both in terms of memory utilization
and data processing performance vs []map[string]string. Also a map[string]string for each row result
makes it more difficult to add supporting methods, etc.

*/

import (
	"fmt"
	"errors"
	"reflect"
)

type Result interface{}
type PropertyPointers []interface{}
type ResultSet []Result

// Base struct for DB query result row objects
type ResultFactory struct {
	prototype	Result
	prototypeType	reflect.Type
	numFields	int
}

// Make a new one of these!
func NewResultFactory(prototype Result) *ResultFactory {
	rf := ResultFactory{
		prototype:	prototype,
		prototypeType:	reflect.TypeOf(prototype),
		numFields:	reflect.TypeOf(prototype).NumField(),
	}
	return &rf
}

// ref: https://stackoverflow.com/questions/29184933/golang-reflect-get-pointer-to-a-struct-field-value
func (rf *ResultFactory) getStructPropertyPointers(u interface{}) PropertyPointers {
	val := reflect.ValueOf(u).Elem()
	v := make(PropertyPointers, (*rf).numFields)
	for i := 0; i < (*rf).numFields; i++ {
		valueField := val.Field(i)
		v[i] = valueField.Addr().Interface()
	}
	return v
}

type blerf struct {
	Bloink	string
	Bleep	string
	Bloop	string
}

// Make a new Result!
// Create a set of pointers to capture query result column values for each row processed with Scan()
// FIXME: https://stackoverflow.com/questions/40512323/golang-cast-interface-back-to-its-original-type
// ref: https://stackoverflow.com/questions/11127723/dynamically-create-variables-of-certain-type-based-on-string-in-go
func (rf *ResultFactory) MakeNewResult() (Result, PropertyPointers, error) {

	// Make a new result object and reflect on it
	//newResult := reflect.New(reflect.TypeOf((*rf).prototype).Elem())
	//newResult := reflect.New(reflect.TypeOf((*rf).prototype))
	// https://groups.google.com/g/golang-dev/c/XWfzNWe4Fy4?pli=1
	// https://www.geeksforgeeks.org/how-to-copy-struct-type-using-value-and-pointer-reference-in-golang/
	//newResult := (*rf).prototype
	// ref: https://stackoverflow.com/questions/11127723/dynamically-create-variables-of-certain-type-based-on-string-in-go
	newResult := blerf{}
	//newResult := reflect.Zero((*rf).prototypeType).Interface()
	propertyPointers := rf.getStructPropertyPointers(&newResult)
	return newResult, propertyPointers, nil

	// Create property pointers for each of the fields in our result object
	// Ref: https://stackoverflow.com/questions/18926303/iterate-through-the-fields-of-a-struct-in-go
	//propertyPointers := make(PropertyPointers, (*rf).numFields)

	// TODO: Reject anything that's not a struct matching our requirements
	//rValue := reflect.ValueOf(newResult)

	// For each of its fields...
	//for i := 0; i < (*rf).numFields; i++ {
		// ref: https://samwize.com/2015/03/20/how-to-use-reflect-to-set-a-struct-field/
		//field := rValue.Field(i)
		//fmt.Printf("Field name: '%s', type: '%s'\n", rValue.Type().Field(i).Name, field.Type())
/*
		newVal, err := rf.newValue(field.Type().String())
		if nil != err {
			// Reject anything that's not one of our supported field types
			return nil, nil, err
		}
*/

		// ref: https://github.com/robertkrimen/otto/issues/83
		//abc := Abc{Def: 3}
		//abc_ := reflect.New(reflect.Indirect(reflect.ValueOf(abc)).Type()).Elem()
		//def_ := abc_.FieldByName("Def")
		//def_.SetInt(4)


		// ref: https://stackoverflow.com/questions/29184933/golang-reflect-get-pointer-to-a-struct-field-value
		//propertyPointers[i] = field.Addr().Interface()

/*
		rValue_ := reflect.New(reflect.Indirect(rValue).Type()).Elem()
		field_ := rValue_.FieldByName(rValue.Type().Field(i).Name)
		switch field.Type().Name() {
			case "string":
				t := field.String()
				//pv.Set(reflect.ValueOf(t))
				//field_.SetString(t)
				field_.Set(t)
				propertyPointers[i] = &t
			case "int":
				t := field.Int()
				//pv.Set(reflect.ValueOf(t))
				//field_.SetInt(t)
				field_.Set(t)
				propertyPointers[i] = &t
			case "uint":
				t := field.Uint()
				pv.Set(reflect.ValueOf(t))
				propertyPointers[i] = &t
			case "bool":
				t := field.Bool()
				pv.Set(reflect.ValueOf(t))
				propertyPointers[i] = &t
			case "float":
				t := field.Float()
				pv.Set(reflect.ValueOf(t))
				propertyPointers[i] = &t
		}
*/


		// https://groups.google.com/g/golang-dev/c/XWfzNWe4Fy4?pli=1
		//pfield := &field
		//var pv, zv reflect.Value
		//pv := reflect.ValueOf(&pfield).Elem()
		//zv = reflect.New(pv.Type().Elem())
		//pv.Set(zv)
		//zvp := zv.Pointer()
		//propertyPointers[i] = zvp

		// TODO: Pointers required?
		//field.Set(reflect.ValueOf(newVal))
/*
		//if field.CanSet() {
		//	field.Set(reflect.ValueOf(newVal))
		//} else {
		//	fmt.Printf("Cannot set for field name: '%s'", rValue.Type().Field(i).Name)
		//}
*/
	//}

	//finalResult := newResult.Interface().(Result)
	//return &finalResult, &propertyPointers, nil
	//return newResult, &propertyPointers, nil
}

// Make a new value based on the specified type
// TODO: Move this to a more generalized library
func (rf *ResultFactory) newValue(datatype string) (interface{}, error) {
	switch datatype {
		case "*string":
			var nv string
			return &nv, nil
		case "string":
			var nv string
			return nv, nil
		case "*[]byte":
			var nv []byte
			return &nv, nil
		case "[]byte":
			var nv []byte
			return nv, nil
		case "*int":
			var nv int
			return &nv, nil
		case "int":
			var nv int
			return nv, nil
		case "*int8":
			var nv int8
			return &nv, nil
		case "int8":
			var nv int8
			return nv, nil
		case "*int16":
			var nv int16
			return &nv, nil
		case "int16":
			var nv int16
			return nv, nil
		case "*int32":
			var nv int32
			return &nv, nil
		case "int32":
			var nv int32
			return nv, nil
		case "*int64":
			var nv int64
			return &nv, nil
		case "int64":
			var nv int64
			return nv, nil
		case "*uint":
			var nv uint
			return &nv, nil
		case "uint":
			var nv uint
			return nv, nil
		case "*uint8":
			var nv uint8
			return &nv, nil
		case "uint8":
			var nv uint8
			return nv, nil
		case "*uint16":
			var nv uint16
			return &nv, nil
		case "uint16":
			var nv uint16
			return nv, nil
		case "*uint32":
			var nv uint32
			return &nv, nil
		case "uint32":
			var nv uint32
			return nv, nil
		case "*uint64":
			var nv uint64
			return &nv, nil
		case "uint64":
			var nv uint64
			return nv, nil
		case "*bool":
			var nv bool
			return &nv, nil
		case "bool":
			var nv bool
			return nv, nil
		case "*float32":
			var nv float32
			return &nv, nil
		case "float32":
			var nv float32
			return nv, nil
		case "*float64":
			var nv float64
			return &nv, nil
		case "float64":
			var nv float64
			return nv, nil
		// TODO: Add support for other types (such as non-pointer version of all of the above)
	}
	return nil, errors.New(fmt.Sprintf("Unsupported type: '%s'", datatype))
}

