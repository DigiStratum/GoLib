package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDependencies_ReturnsSomething(t *testing.T) {
	// Setup
	var sut DependenciesIfc

	// Test
	sut = NewDependencies()

	// Verify
	ExpectNonNil(sut, t)
}

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
		NewDependency(DEP_NAME, DEP_VARIANT, false),
	)

	// Test
	actual := sut.GetUniqueIds()

	// Verify
	ExpectInt(1, len(*actual), t)
	ExpectTrue(sut.Has((*actual)[0]), t)
}

func TestThat_Dependencies_Has_ReturnsFalse_ForMissingDependency(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	actual := sut.Has("bogusdep")

	// Verify
	ExpectFalse(actual, t)
}

func TestThat_Dependencies_Add_AddsDependency(t *testing.T) {
	// Setup
	sut := NewDependencies()

	// Test
	sut.Add(
		NewDependency(DEP_NAME, DEP_VARIANT, false),
	)
	actual := sut.GetUniqueIds()

	// Verify
	ExpectInt(1, len(*actual), t)
	ExpectTrue(sut.Has((*actual)[0]), t)
}

