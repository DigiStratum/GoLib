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
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectTrue(interfaceAssertionOk, t) { return }
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
	if ! ExpectString(DEP_NAME, actualName, t) { return }
	if ! ExpectString(DEP_VARIANT_DEFAULT, actualVariant, t) { return }
	if ! ExpectFalse(actualIsRequired, t) { return }
}

// SetVariant(variant string) *dependency
func TestThat_Dependency_SetVariant_ChangesVariantValue(t *testing.T) {
	// Setup
	sut := NewDependency(DEP_NAME)

	// Test
	sut.SetVariant(DEP_VARIANT)
	actualVariant := sut.GetVariant()

	// Verify
	if ! ExpectString(DEP_VARIANT, actualVariant, t) { return }
}

// SetRequired() *dependency
func TestThat_Dependency_SetRequired_ChangesRequiredValue(t *testing.T) {
	// Setup
	sut := NewDependency(DEP_NAME)

	// Test
	sut.SetRequired()
	actualIsRequired := sut.IsRequired()

	// Verify
	if ! ExpectTrue(actualIsRequired, t) { return }
}

