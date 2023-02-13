package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDependencyInjected_ReturnsSomething_WhenGivenNil(t *testing.T) {
	// Test
	sut := NewDependencyInjected(nil)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewDependencyInjected_ReturnsSomething_WhenGivenDependencies(t *testing.T) {
	// Setup
	deps := NewDependencies()

	// Test
	sut := NewDependencyInjected(deps)

	// Verify
	ExpectNonNil(sut, t)
}


func TestThat_DependencyInjected_GetDeclaredDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.GetDeclaredDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(*(actual.GetUniqueIds())), t)
}

