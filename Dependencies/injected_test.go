package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// NewDependencyInjected(declaredDependencies DependenciesIfc)
func TestThat_NewDependencyInjected_ReturnsSomething_WhenGivenNil(t *testing.T) {
	// Setup
	var sut DependencyInjectedIfc

	// Test
	sut = NewDependencyInjected(nil)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_NewDependencyInjected_InterfaceAssertions_Pass(t *testing.T) {
	// Setup
	var sut DependencyInjectedIfc
	sut = NewDependencyInjected(NewDependencies())

	// Test
	_, discoveryIfcOk := sut.(DependencyDiscoveryIfc)
	_, injectableIfcOk := sut.(DependencyInjectableIfc)
	_, readableIfcOk := sut.(readableDependenciesIfc)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(discoveryIfcOk && injectableIfcOk && readableIfcOk, t)
}

// InjectDependencies(depinst ...DependencyInstanceIfc) error
func TestThat_DependencyInjected_InjectDependencies_ReturnsNoError_ForEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.InjectDependencies()

	// Verify
	ExpectNoError(actual, t)
}

// GetDeclaredDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetDeclaredDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.GetDeclaredDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(actual.GetVariants()), t)
}

// GetDeclaredDependencies() readableDependenciesIfc
// GetRequiredDependencies() readableDependenciesIfc
// GetOptionalDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetDeclaredDependencies_ReturnsExpectedSet(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME).SetVariant(DEP_VARIANT)
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)

	// Test
	actualDec := sut.GetDeclaredDependencies()
	actualReq := sut.GetRequiredDependencies()
	actualOpt := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actualDec, t)
	ExpectNonNil(actualReq, t)
	ExpectNonNil(actualOpt, t)

	// One declared dependency with matching name and variant, not required
	ExpectInt(1, len(actualDec.GetVariants()), t)
	actualDep := actualDec.GetVariant(DEP_NAME, DEP_VARIANT)
	ExpectString(DEP_NAME, actualDep.GetName(), t)
	ExpectString(DEP_VARIANT, actualDep.GetVariant(), t)
	ExpectFalse(actualDep.IsRequired(), t)

	// Zero required dependencies
	ExpectNonNil(actualReq, t)
	ExpectInt(0, len(actualReq.GetVariants()), t)

	// One optional dependency
	ExpectNonNil(actualOpt, t)
	ExpectInt(1, len(actualOpt.GetVariants()), t)
}

// GetRequiredDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetRequiredDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.GetRequiredDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(actual.GetVariants()), t)
}

// GetRequiredDependencies() readableDependenciesIfc
// GetOptionalDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetRequiredDependencies_ReturnsEmptySet_ForOptionalDeps(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME).SetVariant(DEP_VARIANT)
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)

	// Test
	actualReq := sut.GetRequiredDependencies()
	actualOpt := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actualReq, t)
	ExpectInt(0, len(actualReq.GetVariants()), t)

	ExpectNonNil(actualOpt, t)
	ExpectInt(1, len(actualOpt.GetVariants()), t)
}

// GetRequiredDependencies() readableDependenciesIfc
// GetOptionalDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetRequiredDependencies_ReturnsExpectedSet_ForRequiredDeps(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME).SetVariant(DEP_VARIANT).SetRequired()
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)

	// Test
	actualReq := sut.GetRequiredDependencies()
	actualOpt := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actualReq, t)
	ExpectInt(1, len(actualReq.GetVariants()), t)
	actualDep := actualReq.GetVariant(expectedDep.GetName(), expectedDep.GetVariant())
	ExpectString(DEP_NAME, actualDep.GetName(), t)
	ExpectString(DEP_VARIANT, actualDep.GetVariant(), t)
	ExpectTrue(actualDep.IsRequired(), t)

	ExpectNonNil(actualOpt, t)
	ExpectInt(0, len(actualOpt.GetVariants()), t)
}

// GetOptionalDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetOptionalDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actualOpt := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actualOpt, t)
	ExpectInt(0, len(actualOpt.GetVariants()), t)
}

// The ones below are for an instance checking what has been injected:

// GetMissingDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetMissingDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actualMiss := sut.GetMissingDependencies()

	// Verify
	ExpectNonNil(actualMiss, t)
	ExpectInt(0, len(actualMiss.GetVariants()), t)
}

// GetInjectedDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
func TestThat_DependencyInjected_GetInjectedDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actualInj := sut.GetInjectedDependencies()

	// Verify
	ExpectNonNil(actualInj, t)
	ExpectInt(0, len(actualInj.GetVariants()), t)
}

// GetInjectedDependencies() readableDependenciesIfc
// GetVariants() map[string][]string
// GetInstance(name string) interface{}
func TestThat_DependencyInjected_InjectDependencies_InjectsOptionalDependency(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME).SetRequired()
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)
	depInst := NewDependencyInstance(
		expectedDep.GetName(),
		NewDependencies(),		// Arbitrary interface that we already have
	)

	// Test
	actualMiss := sut.GetMissingDependencies()
	actualErr := sut.InjectDependencies(depInst)
	actualInj := sut.GetInjectedDependencies()
	actualDepInst := sut.GetInstance(expectedDep.GetName())

	// Verify
	ExpectInt(1, len(actualMiss.GetVariants()), t)

	ExpectNoError(actualErr, t)

	ExpectNonNil(actualInj, t)
	ExpectInt(1, len(actualInj.GetVariants()), t)
	actualDep := actualInj.Get(expectedDep.GetName())
	ExpectNonNil(actualDep, t)
	ExpectString(DEP_NAME, actualDep.GetName(), t)
	ExpectString(DEP_VARIANT_DEFAULT, actualDep.GetVariant(), t)

	ExpectTrue(sut.Has(expectedDep.GetName()), t)
	ExpectTrue(sut.HasAllRequiredDependencies(), t)

	ExpectNonNil(actualDepInst, t)
	_, ok := actualDepInst.(DependenciesIfc)
	ExpectTrue(ok, t)
}

// HasAllRequiredDependencies() bool
func TestThat_DependencyInjected_HasRequiredDependencies_ReturnsTrue_ForEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.HasAllRequiredDependencies()

	// Verify
	ExpectTrue(actual, t)
}

// ValidateRequiredDependencies() error
func TestThat_DependencyInjected_ValidateRequiredDependencies_ReturnsNoError_ForNothingMissing(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.ValidateRequiredDependencies()

	// Verify
	ExpectNoError(actual, t)
}

func TestThat_DependencyInjected_ValidateRequiredDependencies_ReturnsError_ForMissingRequiredDependency(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(
			NewDependency(DEP_NAME).SetRequired(),
		),
	)

	// Test
	actual := sut.ValidateRequiredDependencies()

	// Verify
	ExpectError(actual, t)
}

// GetInstance(name string) interface{}
func TestThat_DependencyInjected_GetInstance_ReturnsNil(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(
			NewDependency(DEP_NAME),
		),
	)

	// Test
	actual := sut.GetInstance(DEP_NAME)

	// Verify
	ExpectNil(actual, t)
}

func TestThat_DependencyInjected_GetInstance_ReturnsDefaultVariant(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(
			NewDependency(DEP_NAME).SetRequired(),
		),
	)
	sut.InjectDependencies(
		NewDependencyInstance(
			DEP_NAME,
			NewDependencies(),		// Arbitrary interface that we already have
		),
	)

	// Test
	actual := sut.GetInstance(DEP_NAME)

	// Verify
	ExpectNonNil(actual, t)
}

func TestThat_DependencyInjected_GetInstance_ReturnsNonDefaultVariant(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(
			NewDependency(DEP_NAME).SetRequired(),
		),
	)
	sut.InjectDependencies(
		NewDependencyInstance(
			DEP_NAME,
			NewDependencies(),		// Arbitrary interface that we already have
		).SetVariant(DEP_VARIANT_DEFAULT),
	)

	// Test
	actual := sut.GetInstance(DEP_NAME)

	// Verify
	ExpectNonNil(actual, t)
}

