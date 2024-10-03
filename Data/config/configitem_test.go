package config

import(
	"testing"

	//"GoLib/Data"

	. "GoLib/Testing"
)

// Interface

func TestThat_ConfigItem_NewConfigItem_ReturnsInstance(t *testing.T) {
	// Setup
	var sut ConfigItemIfc = NewConfigItem("") // Verifies that result satisfies IFC

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}

// GetSelector()

func TestThat_ConfigItem_GetSelector_Returns_Selector(t *testing.T) {
	// Setup
	expected := "selector"
	sut := NewConfigItem(expected)

	// Verify
	if ! ExpectString(expected, sut.GetSelector(), t) { return }
}

// SetRequired() | IsRequired()

func TestThat_ConfigItem_GetRequired_Returns_false_by_default(t *testing.T) {
	// Setup
	sut := NewConfigItem("selector")

	// Verify
	if ! ExpectFalse(sut.IsRequired(), t) { return }
}

func TestThat_ConfigItem_SetRequired_Changes_Required_State(t *testing.T) {
	// Setup
	sut := NewConfigItem("selector")

	// Verify
	if ! ExpectTrue(sut.SetRequired().IsRequired(), t) { return }
}


