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
	ExpectNonNil(sut, t)
	ExpectTrue(interfaceAssertionOk, t)
}

// Get(name string) *dependency
func TestThat_Dependencies_Get_ReturnsNil_ForMissingDependency(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.Get("bogusdep")

	// Verify
	ExpectNil(actual, t)
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
	ExpectNonNil(actual, t)
	ExpectTrue(actualHas, t)
	ExpectFalse(actualHasNot, t)
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
	ExpectNil(actual, t)
	ExpectTrue(actualHas, t)
	ExpectFalse(actualHasNot, t)
}


// GetAllVariants() map[string][]string
func TestThat_Dependencies_GetAllVariants_ReturnsEmptySet_ForNewDependencies(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.GetAllVariants()

	// Verify
	ExpectInt(0, len(actual), t)
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
	ExpectInt(2, len(actual), t)
	ExpectTrue(okA && okB, t)
	ExpectInt(2, len(actualA), t)
	ExpectInt(2, len(actualB), t)
	ExpectTrue(((actualAVariants == "default:alternateA") || (actualAVariants == "alternateA:default")), t)
	ExpectTrue(((actualBVariants == "default:alternateB") || (actualBVariants == "alternateB:default")), t)
}

