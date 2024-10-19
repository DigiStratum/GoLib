package config

import(
	"fmt"
	"testing"

	"github.com/DigiStratum/GoLib/Data"

	. "github.com/DigiStratum/GoLib/Testing"
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
	sut := NewConfigItem("")

	// Verify
	if ! ExpectFalse(sut.IsRequired(), t) { return }
}

func TestThat_ConfigItem_SetRequired_Changes_Required_State(t *testing.T) {
	// Setup
	sut := NewConfigItem("")

	// Verify
	if ! ExpectTrue(sut.SetRequired().IsRequired(), t) { return }
}

// SetDefault

func TestThat_ConfigItem_SetDefault_Returns_self(t *testing.T) {
	// Setup
	expectedSelector := "selector"
	sut := NewConfigItem(expectedSelector)

	// Test
	actual := sut.SetDefault(data.NewNull())

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expectedSelector, actual.GetSelector(), t) { return }
}

// CanCapture

func TestThat_ConfigItem_CanCapture_Returns_false_by_default(t *testing.T) {
	// Setup
	sut := NewConfigItem("")

	// Verify
	if ! ExpectFalse(sut.CanCapture(), t) { return }
}

// CaptureWith

func TestThat_ConfigItem_CaptureWith_changes_CanCapture_Response_to_true(t *testing.T) {
	// Setup
	expectedSelector := "selector"
	sut := NewConfigItem(expectedSelector)

	// Test
	actual := sut.CaptureWith(func (dataValue data.DataValueIfc) error { return nil })

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expectedSelector, actual.GetSelector(), t) { return }
	if ! ExpectTrue(sut.CanCapture(), t) { return }
}

// Capture

func TestThat_ConfigItem_Capture_Returns_error_without_CaptureFunc(t *testing.T) {
	// Setup
	sut := NewConfigItem("")

	// Verify
	if ! ExpectError(sut.Capture(nil), t) { return }
}

func TestThat_ConfigItem_Capture_Passes_Through_CaptureFunc_Error(t *testing.T) {
	// Setup
	sut := NewConfigItem("")
	expectedError := "yep, error!"
	sut.CaptureWith(func (dataValue data.DataValueIfc) error {
		if nil == dataValue { return nil }
		return fmt.Errorf(expectedError)
	})

	// Test
	actualError := sut.Capture(data.NewNull())

	// Verify
	if ! ExpectError(actualError, t) { return }
	if ! ExpectString(expectedError, actualError.Error(), t) { return }
	if ! ExpectNoError(sut.Capture(nil), t) { return }
}

func TestThat_ConfigItem_Capture_when_required_returns_error_given_nil(t *testing.T) {
	// Setup
	sut := NewConfigItem("")
	sut.CaptureWith(func (dataValue data.DataValueIfc) error { return nil }).SetRequired()

	// Verify
	if ! ExpectError(sut.Capture(nil), t) { return }
	if ! ExpectNoError(sut.Capture(data.NewNull()), t) { return }
}

func TestThat_ConfigItem_Capture_when_has_default_captures_default_given_nil(t *testing.T) {
	// Setup
	expectedString1 := "got 1!"
	expectedString2 := "got 2!"
	var capturedValue data.DataValueIfc
	sut := NewConfigItem("").
		CaptureWith(func (dataValue data.DataValueIfc) error { capturedValue = dataValue; return nil }).
		SetDefault(data.NewString(expectedString1))

	// Test
	actualErr1 := sut.Capture(nil)

	// Verify
	if ! ExpectNoError(actualErr1, t) { return }
	if ! ExpectNonNil(capturedValue, t) { return }
	if ! ExpectTrue(capturedValue.IsString(), t) { return }
	if ! ExpectString(expectedString1, capturedValue.GetString(), t) { return }
	actualErr2 := sut.Capture(data.NewString(expectedString2))
	if ! ExpectNoError(actualErr2, t) { return }
	if ! ExpectNonNil(capturedValue, t) { return }
	if ! ExpectTrue(capturedValue.IsString(), t) { return }
	if ! ExpectString(expectedString2, capturedValue.GetString(), t) { return }
}

// CanValidate

func TestThat_ConfigItem_CanValidate_Returns_false_by_default(t *testing.T) {
	// Setup
	sut := NewConfigItem("")

	// Verify
	if ! ExpectFalse(sut.CanValidate(), t) { return }
}

// ValidateWith

func TestThat_ConfigItem_ValidateWith_changes_CanValidate_Response_to_true(t *testing.T) {
	// Setup
	expectedSelector := "selector"
	sut := NewConfigItem(expectedSelector)

	// Test
	actual := sut.ValidateWith(func (dataValue data.DataValueIfc) error { return nil })

	// Verify
	if ! ExpectNonNil(actual, t) { return }
	if ! ExpectString(expectedSelector, actual.GetSelector(), t) { return }
	if ! ExpectTrue(sut.CanValidate(), t) { return }
}

// Validate

func TestThat_ConfigItem_Validate_Returns_error_without_ValidateFunc(t *testing.T) {
	// Setup
	sut := NewConfigItem("")

	// Verify
	if ! ExpectError(sut.Validate(nil), t) { return }
}

func TestThat_ConfigItem_Validate_Passes_Through_ValidateFunc_Error(t *testing.T) {
	// Setup
	sut := NewConfigItem("")
	expectedError := "yep, error!"
	sut.ValidateWith(func (dataValue data.DataValueIfc) error {
		if nil == dataValue { return nil }
		return fmt.Errorf(expectedError)
	})

	// Test
	actualError := sut.Validate(data.NewNull())

	// Verify
	if ! ExpectError(actualError, t) { return }
	if ! ExpectString(expectedError, actualError.Error(), t) { return }
	if ! ExpectNoError(sut.Validate(nil), t) { return }
}

