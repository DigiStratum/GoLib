package mysql

/*
Query Results - When a query is created, it needs a prototype (ResultIfc) struct with Property Pointers to Scan() result data into
*/

type PropertyPointers []interface{}

// Result public interface
type ResultIfc interface {
	ZeroClone() (ResultIfc, PropertyPointers)
}

// Result set public interface
type ResultSetIfc interface {
	Len() int
	Append(result ResultIfc) ResultSetIfc
	Get(index int) ResultIfc
}

// Result set private data type
type resultSet []ResultIfc

// Make a new one of these
func NewResultSet() ResultSetIfc {
	rs := make(resultSet, 0)
	return &rs
}

// Get the length of the ResultSet
func (rs *resultSet) Len() int {
	return len(*rs)
}

// Make a new ResultSet with the supplied Result appended to the end
func (rs *resultSet) Append(result ResultIfc) ResultSetIfc {
	nrs := append(*rs, result)
	return &nrs
}

// Get the Result from this set at the specified index
func (rs *resultSet) Get(index int) ResultIfc {
	if index >= rs.Len() { return nil }
	return (*rs)[index]
}

