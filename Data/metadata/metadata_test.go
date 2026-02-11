package metadata

import (
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// -------------------------------------------------------------------------------------------------
// Has() Tests
// -------------------------------------------------------------------------------------------------

func TestThat_Metadata_Has_ReturnsTrue_WhenKeyExists(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		Set("key2", "value2").
		GetMetadata()

	// Test & Verify
	ExpectTrue(sut.Has("key1"), t)
	ExpectTrue(sut.Has("key2"), t)
}

func TestThat_Metadata_Has_ReturnsFalse_WhenKeyDoesNotExist(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		GetMetadata()

	// Test & Verify
	ExpectFalse(sut.Has("nonexistent"), t)
}

func TestThat_Metadata_Has_ReturnsTrue_WhenAllKeysExist(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		Set("key2", "value2").
		Set("key3", "value3").
		GetMetadata()

	// Test & Verify
	ExpectTrue(sut.Has("key1", "key2", "key3"), t)
}

func TestThat_Metadata_Has_ReturnsFalse_WhenSomeKeysDoNotExist(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		Set("key2", "value2").
		GetMetadata()

	// Test & Verify
	ExpectFalse(sut.Has("key1", "key2", "key3"), t)
}

func TestThat_Metadata_Has_ReturnsFalse_WhenNoKeysExist(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		GetMetadata()

	// Test & Verify
	ExpectFalse(sut.Has("nonexistent1", "nonexistent2"), t)
}

func TestThat_Metadata_Has_ReturnsTrue_ForEmptyMetadata(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().GetMetadata()

	// Test & Verify - calling Has() with no arguments should return true
	ExpectTrue(sut.Has(), t)
}

// -------------------------------------------------------------------------------------------------
// Get() Tests
// -------------------------------------------------------------------------------------------------

func TestThat_Metadata_Get_ReturnsValue_WhenKeyExists(t *testing.T) {
	// Setup
	expectedValue := "testValue"
	sut := NewMetadataBuilder().
		Set("testKey", expectedValue).
		GetMetadata()

	// Test
	actual := sut.Get("testKey")

	// Verify
	ExpectNonNil(actual, t)
	ExpectString(expectedValue, *actual, t)
}

func TestThat_Metadata_Get_ReturnsNil_WhenKeyDoesNotExist(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		GetMetadata()

	// Test
	actual := sut.Get("nonexistent")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Metadata_Get_ReturnsNil_ForEmptyMetadata(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().GetMetadata()

	// Test
	actual := sut.Get("anyKey")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Metadata_Get_ReturnsCorrectValue_ForMultipleKeys(t *testing.T) {
	// Table-driven test
	tests := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
		{"specialKey", "specialValue"},
		{"emptyValue", ""},
	}

	// Setup
	builder := NewMetadataBuilder()
	for _, tt := range tests {
		builder.Set(tt.key, tt.value)
	}
	sut := builder.GetMetadata()

	// Test & Verify
	for _, tt := range tests {
		actual := sut.Get(tt.key)
		if !ExpectNonNil(actual, t) {
			t.Errorf("Expected non-nil for key '%s'", tt.key)
			continue
		}
		ExpectString(tt.value, *actual, t)
	}
}

// -------------------------------------------------------------------------------------------------
// GetNames() Tests
// -------------------------------------------------------------------------------------------------

func TestThat_Metadata_GetNames_ReturnsEmptySlice_ForEmptyMetadata(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().GetMetadata()

	// Test
	actual := sut.GetNames()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(actual), t)
}

func TestThat_Metadata_GetNames_ReturnsSingleKey_ForSingleEntry(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("onlyKey", "onlyValue").
		GetMetadata()

	// Test
	actual := sut.GetNames()

	// Verify
	ExpectInt(1, len(actual), t)
	ExpectString("onlyKey", actual[0], t)
}

func TestThat_Metadata_GetNames_ReturnsAllKeys_ForMultipleEntries(t *testing.T) {
	// Setup
	expectedKeys := map[string]bool{
		"key1": true,
		"key2": true,
		"key3": true,
	}
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		Set("key2", "value2").
		Set("key3", "value3").
		GetMetadata()

	// Test
	actual := sut.GetNames()

	// Verify
	ExpectInt(3, len(actual), t)
	
	// Verify all expected keys are present
	for _, key := range actual {
		if !expectedKeys[key] {
			t.Errorf("Unexpected key in results: %s", key)
		}
		delete(expectedKeys, key)
	}
	
	// Verify no expected keys are missing
	if len(expectedKeys) > 0 {
		t.Errorf("Missing expected keys: %v", expectedKeys)
	}
}

// -------------------------------------------------------------------------------------------------
// Integration Tests
// -------------------------------------------------------------------------------------------------

func TestThat_Metadata_Operations_WorkTogether(t *testing.T) {
	// Setup - create metadata with multiple entries
	sut := NewMetadataBuilder().
		Set("name", "John Doe").
		Set("email", "john@example.com").
		Set("age", "30").
		GetMetadata()

	// Test Has()
	ExpectTrue(sut.Has("name"), t)
	ExpectTrue(sut.Has("email"), t)
	ExpectTrue(sut.Has("age"), t)
	ExpectTrue(sut.Has("name", "email", "age"), t)
	ExpectFalse(sut.Has("phone"), t)

	// Test Get()
	name := sut.Get("name")
	if ExpectNonNil(name, t) {
		ExpectString("John Doe", *name, t)
	}

	email := sut.Get("email")
	if ExpectNonNil(email, t) {
		ExpectString("john@example.com", *email, t)
	}

	// Test GetNames()
	names := sut.GetNames()
	ExpectInt(3, len(names), t)
}

func TestThat_Metadata_HandlesEmptyValues(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("emptyKey", "").
		Set("normalKey", "normalValue").
		GetMetadata()

	// Test & Verify
	ExpectTrue(sut.Has("emptyKey"), t)
	
	emptyValue := sut.Get("emptyKey")
	if ExpectNonNil(emptyValue, t) {
		ExpectString("", *emptyValue, t)
	}

	names := sut.GetNames()
	ExpectInt(2, len(names), t)
}

func TestThat_Metadata_HandlesSpecialCharacters(t *testing.T) {
	// Table-driven test with special characters
	tests := []struct {
		key   string
		value string
	}{
		{"key-with-dash", "value"},
		{"key_with_underscore", "value"},
		{"key.with.dots", "value"},
		{"key with spaces", "value with spaces"},
		{"key@special#chars", "value!@#$%^&*()"},
	}

	// Setup
	builder := NewMetadataBuilder()
	for _, tt := range tests {
		builder.Set(tt.key, tt.value)
	}
	sut := builder.GetMetadata()

	// Test & Verify
	for _, tt := range tests {
		if !ExpectTrue(sut.Has(tt.key), t) {
			t.Errorf("Expected to have key '%s'", tt.key)
		}
		
		actual := sut.Get(tt.key)
		if ExpectNonNil(actual, t) {
			ExpectString(tt.value, *actual, t)
		}
	}
}
