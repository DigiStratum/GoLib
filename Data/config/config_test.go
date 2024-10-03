package config

/*

TODO:
 * Add test coverage for max depth by setting it low and then hitting it and verifying that the
   number of substitutions is fewer than they would have been with a higher max

*/

import(
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
	actual, actualNum := sut.DereferenceString(expected)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, actualNum, t) { return }
	if ! ExpectString(expected, *actual, t) { return }
}

func TestThat_Config_DereferenceString_Returns_Original_String_with_broken_selector(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Test
	actual, actualNum := sut.DereferenceString("Howdy %name!")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, actualNum, t) { return }
	if ! ExpectString("Howdy %name!", *actual, t) { return }
}

func TestThat_Config_DereferenceString_Returns_Original_String_with_missing_selector(t *testing.T) {
	// Setup
	sut := NewConfig()

	// Test
	actual, actualNum := sut.DereferenceString("Howdy %name%!")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(0, actualNum, t) { return }
	if ! ExpectString("Howdy %name%!", *actual, t) { return }
}

func TestThat_Config_DereferenceString_Returns_String_with_object_property_selector_replaced(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareObject().SetObjectProperty("name", data.NewString("Doody"))

	// Test
	actual, actualNum := sut.DereferenceString("Howdy %name%!")

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectInt(1, actualNum, t) { return }
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
	actualValue := sut.Select("[3]")

	// Verify
	if ! ExpectInt(2, actual, t) { return }
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

func TestThat_Config_Dereference_iterates_for_object_few_reference_configs(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareObject().
		SetObjectProperty("numerics", data.NewObject().
			SetObjectProperty("formatter", data.NewString("%numerics.lucky-number% and %numerics.pi%: %numerics.truth%")).
			SetObjectProperty("lucky-number", data.NewInteger(333)).
			SetObjectProperty("pi", data.NewFloat(3.14159)).
			SetObjectProperty("truth", data.NewBoolean(true)),
		).
		SetObjectProperty("greeting", data.NewString("Greetings, %[0]% %invalid%%[2]%!"))

	ref := NewConfig()
	ref.PrepareArray().
		AppendArrayValue(data.NewString("Earthling - %[1]%")).
		AppendArrayValue(data.NewString("your lucky numbers are %numerics.formatter%"))

	// Test
	actual := sut.Dereference(ref)
	actualValue := sut.GetObjectProperty("greeting")

	// Verify
	if ! ExpectInt(6, actual, t) { return }
	if ! ExpectNonNil(actualValue, t) { return }
	if ! ExpectString("Greetings, Earthling - your lucky numbers are 333 and 3.14159: true %invalid%%[2]%!", actualValue.ToString(), t) { return }
}

func TestThat_Config_MergeConfig_returns_original_for_nil(t *testing.T) {
	// Setup
	sut := NewConfig()
	sut.PrepareObject().
		SetObjectProperty("prop1", data.NewString("Soup"))

	// Test
	actual := sut.MergeConfig(nil)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectTrue(actual.IsObject(), t) { return }
	if ! ExpectTrue(actual.HasObjectProperty("prop1"), t) { return }
	if ! ExpectInt(1, len(actual.GetObjectProperties()), t) { return }
}

func TestThat_Config_MergeConfig_returns_original_for_array_vs_object_datavalue(t *testing.T) {
	// Setup
	expectedName := "prop1"
	expectedValue := "Soup"
	sut := NewConfig()
	sut.PrepareObject().
		SetObjectProperty(expectedName, data.NewString(expectedValue))
	mergeCfg := NewConfig()
	mergeCfg.PrepareArray().
		AppendArrayValue(data.NewString("bogus"))

	// Test
	actual := sut.MergeConfig(mergeCfg)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectError(actual.GetError(), t) { return }
	if ! ExpectTrue(actual.IsObject(), t) { return }
	if ! ExpectTrue(actual.HasObjectProperty(expectedName), t) { return }
	if ! ExpectTrue(actual.GetObjectProperty(expectedName).IsString(), t) { return }
	if ! ExpectString(expectedValue, actual.GetObjectProperty(expectedName).GetString(), t) { return }
	if ! ExpectInt(1, len(actual.GetObjectProperties()), t) { return }
}

func TestThat_Config_MergeConfig_returns_original_for_object_vs_array_datavalue(t *testing.T) {
	// Setup
	expectedValue := "Toast"
	sut := NewConfig()
	sut.PrepareArray().
		AppendArrayValue(data.NewString(expectedValue))
	mergeCfg := NewConfig()
	mergeCfg.PrepareObject().
		SetObjectProperty("bogusprop", data.NewString("Salad"))

	// Test
	actual := sut.MergeConfig(mergeCfg)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectError(actual.GetError(), t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
	if ! ExpectInt(1, actual.GetArraySize(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(0).IsString(), t) { return }
	if ! ExpectString(expectedValue, actual.GetArrayValue(0).GetString(), t) { return }
}

func TestThat_Config_MergeConfig_merges_objects(t *testing.T) {
	// Setup
	expectedName := "prop2"
	expectedValue := "Salad"
	sut := NewConfig()
	sut.PrepareObject().
		SetObjectProperty("prop1", data.NewString("Soup"))
	mergeCfg := NewConfig()
	mergeCfg.PrepareObject().
		SetObjectProperty(expectedName, data.NewString(expectedValue))

	// Test
	actual := sut.MergeConfig(mergeCfg)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
	if ! ExpectTrue(actual.IsObject(), t) { return }
	if ! ExpectTrue(actual.HasObjectProperty(expectedName), t) { return }
	if ! ExpectTrue(actual.GetObjectProperty(expectedName).IsString(), t) { return }
	if ! ExpectString(expectedValue, actual.GetObjectProperty(expectedName).GetString(), t) { return }
	if ! ExpectInt(2, len(actual.GetObjectProperties()), t) { return }
}

func TestThat_Config_MergeConfig_merges_arrays(t *testing.T) {
	// Setup
	expectedValue := "Butter"
	sut := NewConfig()
	sut.PrepareArray().
		AppendArrayValue(data.NewString("Toast"))
	mergeCfg := NewConfig()
	mergeCfg.PrepareArray().
		AppendArrayValue(data.NewString(expectedValue))

	// Test
	actual := sut.MergeConfig(mergeCfg)

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectNoError(actual.GetError(), t) { return }
	if ! ExpectTrue(actual.IsArray(), t) { return }
	if ! ExpectInt(2, actual.GetArraySize(), t) { return }
	if ! ExpectTrue(actual.GetArrayValue(1).IsString(), t) { return }
	if ! ExpectString(expectedValue, actual.GetArrayValue(1).GetString(), t) { return }
}

