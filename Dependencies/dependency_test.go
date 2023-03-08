package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

const DEP_NAME = "depnamegood"
const DEP_VARIANT = "depvariant"

// NewDependency()
func TestThat_NewDependency_ReturnsSomething(t *testing.T) {
	// Setup
	var sut DependencyIfc

	// Test
	sut = NewDependency(DEP_NAME)
	_, interfaceAssertionOk := sut.(DependencyIfc)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(interfaceAssertionOk, t)
}

// GetName() string
// GetVariant() string
// IsRequired() bool
func TestThat_Dependency_DefaultProperties_MatcheExpectations(t *testing.T) {
	// Setup
	sut := NewDependency(DEP_NAME)

	// Test
	actualName := sut.GetName()
	actualVariant := sut.GetVariant()
	actualIsRequired := sut.IsRequired()

	// Verify
	ExpectString(DEP_NAME, actualName, t)
	ExpectString(DEP_VARIANT_DEFAULT, actualVariant, t)
	ExpectFalse(actualIsRequired, t)
}

// SetVariant(variant string) *dependency
func TestThat_Dependency_SetVariant_ChangesVariantValue(t *testing.T) {
	// Setup
	sut := NewDependency(DEP_NAME)

	// Test
	sut.SetVariant(DEP_VARIANT)
	actualVariant := sut.GetVariant()

	// Verify
	ExpectString(DEP_VARIANT, actualVariant, t)
}

// SetRequired() *dependency
func TestThat_Dependency_SetRequired_ChangesRequiredValue(t *testing.T) {
	// Setup
	sut := NewDependency(DEP_NAME)

	// Test
	sut.SetRequired()
	actualIsRequired := sut.IsRequired()

	// Verify
	ExpectTrue(actualIsRequired, t)
}

