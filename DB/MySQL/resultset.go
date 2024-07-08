package mysql

/*

TODO:
 * Add a FromJson()? Any need/value for this? How would it handle the nested structure?
*/

import (
	"encoding/json"

	it "github.com/DigiStratum/GoLib/Data/iterable"
)

type ResultSetIfc interface {
	// Embedded interface(s)
	it.IterableIfc

	// Our own interface
	Get(rowNum int) *ResultRow
	Len() int
	IsEmpty() bool
	Add(result ResultRowIfc)
	IsFinalized() bool
	Finalize()
	ToJson() (*string, error)
}

type resultSet struct {
	results		[]ResultRow
	isFinalized	bool
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewResultSet() *resultSet {
	return &resultSet{
		results:		make([]ResultRow, 0),
	}
}

// -------------------------------------------------------------------------------------------------
// ResultSetIfc
// -------------------------------------------------------------------------------------------------

func (r resultSet) Get(rowNum int) *ResultRow {
	if rowNum >= r.Len() { return nil }
	return &r.results[rowNum]
}

func (r resultSet) Len() int {
	return len(r.results)
}

func (r resultSet) IsEmpty() bool {
	return r.Len() == 0
}

func (r *resultSet) Add(result ResultRowIfc) {
	// No more changes (immutable) after finalization
	if r.IsFinalized() { return }
	resultRow := result.(*ResultRow)
	r.results = append(r.results, *resultRow)
}

func (r resultSet) IsFinalized() bool {
	return r.isFinalized
}

func (r *resultSet) Finalize() {
	r.isFinalized = true
}

func (r resultSet) ToJson() (*string, error) {
	jsonBytes, err := r.MarshalJSON()
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler
// -------------------------------------------------------------------------------------------------

func (r resultSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.results)
}

// -------------------------------------------------------------------------------------------------
// IterableIfc
// -------------------------------------------------------------------------------------------------

func (r resultSet) GetIterator() func () interface{} {
	idx := 0
	var data_len = r.Len()
	return func () interface{} {
		// If we're done iterating, return do nothing
		if idx >= data_len { return nil }
		prev_idx := idx
		idx++
		return &r.results[prev_idx]
	}
}

