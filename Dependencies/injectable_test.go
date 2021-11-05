package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDependencyInjectable_ReturnsNothing_WhenGivenNil(t *testing.T) {
	// Test
	sut := NewDependencyInjectable(nil)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewDependencyInjectable_ReturnsSomething_WhenGivenDependencies(t *testing.T) {
	// Test
	deps := NewDependencies()
	sut := NewDependencyInjectable(deps)

	// Verify
	ExpectNonNil(sut, t)
}
