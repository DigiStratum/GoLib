package hashmap

/*

Unit Tests for HashMap

ref: https://blog.alexellis.io/golang-writing-unit-tests/

TODO: Add some tests to verify that we get a copy of the result rather than a reference;
We don't want to accidentally introduce a change that starts handing out references which would
allow the state to be modified directly by an outside party.

*/

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_HashMap_NewHashMap_ReturnsEmptyHashMap(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Verify
	ExpectNonNil(sut, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_CopyHashMap_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *HashMap	// nil

	// Verify
	actual := CopyHashMap(sut)
	ExpectNil(actual, t)
}

func TestThat_HashMap_CopyHashMap_ReturnsEmpty_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	actual := CopyHashMap(sut)

	// Verify
	ExpectInt(0, actual.Size(), t)
}

func TestThat_HashMap_CopyHashMap_ReturnsNonEmpty_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewHashMap()
	num := 25
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		sut.Set(key, value)
	}

	// Test
	actual := CopyHashMap(sut)

	// Verify
	ExpectInt(num, actual.Size(), t)
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		expectedValue := fmt.Sprintf("value-%d", i)
		actualValue := actual.Get(key)
		ExpectString(expectedValue, *actualValue, t)
	}
}

func TestThat_HashMap_NewHashMapFromJsonString_ReturnsPopulatedHashMap_ForPopulatedJson(t *testing.T) {
	// Setup
	jsonStr := "{\"k1\":\"v1\",\"k2\":\"v2\"}"

	// Test
	sut, err := NewHashMapFromJsonString(&jsonStr)

	// Verify
	ExpectNil(err, t)
	ExpectInt(2, sut.Size(), t)
	ExpectTrue(sut.Has("k1"), t)
	ExpectTrue(sut.Has("k2"), t)
	v1 := sut.Get("k1")
	ExpectString("v1", *v1, t)
	v2 := sut.Get("k2")
	ExpectString("v2", *v2, t)
}

func TestThat_HashMap_NewHashMapFromJsonString_ReturnsEmptyHashMap_ForNonPopulatedJson(t *testing.T) {
	// Setup
	jsonStr := "{}"

	// Test
	sut, err := NewHashMapFromJsonString(&jsonStr)

	// Verify
	ExpectNil(err, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_NewHashMapFromJsonString_ReturnsError_ForInvalidJson(t *testing.T) {
	// Setup
	jsonStr := "}{"

	// Test
	sut, err := NewHashMapFromJsonString(&jsonStr)

	// Verify
	ExpectNonNil(err, t)
	ExpectNil(sut, t)
}

func TestThat_HashMap_NewHashMapFromJsonFile_ReturnsPopulatedMap_ForPopulatedJson(t *testing.T) {
	// Test
	sut, err := NewHashMapFromJsonFile("hashmap_test.json")

	// Verify
	ExpectNonNil(sut, t)
	ExpectNil(err, t)
	ExpectInt(2, sut.Size(), t)
	ExpectTrue(sut.Has("k1"), t)
	ExpectTrue(sut.Has("k2"), t)
	v1 := sut.Get("k1")
	ExpectString("v1", *v1, t)
	v2 := sut.Get("k2")
	ExpectString("v2", *v2, t)
}

func TestThat_HashMap_NewHashMapFromJsonFile_ReturnsEmptyMap_ForNonPopulatedJson(t *testing.T) {
	// Test
	sut, err := NewHashMapFromJsonFile("hashmap_test_empty.json")

	// Verify
	ExpectNonNil(sut, t)
	ExpectNil(err, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_NewHashMapFromJsonFile_ReturnsError_ForInvalidJson(t *testing.T) {
	// Test
	sut, err := NewHashMapFromJsonFile("hashmap_test_invalid.json")

	// Verify
	ExpectNil(sut, t)
	ExpectNonNil(err, t)
}

func TestThat_HashMap_NewHashMapFromJsonFile_ReturnsError_ForMissingJson(t *testing.T) {
	// Test
	sut, err := NewHashMapFromJsonFile("missingfile.json")

	// Verify
	ExpectNil(sut, t)
	ExpectNonNil(err, t)
}

func TestThat_HashMap_LoadFromJsonString_PopulatesMap_ForPopulatedJson(t *testing.T) {
	// Setup
	sut := NewHashMap()
	jsonStr := "{\"k1\":\"v1\",\"k2\":\"v2\"}"

	// Test
	err := sut.LoadFromJsonString(&jsonStr)

	// Verify
	ExpectNil(err, t)
	ExpectInt(2, sut.Size(), t)
	ExpectTrue(sut.Has("k1"), t)
	ExpectTrue(sut.Has("k2"), t)
	v1 := sut.Get("k1")
	ExpectString("v1", *v1, t)
	v2 := sut.Get("k2")
	ExpectString("v2", *v2, t)
}

func TestThat_HashMap_LoadFromJsonString_ChangesNothing_ForNonPopulatedJson(t *testing.T) {
	// Setup
	sut := NewHashMap()
	jsonStr := "{}"

	// Test
	err := sut.LoadFromJsonString(&jsonStr)

	// Verify
	ExpectNil(err, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_LoadFromJsonString_ReturnsError_ForInvalidJson(t *testing.T) {
	// Setup
	sut := NewHashMap()
	jsonStr := "}{"

	// Test
	err := sut.LoadFromJsonString(&jsonStr)

	// Verify
	ExpectNonNil(err, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_LoadFromJsonFile_PopulatesMap_ForPopulatedJson(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	err := sut.LoadFromJsonFile("hashmap_test.json")

	// Verify
	ExpectNil(err, t)
	ExpectInt(2, sut.Size(), t)
	ExpectTrue(sut.Has("k1"), t)
	ExpectTrue(sut.Has("k2"), t)
	v1 := sut.Get("k1")
	ExpectString("v1", *v1, t)
	v2 := sut.Get("k2")
	ExpectString("v2", *v2, t)
}

func TestThat_HashMap_LoadFromJsonFile_ChangesNothing_ForNonPopulatedJson(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	err := sut.LoadFromJsonFile("hashmap_test_empty.json")

	// Verify
	ExpectNil(err, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_LoadFromJsonFile_ReturnsError_ForInvalidJson(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	err := sut.LoadFromJsonFile("hashmap_test_invalid.json")

	// Verify
	ExpectNonNil(err, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_LoadFromJsonFile_ReturnsError_ForMissingJson(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	err := sut.LoadFromJsonFile("missingfile.json")

	// Verify
	ExpectNonNil(err, t)
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_Size_Is1_WithOneSet(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	sut.Set("rosie", "posey")

	// Verify
	ExpectInt(1, sut.Size(), t)
}

func TestThat_HashMap_IsEmpty_IsTrue_WhenNew(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Verify
	ExpectBool(true, sut.IsEmpty(), t)
}

func TestThat_HashMap_IsEmpty_IsFalse_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewHashMap()
	sut.Set("scooby", "dooby")

	// Verify
	ExpectBool(false, sut.IsEmpty(), t)
}

func TestThat_HashMap_Size_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Verify
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_Merge_AddsNothing_ForEmptyMaps(t *testing.T) {
	// Setup
	sut := NewHashMap()
	otherMap := NewHashMap()

	// Test
	sut.Merge(otherMap)

	// Verify
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_Merge_AddsEntries_ForNonEmptyMaps(t *testing.T) {
	// Setup
	key := "tweedle"
	expected := "deedle"
	other := NewHashMap()
	other.Set(key, expected)
	sut := NewHashMap()

	// Test
	sut.Merge(other)

	// Verify
	ExpectBool(false, sut.IsEmpty(), t)
	ExpectInt(1, sut.Size(), t)
	ExpectBool(true, sut.Has(key), t)
	actual := sut.Get(key)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_HashMap_Merge_DoesNothing_WhenNilForEmptySet(t *testing.T) {
	// Setup
	var sut *HashMap	// nil

	// Test
	sut.Merge(NewHashMap())

	// Verify
	ExpectNil(sut, t)
}

func TestThat_HashMap_Merge_DoesNothing_WhenNilForNonEmptySet(t *testing.T) {
	// Setup
	//defer ExpectPanic(t)
	other := NewHashMap()
	other.Set("beep", "boop")
	var sut *HashMap	// nil

	// Test
	sut.Merge(other)
}

func TestThat_HashMap_Set_DoesNothing_WhenNil(t *testing.T) {
	// Setup
	//defer ExpectPanic(t)
	var sut *HashMap	// nil

	// Test
	sut.Set("testkey", "testvalue")
}

func TestThat_HashMap_Set_OverwritesValue_ForExistingKey(t *testing.T) {
	// Setup
	sut := NewHashMap()
	key := "somekey"
	sut.Set(key, "unexpected")

	// Test
	expected := "expected"
	sut.Set(key, expected)

	// Verify
	actual := sut.Get(key)
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_HashMap_Get_ReturnsValue_ForSetKey(t *testing.T) {
	// Setup
	sut := NewHashMap()
	key := "testkey"
	expected := "testvalue"

	// Test
	sut.Set(key, expected)
	actual := sut.Get(key)

	// Verify
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_HashMap_Get_ReturnsNil_ForUnsetKey(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	actual := sut.Get("boguskey")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_HashMap_GetInt64_ReturnsValue_ForSetKeyWithParseableInt(t *testing.T) {
	// Setup
	sut := NewHashMap()
	key := "testkey"
	var expected int64 = 1234567

	// Test
	sut.Set(key, fmt.Sprintf("%d", expected))
	actual := sut.GetInt64(key)

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt64(expected, *actual, t)
}

func TestThat_HashMap_GetInt64_ReturnsNil_ForSetKeyWithoutParseableInt(t *testing.T) {
	// Setup
	sut := NewHashMap()
	key := "testkey"

	// Test
	sut.Set(key, "bogusvalue")
	actual := sut.GetInt64(key)

	// Verify
	ExpectNil(actual, t)
}

func TestThat_HashMap_GetInt64_ReturnsNil_ForUnsetKey(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	actual := sut.GetInt64("boguskey")

	// Verify
	ExpectNil(actual, t)
}

func TestThat_HashMap_GetKeys_ReturnsNothing_ForEmptyMap(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	actual := sut.GetKeys()

	// Verify
	ExpectInt(0, len(actual), t)
}

func TestThat_HashMap_GetKeys_ReturnsKeys_ForNonEmptyMap(t *testing.T) {
	// Setup
	sut := NewHashMap()
	sut.Set("k1", "v1")
	sut.Set("k2", "v2")

	// Test
	actual := sut.GetKeys()

	// Verify
	ExpectInt(2, len(actual), t)
	ExpectTrue(((actual[0]=="k1")&&(actual[1]=="k2"))||((actual[1]=="k1")&&(actual[0]=="k2")), t)
}

func TestThat_HashMap_Has_IsFalse_WhenKeyMissing(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Verify
	ExpectBool(false, sut.Has("boguskey"), t)
}

func TestThat_HashMap_Has_IsTrue_WhenKeyExists(t *testing.T) {
	// Setup
	sut := NewHashMap()
	key := "testkey"

	// Test
	sut.Set(key, "testvalue")

	// Verify
	ExpectBool(true, sut.Has(key), t)
}

func TestThat_HashMap_HasAll_ReturnsTrue_WhenKeysEmptySet(t *testing.T) {
	// Setup
	sut := NewHashMap()
	keys := make([]string, 0)

	// Test
	ExpectBool(true, sut.HasAll(&keys), t)
}

func TestThat_HashMap_HasAll_ReturnsFalse_WhenAnyKeyMissing(t *testing.T) {
	// Setup
	sut := NewHashMap()
	keys := make([]string, 1)
	keys[0] = "missingkey"

	// Test
	ExpectBool(false, sut.HasAll(&keys), t)
}

func TestThat_HashMap_HasAll_ReturnsTrue_WhenAllKeysExist(t *testing.T) {
	// Setup
	sut := NewHashMap()
	sut.Set("alrighty", "then")
	sut.Set("okey", "dokey")
	keys := make([]string, 2)
	keys[0] = "alrighty"
	keys[1] = "okey"

	// Test
	ExpectBool(true, sut.HasAll(&keys), t)
}

func TestThat_HashMap_GetIterator_ReturnsIterator_WhenEmpty(t *testing.T) {
	// Setup
	sut := NewHashMap()

	// Test
	it := sut.GetIterator()

	// Verify
	ExpectNonNil(it, t)
	item := it()
	ExpectNil(item, t)
}

func TestThat_HashMap_GetIterator_ReturnsWorkingIterator_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := NewHashMap()
	expectedItems := make(map[string]string)
	numItems := 10
	for i :=0; i < numItems; i++ {
		key := fmt.Sprintf("k%d", i)
		value := fmt.Sprintf("v%d", i)
		sut.Set(key, value)
		expectedItems[key] = value
	}

	// Test
	it := sut.GetIterator()

	// Verify
	ExpectNonNil(it, t)
	// Iterator will return items in whatever natural order the map sees fit; we must look up and eliminate each one
	for itemi := it(); itemi != nil; itemi = it() {
		ExpectNonNil(itemi, t)
		if kvp, ok := itemi.(*KeyValuePair); ok {
			ExpectTrue(ok, t)
			actualKey := kvp.Key
			actualValue := kvp.Value
			if expectedValue, isExpectedKey := expectedItems[actualKey]; isExpectedKey {
				ExpectTrue(isExpectedKey, t)
				ExpectString(expectedValue, actualValue, t)
				delete(expectedItems, actualKey)
			} else {
				ExpectTrue(false, t)
			}
		} else {
			ExpectTrue(false, t)
		}
	}
	ExpectInt(0, len(expectedItems), t)
}

func TestThat_HashMap_ToJson_ReturnsEmptyJson_ForEmptyMap(t *testing.T) {
	// Setup
	sut := NewHashMap()
	expected := "{}"

	// Test
	actual, err := sut.ToJson()

	// Verify
	ExpectNil(err, t)
	ExpectString(expected, *actual, t)
}

func TestThat_HashMap_ToJson_ReturnsPopulatedJson_ForNonEmptyMap(t *testing.T) {
	// Setup
	sut := NewHashMap()
	sut.Set("k1", "v1")
	sut.Set("k2", "v2")
	expected := "{\"k1\":\"v1\",\"k2\":\"v2\"}"

	// Test
	actual, err := sut.ToJson()

	// Verify
	ExpectNil(err, t)
	ExpectString(expected, *actual, t)
}
