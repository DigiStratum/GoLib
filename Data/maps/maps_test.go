package maps

/*

Unit Tests for maps package

*/

import (
	"fmt"
	"sort"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// TestThat_Strkeys_ReturnsEmptySlice_WhenNil tests that Strkeys returns empty slice for nil input
func TestThat_Strkeys_ReturnsEmptySlice_WhenNil(t *testing.T) {
	// Test
	actual := Strkeys(nil)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, len(actual), t)
}

// TestThat_Strkeys_ReturnsEmptySlice_WhenNotAMap tests that Strkeys returns empty slice for non-map input
func TestThat_Strkeys_ReturnsEmptySlice_WhenNotAMap(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"string input", "not a map"},
		{"int input", 42},
		{"slice input", []string{"a", "b", "c"}},
		{"struct input", struct{ Name string }{"test"}},
		{"pointer input", new(int)},
		{"bool input", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			actual := Strkeys(tt.input)

			// Verify
			ExpectNonNil(actual, t)
			ExpectInt(0, len(actual), t)
		})
	}
}

// TestThat_Strkeys_ReturnsEmptySlice_WhenEmptyMap tests that Strkeys returns empty slice for empty map
func TestThat_Strkeys_ReturnsEmptySlice_WhenEmptyMap(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"empty map[string]interface{}", map[string]interface{}{}},
		{"empty map[string]string", map[string]string{}},
		{"empty map[string]int", map[string]int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			actual := Strkeys(tt.input)

			// Verify
			ExpectNonNil(actual, t)
			ExpectInt(0, len(actual), t)
		})
	}
}

// TestThat_Strkeys_ReturnsKeys_WhenMapWithStringKeys tests that Strkeys returns correct keys for maps with string keys
func TestThat_Strkeys_ReturnsKeys_WhenMapWithStringKeys(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []string
	}{
		{
			name:     "map[string]interface{} with one key",
			input:    map[string]interface{}{"key1": "value1"},
			expected: []string{"key1"},
		},
		{
			name:     "map[string]interface{} with multiple keys",
			input:    map[string]interface{}{"key1": "value1", "key2": "value2", "key3": "value3"},
			expected: []string{"key1", "key2", "key3"},
		},
		{
			name:     "map[string]string with multiple keys",
			input:    map[string]string{"alpha": "a", "beta": "b", "gamma": "c"},
			expected: []string{"alpha", "beta", "gamma"},
		},
		{
			name:     "map[string]int with numeric values",
			input:    map[string]int{"one": 1, "two": 2, "three": 3},
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "map[string]bool with boolean values",
			input:    map[string]bool{"enabled": true, "disabled": false},
			expected: []string{"enabled", "disabled"},
		},
		{
			name:     "map with empty string key",
			input:    map[string]string{"": "empty", "key": "value"},
			expected: []string{"", "key"},
		},
		{
			name:     "map with special characters in keys",
			input:    map[string]interface{}{"key-1": "val", "key_2": "val", "key.3": "val"},
			expected: []string{"key-1", "key_2", "key.3"},
		},
		{
			name:     "map with unicode keys",
			input:    map[string]string{"你好": "hello", "世界": "world"},
			expected: []string{"你好", "世界"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			actual := Strkeys(tt.input)

			// Verify
			ExpectNonNil(actual, t)
			ExpectInt(len(tt.expected), len(actual), t)

			// Sort both slices to compare (map iteration order is not guaranteed)
			sort.Strings(actual)
			expectedSorted := make([]string, len(tt.expected))
			copy(expectedSorted, tt.expected)
			sort.Strings(expectedSorted)

			// Verify each key is present
			for i, key := range expectedSorted {
				ExpectString(key, actual[i], t)
			}
		})
	}
}

// TestThat_Strkeys_ReturnsEmptySlice_WhenNonStringKeyMap tests that Strkeys handles maps with non-string keys
func TestThat_Strkeys_ReturnsEmptySlice_WhenNonStringKeyMap(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{"map[int]string", map[int]string{1: "one", 2: "two"}},
		{"map[int]interface{}", map[int]interface{}{1: "one", 2: "two"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			actual := Strkeys(tt.input)

			// Verify - int keys are converted to their string representation
			ExpectNonNil(actual, t)
			ExpectTrue(len(actual) > 0, t)
		})
	}
}

// TestThat_Strkeys_ReturnsKeys_WhenMapWithNilValues tests that Strkeys handles maps with nil values
func TestThat_Strkeys_ReturnsKeys_WhenMapWithNilValues(t *testing.T) {
	// Setup
	input := map[string]interface{}{
		"key1": nil,
		"key2": "value",
		"key3": nil,
	}
	expected := []string{"key1", "key2", "key3"}

	// Test
	actual := Strkeys(input)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(len(expected), len(actual), t)

	// Sort both slices to compare
	sort.Strings(actual)
	sort.Strings(expected)

	for i, key := range expected {
		ExpectString(key, actual[i], t)
	}
}

// TestThat_Strkeys_HandlesLargeMaps tests that Strkeys can handle large maps
func TestThat_Strkeys_HandlesLargeMaps(t *testing.T) {
	// Setup
	largeMap := make(map[string]interface{})
	expectedSize := 1000
	for i := 0; i < expectedSize; i++ {
		key := fmt.Sprintf("key%d", i)
		largeMap[key] = i
	}

	// Test
	actual := Strkeys(largeMap)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(expectedSize, len(actual), t)

	// Verify all keys are unique
	keySet := make(map[string]bool)
	for _, key := range actual {
		keySet[key] = true
	}
	ExpectInt(expectedSize, len(keySet), t)
}

// TestThat_Strkeys_ReturnsKeys_WhenNestedMapValues tests that Strkeys handles nested maps as values
func TestThat_Strkeys_ReturnsKeys_WhenNestedMapValues(t *testing.T) {
	// Setup
	nestedMap := map[string]interface{}{
		"outer1": map[string]string{"inner1": "value1"},
		"outer2": "simple value",
		"outer3": map[string]interface{}{"inner2": 42},
	}
	expected := []string{"outer1", "outer2", "outer3"}

	// Test
	actual := Strkeys(nestedMap)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(len(expected), len(actual), t)

	// Sort both slices to compare
	sort.Strings(actual)
	sort.Strings(expected)

	for i, key := range expected {
		ExpectString(key, actual[i], t)
	}
}
