package bytemap

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// -------------------------------------------------------------------------------------------------
// Factory Function Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_NewByteMap_ReturnsEmptyByteMap(t *testing.T) {
	// Test
	sut := NewByteMap()

	// Verify
	ExpectNonNil(sut, t)
	ExpectInt(0, sut.Size(), t)
	ExpectTrue(sut.IsEmpty(), t)
}

// -------------------------------------------------------------------------------------------------
// Copy Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_Copy_ReturnsEmpty_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	actual := sut.Copy()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, actual.Size(), t)
	ExpectTrue(actual.IsEmpty(), t)
}

func TestThat_ByteMap_Copy_ReturnsNonEmpty_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewByteMap()
	num := 25
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := []byte(fmt.Sprintf("value-%d", i))
		sut.Set(key, value)
	}

	// Test
	actual := sut.Copy()

	// Verify
	ExpectInt(num, actual.Size(), t)
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		expectedValue := []byte(fmt.Sprintf("value-%d", i))
		actualValue := actual.Get(key)
		ExpectNonNil(actualValue, t)
		ExpectString(string(expectedValue), string(*actualValue), t)
	}
}

func TestThat_ByteMap_Copy_CreatesDeepCopy(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "testkey"
	originalValue := []byte("original")
	sut.Set(key, originalValue)

	// Test
	copy := sut.Copy()
	copy.Set(key, []byte("modified"))

	// Verify - original should be unchanged
	sutValue := sut.Get(key)
	ExpectNonNil(sutValue, t)
	ExpectString("original", string(*sutValue), t)
}

// -------------------------------------------------------------------------------------------------
// Merge Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_Merge_AddsNothing_ForEmptyMaps(t *testing.T) {
	// Setup
	sut := NewByteMap()
	otherMap := NewByteMap()

	// Test
	sut.Merge(otherMap)

	// Verify
	ExpectInt(0, sut.Size(), t)
}

func TestThat_ByteMap_Merge_AddsEntries_ForNonEmptyMaps(t *testing.T) {
	// Setup
	key := "testkey"
	expected := []byte("testvalue")
	other := NewByteMap()
	other.Set(key, expected)
	sut := NewByteMap()

	// Test
	sut.Merge(other)

	// Verify
	ExpectFalse(sut.IsEmpty(), t)
	ExpectInt(1, sut.Size(), t)
	ExpectTrue(sut.Has(key), t)
	actual := sut.Get(key)
	ExpectNonNil(actual, t)
	ExpectString(string(expected), string(*actual), t)
}

func TestThat_ByteMap_Merge_OverwritesExistingKeys(t *testing.T) {
	// Setup
	key := "samekey"
	sut := NewByteMap()
	sut.Set(key, []byte("original"))
	other := NewByteMap()
	expected := []byte("updated")
	other.Set(key, expected)

	// Test
	sut.Merge(other)

	// Verify
	ExpectInt(1, sut.Size(), t)
	actual := sut.Get(key)
	ExpectNonNil(actual, t)
	ExpectString(string(expected), string(*actual), t)
}

func TestThat_ByteMap_Merge_DoesNothing_WhenNilReceiver(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test
	sut.Merge(NewByteMap())

	// Verify
	ExpectNil(sut, t)
}

// -------------------------------------------------------------------------------------------------
// IsEmpty and Size Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_IsEmpty_IsTrue_WhenNew(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Verify
	ExpectTrue(sut.IsEmpty(), t)
}

func TestThat_ByteMap_IsEmpty_IsFalse_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("key", []byte("value"))

	// Verify
	ExpectFalse(sut.IsEmpty(), t)
}

func TestThat_ByteMap_Size_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Verify
	ExpectInt(0, sut.Size(), t)
}

func TestThat_ByteMap_Size_Is1_WithOneSet(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	sut.Set("key", []byte("value"))

	// Verify
	ExpectInt(1, sut.Size(), t)
}

func TestThat_ByteMap_Size_Returns0_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Verify
	ExpectInt(0, sut.Size(), t)
}

// -------------------------------------------------------------------------------------------------
// Set and Get Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_Set_SetsValue(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "testkey"
	expected := []byte("testvalue")

	// Test
	sut.Set(key, expected)

	// Verify
	actual := sut.Get(key)
	ExpectNonNil(actual, t)
	ExpectString(string(expected), string(*actual), t)
}

func TestThat_ByteMap_Set_OverwritesValue_ForExistingKey(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "somekey"
	sut.Set(key, []byte("unexpected"))

	// Test
	expected := []byte("expected")
	sut.Set(key, expected)

	// Verify
	actual := sut.Get(key)
	ExpectNonNil(actual, t)
	ExpectString(string(expected), string(*actual), t)
}

func TestThat_ByteMap_Set_DoesNothing_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test (should not panic)
	sut.Set("testkey", []byte("testvalue"))

	// Verify
	ExpectNil(sut, t)
}

func TestThat_ByteMap_Get_ReturnsValue_ForSetKey(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "testkey"
	expected := []byte("testvalue")

	// Test
	sut.Set(key, expected)
	actual := sut.Get(key)

	// Verify
	ExpectNonNil(actual, t)
	ExpectString(string(expected), string(*actual), t)
}

func TestThat_ByteMap_Get_ReturnsNil_ForUnsetKey(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	actual := sut.Get("boguskey")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_ByteMap_Get_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test
	actual := sut.Get("anykey")

	// Verify
	ExpectNil(actual, t)
}

// -------------------------------------------------------------------------------------------------
// GetInt64 Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_GetInt64_ReturnsValue_ForParseableInt(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected int64
	}{
		{"positive decimal", "1234567", 1234567},
		{"negative decimal", "-9876", -9876},
		{"zero", "0", 0},
		{"hex with 0x prefix", "0x10", 16},
		{"octal with 0 prefix", "010", 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			sut := NewByteMap()
			key := "testkey"
			sut.Set(key, []byte(tt.value))

			// Test
			actual := sut.GetInt64(key)

			// Verify
			ExpectNonNil(actual, t)
			ExpectInt64(tt.expected, *actual, t)
		})
	}
}

func TestThat_ByteMap_GetInt64_ReturnsNil_ForNonParseableValue(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "testkey"
	sut.Set(key, []byte("notanumber"))

	// Test
	actual := sut.GetInt64(key)

	// Verify
	ExpectNil(actual, t)
}

func TestThat_ByteMap_GetInt64_ReturnsNil_ForUnsetKey(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	actual := sut.GetInt64("boguskey")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_ByteMap_GetInt64_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test
	actual := sut.GetInt64("anykey")

	// Verify
	ExpectNil(actual, t)
}

// -------------------------------------------------------------------------------------------------
// GetBool Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_GetBool_ReturnsTrue_ForTruthyValues(t *testing.T) {
	truthyValues := []string{
		"true", "TRUE", "t", "T", "1",
		"on", "ON", "yes", "YES",
	}

	for _, value := range truthyValues {
		t.Run(fmt.Sprintf("value=%s", value), func(t *testing.T) {
			// Setup
			sut := NewByteMap()
			key := "testkey"
			sut.Set(key, []byte(value))

			// Test
			actual := sut.GetBool(key)

			// Verify
			ExpectTrue(actual, t)
		})
	}
}

func TestThat_ByteMap_GetBool_ReturnsTrue_ForNonZeroIntegers(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "testkey"
	sut.Set(key, []byte("42"))

	// Test
	actual := sut.GetBool(key)

	// Verify
	ExpectTrue(actual, t)
}

func TestThat_ByteMap_GetBool_ReturnsFalse_ForFalsyValues(t *testing.T) {
	falsyValues := []string{
		"false", "FALSE", "f", "F", "0",
		"off", "OFF", "no", "NO", "bogus",
	}

	for _, value := range falsyValues {
		t.Run(fmt.Sprintf("value=%s", value), func(t *testing.T) {
			// Setup
			sut := NewByteMap()
			key := "testkey"
			sut.Set(key, []byte(value))

			// Test
			actual := sut.GetBool(key)

			// Verify
			ExpectFalse(actual, t)
		})
	}
}

func TestThat_ByteMap_GetBool_ReturnsFalse_ForUnsetKey(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	actual := sut.GetBool("boguskey")

	// Verify
	ExpectFalse(actual, t)
}

func TestThat_ByteMap_GetBool_ReturnsFalse_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test
	actual := sut.GetBool("anykey")

	// Verify
	ExpectFalse(actual, t)
}

// -------------------------------------------------------------------------------------------------
// Has and HasAll Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_Has_IsFalse_WhenKeyMissing(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Verify
	ExpectFalse(sut.Has("boguskey"), t)
}

func TestThat_ByteMap_Has_IsTrue_WhenKeyExists(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "testkey"

	// Test
	sut.Set(key, []byte("testvalue"))

	// Verify
	ExpectTrue(sut.Has(key), t)
}

func TestThat_ByteMap_Has_ReturnsFalse_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Verify
	ExpectFalse(sut.Has("anykey"), t)
}

func TestThat_ByteMap_HasAll_ReturnsTrue_WhenKeysEmptySet(t *testing.T) {
	// Setup
	sut := NewByteMap()
	keys := make([]string, 0)

	// Test
	actual := sut.HasAll(&keys)

	// Verify
	ExpectTrue(actual, t)
}

func TestThat_ByteMap_HasAll_ReturnsFalse_WhenAnyKeyMissing(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("key1", []byte("value1"))
	keys := []string{"key1", "missingkey"}

	// Test
	actual := sut.HasAll(&keys)

	// Verify
	ExpectFalse(actual, t)
}

func TestThat_ByteMap_HasAll_ReturnsTrue_WhenAllKeysExist(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("key1", []byte("value1"))
	sut.Set("key2", []byte("value2"))
	keys := []string{"key1", "key2"}

	// Test
	actual := sut.HasAll(&keys)

	// Verify
	ExpectTrue(actual, t)
}

func TestThat_ByteMap_HasAll_ReturnsFalse_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil
	keys := []string{"key1"}

	// Test
	actual := sut.HasAll(&keys)

	// Verify
	ExpectFalse(actual, t)
}

// -------------------------------------------------------------------------------------------------
// GetKeys Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_GetKeys_ReturnsEmpty_ForEmptyMap(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	actual := sut.GetKeys()

	// Verify
	ExpectInt(0, len(actual), t)
}

func TestThat_ByteMap_GetKeys_ReturnsKeys_ForNonEmptyMap(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	sut.Set("k2", []byte("v2"))

	// Test
	actual := sut.GetKeys()

	// Verify
	ExpectInt(2, len(actual), t)
	// Keys can be in any order
	hasK1 := false
	hasK2 := false
	for _, key := range actual {
		if key == "k1" {
			hasK1 = true
		}
		if key == "k2" {
			hasK2 = true
		}
	}
	ExpectTrue(hasK1, t)
	ExpectTrue(hasK2, t)
}

func TestThat_ByteMap_GetKeys_ReturnsEmpty_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test
	actual := sut.GetKeys()

	// Verify
	ExpectInt(0, len(actual), t)
}

// -------------------------------------------------------------------------------------------------
// GetSubset Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_GetSubset_ReturnsEmpty_ForEmptyKeySet(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	sut.Set("k2", []byte("v2"))
	keys := make([]string, 0)

	// Test
	actual := sut.GetSubset(&keys)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, actual.Size(), t)
}

func TestThat_ByteMap_GetSubset_ReturnsEmpty_ForNilKeySet(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))

	// Test
	actual := sut.GetSubset(nil)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(0, actual.Size(), t)
}

func TestThat_ByteMap_GetSubset_ReturnsSubset_ForValidKeys(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	sut.Set("k2", []byte("v2"))
	sut.Set("k3", []byte("v3"))
	keys := []string{"k1", "k3"}

	// Test
	actual := sut.GetSubset(&keys)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(2, actual.Size(), t)
	ExpectTrue(actual.Has("k1"), t)
	ExpectTrue(actual.Has("k3"), t)
	ExpectFalse(actual.Has("k2"), t)
}

func TestThat_ByteMap_GetSubset_IgnoresMissingKeys(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	keys := []string{"k1", "missingkey"}

	// Test
	actual := sut.GetSubset(&keys)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt(1, actual.Size(), t)
	ExpectTrue(actual.Has("k1"), t)
}

func TestThat_ByteMap_GetSubset_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil
	keys := []string{"k1"}

	// Test
	actual := sut.GetSubset(&keys)

	// Verify
	ExpectNil(actual, t)
}

// -------------------------------------------------------------------------------------------------
// Drop and DropSet Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_Drop_RemovesKey_WhenExists(t *testing.T) {
	// Setup
	sut := NewByteMap()
	key := "testkey"
	sut.Set(key, []byte("testvalue"))

	// Test
	result := sut.Drop(key)

	// Verify
	ExpectNonNil(result, t)
	ExpectFalse(sut.Has(key), t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_ByteMap_Drop_DoesNothing_WhenKeyMissing(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("existingkey", []byte("value"))

	// Test
	result := sut.Drop("missingkey")

	// Verify
	ExpectNonNil(result, t)
	ExpectInt(1, sut.Size(), t)
}

func TestThat_ByteMap_Drop_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test
	result := sut.Drop("anykey")

	// Verify
	ExpectNil(result, t)
}

func TestThat_ByteMap_DropSet_RemovesMultipleKeys(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	sut.Set("k2", []byte("v2"))
	sut.Set("k3", []byte("v3"))
	keys := []string{"k1", "k3"}

	// Test
	result := sut.DropSet(&keys)

	// Verify
	ExpectNonNil(result, t)
	ExpectInt(1, sut.Size(), t)
	ExpectTrue(sut.Has("k2"), t)
	ExpectFalse(sut.Has("k1"), t)
	ExpectFalse(sut.Has("k3"), t)
}

func TestThat_ByteMap_DropSet_DoesNothing_ForEmptyKeySet(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	keys := make([]string, 0)

	// Test
	result := sut.DropSet(&keys)

	// Verify
	ExpectNonNil(result, t)
	ExpectInt(1, sut.Size(), t)
}

func TestThat_ByteMap_DropSet_DoesNothing_ForNilKeySet(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))

	// Test
	result := sut.DropSet(nil)

	// Verify
	ExpectNonNil(result, t)
	ExpectInt(1, sut.Size(), t)
}

func TestThat_ByteMap_DropSet_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil
	keys := []string{"k1"}

	// Test
	result := sut.DropSet(&keys)

	// Verify
	ExpectNil(result, t)
}

// -------------------------------------------------------------------------------------------------
// DropAll Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_DropAll_ClearsAllEntries(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	sut.Set("k2", []byte("v2"))

	// Test
	sut.DropAll()

	// Verify
	ExpectInt(0, sut.Size(), t)
	ExpectTrue(sut.IsEmpty(), t)
}

func TestThat_ByteMap_DropAll_DoesNothing_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	sut.DropAll()

	// Verify
	ExpectInt(0, sut.Size(), t)
	ExpectTrue(sut.IsEmpty(), t)
}

func TestThat_ByteMap_DropAll_DoesNothing_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test (should not panic)
	sut.DropAll()

	// Verify
	ExpectNil(sut, t)
}

// -------------------------------------------------------------------------------------------------
// GetIterator Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_GetIterator_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *ByteMap // nil

	// Test
	it := sut.GetIterator()

	// Verify - GetIterator returns nil for a nil ByteMap
	if it != nil {
		t.Errorf("Expected nil iterator for nil ByteMap, got non-nil")
	}
}

func TestThat_ByteMap_GetIterator_ReturnsIterator_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewByteMap()

	// Test
	it := sut.GetIterator()

	// Verify
	ExpectNonNil(it, t)
	item := it()
	ExpectNil(item, t)
}

func TestThat_ByteMap_GetIterator_ReturnsWorkingIterator_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewByteMap()
	expectedItems := make(map[string][]byte)
	numItems := 10
	for i := 0; i < numItems; i++ {
		key := fmt.Sprintf("k%d", i)
		value := []byte(fmt.Sprintf("v%d", i))
		sut.Set(key, value)
		expectedItems[key] = value
	}

	// Test
	it := sut.GetIterator()

	// Verify
	ExpectNonNil(it, t)
	count := 0
	for itemi := it(); itemi != nil; itemi = it() {
		count++
		ExpectNonNil(itemi, t)
		if kvp, ok := itemi.(*KeyValuePair); ok {
			ExpectTrue(ok, t)
			actualKey := kvp.Key
			actualValue := kvp.Value
			if expectedValue, isExpectedKey := expectedItems[actualKey]; isExpectedKey {
				ExpectTrue(isExpectedKey, t)
				ExpectString(string(expectedValue), string(actualValue), t)
				delete(expectedItems, actualKey)
			} else {
				t.Errorf("Unexpected key: %s", actualKey)
			}
		} else {
			t.Errorf("Expected *KeyValuePair, got something else")
		}
	}
	ExpectInt(numItems, count, t)
	ExpectInt(0, len(expectedItems), t)
}

func TestThat_ByteMap_GetIterator_IteratesExactlyOnce_PerItem(t *testing.T) {
	// Setup
	sut := NewByteMap()
	sut.Set("k1", []byte("v1"))
	sut.Set("k2", []byte("v2"))

	// Test
	it := sut.GetIterator()
	item1 := it()
	item2 := it()
	item3 := it()

	// Verify
	ExpectNonNil(item1, t)
	ExpectNonNil(item2, t)
	ExpectNil(item3, t)
}

// -------------------------------------------------------------------------------------------------
// Thread Safety Tests
// -------------------------------------------------------------------------------------------------

func TestThat_ByteMap_Merge_IsThreadSafe(t *testing.T) {
	// Setup
	sut := NewByteMap()
	numGoroutines := 10
	itemsPerGoroutine := 100
	var wg sync.WaitGroup

	// Test - concurrent merges (Merge uses mutex for thread safety)
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			other := NewByteMap()
			for j := 0; j < itemsPerGoroutine; j++ {
				key := fmt.Sprintf("g%d-k%d", goroutineID, j)
				value := []byte(fmt.Sprintf("g%d-v%d", goroutineID, j))
				other.Set(key, value)
			}
			sut.Merge(other)
		}(i)
	}
	wg.Wait()

	// Verify
	expectedSize := numGoroutines * itemsPerGoroutine
	ExpectInt(expectedSize, sut.Size(), t)
}
