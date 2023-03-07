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
	"math"
	"testing"
	"runtime"
	"reflect"
	"regexp"
)

func getCaller() string {
	_, file, no, ok := runtime.Caller(2)
	if ok { return fmt.Sprintf("@%s:%d", file, no) }
	return ""
}

func ExpectEmptyString(actual string, t *testing.T) bool {
	if 0 == len(actual) { return true }
	t.Errorf("\n\n%s:\nExpect len('%s') == 0, Actual: %d", getCaller(), actual, len(actual))
	return false
}

func ExpectNonEmptyString(actual string, t *testing.T) bool {
	if 0 < len(actual) { return true }
	t.Errorf("\n\n%s:\nExpect len('%s') > 0, Actual: %d", getCaller(), actual, len(actual))
	return false
}

func ExpectString(expect, actual string, t *testing.T) bool {
	if expect == actual { return true }
	t.Errorf("\n\n%s:\nExpect: '%s', Actual: '%s'", getCaller(), expect, actual)
	return false
}

func ExpectInt(expect, actual int, t *testing.T) bool {
	if expect == actual { return true }
	t.Errorf("\n\n%s\nExpect: '%d', Actual: '%d'", getCaller(), expect, actual)
	return false
}

func ExpectInt64(expect, actual int64, t *testing.T) bool {
	if expect == actual { return true }
	t.Errorf("\n\n%s\nExpect: '%d', Actual: '%d'", getCaller(), expect, actual)
	return false
}

func ExpectFloat64(expect, actual float64, t *testing.T) bool {
	if math.Abs(expect - actual) < 0.0001 { return true }
	t.Errorf("\n\n%s\nExpect: '%f', Actual: '%f'", getCaller(), expect, actual)
	return false
}

func ExpectBool(expect, actual bool, t *testing.T) bool {
	if expect == actual { return true }
	t.Errorf("\n\n%s:\nExpect: '%t', Actual: '%t'", getCaller(), expect, actual)
	return false
}

func ExpectTrue(actual bool, t *testing.T) bool {
	if true == actual { return true }
	t.Errorf("\n\n%s:\nExpect: 'true', Actual: '%t'", getCaller(), actual)
	return false
}

func ExpectFalse(actual bool, t *testing.T) bool {
	if false == actual { return true }
	t.Errorf("\n\n%s:\nExpect: 'false', Actual: '%t'", getCaller(), actual)
	return false
}

// ref: https://stackoverflow.com/questions/31595791/how-to-test-panics
func ExpectPanic(t *testing.T) bool {
	if r := recover(); r == nil {
		t.Errorf("\n\n%s:\nExpect: Panic, Actual: Nope!", getCaller())
		return false
	}
	return true
}

func ExpectNil(value interface{}, t *testing.T) bool {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		return true
	}
	t.Errorf("\n\n%s:\nExpect: nil, Actual: non-nil", getCaller())
	return false
}

func ExpectNonNil(value interface{}, t *testing.T) bool {
	if value == nil || (reflect.ValueOf(value).Kind() == reflect.Ptr && reflect.ValueOf(value).IsNil()) {
		t.Errorf("\n\n%s:\nExpect: non-nil, Actual: nil", getCaller())
		return false
	}
	return true
}

func ExpectError(err error, t *testing.T) bool {
	if nil != err {	return true }
	t.Errorf("\n\n%s:\nExpect: error, Actual: nil error", getCaller())
	return false
}

func ExpectNoError(err error, t *testing.T) bool {
	if nil == err {	return true }
	t.Errorf("\n\n%s:\nExpect: nil error, Actual: error('%s')", getCaller(), err.Error())
	return false
}

func ExpectMatch(expect, actual string, t *testing.T) bool {
	matched, err := regexp.MatchString(expect, actual)
	if nil == err && matched { return true }
	t.Errorf("\n\n%s:\nExpect: pattern [%s] match actual [%s], Actual: no match", getCaller(), expect, actual)
	return false
}

