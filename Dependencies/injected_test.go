package dependencies

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_NewDependencyInjected_ReturnsSomething_WhenGivenNil(t *testing.T) {
	// Setup
	var sut DependencyInjectedIfc

	// Test
	sut = NewDependencyInjected(nil)

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

func TestThat_DependencyInjected_GetDeclaredDependencies_ReturnsExpectedSet(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME, DEP_VARIANT, false)
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)

	// Test
	actual := sut.GetDeclaredDependencies()
	actualReq := sut.GetRequiredDependencies()
	actualOpt := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(1, len(*(actual.GetUniqueIds())), t)
	actualDep := actual.Get(expectedDep.GetUniqueId())
	ExpectString(DEP_NAME, actualDep.GetName(), t)
	ExpectString(DEP_VARIANT, actualDep.GetVariant(), t)
	ExpectFalse(actualDep.IsRequired(), t)

	ExpectNonNil(actualReq, t)
	ExpectInt(0, len(*(actualReq.GetUniqueIds())), t)

	ExpectNonNil(actualOpt, t)
	ExpectInt(1, len(*(actualOpt.GetUniqueIds())), t)
}

func TestThat_DependencyInjected_GetRequiredDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.GetRequiredDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(*(actual.GetUniqueIds())), t)
}

func TestThat_DependencyInjected_GetRequiredDependencies_ReturnsEmptySet_ForOptionalDeps(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME, DEP_VARIANT, false)
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)

	// Test
	actual := sut.GetRequiredDependencies()
	actualOpt := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(*(actual.GetUniqueIds())), t)

	ExpectNonNil(actualOpt, t)
	ExpectInt(1, len(*(actualOpt.GetUniqueIds())), t)
}

func TestThat_DependencyInjected_GetRequiredDependencies_ReturnsExpectedSet_ForRequiredDeps(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME, DEP_VARIANT, true)
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)

	// Test
	actual := sut.GetRequiredDependencies()
	actualOpt := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(1, len(*(actual.GetUniqueIds())), t)
	actualDep := actual.Get(expectedDep.GetUniqueId())
	ExpectString(DEP_NAME, actualDep.GetName(), t)
	ExpectString(DEP_VARIANT, actualDep.GetVariant(), t)
	ExpectTrue(actualDep.IsRequired(), t)

	ExpectNonNil(actualOpt, t)
	ExpectInt(0, len(*(actualOpt.GetUniqueIds())), t)
}

func TestThat_DependencyInjected_GetOptionalDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.GetOptionalDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(*(actual.GetUniqueIds())), t)
}

// The ones below are for an instance checking what has been injected:

func TestThat_DependencyInjected_GetMissingDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.GetMissingDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(*(actual.GetUniqueIds())), t)
}

func TestThat_DependencyInjected_GetInjectedDependencies_ReturnsEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.GetInjectedDependencies()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(*(actual.GetUniqueIds())), t)
}

func TestThat_DependencyInjected_HasRequiredDependencies_ReturnsTrue_ForEmptySet(t *testing.T) {
	// Setup
	sut := NewDependencyInjected(
		NewDependencies(),
	)

	// Test
	actual := sut.HasRequiredDependencies()

	// Verify
	ExpectTrue(actual, t)
}

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

func TestThat_DependencyInjected_InjectDependencies_InjectsOptionalDependency(t *testing.T) {
	// Setup
	expectedDep := NewDependency(DEP_NAME, DEP_VARIANT, true)
	sut := NewDependencyInjected(
		NewDependencies(
			expectedDep,
		),
	)
	depInst := NewDependencyInstance(
		expectedDep,
		NewDependencies(), // Arbitrary interface that we already have
	)

	// Test
	missing := sut.GetMissingDependencies()
	actual := sut.InjectDependencies(
		depInst,
	)
	actualInj := sut.GetInjectedDependencies()
	actualDepInst := sut.GetInstance(expectedDep.GetUniqueId())

	// Verify
	ExpectInt(1, len(*(missing.GetUniqueIds())), t)

	ExpectNoError(actual, t)

	ExpectNonNil(actualInj, t)
	ExpectInt(1, len(*(actualInj.GetUniqueIds())), t)
	actualDep := actualInj.Get(expectedDep.GetUniqueId())
	ExpectString(DEP_NAME, actualDep.GetName(), t)
	ExpectString(DEP_VARIANT, actualDep.GetVariant(), t)
	ExpectTrue(actualDep.IsRequired(), t)

	ExpectTrue(sut.Has(expectedDep.GetUniqueId()), t)
	ExpectTrue(sut.HasRequiredDependencies(), t)


	ExpectNonNil(actualDepInst, t)
	_, ok := actualDepInst.(DependenciesIfc)
	ExpectTrue(ok, t)
}

