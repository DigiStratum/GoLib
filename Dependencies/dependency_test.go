package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)


// NewDependency(name string, dep interface{}) *Dependency
func TestThat_NewDependency_ReturnsSomething(t *testing.T) {
	// Setup
	var sut *Dependency
	expectedName := "depname"
	expectedValue := "depvalue"

	// Test
	sut = NewDependency(expectedName, expectedValue)

	// Verify
	ExpectNonNil(sut, t)
	actualName, actualValuei := sut.GetDep()
	actualValue := actualValuei.(string)
	ExpectString(expectedName, actualName, t)
	ExpectString(expectedValue, actualValue, t)
}

