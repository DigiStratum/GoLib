package config

import(
	//"fmt"
	//"strings"
	"testing"

	"GoLib/Data"

	. "GoLib/Testing"
)

// Interface

func TestThat_Config_NewConfig_ReturnsInstance(t *testing.T) {
	// Setup
	var sut ConfigIfc = NewConfig() // Verifies that result satisfies IFC

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}

func TestThat_Config_DereferenceString_Returns_Original_String_without_selectors(t *testing.T) {
	// Setup
	sut := NewConfig()
	expected := "Howdy!"

	// Test
	actual := sut.DereferenceString(expected)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expected, *actual, t) { return }
}

func TestThat_Config_DereferenceString_Returns_Original_String_with_broken_selector(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Test
	actual := sut.DereferenceString("Howdy %name!")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString("Howdy %name!", *actual, t) { return }
}

func TestThat_Config_DereferenceString_Returns_Original_String_with_missing_selector(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Test
	actual := sut.DereferenceString("Howdy %name%!")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString("Howdy %name%!", *actual, t) { return }
}

func TestThat_Config_DereferenceString_Returns_String_with_object_property_selector_replaced(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareObject().SetObjectProperty("name", data.NewString("Doody"))

	// Test
	actual := sut.DereferenceString("Howdy %name%!")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString("Howdy Doody!", *actual, t) { return }
}

// Non-Structures

func TestThat_Config_Dereference_Does_nothing_for_non_structure_no_reference_config(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Test
	actual := sut.Dereference()

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

func TestThat_Config_Dereference_Does_nothing_for_non_structure_one_reference_config(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Test
	actual := sut.Dereference(NewConfig())

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

func TestThat_Config_Dereference_Does_nothing_for_non_structure_few_reference_configs(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Test
	actual := sut.Dereference(NewConfig(), NewConfig(), NewConfig())

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

// Arrays

func TestThat_Config_Dereference_Does_nothing_for_array_no_reference_config(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareArray()

	// Test
	actual := sut.Dereference()

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

func TestThat_Config_Dereference_Does_nothing_for_array_one_reference_config(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareArray()

	// Test
	actual := sut.Dereference(NewConfig())

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

func TestThat_Config_Dereference_Does_nothing_for_array_few_reference_configs(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareArray()

	// Test
	actual := sut.Dereference(NewConfig(), NewConfig(), NewConfig())

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

func TestThat_Config_Dereference_iterates_for_array_few_reference_configs(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareArray().
		AppendArrayValue(data.NewInteger(333)).
		AppendArrayValue(data.NewFloat(3.14159)).
		AppendArrayValue(data.NewBoolean(true)).
		AppendArrayValue(data.NewString("Greetings, %label%!"))

	ref := NewConfig()
	ref.PrepareObject().
		SetObjectProperty("label", data.NewString("Earthling number %[0]%"))


	// Test
	actual := sut.Dereference(ref)
	actualValue, err := sut.Select("[3]")

	// Verify
	if ! ExpectInt(2, actual, t) { return }
	if ! ExpectNoError(err, t) { return }
	if ! ExpectString("Greetings, Earthling number 333!", actualValue.ToString(), t) { return }
}

// Objects

func TestThat_Config_Dereference_Does_nothing_for_object_no_reference_config(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareObject()

	// Test
	actual := sut.Dereference()

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

func TestThat_Config_Dereference_Does_nothing_for_object_one_reference_config(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareObject()

	// Test
	actual := sut.Dereference(NewConfig())

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

func TestThat_Config_Dereference_Does_nothing_for_object_few_reference_configs(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareObject()

	// Test
	actual := sut.Dereference(NewConfig(), NewConfig(), NewConfig())

	// Verify
	if ! ExpectInt(0, actual, t) { return }
}

