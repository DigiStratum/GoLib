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

// A Query Result Row Object Interface
type ResultIfc interface {
	GetZeroClone() ResultIfc
	GetPropertyPointers() (*[]interface{}, error)
}

// Base struct for DB query result row objects
type Result struct {
	propertyPointers	*[]interface{}
}

// Make a "Zero" clone (all properties will be zero value) of the Result structure
func (r *Result) GetZeroClone() ResultIfc {
	thisResult := *r
	newResult := reflect.New(reflect.TypeOf(thisResult).Elem())
	return newResult
}

// Create a set of pointers to capture query result column values for each row processed with Scan()
func (r *Result) GetPropertyPointers() (*[]interface{}, error) {
	// If we have already figured this out, just return the result
	if nil != r.propertyPointers { return r.propertyPointers, nil }

	// Reflect on ourselves
	result := *r

	// Create property pointers for each of the fields in our result object
	// Ref: https://stackoverflow.com/questions/18926303/iterate-through-the-fields-of-a-struct-in-go
	numFields := reflect.TypeOf(result).NumField()
	pp := make([]interface{}, numFields)
	r.propertyPointers = &pp

	// TODO: Reject anything that's not a struct matching our requirements
	rValue := reflect.ValueOf(result)

	// For each of its fields...
	for i := 0; i < numFields; i++ {
		// ref: https://samwize.com/2015/03/20/how-to-use-reflect-to-set-a-struct-field/
		//fmt.Printf("Field name: '%s', type: '%s'\n", voPrototype.Type().Field(i).Name, field.Type())
		field:= rValue.Field(i)
		newVal, err := r.newValue(field.Type().String())
		if nil != err {
			// Reject anything that's not one of our supported field types
			return nil, err
		}

		// TODO: Pointers required?
		field.Set(&newVal)
		propertyPointers[i] = &newVal
	}
	return r.propertyPointers, nil
}

// Make a new value based on the specified type
// TODO: Move this to a more generalized library
func (r *Result) newValue(datatype string) (interface{}, error) {
	switch datatype {
		case "*string":
			nv := ""
			return &nv, nil
		case "*[]byte":
			nv := []byte{}
			return &nv, nil
		case "*int":
			nv := 0.(int)
			return &nv, nil
		case "*int8":
			nv := 0.(int8)
			return &nv, nil
		case "*int16":
			nv := 0.(int16)
			return &nv, nil
		case "*int32":
			nv := 0.(int32)
			return &nv, nil
		case "*int64":
			nv := 0.(int64)
			return &nv, nil
		case "*uint":
			nv := 0.(uint)
			return &nv, nil
		case "*uint8":
			nv := 0.(uint8)
			return &nv, nil
		case "*uint16":
			nv := 0.(uint16)
			return &nv, nil
		case "*uint32":
			nv := 0.(uint32)
			return &nv, nil
		case "*uint64":
			nv := 0.(uint64)
			return &nv, nil
		case "*bool":
			nv := true
			return &nv, nil
		case "*float32":
			nv := 0.(float32)
			return &nv, nil
		case "*float64":
			nv := 0.(float64)
			return &nv, nil
		// TODO: Add support for other types (such as non-pointer version of all of the above)
	}
	return nil, errors.New(fmt.Sprintf("Unsupported type: '%s'", datatype))
}

type ResultSet []ResultIfc


