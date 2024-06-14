package json

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Interface

func TestThat_JsonValue_NewJsonValue_ReturnsInstance(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(
		sut.IsNull() || sut.IsString() || sut.IsObject() || sut.IsArray() || sut.IsInteger() || sut.IsFloat() || sut.IsValid(),
		t,
	) { return }
}

// Validity

func TestThat_JsonValue_IsValid_Returns_false_for_new_value(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(sut.IsValid(), t) { return }
}

// Nulls

func TestThat_JsonValue_IsNull_Returns_false_for_new_value(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(sut.IsNull(), t) { return }
}

func TestThat_JsonValue_IsNull_Returns_true_after_setting_null(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()
	sut.SetNull()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectTrue(sut.IsNull(), t) { return }
}

// Strings

func TestThat_JsonValue_IsString_Returns_false_for_new_value(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectFalse(sut.IsString(), t) { return }
}

func TestThat_JsonValue_IsString_Returns_true_after_setting_string(t *testing.T) {
	// Setup
	var sut JsonValueIfc = NewJsonValue()
	expected := "hiyee!"
	sut.SetString(expected)

	// Verify
	if ! ExpectNonNil(sut, t) { return }
	if ! ExpectTrue(sut.IsString(), t) { return }
	if ! ExpectString(expected, sut.GetString(), t) { return }
}

// Objects


