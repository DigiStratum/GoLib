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

TODO:
 * Consider using t.Fatalf() instead of t.Errorf() to stop test execution; this would eliminate the
   need for the bool return and the caller's responsibility to check it. It would, however force
   every expectation check to abort the test which means that only one failure per pass can be
   reported and corrected before running another pass to find the next failure - too annoying? Let's
   try it and find out! ref: https://ieftimov.com/posts/testing-in-go-failing-tests/

*/

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"testing"
)

func ExpectUnequal[T comparable](expected, actual T, t *testing.T) bool {
	if expected != actual {
		return true
	}
	t.Errorf("\n\n%s\nExpectUnequal: Expected: '%v', Actual: '%v'", getCaller(), expected, actual)
	return false
}

func ExpectEqual[T comparable](expected, actual T, t *testing.T) bool {
	if expected == actual {
		return true
	}
	t.Errorf("\n\n%s\nExpectEqual: Expected: '%v', Actual: '%v'", getCaller(), expected, actual)
	return false
}

func ExpectFloat[T float32 | float64](expected, actual T, t *testing.T) bool {
	// Note: Can't use Math.Abs() here because it doesn't work with generic types
	var diff T = expected - actual
	if diff < 0 {
		diff = -diff
	}
	if diff <= 0.0001 {
		return true
	}
	t.Errorf("\n\n%s\nExpectFloat: Expected: '%v', Actual: '%v'", getCaller(), expected, actual)
	return false
}

// Deprecated: Use ExpectEqual instead
func ExpectEmptyString(actual string, t *testing.T) bool {
	return ExpectEqual("", actual, t)
}

// Deprecated: Use ExpectUnequal instead
func ExpectNonEmptyString(actual string, t *testing.T) bool {
	return ExpectUnequal("", actual, t)
}

// Deprecated: Use ExpectEqual instead
func ExpectString(expect, actual string, t *testing.T) bool {
	return ExpectEqual(expect, actual, t)
}

// Deprecated: Use ExpectEqual instead
func ExpectInt(expect, actual int, t *testing.T) bool {
	return ExpectEqual(expect, actual, t)
}

// Deprecated: Use ExpectEqual instead
func ExpectInt64(expect, actual int64, t *testing.T) bool {
	return ExpectEqual(expect, actual, t)
}

// Deprecated: Use ExpectEqual instead
func ExpectBool(expect, actual bool, t *testing.T) bool {
	return ExpectEqual(expect, actual, t)
}

// Deprecated: Use ExpectEqual instead
func ExpectTrue(actual bool, t *testing.T) bool {
	return ExpectEqual(true, actual, t)
}

// Deprecated: Use ExpectEqual instead
func ExpectFalse(actual bool, t *testing.T) bool {
	return ExpectEqual(false, actual, t)
}

// Deprecated: Use ExpectFloat instead
func ExpectFloat64(expect, actual float64, t *testing.T) bool {
	return ExpectFloat(expect, actual, t)
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
	if nil != err {
		return true
	}
	t.Errorf("\n\n%s:\nExpect: error, Actual: nil error", getCaller())
	return false
}

func ExpectNoError(err error, t *testing.T) bool {
	if nil == err {
		return true
	}
	t.Errorf("\n\n%s:\nExpect: nil error, Actual: error('%s')", getCaller(), err.Error())
	return false
}

func ExpectMatch(pattern, actual string, t *testing.T) bool {
	matched, err := regexp.MatchString(pattern, actual)
	if nil == err && matched {
		return true
	}
	t.Errorf("\n\n%s:\nExpect: pattern [%s] match actual [%s], Actual: no match", getCaller(), pattern, actual)
	return false
}

func getCaller() string {
	_, file, no, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("@%s:%d", file, no)
	}
	return ""
}
