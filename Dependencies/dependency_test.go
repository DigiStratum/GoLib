package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// NewDependency(name string, dep interface{}) *Dependency
func TestThat_NewDependency_ReturnsSomething(t *testing.T) {
	// Setup
	var sut DependencyIfc

	// Test
	sut = NewDependency(DEP_NAME, DEP_VARIANT, true)

	// Verify
	ExpectNonNil(sut, t)
}

// NewDependency(name string, dep interface{}) *Dependency
func TestThat_NewDependency_MatchesExpectedProperties(t *testing.T) {
	// Setup
	var sut DependencyIfc
	sut = NewDependency(DEP_NAME, DEP_VARIANT, true)

	// Test
	actualName := sut.GetName()
	actualVariant := sut.GetVariant()
	actualIsRequired := sut.IsRequired()

	// Verify
	ExpectString(DEP_NAME, actualName, t)
	ExpectString(DEP_VARIANT, actualVariant, t)
	ExpectTrue(actualIsRequired, t)
}
