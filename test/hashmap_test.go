package golib_test

/*

Unit Tests for HasMap

ref: https://blog.alexellis.io/golang-writing-unit-tests/

*/
import(
	"testing"

	lib "github.com/DigiStratum/GoLib"
)

func TestThat_HashMap_GetsSetValues(t *testing.T) {
	sut := lib.HashMap{}
	key := "test_key"
	value := "test_value"
	sut.Set(key, value)
	res := sut.Get(key)
	if res != value {
		t.Errorf("Expected: '%s', Actual: '%s'", value, res)
	}
}

