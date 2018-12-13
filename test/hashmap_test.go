package golib_test

/*

Unit Tests for HasMap

ref: https://blog.alexellis.io/golang-writing-unit-tests/

*/
import(
	"fmt"
	"testing"
	"runtime"

	lib "github.com/DigiStratum/GoLib"
)

// ----------------------------------------------------------------------------

func GetCaller() string {
	_, file, no, ok := runtime.Caller(2)
	if ok { return fmt.Sprintf("@%s:%d", file, no) }
	return ""
}

func ExpectString(expect, actual string, t *testing.T) {
	if expect == actual { return }
	t.Errorf("%s:\nExpect: '%s', Actual: '%s'", GetCaller(), expect, actual)
}

func ExpectInt(expect, actual int, t *testing.T) {
	if expect == actual { return }
	t.Errorf("%s\nExpect: '%d', Actual: '%d'", GetCaller(), expect, actual)
}

func ExpectBool(expect, actual bool, t *testing.T) {
	if expect == actual { return }
	t.Errorf("%s:\nExpect: '%t', Actual: '%t'", GetCaller(), expect, actual)
}

// ----------------------------------------------------------------------------

func TestThat_HashMap_Size_Is0_WhenNew(t *testing.T) {
	sut := lib.NewHashMap()
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_Size_Is1_WithOneSet(t *testing.T) {
	sut := lib.NewHashMap()
	sut.Set("rosie", "posey")
	ExpectInt(1, sut.Size(), t)
}

func TestThat_HashMap_IsEmpty_IsTrue_WhenNew(t *testing.T) {
	sut := lib.NewHashMap()
	ExpectBool(true, sut.IsEmpty(), t)
}

func TestThat_HashMap_IsEmpty_IsFalse_WhenNonEmpty(t *testing.T) {
	sut := lib.NewHashMap()
	sut.Set("scooby", "dooby")
	ExpectBool(false, sut.IsEmpty(), t)
}

func TestThat_HashMap_Get_ReturnsValue_ForSetKey(t *testing.T) {
	sut := lib.NewHashMap()
	key := "test_key"
	value := "test_value"
	sut.Set(key, value)
	ExpectString(value, sut.Get(key), t)
}

func TestThat_HashMap_Get_ReturnsEmptyString_ForUnsetKey(t *testing.T) {
	sut := lib.NewHashMap()
	ExpectString("", sut.Get("boguskey"), t)
}

func TestThat_HashMap_Merge_AddsNothing_ForEmptyMaps(t *testing.T) {
	sut := lib.NewHashMap()
	otherMap := lib.NewHashMap()
	sut.Merge(otherMap)
	ExpectInt(0, sut.Size(), t)
}

