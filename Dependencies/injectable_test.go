package dependencies

import(
	"fmt"
	"strings"
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
	if ! ExpectNonNil(sut, t) { return }
}

// IsStarted()
func TestThat_NewDependencyInjectable_IsStarted_ReturnsFalse_BeforeStarted(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.IsStarted()

	// Verify
	if ! ExpectFalse(actual, t) { return }
}

// Start()
func TestThat_NewDependencyInjectable_Start_ReturnsNoError_WhenNoRequiredDeps(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	err := sut.Start()
	actual := sut.IsStarted()

	// Verify
	if ! ExpectNoError(err, t) { return }
	if ! ExpectTrue(actual, t) { return }
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
	if ! ExpectNoError(err, t) { return }
	if ! ExpectTrue(actual, t) { return }
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
	if ! ExpectError(err, t) { return }
	if ! ExpectFalse(actual, t) { return }
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
	if ! ExpectNoError(err, t) { return }
	if ! ExpectTrue(actual, t) { return }
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
	if ! ExpectError(err, t) { return }
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
	if ! ExpectNoError(err, t) { return }
}

// GetInstance(name string) interface{}
func TestThat_NewDependencyInjectable_GetInstance_ReturnsNil_ForInvalidDependencyName(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetInstance("bogusdep")

	// Verify
	if ! ExpectNil(actual, t) { return }
}

func TestThat_NewDependencyInjectable_GetInstance_ReturnsNonNil_ForValidDependencyName(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	sut.InjectDependencies(
		NewDependencyInstance("depname", sut), // sut is as good as any other interface to use here
	)
	actual := sut.GetInstance("depname")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
}

func TestThat_NewDependencyInjectable_GetInstance_ReturnsNonNil_ForAnyVariant(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	sut.InjectDependencies(
		NewDependencyInstance("depname", sut).SetVariant("goodvariant"),
	)
	actual := sut.GetInstance("depname")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
}

// GetInstanceVariant(name, variant string) interface{}
func TestThat_NewDependencyInjectable_GetInstanceVariant_ReturnsNil_ForInvalidVariant(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	sut.InjectDependencies(
		NewDependencyInstance("depname", sut),
	)
	actual := sut.GetInstanceVariant("depname", "badvariant")

	// Verify
	if ! ExpectNil(actual, t) { return }
}

func TestThat_NewDependencyInjectable_GetInstanceVariant_ReturnsNonNil_ForValidVariant(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	sut.InjectDependencies(
		NewDependencyInstance("depname", sut).SetVariant("goodvariant"),
	)
	actual := sut.GetInstanceVariant("depname", "goodvariant")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
}


// HasAllRequiredDependencies() bool
func TestThat_NewDependencyInjectable_HasAllRequiredDependencies_ReturnsTrue_ForNoDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.HasAllRequiredDependencies()

	// Verify
	if ! ExpectTrue(actual, t) { return }
}

func TestThat_NewDependencyInjectable_HasAllRequiredDependencies_ReturnsTrue_ForOptionalDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
	)

	// Test
	actual := sut.HasAllRequiredDependencies()

	// Verify
	if ! ExpectTrue(actual, t) { return }
}

func TestThat_NewDependencyInjectable_HasAllRequiredDependencies_ReturnsFalse_ForMissingRequiredDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("requireddep").SetRequired(),
	)

	// Test
	actual := sut.HasAllRequiredDependencies()

	// Verify
	if ! ExpectFalse(actual, t) { return }
}

func TestThat_NewDependencyInjectable_HasAllRequiredDependencies_ReturnsTrue_ForInjectedRequiredDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("requireddep").SetRequired(),
	)
	sut.InjectDependencies(
		NewDependencyInstance("requireddep", sut),
	)

	// Test
	actual := sut.HasAllRequiredDependencies()

	// Verify
	if ! ExpectTrue(actual, t) { return }
}

// GetDeclaredDependencies() DependenciesIfc
func TestThat_NewDependencyInjectable_GetDeclaredDependencies_ReturnsEmpty_ForNoDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetDeclaredDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.GetAllVariants()), t) { return }
}

func TestThat_NewDependencyInjectable_GetDeclaredDependencies_ReturnsExpectedSet_ForDeclaredDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
		NewDependency("optionaldep").SetVariant("vname"),
		NewDependency("requireddep").SetRequired(),
		NewDependency("requireddep").SetRequired().SetVariant("vname"),
	)

	// Test & Verify
	actual := sut.GetDeclaredDependencies()
	if ! ExpectNonNil(actual, t) { return }

	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(2, len(actualVariants), t) { return }

	actualOptional, okOptional := actualVariants["optionaldep"]
	if ! ExpectTrue(okOptional, t) { return }
	if ! ExpectInt(2, len(actualOptional), t) { return }
	optionalVariants := strings.Join(actualOptional, ",")
	optionalVariantsOk := (optionalVariants == DEP_VARIANT_DEFAULT+",vname") || (optionalVariants == "vname,"+DEP_VARIANT_DEFAULT)
	if ! ExpectTrue(optionalVariantsOk, t) { return }

	actualRequired, okRequired:= actualVariants["requireddep"]
	if ! ExpectTrue(okRequired, t) { return }
	if ! ExpectInt(2, len(actualRequired), t) { return }
	requiredVariants := strings.Join(actualRequired, ",")
	requiredVariantsOk := (requiredVariants == DEP_VARIANT_DEFAULT+",vname") || (requiredVariants == "vname,"+DEP_VARIANT_DEFAULT)
	if ! ExpectTrue(requiredVariantsOk, t) { return }
}

// GetRequiredDependencies() DependenciesIfc
func TestThat_NewDependencyInjectable_GetRequiredDependencies_ReturnsEmpty_ForNoDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetRequiredDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.GetAllVariants()), t) { return }
}

func TestThat_NewDependencyInjectable_GetReuiredDependencies_ReturnsExpectedSet_ForRequiredDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
		NewDependency("optionaldep").SetVariant("vname"),
		NewDependency("requireddep").SetRequired(),
		NewDependency("requireddep").SetRequired().SetVariant("vname"),
	)

	// Test & Verify
	actual := sut.GetRequiredDependencies()
	if ! ExpectNonNil(actual, t) { return }

	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(1, len(actualVariants), t) { return }

	actualRequired, okRequired:= actualVariants["requireddep"]
	if ! ExpectTrue(okRequired, t) { return }
	if ! ExpectInt(2, len(actualRequired), t) { return }
	requiredVariants := strings.Join(actualRequired, ",")
	requiredVariantsOk := (requiredVariants == DEP_VARIANT_DEFAULT+",vname") || (requiredVariants == "vname,"+DEP_VARIANT_DEFAULT)
	if ! ExpectTrue(requiredVariantsOk, t) { return }
}


// GetOptionalDependencies() DependenciesIfc
func TestThat_NewDependencyInjectable_GetOptionalDependencies_ReturnsEmpty_ForNoDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetOptionalDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.GetAllVariants()), t) { return }
}

func TestThat_NewDependencyInjectable_GetOptionalDependencies_ReturnsExpectedSet_ForOptionalDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
		NewDependency("optionaldep").SetVariant("vname"),
		NewDependency("requireddep").SetRequired(),
		NewDependency("requireddep").SetRequired().SetVariant("vname"),
	)

	// Test & Verify
	actual := sut.GetOptionalDependencies()
	if ! ExpectNonNil(actual, t) { return }

	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(1, len(actualVariants), t) { return }

	actualOptional, okOptional := actualVariants["optionaldep"]
	if ! ExpectTrue(okOptional, t) { return }
	if ! ExpectInt(2, len(actualOptional), t) { return }
	optionalVariants := strings.Join(actualOptional, ",")
	optionalVariantsOk := (optionalVariants == DEP_VARIANT_DEFAULT+",vname") || (optionalVariants == "vname,"+DEP_VARIANT_DEFAULT)
	if ! ExpectTrue(optionalVariantsOk, t) { return }
}

// GetInjectedDependencies() DependenciesIfc

func TestThat_NewDependencyInjectable_GetInjectedDependencies_ReturnsEmpty_ForNoDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetInjectedDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.GetAllVariants()), t) { return }
}

func TestThat_NewDependencyInjectable_GetInjectedDependencies_ReturnsExpectedSet_ForInjectedDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()
	sut.InjectDependencies(
		NewDependencyInstance("requireddep", sut),
		NewDependencyInstance("requireddep", sut).SetVariant("vname"),
	)

	// Test
	actual := sut.GetInjectedDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(1, len(actualVariants), t) { return }

	actualInjected, okInjected:= actualVariants["requireddep"]
	if ! ExpectTrue(okInjected, t) { return }
	if ! ExpectInt(2, len(actualInjected), t) { return }
	injectedVariants := strings.Join(actualInjected, ",")
	injectedVariantsOk := (injectedVariants == DEP_VARIANT_DEFAULT+",vname") || (injectedVariants == "vname,"+DEP_VARIANT_DEFAULT)
	if ! ExpectTrue(injectedVariantsOk, t) { return }
}


// GetMissingDependencies() DependenciesIfc
func TestThat_NewDependencyInjectable_GetMissingDependencies_ReturnsEmpty_ForNoDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetMissingDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.GetAllVariants()), t) { return }
}

func TestThat_NewDependencyInjectable_GetMissingDependencies_ReturnsEmptySet_ForOptionalDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
		NewDependency("optionaldep").SetVariant("vname"),
	)

	// Test
	actual := sut.GetMissingDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(0, len(actualVariants), t) { return }
}

func TestThat_NewDependencyInjectable_GetMissingDependencies_ReturnsExpectedSet_ForRequiredDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
		NewDependency("optionaldep").SetVariant("vname"),
		NewDependency("requireddep").SetRequired(),
		NewDependency("requireddep").SetRequired().SetVariant("vname"),
	)

	// Test
	actual := sut.GetMissingDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(1, len(actualVariants), t) { return }

	actualMissing, okMissing:= actualVariants["requireddep"]
	if ! ExpectTrue(okMissing, t) { return }
	if ! ExpectInt(2, len(actualMissing), t) { return }
	missingVariants := strings.Join(actualMissing, ",")
	missingVariantsOk := (missingVariants == DEP_VARIANT_DEFAULT+",vname") || (missingVariants == "vname,"+DEP_VARIANT_DEFAULT)
	if ! ExpectTrue(missingVariantsOk, t) { return }
}

func TestThat_NewDependencyInjectable_GetMissingDependencies_ReturnsEmptySet_ForInjectedDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
		NewDependency("optionaldep").SetVariant("vname"),
		NewDependency("requireddep").SetRequired(),
		NewDependency("requireddep").SetRequired().SetVariant("vname"),
	)
	sut.InjectDependencies(
		NewDependencyInstance("requireddep", sut),
		NewDependencyInstance("requireddep", sut).SetVariant("vname"),
	)

	// Test
	actual := sut.GetMissingDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(0, len(actualVariants), t) { return }
}

// GetUnknownDependencies() DependenciesIfc
func TestThat_NewDependencyInjectable_GetUnknownDependencies_ReturnsEmpty_ForNoDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable()

	// Test
	actual := sut.GetUnknownDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, len(actual.GetAllVariants()), t) { return }
}

func TestThat_NewDependencyInjectable_GetUnknownDependencies_ReturnsExpectedSet_ForUnknownDependencies(t *testing.T) {
	// Setup
	sut := NewDependencyInjectable(
		NewDependency("optionaldep"),
		NewDependency("optionaldep").SetVariant("vname"),
		NewDependency("requireddep").SetRequired(),
		NewDependency("requireddep").SetRequired().SetVariant("vname"),
	)
	sut.InjectDependencies(
		NewDependencyInstance("unknowndep", sut),
		NewDependencyInstance("unknowndep", sut).SetVariant("vname"),
		NewDependencyInstance("requireddep", sut),
		NewDependencyInstance("requireddep", sut).SetVariant("vname"),
		NewDependencyInstance("optionaldep", sut),
		NewDependencyInstance("optionaldep", sut).SetVariant("vname"),
	)

	// Test
	actual := sut.GetUnknownDependencies()

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	actualVariants := actual.GetAllVariants()
	if ! ExpectInt(1, len(actualVariants), t) { return }

	actualUnknown, okUnknown:= actualVariants["unknowndep"]
	if ! ExpectTrue(okUnknown, t) { return }
	if ! ExpectInt(2, len(actualUnknown), t) { return }
	unknownVariants := strings.Join(actualUnknown, ",")
	unknownVariantsOk := (unknownVariants == DEP_VARIANT_DEFAULT+",vname") || (unknownVariants == "vname,"+DEP_VARIANT_DEFAULT)
	if ! ExpectTrue(unknownVariantsOk, t) { return }
}

