package golib_test

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

	test "github.com/DigiStratum/GoTools/test"
	lib "github.com/DigiStratum/GoLib"
)

func TestThat_HashMap_Size_Is0_WhenNew(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()

	// Verify
	test.ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_Size_Is1_WithOneSet(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()

	// Test
	sut.Set("rosie", "posey")

	// Verify
	test.ExpectInt(1, sut.Size(), t)
}

func TestThat_HashMap_IsEmpty_IsTrue_WhenNew(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()

	// Verify
	test.ExpectBool(true, sut.IsEmpty(), t)
}

func TestThat_HashMap_IsEmpty_IsFalse_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	sut.Set("scooby", "dooby")

	// Verify
	test.ExpectBool(false, sut.IsEmpty(), t)
}

func TestThat_HashMap_Has_IsFalse_WhenKeyMissing(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()

	// Verify
	test.ExpectBool(false, sut.Has("boguskey"), t)
}

func TestThat_HashMap_Has_IsTrue_WhenKeyExists(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	key := "testkey"

	// Test
	sut.Set(key, "testvalue")

	// Verify
	test.ExpectBool(true, sut.Has(key), t)
}

func TestThat_HashMap_Get_ReturnsValue_ForSetKey(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	key := "testkey"
	value := "testvalue"

	// Test
	sut.Set(key, value)

	// Verify
	test.ExpectString(value, sut.Get(key), t)
}

func TestThat_HashMap_Get_ReturnsEmptyString_ForUnsetKey(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()

	// Verify
	test.ExpectString("", sut.Get("boguskey"), t)
}

func TestThat_HashMap_Merge_AddsNothing_ForEmptyMaps(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	otherMap := lib.NewHashMap()

	// Test
	sut.Merge(otherMap)

	// Verify
	test.ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_Merge_AddsEntries_ForNonEmptyMaps(t *testing.T) {
	// Setup
	key := "tweedle"
	value := "deedle"
	other := lib.NewHashMap()
	other.Set(key, value)
	sut := lib.NewHashMap()

	// Test
	sut.Merge(other)

	// Verify
	test.ExpectBool(false, sut.IsEmpty(), t)
	test.ExpectInt(1, sut.Size(), t)
	test.ExpectBool(true, sut.Has(key), t)
	test.ExpectString(value, sut.Get(key), t)
}

func TestThat_HashMap_Merge_ChangesNothing_WhenNilForEmptySet(t *testing.T) {
	// Setup
	var sut *lib.HashMap	// nil

	// Test
	sut.Merge(lib.NewHashMap())

	// Verify
	test.ExpectNil(sut, t)
}

func TestThat_HashMap_Merge_Panics_WhenNilForNonEmptySet(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	other := lib.NewHashMap()
	other.Set("beep", "boop")
	var sut *lib.HashMap	// nil

	// Test
	sut.Merge(other)
}

func TestThat_HashMap_Set_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	sut.Set("testkey", "testvalue")
}

func TestThat_HashMap_IsEmpty_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	sut.IsEmpty()
}

func TestThat_HashMap_Size_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	sut.Size()
}

func TestThat_HashMap_Get_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	sut.Get("boguskey")
}

func TestThat_HashMap_Has_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	sut.Has("boguskey")
}

func TestThat_HashMap_HasAll_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	sut.HasAll(nil)
}

func TestThat_HashMap_HasAll_Panics_WhenKeysNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	sut := lib.NewHashMap()

	// Test
	sut.HasAll(nil)
}

func TestThat_HashMap_HasAll_ReturnsTrue_WhenKeysEmptySet(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	keys := make([]string, 0)

	// Test
	test.ExpectBool(true, sut.HasAll(&keys), t)
}

func TestThat_HashMap_HasAll_ReturnsFalse_WhenAnyKeyMissing(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	keys := make([]string, 1)
	keys[0] = "missingkey"

	// Test
	test.ExpectBool(false, sut.HasAll(&keys), t)
}

func TestThat_HashMap_HasAll_ReturnsTrue_WhenAllKeysExist(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	sut.Set("alrighty", "then")
	sut.Set("okey", "dokey")
	keys := make([]string, 2)
	keys[0] = "alrighty"
	keys[1] = "okey"

	// Test
	test.ExpectBool(true, sut.HasAll(&keys), t)
}

func TestThat_HashMap_Set_OverwritesValue_ForExistingKey(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	key := "somekey"
	sut.Set(key, "unexpected")

	// Test
	expected := "expected"
	sut.Set(key, expected)

	// Verify
	test.ExpectString(expected, sut.Get(key), t)
}

func TestThat_HashMap_GetCopy_ReturnsNil_WhenNil(t *testing.T) {
	// Setup
	var sut *lib.HashMap	// nil

	// Verify
	test.ExpectNil(sut.GetCopy(), t)
}

func TestThat_HashMap_GetCopy_ReturnsEmpty_WhenEmpty(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()

	// Test
	res := sut.GetCopy()

	// Verify
	test.ExpectInt(0, res.Size(), t)
}

func TestThat_HashMap_GetCopy_ReturnsNonEmpty_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	num := 25
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		sut.Set(key, value)
	}

	// Test
	res := sut.GetCopy()

	// Verify
	test.ExpectInt(num, res.Size(), t)
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		test.ExpectString(value, res.Get(key), t)
	}
}

func TestThat_HashMap_IterateCallback_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	sut.IterateCallback(func (kvp lib.KeyValuePair) {})
}

func TestThat_HashMap_IterateCallback_MakesNoCalls_WhenEmpty(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	callbackCounter := 0

	// Test
	sut.IterateCallback(func (kvp lib.KeyValuePair) { callbackCounter++ })

	// Verify
	test.ExpectInt(0, callbackCounter, t)
}

func TestThat_HashMap_IterateCallback_MakesOneCallPerKey_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	num := 25
	keys := make([]string, num)
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		sut.Set(key, value)
		keys[i] = key
	}
	type kvpdata struct {
		Value	string
		Num	int
	}
	callbacks := make(map[string]*kvpdata)

	// Test
	sut.IterateCallback(func (kvp lib.KeyValuePair) {
		if _, ok := callbacks[kvp.Key]; ! ok {
			callbacks[kvp.Key] = &kvpdata{ Value: kvp.Value, Num: 0 }
		}
		(*callbacks[kvp.Key]).Num++
	})

	// Verify
	test.ExpectInt(num, len(callbacks), t)
	for i := 0; i < num; i++ {
		test.ExpectInt(1, callbacks[keys[i]].Num, t)
		test.ExpectString(sut.Get(keys[i]), callbacks[keys[i]].Value, t)
	}
}

// TODO: Add HashMap.IterateChannel() coverage

func TestThat_HashMap_IterateChannel_Panics_WhenNil(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	var sut *lib.HashMap	// nil

	// Test
	for _ = range sut.IterateChannel() {}
}

func TestThat_HashMap_IterateChannel_YieldsNoEntries_WhenEmpty(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	entryCounter := 0

	// Test
	for _ = range sut.IterateChannel() { entryCounter++ }

	// Verify
	test.ExpectInt(0, entryCounter, t)
}

func TestThat_HashMap_IterateChannel_YieldsOneEntryPerKey_WhenNonEmpty(t *testing.T) {
	// Setup
	sut := lib.NewHashMap()
	num := 25
	keys := make([]string, num)
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("key-%d", i)
		value := fmt.Sprintf("value-%d", i)
		sut.Set(key, value)
		keys[i] = key
	}
	type kvpdata struct {
		Value	string
		Num	int
	}
	entries := make(map[string]*kvpdata)

	// Test
	for kvp := range sut.IterateChannel() {
		if _, ok := entries[kvp.Key]; ! ok {
			entries[kvp.Key] = &kvpdata{ Value: kvp.Value, Num: 0 }
		}
		(*entries[kvp.Key]).Num++
	}

	// Verify
	test.ExpectInt(num, len(entries), t)
	for i := 0; i < num; i++ {
		test.ExpectInt(1, entries[keys[i]].Num, t)
		test.ExpectString(sut.Get(keys[i]), entries[keys[i]].Value, t)
	}
}

