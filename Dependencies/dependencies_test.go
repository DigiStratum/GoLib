package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDependencies_ReturnsSomething(t *testing.T) {
	// Setup
	var sut *Dependencies

	// Test
	sut = NewDependencies()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Dependencies_Set_AddsNamedDependency(t *testing.T) {
	// Setup
	var sut *Dependencies = NewDependencies()
	expectedName := "bogusname"
	expectedValue := "bogusvalue"

	// Test & Verify
	sut.Set(expectedName, expectedValue)
	hasIt := sut.Has(expectedName)
	ExpectTrue(hasIt, t)
	actual := sut.Get(expectedName)
	ExpectNonNil(actual, t)
	actualValue := actual.(string)
	ExpectString(expectedValue, actualValue, t)
	names := sut.GetNames()
	ExpectTrue(sut.HasAll(names), t)
	ExpectFalse(sut.Has("unexpectedname"), t)
}

// GetIterator() func () *Dependency

func TestThat_Dependencies_GetIterator_ReturnsGoodIterator(t *testing.T) {
	// Setup
	var sut *Dependencies = NewDependencies()
	sut.Set("name0", "value0")
	sut.Set("name1", "value1")
	sut.Set("name2", "value2")

	// Test
	var it func () *Dependency = sut.GetIterator()

	// Verify
	ExpectNonNil(it, t)
	num := 0
	for dep := it(); nil != dep; dep = it() {
		num++
		if 3 < num { break }
	}
	ExpectInt(3, num, t)

}

