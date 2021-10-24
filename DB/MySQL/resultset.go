package mysql

/*

TODO:
 * Add a FromJson()? Any need/value for this? How would it handle the nested structure?
*/

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

// Non-exported structure with exported properties that we can serialize
type resultSetSerializableProperties struct {
	Results		[]ResultRow
	IsFinalized	bool
}

// Exported structure with non-exported properties to prevent consumer from accessing directly
type ResultSet struct {
	props		resultSetSerializableProperties
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewResultSet() *ResultSet {
	return &ResultSet{
		props:	resultSetSerializableProperties{
			Results:	make([]ResultRow, 0),
		},
	}
}

// -------------------------------------------------------------------------------------------------
// ResultSetIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r ResultSet) Get(rowNum int) *ResultRow {
	if rowNum >= r.Len() { return nil }
	return &r.props.Results[rowNum]
}

func (r ResultSet) Len() int {
	return len(r.props.Results)
}

func (r ResultSet) IsEmpty() bool {
	return r.Len() == 0
}

func (r *ResultSet) Add(result ResultRowIfc) {
	// No more changes (immutable) after finalization
	if r.IsFinalized() { return }
	resultRow := result.(*ResultRow)
	r.props.Results = append(r.props.Results, *resultRow)
}

func (r ResultSet) IsFinalized() bool {
	return r.props.IsFinalized
}

func (r *ResultSet) Finalize() {
	r.props.IsFinalized = true
}

func (r ResultSet) ToJson() (*string, error) {
	jsonBytes, err := r.MarshalJSON()
	if nil != err { return nil, err }
	jsonString := string(jsonBytes[:])
	return &jsonString, nil
}

// -------------------------------------------------------------------------------------------------
// encoding/json.Marshaler Public Interface
// -------------------------------------------------------------------------------------------------

func (r ResultSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.props)
}
