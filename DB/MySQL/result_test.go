package mysql

import(
	"fmt"
	"testing"

        "github.com/DATA-DOG/go-sqlmock"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewResult_ReturnsNothing_WhenGivenNothing(t *testing.T) {
	// Test
	var sut ResultIfc = NewResult(nil)	// <- ensures that we satisfy our interface

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewResult_ReturnsSomething_WhenGivenNonNilValue(t *testing.T) {
	// Setup
	r := sqlmock.NewResult(222, 333)

	// Test
	sut := NewResult(r)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Result_GetLastInsertId_ReturnsLastInsertId(t *testing.T) {
	// Setup
	var expectedLastInsertId int64 = 222
	var expectedRowsAffected int64 = 333
	r := sqlmock.NewResult(expectedLastInsertId, expectedRowsAffected)

	// Test
	sut := NewResult(r)
	actualLastInsertId, err := sut.GetLastInsertId()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualLastInsertId, t)
	ExpectInt64(expectedLastInsertId, *actualLastInsertId, t)
}

func TestThat_Result_GetLastInsertId_ReturnsError(t *testing.T) {
	// Setup
	r := sqlmock.NewErrorResult(fmt.Errorf("result error"))

	// Test
	sut := NewResult(r)
	actualLastInsertId, err := sut.GetLastInsertId()

	// Verify
	ExpectError(err, t)
	ExpectNil(actualLastInsertId, t)
}

func TestThat_Result_GetRowsAffected_ReturnsRowsAffected(t *testing.T) {
	// Setup
	var expectedLastInsertId int64 = 222
	var expectedRowsAffected int64 = 333
	r := sqlmock.NewResult(expectedLastInsertId, expectedRowsAffected)

	// Test
	sut := NewResult(r)
	actualRowsAffected, err := sut.GetRowsAffected()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(actualRowsAffected, t)
	ExpectInt64(expectedRowsAffected, *actualRowsAffected, t)
}

func TestThat_Result_GetRowsAffected_ReturnsError(t *testing.T) {
	// Setup
	r := sqlmock.NewErrorResult(fmt.Errorf("result error"))

	// Test
	sut := NewResult(r)
	actualRowsAffected, err := sut.GetRowsAffected()

	// Verify
	ExpectError(err, t)
	ExpectNil(actualRowsAffected, t)
}
