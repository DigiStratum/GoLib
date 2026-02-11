package metadata

import (
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// -------------------------------------------------------------------------------------------------
// Factory Function Tests
// -------------------------------------------------------------------------------------------------

func TestThat_MetadataBuilder_NewMetadataBuilder_ReturnsNonNil(t *testing.T) {
	// Test
	sut := NewMetadataBuilder()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_MetadataBuilder_NewMetadataBuilder_ReturnsEmptyMetadata(t *testing.T) {
	// Test
	sut := NewMetadataBuilder()
	metadata := sut.GetMetadata()

	// Verify
	ExpectNonNil(metadata, t)
	names := metadata.GetNames()
	ExpectInt(0, len(names), t)
}

// -------------------------------------------------------------------------------------------------
// Set() Tests
// -------------------------------------------------------------------------------------------------

func TestThat_MetadataBuilder_Set_AddsKeyValuePair(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()

	// Test
	result := sut.Set("testKey", "testValue")

	// Verify - Set should return the builder for chaining
	ExpectNonNil(result, t)
	
	// Verify the value was set
	metadata := sut.GetMetadata()
	value := metadata.Get("testKey")
	if ExpectNonNil(value, t) {
		ExpectString("testValue", *value, t)
	}
}

func TestThat_MetadataBuilder_Set_ReturnsSelfForChaining(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()

	// Test
	result := sut.Set("key1", "value1")

	// Verify - the returned builder should be the same instance
	if sut != result {
		t.Error("Expected Set() to return the same builder instance for method chaining")
	}
}

func TestThat_MetadataBuilder_Set_SupportsMethodChaining(t *testing.T) {
	// Setup & Test
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		Set("key2", "value2").
		Set("key3", "value3")

	// Verify
	metadata := sut.GetMetadata()
	ExpectTrue(metadata.Has("key1", "key2", "key3"), t)
}

func TestThat_MetadataBuilder_Set_UpdatesExistingKey(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()
	sut.Set("key", "originalValue")

	// Test - update the same key
	sut.Set("key", "newValue")

	// Verify
	metadata := sut.GetMetadata()
	value := metadata.Get("key")
	if ExpectNonNil(value, t) {
		ExpectString("newValue", *value, t)
	}
}

func TestThat_MetadataBuilder_Set_HandlesEmptyValue(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()

	// Test
	sut.Set("emptyKey", "")

	// Verify
	metadata := sut.GetMetadata()
	ExpectTrue(metadata.Has("emptyKey"), t)
	value := metadata.Get("emptyKey")
	if ExpectNonNil(value, t) {
		ExpectString("", *value, t)
	}
}

func TestThat_MetadataBuilder_Set_HandlesMultipleValues(t *testing.T) {
	// Table-driven test
	tests := []struct {
		key   string
		value string
	}{
		{"key1", "value1"},
		{"key2", "value2"},
		{"key3", "value3"},
		{"key4", "value4"},
		{"key5", "value5"},
	}

	// Setup
	sut := NewMetadataBuilder()

	// Test
	for _, tt := range tests {
		sut.Set(tt.key, tt.value)
	}

	// Verify
	metadata := sut.GetMetadata()
	for _, tt := range tests {
		value := metadata.Get(tt.key)
		if !ExpectNonNil(value, t) {
			t.Errorf("Expected non-nil value for key '%s'", tt.key)
			continue
		}
		ExpectString(tt.value, *value, t)
	}
}

// -------------------------------------------------------------------------------------------------
// GetMetadata() Tests
// -------------------------------------------------------------------------------------------------

func TestThat_MetadataBuilder_GetMetadata_ReturnsNonNil(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()

	// Test
	metadata := sut.GetMetadata()

	// Verify
	ExpectNonNil(metadata, t)
}

func TestThat_MetadataBuilder_GetMetadata_ReturnsEmptyMetadata_WhenNothingSet(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()

	// Test
	metadata := sut.GetMetadata()

	// Verify
	names := metadata.GetNames()
	ExpectInt(0, len(names), t)
}

func TestThat_MetadataBuilder_GetMetadata_ReturnsPopulatedMetadata_WhenValuesSet(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().
		Set("key1", "value1").
		Set("key2", "value2")

	// Test
	metadata := sut.GetMetadata()

	// Verify
	names := metadata.GetNames()
	ExpectInt(2, len(names), t)
	ExpectTrue(metadata.Has("key1", "key2"), t)
}

func TestThat_MetadataBuilder_GetMetadata_ReturnsSameInstance(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().Set("key", "value")

	// Test - call GetMetadata multiple times
	metadata1 := sut.GetMetadata()
	metadata2 := sut.GetMetadata()

	// Verify - should return the same instance
	if metadata1 != metadata2 {
		t.Error("Expected GetMetadata() to return the same instance on multiple calls")
	}
}

func TestThat_MetadataBuilder_GetMetadata_ReflectsSubsequentChanges(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder().Set("key1", "value1")
	metadata := sut.GetMetadata()

	// Test - add more values after getting metadata
	sut.Set("key2", "value2")

	// Verify - the metadata should reflect the new value
	// (since it's the same instance)
	ExpectTrue(metadata.Has("key2"), t)
}

// -------------------------------------------------------------------------------------------------
// Integration Tests
// -------------------------------------------------------------------------------------------------

func TestThat_MetadataBuilder_BuildsComplexMetadata(t *testing.T) {
	// Setup & Test - build complex metadata with method chaining
	sut := NewMetadataBuilder().
		Set("name", "John Doe").
		Set("email", "john@example.com").
		Set("age", "30").
		Set("city", "New York").
		Set("country", "USA")

	// Verify
	metadata := sut.GetMetadata()
	
	// Check all keys exist
	ExpectTrue(metadata.Has("name", "email", "age", "city", "country"), t)
	
	// Check values
	name := metadata.Get("name")
	if ExpectNonNil(name, t) {
		ExpectString("John Doe", *name, t)
	}
	
	email := metadata.Get("email")
	if ExpectNonNil(email, t) {
		ExpectString("john@example.com", *email, t)
	}
	
	// Check count
	names := metadata.GetNames()
	ExpectInt(5, len(names), t)
}

func TestThat_MetadataBuilder_SupportsIterativeBuilding(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()

	// Test - build iteratively (not chained)
	sut.Set("key1", "value1")
	sut.Set("key2", "value2")
	sut.Set("key3", "value3")

	// Verify
	metadata := sut.GetMetadata()
	ExpectTrue(metadata.Has("key1", "key2", "key3"), t)
}

func TestThat_MetadataBuilder_HandlesSpecialCharacters(t *testing.T) {
	// Table-driven test with special characters
	tests := []struct {
		key   string
		value string
	}{
		{"key-with-dash", "value-with-dash"},
		{"key_with_underscore", "value_with_underscore"},
		{"key.with.dots", "value.with.dots"},
		{"key with spaces", "value with spaces"},
		{"unicode-key-日本語", "unicode-value-こんにちは"},
	}

	// Setup
	sut := NewMetadataBuilder()

	// Test
	for _, tt := range tests {
		sut.Set(tt.key, tt.value)
	}

	// Verify
	metadata := sut.GetMetadata()
	for _, tt := range tests {
		value := metadata.Get(tt.key)
		if !ExpectNonNil(value, t) {
			t.Errorf("Expected non-nil value for key '%s'", tt.key)
			continue
		}
		ExpectString(tt.value, *value, t)
	}
}

func TestThat_MetadataBuilder_HandlesLargeNumberOfEntries(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()
	count := 100

	// Test - add many entries
	for i := 0; i < count; i++ {
		key := "key" + string(rune('A'+i%26)) + string(rune('0'+i/26))
		value := "value" + string(rune('A'+i%26)) + string(rune('0'+i/26))
		sut.Set(key, value)
	}

	// Verify
	metadata := sut.GetMetadata()
	names := metadata.GetNames()
	ExpectInt(count, len(names), t)
}

func TestThat_MetadataBuilder_UpdatesAndOverwrites(t *testing.T) {
	// Setup
	sut := NewMetadataBuilder()

	// Test - set initial values
	sut.Set("key1", "original1")
	sut.Set("key2", "original2")
	sut.Set("key3", "original3")

	// Update some values
	sut.Set("key1", "updated1")
	sut.Set("key3", "updated3")

	// Verify
	metadata := sut.GetMetadata()
	
	// Check that we still have 3 keys (not 5)
	names := metadata.GetNames()
	ExpectInt(3, len(names), t)
	
	// Check updated values
	value1 := metadata.Get("key1")
	if ExpectNonNil(value1, t) {
		ExpectString("updated1", *value1, t)
	}
	
	value2 := metadata.Get("key2")
	if ExpectNonNil(value2, t) {
		ExpectString("original2", *value2, t)
	}
	
	value3 := metadata.Get("key3")
	if ExpectNonNil(value3, t) {
		ExpectString("updated3", *value3, t)
	}
}
