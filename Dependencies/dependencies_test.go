package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

/*
type readableDependenciesIfc interface {
        // Get a dependency by name/variant
        Get(name string) *dependency
        GetVariant(name, variant string) *dependency

        // Check whether a dependency is in the set by name/variant
        Has(name string) bool
        HasVariant(name, variant string) bool

        // Get the list of currently set dependencies
        GetVariants() map[string][]string
}

type DependenciesIfc interface {
        // Embed all the readableDependenciesIfc requirements
        readableDependenciesIfc
        // Add a Dependency to the set
        Add(dep ...*dependency)
}

*/
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

func TestThat_Dependencies_Get_ReturnsDependency_From_Add(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	sut.Add(NewDependency("sampledep"))
	actual := sut.Get("sampledep")

	// Verify
	ExpectNonNil(actual, t)
}


/*
func TestThat_Dependencies_GetUniqueIds_IsEmpty_ForNewDependencies(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.GetUniqueIds()

	// Verify
	ExpectInt(0, len(*actual), t)
}

func TestThat_Dependencies_GetUniqueIds_HasExpectedOneDependency(t *testing.T) {
	// Setup
	sut := NewDependencies(
		NewDependency(DEP_NAME).SetVariant(DEP_VARIANT),
	)

	// Test
	actual := sut.GetUniqueIds()

	// Verify
	ExpectInt(1, len(*actual), t)
	ExpectTrue(sut.Has((*actual)[0]), t)
}

func TestThat_Dependencies_Add_AddsDependency(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	sut.Add(
		NewDependency(DEP_NAME).SetVariant(DEP_VARIANT),
	)
	actual := sut.GetUniqueIds()

	// Verify
	ExpectInt(1, len(*actual), t)
	ExpectTrue(sut.Has((*actual)[0]), t)
}
*/

func TestThat_Dependencies_Has_ReturnsFalse_ForMissingDependency(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.Has("bogusdep")

	// Verify
	ExpectFalse(actual, t)
}

