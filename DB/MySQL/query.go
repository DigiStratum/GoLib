package mysql

/*
TODO: Add some sort of query builder - this will allow us to ditch writing SQL for most needs.
*/

import (
	_ "github.com/go-sql-driver/mysql"
)

// The spec for a prepared statement query. Single '?' substitution is handled by db.Query()
// automatically. '???' expands to include enough placeholders (as with an IN () list for any count
// of keys > min. max must be >= min unless max == 0.
type Query struct {
	query		string          // The query to execute as prepared statement
	//resultPrototype	ResultIfc       // Object to use as a prototype to produce query Result row objects
	resultFactory	*ResultFactory	// Object to use as a prototype to produce query Result row objects
}

// Make a new one of these
//func NewQuery(query string, prototype ResultIfc) *Query {
//func NewQuery(query string, resultFactory ResultFactory) *Query {
func NewQuery(query string, prototype Result) *Query {
	resultFactory := NewResultFactory(prototype)
	q := Query{
		query:		query,
		//resultPrototype:	prototype,
		resultFactory:	resultFactory,
	}
	return &q
}

// Run this query against the supplied database Connection with the provided query arguments
func (q *Query) Run(conn *Connection, args ...interface{}) (*ResultSet, error) {
	results := ResultSet{}
	protoQuery := (*q).query
	// TODO: expand query '???' placeholders
	finalQuery := protoQuery

	// Execute the Query
	rows, err := conn.GetConnection().Query(finalQuery, args...)
	if err != nil { return nil, err }

	// Process the result rows
	for rows.Next() {
		// Make a new result object for this row
		// (... and get pointers to all the all the result object members)
		//result := (*q).resultPrototype.GetZeroClone()
		result, resultProperties, err := (*q).resultFactory.MakeNewResult()
		if nil != err { return nil, err }

		// Read MySQL columns for this row into the result object member pointers
		err = rows.Scan(*resultProperties...)
		if nil != err { return nil, err }

		results = append(results, result)
	}

	return &results, nil
}



/*
// ------------------------------------------------------------------------------------------------
// Query Bits

// ------------------------------------------------------------------------------------------------
// Execute a query with varargs for substitution and Structured results
// ref: https://appliedgo.net/generics/
// ref: https://stackoverflow.com/questions/37851500/how-to-copy-an-interface-value-in-go
// ref: https://forum.golangbridge.org/t/database-rows-scan-unknown-number-of-columns-json/7378/2
// ref: https://stackoverflow.com/questions/26744873/converting-map-to-struct/26746461
// ref: https://stackoverflow.com/questions/29184933/golang-reflect-get-pointer-to-a-struct-field-value
func (dbc *Connection) Query(querySpec *QuerySpec, args ...interface{}) (*[]interface{}, error) {
	results := []interface{}{}
	//template := querySpec.Template
	template := querySpec.ResultFactory()

	// Ref: https://stackoverflow.com/questions/18926303/iterate-through-the-fields-of-a-struct-in-go
	numFields := reflect.TypeOf(template).NumField()
	fmt.Printf("SQuery() QuerySpec.Template has %d Fields\n", numFields)
	values := make([]interface{}, numFields)
	templateValue := reflect.ValueOf(template)
	tvType := templateValue.Type()
	//templateValue := reflect.ValueOf(querySpec.Template).Elem()
	for i := 0; i < numFields; i++ {
		//fi := templateValue.Field(i).Interface()
		//values[i] = &fi
		// ref: https://stackoverflow.com/questions/27992821/how-get-pointer-of-structs-member-from-interface
		//values[i] = templateValue.Field(i).Addr().Interface()
		//values[i] = templateValue.Elem().FieldByIndex(i).Addr().Interface()
		//values[i] = templateValue.Field(i).Addr()
		//values[i] = templateValue.Field(i).Interface().Addr()
		//values[i] = templateValue.Field(i).Interface()

		// ref: https://stackoverflow.com/questions/29184933/golang-reflect-get-pointer-to-a-struct-field-value
		//valueField := templateValue.Field(i)
		//values[i] = valueField.Addr().Interface()

		// ref: https://samwize.com/2015/03/20/how-to-use-reflect-to-set-a-struct-field/
		fieldName := tvType.Field(i).Name
		field:= templateValue.Field(i)
		fmt.Printf("Field name: '%s', type: '%s'\n", fieldName, field.Type())
		//v := templateValue.FieldByName(fieldName).Interface()
		v := field.Interface()
		switch field.Type().String() {
			case "*int":
				values[i] = v.(*int)
			case "*string":
				values[i] = v.(*string)
		}
		//values[i] = templateValue.Elem().FieldByName(fieldName).Addr().Interface()
	}

	protoQuery := querySpec.Query
	// TODO: expand querySpec.Query '???' placeholders
	finalQuery := protoQuery

	// Convert string args to interface{} for Query()
	//iArgs := make([]interface{}, len(args))
	//for i, v := range args { iArgs[i] = v }

	// Execute the Query
	//rows, err := dbc.Conn.Query(finalQuery, iArgs...)
	rows, err := dbc.Conn.Query(finalQuery, args...)
	if err != nil { return nil, err }

	// Process the result rows
	for rows.Next() {

		// figure out what columns were returned
		// the column names will be the JSON object field keys
		//columns, err := rows.ColumnTypes()
		//if err != nil { return nil, err }

		// Scan needs an array of pointers to the values it is setting
		// This creates the object and sets the values correctly
		//values := make([]interface{}, len(columns))
		//object := map[string]interface{}{}
		//for i, column := range columns {
		//	object[column.Name()] = reflect.New(column.ScanType()).Interface()
		//	values[i] = object[column.Name()]
		//}
		//err = rows.Scan(values...)
		//if err != nil { return nil, err }
		//results = append(results, template)
		// Make a new result object for this row
		result, err := querySpec.ResultFactory.NewResult()
		if nil != err { return nil, err }

		// Get pointers to all the all the result object members
		resultProperties := result.GetPropertyPointers()

		// Read MySQL columns for this row into the result object member pointers
		err = rows.Scan(resultProperties...)
		if nil != err { return nil, err }

		results = append(results, result)
	}

	return &results, nil
}
*/
