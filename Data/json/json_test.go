package json

/*

Unit Tests for Json

*/

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Json_NewJson_ReturnsInstance(t *testing.T) {
	// Setup
	sut := NewJson(nil)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Json_NewJsonFromFile_ReturnsInstance(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("")

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Json_Load_ReturnsError_ForNilJsonString(t *testing.T) {
	// Setup
	sut := NewJson(nil)
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	ExpectNonNil(err, t)
}

func TestThat_Json_Load_ReturnsError_ForBadJsonString(t *testing.T) {
	// Setup
	jsonString := "</not json>"
	sut := NewJson(&jsonString)
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	ExpectNonNil(err, t)
}

func TestThat_Json_Load_Works_ForGoodJsonString(t *testing.T) {
	// Setup
	name := "boogie"
	value := "woogie"
	jsonString := fmt.Sprintf("{\"%s\": \"%s\"}", name, value)
	sut := NewJson(&jsonString)
	var target map[string]string

	// Test
	err := sut.Load(&target)
	res, ok := target[name]

	// Verify
	ExpectNil(err, t)
	ExpectNonNil(target, t)
	ExpectBool(true, ok, t)
	ExpectString(value, res, t)
}

func TestThat_Json_Load_ReturnsError_ForBadFilePath(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	ExpectNonNil(err, t)
}

func TestThat_Json_Load_ReturnsError_ForBadJsonFile(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("json_test.bad.json")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	ExpectNonNil(err, t)
}

func TestThat_Json_Load_Works_ForGoodJsonFile(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("json_test.good.json")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	ExpectNil(err, t)
	ExpectNonNil(target, t)
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		res, ok := target[name]
		ExpectBool(true, ok, t)
		ExpectString(value, res, t)
	}
}
