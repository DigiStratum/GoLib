package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDependencyInjected_ReturnsNothing_WhenGivenNil(t *testing.T) {
	// Test
	sut := NewDependencyInjected(nil)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewDependencyInjected_ReturnsSomething_WhenGivenDependencies(t *testing.T) {
	// Setup
	deps := NewDependencies()

	// Test
	sut := NewDependencyInjected(deps)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_DependencyInjected_SetRequired_ReturnsSomething_WhenGivenDependencies(t *testing.T) {
	// Setup
	deps := NewDependencies()
	sut := NewDependencyInjected(deps)

	// Test

	// Verify
	ExpectNonNil(sut, t)
}
