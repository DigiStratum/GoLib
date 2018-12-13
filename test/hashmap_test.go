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

func TestThat_HashMap_IsEmptyWhenNew(t *testing.T) {
	sut := lib.HashMap{}
	ExpectInt(0, sut.Size(), t)
}

func TestThat_HashMap_GetsSetValues(t *testing.T) {
	sut := lib.HashMap{}
	key := "test_key"
	value := "test_value"
	sut.Set(key, value)
	ExpectString(value, sut.Get(key), t)
}

