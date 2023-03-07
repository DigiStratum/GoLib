package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// NewDependencyInjectable()
func TestThat_NewDependencyInjectable_ReturnsSomething_WithoutArguments(t *testing.T) {
	// Setup
	var sut DependencyInjectableIfc

	// Test
	sut = NewDependencyInjectable()

	// Verify
	ExpectNonNil(sut, t)
}

// IsStarted()
func TestThat_NewDependencyInjectable_IsStarted_ReturnsFalse_BeforeStarted(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.IsStarted()

	// Verify
	ExpectFalse(actual, t)
}

// Start()
func TestThat_NewDependencyInjectable_Start_ReturnsNoError_WhenNoRequiredDeps(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	err := sut.Start()
	actual := sut.IsStarted()

	// Verify
	ExpectNoError(err, t)
	ExpectTrue(actual, t)
}

func TestThat_NewDependencyInjectable_Start_ReturnsNoError_WhenDepsOptional(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
	)

	// Test
	err := sut.Start()
	actual := sut.IsStarted()

	// Verify
	ExpectNoError(err, t)
	ExpectTrue(actual, t)
}

func TestThat_NewDependencyInjectable_Start_ReturnsError_WhenMissingRequiredDeps(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("requireddep").SetRequired(),
	)

	// Test
	err := sut.Start()
	actual := sut.IsStarted()

	// Verify
	ExpectError(err, t)
	ExpectFalse(actual, t)
}

func TestThat_NewDependencyInjectable_Start_ReturnsNoError_WhenRequiredDepsInjected(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("requireddep").SetRequired(),
	)
	var ifc interface{}

	// Test
	sut.InjectDependencies(
		NewDependencyInstance("requireddep", ifc),
	)
	err := sut.Start()
	actual := sut.IsStarted()

	// Verify
	ExpectNoError(err, t)
	ExpectTrue(actual, t)
}

