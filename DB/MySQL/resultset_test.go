package mysql

import(
	"testing"

	"github.com/DigiStratum/GoLib/DB/MySQL/nullables"
	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewResultSet_ReturnsEmptyResultSet(t *testing.T) {
	// Test
	sut := NewResultSet()

	// Verify
	ExpectNonNil(sut, t)
	ExpectInt(0, sut.Len(), t)
	ExpectTrue(sut.IsEmpty(), t)
}

func TestThat_ResultSet_Add_AddsResultRowToSet(t *testing.T) {
	// Setup
	sut := NewResultSet()
	expectedLen := 5

	// Test
	for i := 0; i < expectedLen; i++ {
		sut.Add(NewResultRow())
	}

	// Verify
	ExpectInt(expectedLen, sut.Len(), t)
	ExpectFalse(sut.IsEmpty(), t)
}

func TestThat_ResultSet_Finalize_PreventsMoreAdditions(t *testing.T) {
	// Setup
	sut := NewResultSet()
	expectedLen := 5

	// Test
	for i := 0; i < expectedLen; i++ {
		sut.Add(NewResultRow())
	}
	ExpectFalse(sut.IsFinalized(), t)
	sut.Finalize()
	ExpectTrue(sut.IsFinalized(), t)
	sut.Add(NewResultRow()) // This one should be prevented

	// Verify
	ExpectInt(expectedLen, sut.Len(), t)
}

func TestThat_ResultSet_Get_ReturnsNil_ForIndexOutOfBounds(t *testing.T) {
	// Setup
	sut := NewResultSet()

	// Test / Verify
	ExpectNil(sut.Get(0), t)
	sut.Add(NewResultRow())
	ExpectNonNil(sut.Get(0), t)
	ExpectNil(sut.Get(1), t)
}

func TestThat_ResultSet_Get_ReturnsResultRowsInExpectedOrder(t *testing.T) {
	// Setup
	sut := NewResultSet()
	var expectedLen int64 = 5
	var i int64
	for i = 0; i < expectedLen; i++ {
		resultRow := NewResultRow()
		resultRow.Set("row", *nullables.NewNullable(i))
		sut.Add(resultRow)
	}

	// Test / Verify
	for i = 0; i < expectedLen; i++ {
		resultRow := sut.Get(int(i))
		rowIdNullable := resultRow.Get("row")
		rowId := rowIdNullable.GetInt64()
		ExpectNonNil(rowId, t)
		ExpectInt64(i, *rowId, t)
	}
}

func TestThat_ResultSet_ToJson_ReturnsEmptyStateObject_WhenNoRowsSet(t *testing.T) {
	// Setup
	sut := NewResultSet()

	// Test
	actualJson, err := sut.ToJson()

	// Verify
	ExpectNonNil(actualJson, t)
	ExpectNoError(err, t)
	ExpectString("{\"Results\":[],\"IsFinalized\":false}", *actualJson, t)
}

func TestThat_ResultSet_ToJson_ReturnsFinalizedEmptyStateObject_WhenNoRowsSetAndFinalized(t *testing.T) {
	// Setup
	sut := NewResultSet()
	sut.Finalize()

	// Test
	actualJson, err := sut.ToJson()

	// Verify
	ExpectNonNil(actualJson, t)
	ExpectNoError(err, t)
	ExpectString("{\"Results\":[],\"IsFinalized\":true}", *actualJson, t)
}

func TestThat_ResultSet_ToJson_ReturnsPopulatedStateObject_WhenRowsSet(t *testing.T) {
	// Setup
	sut := NewResultSet()
	resultRow := NewResultRow()
	resultRow.Set("row", *nullables.NewNullable(111))
	sut.Add(resultRow)

	// Test
	actualJson, err := sut.ToJson()

	// Verify
	ExpectNonNil(actualJson, t)
	ExpectNoError(err, t)
	ExpectString("{\"Results\":[{\"row\":111}],\"IsFinalized\":false}", *actualJson, t)
}
