package dependencies

import(
	"fmt"
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

// InjectDependencies(depinst ...DependencyInstanceIfc) error
func TestThat_NewDependencyInjectable_InjectDependencies_ReturnsError_WhenCaptureFuncReturnsError(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("requireddep").SetRequired().CaptureWith(
			func (v interface{}) error { return fmt.Errorf("capture error!") },
		),
	)
	var ifc interface{}

	// Test
	err := sut.InjectDependencies(
		NewDependencyInstance("requireddep", ifc),
	)

	// Verify
	ExpectError(err, t)
}

func TestThat_NewDependencyInjectable_InjectDependencies_ReturnsNoError_WhenCaptureFuncReturnsNoError(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("requireddep").SetRequired().CaptureWith(
			func (v interface{}) error { return nil },
		),
	)
	var ifc interface{}

	// Test
	err := sut.InjectDependencies(
		NewDependencyInstance("requireddep", ifc),
	)

	// Verify
	ExpectNoError(err, t)
}

// GetInstance(name string) interface{}
func TestThat_NewDependencyInjectable_GetInstance_ReturnsNil_ForInvalidDependency(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetInstance("bogusdep")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_NewDependencyInjectable_GetInstance_ReturnsNonNil_ForValidDependency(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("requireddep").SetRequired(),
	)

	// Test
	sut.InjectDependencies(
		NewDependencyInstance("requireddep", sut), // sut is as good as any other interface to use here
	)
	actual := sut.GetInstance("requireddep")

	// Verify
	ExpectNonNil(actual, t)
}

// TODO: Test more of these things:

// GetInstanceVariant(name, variant string) interface{}
// HasAllRequiredDependencies() bool
// GetDeclaredDependencies() DependenciesIfc
// GetRequiredDependencies() DependenciesIfc
// GetOptionalDependencies() DependenciesIfc
// GetInjectedDependencies() DependenciesIfc
// GetMissingDependencies() DependenciesIfc
// GetUnknownDependencies() DependenciesIfc

