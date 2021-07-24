package mysql

import (
	"encoding/json"
)

type ResultSetIfc interface {
	// Public
	Get(rowNum int) ResultRowIfc
	Len() int
	IsEmpty() bool
	ToJson() (*string, error)
	// Private
	add(result ResultRowIfc)
}

type resultSet struct {
	results		[]ResultRowIfc
}

func NewResultSet() ResultSetIfc {
	rs := resultSet{
		results:	make([]ResultRowIfc, 0),
	}
	return &rs
}

// -------------------------------------------------------------------------------------------------
// ResultSetIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (rs *resultSet) Get(rowNum int) ResultRowIfc {
	if rowNum >= rs.Len() { return nil }
	return (*rs).results[rowNum]
}

func (rs *resultSet) Len() int {
	return len((*rs).results)
}

func (rs *resultSet) IsEmpty() bool {
	return rs.Len() == 0
}

func (rs *resultSet) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(rs)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// ResultSetIfc Private Interface
// -------------------------------------------------------------------------------------------------

func (rs *resultSet) add(result ResultRowIfc) {
	(*rs).results = append((*rs).results, result)
}
