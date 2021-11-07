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

func TestThat_DependencyInjected_SetRequired_SetsRequired_And_ReturnsInstance(t *testing.T) {
	// Setup
	deps := NewDependencies()
	sut := NewDependencyInjected(deps)
	requiredNames := []string{"one", "two"}
	expected := len(requiredNames)

	// Test
	res := sut.SetRequired(&requiredNames)
	ExpectNonNil(res, t)
	actual := res.NumRequired()
	actualNames := sut.GetRequired()
	actualNameMap := make(map[string]bool)
	for _, actualName := range *actualNames {
		actualNameMap[actualName] = true
	}
	for _, requiredName := range requiredNames {
		_, ok := actualNameMap[requiredName]
		ExpectTrue(ok, t)
	}

	// Verify
	ExpectInt(expected, actual, t)
}

func TestThat_DependencyInjected_SetOptional_SetsOptional_And_ReturnsInstance(t *testing.T) {
	// Setup
	deps := NewDependencies()
	sut := NewDependencyInjected(deps)
	optionalNames := []string{"one", "two"}
	expected := len(optionalNames)

	// Test
	res := sut.SetOptional(&optionalNames)
	ExpectNonNil(res, t)
	actual := res.NumOptional()
	actualNames := sut.GetOptional()
	actualNameMap := make(map[string]bool)
	for _, actualName := range *actualNames {
		actualNameMap[actualName] = true
	}
	for _, optionalName := range optionalNames {
		_, ok := actualNameMap[optionalName]
		ExpectTrue(ok, t)
	}

	// Verify
	ExpectInt(expected, actual, t)
}

func TestThat_DependencyInjected_GetMissingRequiredDependencyNames_ReturnsNil_WhenNotInstantiated(t *testing.T) {
	// Setup
	var sut DependencyInjected

	// Test
	res := sut.GetMissingRequiredDependencyNames()

	// Verify
	ExpectNil(res, t)
}

func TestThat_DependencyInjected_GetMissingRequiredDependencyNames_ReturnsMissingNames(t *testing.T) {
	// Setup
	deps := NewDependencies()
	deps.Set("one", "sillystring!")
	sut := NewDependencyInjected(deps)
	requiredNames := []string{"one", "two"}
	sut.SetRequired(&requiredNames)
	expected := len(requiredNames)
	depNames := deps.GetNames()
	if nil != depNames { expected -= len(*depNames) }

	// Test
	missingNames := sut.GetMissingRequiredDependencyNames()
	ExpectNonNil(missingNames, t)
	actual := len(*missingNames)

	// Verify
	ExpectInt(expected, actual, t)
}


func TestThat_DependencyInjected_GetInvalidDependencyNames_ReturnsInvalidNames(t *testing.T) {
	// Setup
	deps := NewDependencies()
	deps.Set("three", "sillystring!")
	sut := NewDependencyInjected(deps)
	optionalNames := []string{"one", "two"}
	sut.SetOptional(&optionalNames)
	depNames := deps.GetNames()
	expected := len(*depNames)

	// Test
	invalidNames := sut.GetInvalidDependencyNames()
	ExpectNonNil(invalidNames, t)
	actual := len(*invalidNames)

	// Verify
	ExpectInt(expected, actual, t)
}
