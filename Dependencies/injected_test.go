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

func TestThat_DependencyInjected_IsValid_ReturnsFalse(t *testing.T) {
	var optionalNames, requiredNames []string
	var deps *Dependencies
	var sut *DependencyInjected

	// Setup / Test : One required, nothing optional/injected
	deps = NewDependencies()
	sut = NewDependencyInjected(deps)
	requiredNames = []string{"one"}
	sut.SetRequired(&requiredNames)
	res1 := sut.IsValid()
	ExpectFalse(res1, t)

	// Setup / Test : Two required, one required injected, nothing optional
	deps = NewDependencies()
	deps.Set("one", "sillystring!")
	sut = NewDependencyInjected(deps)
	requiredNames = []string{"one", "two"}
	sut.SetRequired(&requiredNames)
	res2 := sut.IsValid()
	ExpectFalse(res2, t)

	// Setup / Test : Nothing required, one optional, one invalid optional injected
	deps = NewDependencies()
	deps.Set("two", "sillystring!")
	sut = NewDependencyInjected(deps)
	optionalNames = []string{"one"}
	sut.SetOptional(&optionalNames)
	res3 := sut.IsValid()
	ExpectFalse(res3, t)

	// Setup / Test : Two required, one optional, one required injected, one invalid optional injected
	deps = NewDependencies()
	deps.Set("one", "sillystring!")
	deps.Set("four", "sillystring!")
	sut = NewDependencyInjected(deps)
	requiredNames = []string{"one", "two"}
	sut.SetRequired(&requiredNames)
	optionalNames = []string{"three"}
	sut.SetOptional(&optionalNames)
	res4 := sut.IsValid()
	ExpectFalse(res4, t)
}

func TestThat_DependencyInjected_IsValid_ReturnsTrue(t *testing.T) {
	var optionalNames, requiredNames []string
	var deps *Dependencies
	var sut *DependencyInjected

	// Setup / Test : Nothing required/optional/injected
	deps = NewDependencies()
	sut = NewDependencyInjected(deps)
	res1 := sut.IsValid()
	ExpectTrue(res1, t)

	// Setup / Test : Nothing required/injected, two optional
	deps = NewDependencies()
	sut = NewDependencyInjected(deps)
	optionalNames = []string{"one", "two"}
	sut.SetOptional(&optionalNames)
	res2 := sut.IsValid()
	ExpectTrue(res2, t)

	// Setup / Test : Nothing required, two optional, one optional injected
	deps = NewDependencies()
	deps.Set("one", "sillystring!")
	sut = NewDependencyInjected(deps)
	optionalNames = []string{"one", "two"}
	sut.SetOptional(&optionalNames)
	res3 := sut.IsValid()
	ExpectTrue(res3, t)

	// Setup / Test : Nothing optional, one required, one required injected
	deps = NewDependencies()
	deps.Set("one", "sillystring!")
	sut = NewDependencyInjected(deps)
	optionalNames = []string{"one"}
	sut.SetOptional(&optionalNames)
	res4 := sut.IsValid()
	ExpectTrue(res4, t)

	// Setup / Test : One optional, one required, one optional+required injected
	deps = NewDependencies()
	deps.Set("one", "sillystring!")
	deps.Set("two", "sillystring!")
	sut = NewDependencyInjected(deps)
	optionalNames = []string{"one"}
	sut.SetOptional(&optionalNames)
	requiredNames = []string{"two"}
	sut.SetRequired(&requiredNames)
	res5 := sut.IsValid()
	ExpectTrue(res5, t)
}
