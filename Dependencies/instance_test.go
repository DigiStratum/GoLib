package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// NewDependencyInstance()
func TestThat_NewDependencyInstance_ReturnsSomething(t *testing.T) {
	// Setup
	var sut DependencyInstanceIfc

	// Test
	sut = NewDependencyInstance("depname", nil)

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}

// GetName() string
func TestThat_DependencyInstance_GetName_ReturnsName(t *testing.T) {
	// Setup
	sut := NewDependencyInstance("depname", nil)

	// Test
	actual := sut.GetName()

	// Verify
	if ! ExpectString("depname", actual, t) { return }
}

// GetVariant() string
func TestThat_DependencyInstance_GetVariant_ReturnsDefault(t *testing.T) {
	// Setup
	sut := NewDependencyInstance("depname", nil)

	// Test
	actual := sut.GetVariant()

	// Verify
	if ! ExpectString(DEP_VARIANT_DEFAULT, actual, t) { return }
}

// SetVariant(variant string) *dependencyInstance
func TestThat_DependencyInstance_GetVariant_ReturnsNameSet(t *testing.T) {
	// Setup
	sut := NewDependencyInstance("depname", nil).SetVariant("vname")

	// Test
	actual := sut.GetVariant()

	// Verify
	if ! ExpectString("vname", actual, t) { return }
}

// GetInstance() interface{}
func TestThat_DependencyInstance_GetInstance_ReturnsMatchingInstance(t *testing.T) {
	// Setup
	sut := NewDependencyInstance("depname", NewDependency("depname"))

	// Test
	actual := sut.GetInstance()
	_, okIfc := actual.(DependencyIfc)

	// Verify
	if ! ExpectTrue(okIfc, t) { return }
}

