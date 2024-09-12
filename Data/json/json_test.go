package json

/*

Unit Tests for Json

*/

import(
	"fmt"
	"testing"

	. "GoLib/Testing"
)

func TestThat_Json_NewJson_ReturnsInstance(t *testing.T) {
	// Setup
	sut := NewJson(nil)

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}

func TestThat_Json_NewJsonFromFile_ReturnsInstance(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("")

	// Verify
	if ! ExpectNonNil(sut, t) { return }
}

func TestThat_Json_Load_ReturnsError_ForNilJsonString(t *testing.T) {
	// Setup
	sut := NewJson(nil)
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	if ! ExpectNonNil(err, t) { return }
}

func TestThat_Json_Load_ReturnsError_ForBadJsonString(t *testing.T) {
	// Setup
	jsonString := "</not json>"
	sut := NewJson(&jsonString)
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	if ! ExpectNonNil(err, t) { return }
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
	if ! ExpectNil(err, t) { return }
	if ! ExpectNonNil(target, t) { return }
	if ! ExpectTrue(ok, t) { return }
	if ! ExpectString(value, res, t) { return }
}

func TestThat_Json_Load_ReturnsError_ForBadFilePath(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	if ! ExpectNonNil(err, t) { return }
}

func TestThat_Json_Load_ReturnsError_ForBadJsonFile(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("json_test.bad.json")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	if ! ExpectNonNil(err, t) { return }
}

func TestThat_Json_Load_Works_ForGoodJsonFile(t *testing.T) {
	// Setup
	sut := NewJsonFromFile("json_test.good.json")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	if ! ExpectNil(err, t) { return }
	if ! ExpectNonNil(target, t) { return }
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		res, ok := target[name]
		if ! ExpectTrue(ok, t) { return }
		if ! ExpectString(value, res, t) { return }
	}
}
