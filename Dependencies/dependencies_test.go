package dependencies

import(
	"strings"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// NewDependencies()
func TestThat_NewDependencies_ReturnsSomething(t *testing.T) {
	// Setup
	var sut DependenciesIfc

	// Test
	sut = NewDependencies()
	_, interfaceAssertionOk := sut.(DependenciesIfc)

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectTrue(interfaceAssertionOk, t) { return }
}

// Get(name string) *dependency
func TestThat_Dependencies_Get_ReturnsNil_ForMissingDependency(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.Get("bogusdep")

	// Verify
	if ! ExpectNil(actual, t) { return }
}

// Add(deps ...*dependency)
// Has(name string) bool
func TestThat_Dependencies_Get_ReturnsDependency_From_Add(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	sut.Add(NewDependency("sampledep"))
	actual := sut.Get("sampledep")
	actualHas := sut.Has("sampledep")
	actualHasNot := sut.Has("bogusdep")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actualHas, t) { return }
	if ! ExpectFalse(actualHasNot, t) { return }
}

// GetVariant(name, variant string) *dependency
// HasVariant(name, variant string) bool
func TestThat_Dependencies_GetVariant_ReturnsNil_ForMissingVariant(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	sut.Add(NewDependency("sampledep"))
	actual := sut.GetVariant("sampledep", "bogusvariant")
	actualHas := sut.HasVariant("sampledep", DEP_VARIANT_DEFAULT)
	actualHasNot := sut.HasVariant("sampledep", "bogusvariant")

	// Verify
	if ! ExpectNil(actual, t) { return }
	if ! ExpectTrue(actualHas, t) { return }
	if ! ExpectFalse(actualHasNot, t) { return }
}


// GetVariants(name string) []string
func TestThat_Dependencies_GetVariants_ReturnsEmptySet_ForNewDependencies(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.GetVariants("bogusdep")

	// Verify
	if ! ExpectInt(0, len(actual), t) { return }
}

func TestThat_Dependencies_GetVariants_ReturnsExpectedSet(t *testing.T) {
	// Setup
	sut := NewDependencies(
		NewDependency("depname"),
	)

	// Test
	actual := sut.GetVariants("depname")

	// Verify
	if ! ExpectInt(1, len(actual), t) { return }
	if ! ExpectString(DEP_VARIANT_DEFAULT, actual[0], t) { return }
}


// GetAllVariants() map[string][]string
func TestThat_Dependencies_GetAllVariants_ReturnsEmptySet_ForNewDependencies(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.GetAllVariants()

	// Verify
	if ! ExpectInt(0, len(actual), t) { return }
}

func TestThat_Dependencies_GetAllVariants_Returns4_For2Names2Variants(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	sut.Add(
		NewDependency("sampledepA"),
		NewDependency("sampledepA").SetVariant("alternateA"),
		NewDependency("sampledepB"),
		NewDependency("sampledepB").SetVariant("alternateB"),
	)
	actual := sut.GetAllVariants()
	actualA, okA := actual["sampledepA"]
	actualAVariants := strings.Join(actualA[:], ":")
	actualB, okB := actual["sampledepB"]
	actualBVariants := strings.Join(actualB[:], ":")

	// Verify
	if ! ExpectInt(2, len(actual), t) { return }
	if ! ExpectTrue(okA && okB, t) { return }
	if ! ExpectInt(2, len(actualA), t) { return }
	if ! ExpectInt(2, len(actualB), t) { return }
	if ! ExpectTrue(((actualAVariants == "default:alternateA") || (actualAVariants == "alternateA:default")), t) { return }
	if ! ExpectTrue(((actualBVariants == "default:alternateB") || (actualBVariants == "alternateB:default")), t) { return }
}

