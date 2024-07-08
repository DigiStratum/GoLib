package fileio

import(
	"testing"
	//"runtime"
	//"path/filepath"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewFileSet_Returns_FileSetIfc(t *testing.T) {
	// Test
	var sut FileSetIfc = NewFileSet()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_FileSet_Len_Returns_Zero_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewFileSet()

	// Test
	actual := sut.Len()

	// Verify
	ExpectInt(0, actual, t)
}

func TestThat_FileSet_Len_Returns_NonZero_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewFileSet()
	sut.AddFile("missingfile.txt")

	// Test
	actual := sut.Len()

	// Verify
	ExpectInt(1, actual, t)
}

func TestThat_FileSet_GetIterator_Returns_NonZero_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewFileSet()
	sut.AddFile("missingfile.txt")

	// Test
	it := sut.GetIterator()
	actual1 := it()
	actual2 := it()

	// Verify
	ExpectNonNil(actual1, t)
	_, ok := actual1.(FileIfc)
	ExpectTrue(ok, t)
	ExpectNil(actual2, t)
}

