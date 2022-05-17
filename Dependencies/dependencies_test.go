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

// Set(name string, dep interface{})
// Get(name string) interface{}
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
}

