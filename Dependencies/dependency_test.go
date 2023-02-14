package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDependency_ReturnsSomething(t *testing.T) {
	// Setup
	var sut DependencyIfc

	// Test
	sut = NewDependency(DEP_NAME).SetVariant(DEP_VARIANT).SetRequired()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewDependency_MatchesExpectedProperties(t *testing.T) {
	// Setup
	var sut DependencyIfc
	sut = NewDependency(DEP_NAME).SetVariant(DEP_VARIANT).SetRequired()

	// Test
	actualName := sut.GetName()
	actualVariant := sut.GetVariant()
	actualIsRequired := sut.IsRequired()

	// Verify
	ExpectString(DEP_NAME, actualName, t)
	ExpectString(DEP_VARIANT, actualVariant, t)
	ExpectTrue(actualIsRequired, t)
}
