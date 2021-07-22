package mysql

import (
	"fmt"
	"encoding/json"
)

type NewResultIfc interface {
	Get(field string) NullableIfc
	Fields() []string
	ToJson() (*string, error)
}

type resultRow map[string]NullableIfc

type newResult struct {
	//result		map[string]interface{}
	//result		map[string]string
	result		resultRow
}

//func newNewResult(result map[string]interface{}) NewResultIfc {
//func newNewResult(result map[string]string) NewResultIfc {
func newNewResult(result resultRow) NewResultIfc {
	r := newResult{
		result:		result,
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// NewResultIfc Public Interface
// -------------------------------------------------------------------------------------------------

//func (r *newResult) Get(field string) interface{} {
func (r *newResult) Get(field string) NullableIfc {
	if value, ok := (*r).result[field]; ok { return value }
	return nil
}

// Pluck the fields out of the result set and just return them so that caller can iterate with Get()
func (r *newResult) Fields() []string {
	fields := make([]string, 0)
	for field, _ := range (*r).result {
		fields = append(fields, field)
	}
	return fields
}

func (r *newResult) ToJson() (*string, error) {
	// TODO: See if there is a way to encode each Nullable value as it's native JSON data type instead of making them all strings
	for field, value := range (*r).result {
		var svalue string = "nil"
		if ! value.IsNil() {
			sptr := value.GetString()
			if nil != sptr {
				svalue = fmt.Sprintf("%s::%s", GetNullableTypeString(value.GetType()), *sptr)
			}
		}
		fmt.Printf("Field['%s'] = '%s'\n", field, svalue)
	}
	fmt.Println("")
	jsonBytes, err := json.Marshal((*r).result)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------

type NewResultSetIfc interface {
	// Public
	Get(resultNum int) NewResultIfc
	Len() int
	IsEmpty() bool
	// Private
	add(result NewResultIfc)
}

type newResultSet struct {
	results		[]NewResultIfc
}

func newNewResultSet() NewResultSetIfc {
	rs := newResultSet{
		results:	make([]NewResultIfc, 0),
	}
	return &rs
}

// -------------------------------------------------------------------------------------------------
// NewResultSetIfc Public Interface
// -------------------------------------------------------------------------------------------------
func (rs *newResultSet) Get(resultNum int) NewResultIfc {
	if resultNum >= rs.Len() { return nil }
	return (*rs).results[resultNum]
}

func (rs *newResultSet) Len() int {
	return len((*rs).results)
}

func (rs *newResultSet) IsEmpty() bool {
	return rs.Len() == 0
}

// -------------------------------------------------------------------------------------------------
// NewResultSetIfc Private Interface
// -------------------------------------------------------------------------------------------------

func (rs *newResultSet) add(result NewResultIfc) {
	(*rs).results = append((*rs).results, result)
}

// -------------------------------------------------------------------------------------------------

type NewQueryIfc interface {
	Run(args ...interface{}) (NewResultSetIfc, error)
}

type newQuery struct {
	connection	ConnectionIfc
	query		string
}

func NewNewQuery(connection ConnectionIfc, query string) NewQueryIfc {
	q := newQuery{
		connection:	connection,
		query:		query,
	}
	return &q
}

// -------------------------------------------------------------------------------------------------
// NewQueryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (nq *newQuery) Run(args ...interface{}) (NewResultSetIfc, error) {
	results := newNewResultSet()
	// ref: https://kylewbanks.com/blog/query-result-to-map-in-golang
	rows, err := (*nq).connection.GetConnection().Query((*nq).query)
	if nil != err { return nil, err }
	cols, _ := rows.Columns()
	for rows.Next() {
		columnPointers := make([]interface{}, len(cols))
		//columns := make([]interface{}, len(cols))
		columns := make([]string, len(cols))
		//columnPointers := make([]*string, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err = rows.Scan(columnPointers...); err != nil { return nil, err }

		// Create our map, and retrieve the value for each column from the pointers,
		// slice, storing it in the map with the name of the column as the key.
		result := make(resultRow)
		for i, colName := range cols {
			//val := columnPointers[i].(*interface{})
			//val := columnPointers[i].(*string)
			val := columns[i]
fmt.Printf("Run().Row(%T) -> '%s'\n", val, val)

			//m[colName] = *val
			//m[colName] = fmt.Sprintf("%v", *val)
			//m[colName] = fmt.Sprintf("%v", (*val).(string))
			//m[colName] = fmt.Sprintf("%v", string(*val))
			result[colName] = NewNullable(val)
		}
		results.add(newNewResult(result))
	}
	return results, nil
}
