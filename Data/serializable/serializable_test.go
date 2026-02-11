package serializable

import (
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// -------------------------------------------------------------------------------------------------
// Interface Implementation Tests
// -------------------------------------------------------------------------------------------------

// Test that our mock implements SerializableIfc
func TestThat_MockSerializable_ImplementsSerializableIfc(t *testing.T) {
	var _ SerializableIfc = (*MockSerializable)(nil)
}

// Test that our mock implements DeserializableIfc
func TestThat_MockDeserializable_ImplementsDeserializableIfc(t *testing.T) {
	var _ DeserializableIfc = (*MockDeserializable)(nil)
}

// Test that an object can implement both interfaces
func TestThat_MockSerializableDeserializable_ImplementsBothInterfaces(t *testing.T) {
	var _ SerializableIfc = (*MockSerializableDeserializable)(nil)
	var _ DeserializableIfc = (*MockSerializableDeserializable)(nil)
}

// -------------------------------------------------------------------------------------------------
// Mock Implementation Tests - SerializableIfc
// -------------------------------------------------------------------------------------------------

func TestThat_MockSerializable_Serialize_ReturnsSerializedString(t *testing.T) {
	// Setup
	sut := &MockSerializable{Data: "test data"}

	// Test
	result, err := sut.Serialize()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(result, t)
	ExpectString("serialized:test data", *result, t)
}

func TestThat_MockSerializable_Serialize_HandlesEmptyData(t *testing.T) {
	// Setup
	sut := &MockSerializable{Data: ""}

	// Test
	result, err := sut.Serialize()

	// Verify
	ExpectNoError(err, t)
	ExpectNonNil(result, t)
	ExpectString("serialized:", *result, t)
}

func TestThat_MockSerializable_Serialize_CanReturnError(t *testing.T) {
	// Setup
	sut := &MockSerializable{Data: "error", ShouldError: true}

	// Test
	result, err := sut.Serialize()

	// Verify
	ExpectError(err, t)
	ExpectNil(result, t)
}

// -------------------------------------------------------------------------------------------------
// Mock Implementation Tests - DeserializableIfc
// -------------------------------------------------------------------------------------------------

func TestThat_MockDeserializable_Deserialize_RestoresData(t *testing.T) {
	// Setup
	sut := &MockDeserializable{}
	serialized := "deserialized:test data"

	// Test
	err := sut.Deserialize(&serialized)

	// Verify
	ExpectNoError(err, t)
	ExpectString("test data", sut.Data, t)
}

func TestThat_MockDeserializable_Deserialize_HandlesEmptyData(t *testing.T) {
	// Setup
	sut := &MockDeserializable{}
	serialized := "deserialized:"

	// Test
	err := sut.Deserialize(&serialized)

	// Verify
	ExpectNoError(err, t)
	ExpectString("", sut.Data, t)
}

func TestThat_MockDeserializable_Deserialize_ReturnsError_WithNilData(t *testing.T) {
	// Setup
	sut := &MockDeserializable{}

	// Test
	err := sut.Deserialize(nil)

	// Verify
	ExpectError(err, t)
}

func TestThat_MockDeserializable_Deserialize_CanReturnError(t *testing.T) {
	// Setup
	sut := &MockDeserializable{ShouldError: true}
	serialized := "test"

	// Test
	err := sut.Deserialize(&serialized)

	// Verify
	ExpectError(err, t)
}

// -------------------------------------------------------------------------------------------------
// Mock Implementation Tests - Combined SerializableIfc and DeserializableIfc
// -------------------------------------------------------------------------------------------------

func TestThat_MockSerializableDeserializable_RoundTrip_PreservesData(t *testing.T) {
	// Setup
	original := &MockSerializableDeserializable{Data: "original data"}

	// Test - Serialize
	serialized, err := original.Serialize()
	if !ExpectNoError(err, t) { return }
	if !ExpectNonNil(serialized, t) { return }

	// Test - Deserialize into new instance
	restored := &MockSerializableDeserializable{}
	err = restored.Deserialize(serialized)
	if !ExpectNoError(err, t) { return }

	// Verify
	ExpectString(original.Data, restored.Data, t)
}

func TestThat_MockSerializableDeserializable_RoundTrip_WithEmptyData(t *testing.T) {
	// Setup
	original := &MockSerializableDeserializable{Data: ""}

	// Test - Serialize
	serialized, err := original.Serialize()
	if !ExpectNoError(err, t) { return }
	if !ExpectNonNil(serialized, t) { return }

	// Test - Deserialize into new instance
	restored := &MockSerializableDeserializable{}
	err = restored.Deserialize(serialized)
	if !ExpectNoError(err, t) { return }

	// Verify
	ExpectString("", restored.Data, t)
}

// -------------------------------------------------------------------------------------------------
// Mock Implementations
// -------------------------------------------------------------------------------------------------

// MockSerializable implements SerializableIfc for testing
type MockSerializable struct {
	Data        string
	ShouldError bool
}

func (m *MockSerializable) Serialize() (*string, error) {
	if m.ShouldError {
		return nil, &MockError{"serialization error"}
	}
	result := "serialized:" + m.Data
	return &result, nil
}

// MockDeserializable implements DeserializableIfc for testing
type MockDeserializable struct {
	Data        string
	ShouldError bool
}

func (m *MockDeserializable) Deserialize(data *string) error {
	if nil == data {
		return &MockError{"nil data"}
	}
	if m.ShouldError {
		return &MockError{"deserialization error"}
	}
	// Remove "deserialized:" prefix
	if len(*data) >= 13 && (*data)[:13] == "deserialized:" {
		m.Data = (*data)[13:]
	} else {
		m.Data = *data
	}
	return nil
}

// MockSerializableDeserializable implements both interfaces for testing
type MockSerializableDeserializable struct {
	Data        string
	ShouldError bool
}

func (m *MockSerializableDeserializable) Serialize() (*string, error) {
	if m.ShouldError {
		return nil, &MockError{"serialization error"}
	}
	result := "data:" + m.Data
	return &result, nil
}

func (m *MockSerializableDeserializable) Deserialize(data *string) error {
	if nil == data {
		return &MockError{"nil data"}
	}
	if m.ShouldError {
		return &MockError{"deserialization error"}
	}
	// Remove "data:" prefix
	if len(*data) >= 5 && (*data)[:5] == "data:" {
		m.Data = (*data)[5:]
	} else {
		m.Data = *data
	}
	return nil
}

// MockError is a simple error implementation for testing
type MockError struct {
	message string
}

func (e *MockError) Error() string {
	return e.message
}
