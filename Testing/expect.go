package testing

/*

Expect*() helpers for unit tests

Use like so in unit tests:

// --------------------------
import(
	"testing"
	. "github.com/DigiStratum/GoLib/Testing"
)
// ...
func TestThat_Class_Method_DoesWhatever_WhenCondition(t *testing.T) {
	// Setup
	sut := NewClassInstance()

	// Test
	res := sut.Method()

	// Verify
	ExpectInt(0, res, t)
}
// --------------------------

*/

import(
	"fmt"
	"testing"
	"runtime"
	"reflect"
)

func getCaller() string {
	_, file, no, ok := runtime.Caller(2)
	if ok { return fmt.Sprintf("@%s:%d", file, no) }
	return ""
}

func ExpectString(expect, actual string, t *testing.T) {
	if expect == actual { return }
	t.Errorf("\n\n%s:\nExpect: '%s', Actual: '%s'", getCaller(), expect, actual)
}

func ExpectInt(expect, actual int, t *testing.T) {
	if expect == actual { return }
	t.Errorf("\n\n%s\nExpect: '%d', Actual: '%d'", getCaller(), expect, actual)
}

func ExpectInt64(expect, actual int64, t *testing.T) {
	if expect == actual { return }
	t.Errorf("\n\n%s\nExpect: '%d', Actual: '%d'", getCaller(), expect, actual)
}

func ExpectBool(expect, actual bool, t *testing.T) {
	if expect == actual { return }
	t.Errorf("\n\n%s:\nExpect: '%t', Actual: '%t'", getCaller(), expect, actual)
}

func ExpectTrue(actual bool, t *testing.T) {
	if true == actual { return }
	t.Errorf("\n\n%s:\nExpect: 'true', Actual: '%t'", getCaller(), actual)
}

func ExpectFalse(actual bool, t *testing.T) {
	if false == actual { return }
	t.Errorf("\n\n%s:\nExpect: 'false', Actual: '%t'", getCaller(), actual)
}

// ref: https://stackoverflow.com/questions/31595791/how-to-test-panics
func ExpectPanic(t *testing.T) {
	if r := recover(); r == nil {
		t.Errorf("\n\n%s:\nExpect: Panic, Actual: Nope!", getCaller())
	}
}

func ExpectNil(value interface{}, t *testing.T) {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		return
	}
	t.Errorf("\n\n%s:\nExpect: nil, Actual: non-nil", getCaller())
}

func ExpectNonNil(value interface{}, t *testing.T) {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		t.Errorf("\n\n%s:\nExpect: non-nil, Actual: nil", getCaller())
	}
}
