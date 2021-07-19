package mysql

type NewResultIfc interface {
	Get(field string) interface{}
}

type writableResultIfc interface {
	NewResultIfc
	ToImmutable() NewResultIfc
}

type newResult struct {
	result		map[string]interface{}
}

func newNewResult(result map[string]interface{}) NewResultIfc {
	result := newResult{
		result:		result,
	}
	return &result
}

// -------------------------------------------------------------------------------------------------
// NewResultIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *newResult) Get(field string) interface{} {
	if value, ok := (*r).result[field]; ok { return value }
	return nil
}

// -------------------------------------------------------------------------------------------------
// writableResultIfc Private Interface
// -------------------------------------------------------------------------------------------------

func (r *newResult) ToImmutable() NewResultIfc {
	immutable, _ := r.(NewResultIfc)
	return immutable
}

// -------------------------------------------------------------------------------------------------

type NewResultSetIfc interface {
	Get(resultNum int) NewResultIfc
	Num() int
	IsEmpty() bool
}

type writableResultIfc interface {
	NewResultSetIfc
	Add(result NewResultIfc)
	ToImmutable() NewResultSetIfc
}

type newResultSet struct {
	results		[]NewResultIfc
}

func newWritableResultSet() writableResultIfc {
	rs := newResultSet{
		results:	make([]NewResultIfc, 0)
	}
	return &rs
}

// -------------------------------------------------------------------------------------------------
// NewResultSetIfc Public Interface
// -------------------------------------------------------------------------------------------------
func (rs *newResultSet) Get(resultNum int) NewResultIfc {
	if resultNum >= rs.Num() { return nil }
	return (*rs).results[resultNum]
}

func (rs *newResultSet) Num() int {
	return len((*rs).results)
}

func (rs *newResultSet) IsEmpty() bool {
	return rs.Num() == 0
}

// -------------------------------------------------------------------------------------------------
// writableResultIfc Private Interface
// -------------------------------------------------------------------------------------------------

func (rs *newResultSet) Add(result NewResultIfc) {
	(*rs).results = append((*rs.results), result)
}

func (rs *newResultSet) ToImmutable() NewResultSetIfc {
	immutable, _ := rs.(NewResultSetIfc)
	return immutable
}

type NewQueryIfc interface {
	Exec(args ...interface{}) error
	GetResults() NewResultSetIfc
}

type newQuery struct {
	query		string
	results		NewResultSetIfc
}

func NewNewQuery(query string) NewQueryIfc {
	q := newQuery{
		query:		query,
		results:	newWritableResultSet(),
	}
	return &q
}

// -------------------------------------------------------------------------------------------------
// NewQueryIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (nq *newQuery) Exec(args ...interface{}) error {
	rows, _ := db.Query((*nq).query)
	cols, _ := rows.Columns()
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		(*nq).results.add(newNewResult(m))
	}
	return nil // no error
}

func (nq *newQuery) GetResults() NewResultSetIfc {
	return (*nq).results
}