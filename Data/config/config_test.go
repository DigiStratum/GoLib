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

