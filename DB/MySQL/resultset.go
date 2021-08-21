package mysql

import (
	"encoding/json"
)

type ResultSetIfc interface {
	Get(rowNum int) *ResultRow
	Len() int
	IsEmpty() bool
	Add(result ResultRowIfc)
	IsFinalized() bool
	Finalize()
	ToJson() (*string, error)
}

type ResultSet struct {
	results		[]ResultRow
	isFinalized	bool
}

func NewResultSet() *ResultSet {
	return &ResultSet{
		results:	make([]ResultRow, 0),
	}
}

// -------------------------------------------------------------------------------------------------
// ResultSetIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r ResultSet) Get(rowNum int) *ResultRow {
	if rowNum >= r.Len() { return nil }
	return &r.results[rowNum]
}

func (r ResultSet) Len() int {
	return len(r.results)
}

func (r ResultSet) IsEmpty() bool {
	return r.Len() == 0
}

func (r *ResultSet) Add(result ResultRowIfc) {
	// No more changes (immutable) after finalization
	if r.IsFinalized() { return }
	resultRow := result.(*ResultRow)
	(*r).results = append((*r).results, *resultRow)
}

func (r ResultSet) IsFinalized() bool {
	return r.isFinalized
}

func (r *ResultSet) Finalize() {
	r.isFinalized = true
}

func (r ResultSet) ToJson() (*string, error) {
	jsonBytes, err := json.Marshal(r)
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}
