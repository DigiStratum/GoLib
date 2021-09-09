package json

/*

Unit Tests for Json

*/

import(
	"fmt"
	"testing"

	test "github.com/DigiStratum/GoTools/test"
	lib "github.com/DigiStratum/GoLib"
)

func TestThat_Json_NewJson_ReturnsInstance(t *testing.T) {
	// Setup
	sut := lib.NewJson(nil)

	// Verify
	test.ExpectNonNil(sut, t)
}

func TestThat_Json_NewJsonFromFile_ReturnsInstance(t *testing.T) {
	// Setup
	sut := lib.NewJsonFromFile("")

	// Verify
	test.ExpectNonNil(sut, t)
}

func TestThat_Json_Load_ReturnsError_ForNilJsonString(t *testing.T) {
	// Setup
	sut := lib.NewJson(nil)
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	test.ExpectNonNil(err, t)
}

func TestThat_Json_Load_ReturnsError_ForBadJsonString(t *testing.T) {
	// Setup
	jsonString := "</not json>"
	sut := lib.NewJson(&jsonString)
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	test.ExpectNonNil(err, t)
}

func TestThat_Json_Load_Works_ForGoodJsonString(t *testing.T) {
	// Setup
	name := "boogie"
	value := "woogie"
	jsonString := fmt.Sprintf("{\"%s\": \"%s\"}", name, value)
	sut := lib.NewJson(&jsonString)
	var target map[string]string

	// Test
	err := sut.Load(&target)
	res, ok := target[name]

	// Verify
	test.ExpectNil(err, t)
	test.ExpectNonNil(target, t)
	test.ExpectBool(true, ok, t)
	test.ExpectString(value, res, t)
}

func TestThat_Json_Load_ReturnsError_ForBadFilePath(t *testing.T) {
	// Setup
	sut := lib.NewJsonFromFile("")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	test.ExpectNonNil(err, t)
}

func TestThat_Json_Load_ReturnsError_ForBadJsonFile(t *testing.T) {
	// Setup
	sut := lib.NewJsonFromFile("json_test.bad.json")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	test.ExpectNonNil(err, t)
}

func TestThat_Json_Load_Works_ForGoodJsonFile(t *testing.T) {
	// Setup
	sut := lib.NewJsonFromFile("json_test.good.json")
	var target map[string]string

	// Test
	err := sut.Load(&target)

	// Verify
	test.ExpectNil(err, t)
	test.ExpectNonNil(target, t)
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		res, ok := target[name]
		test.ExpectBool(true, ok, t)
		test.ExpectString(value, res, t)
	}
}

func TestThat_Json_LoadOrPanic_Panics_ForBadJsonString(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	jsonString := "</not json>"
	sut := lib.NewJson(&jsonString)
	var target map[string]string

	// Test
	sut.LoadOrPanic(&target)
}

func TestThat_Json_LoadOrPanic_Works_ForGoodJsonString(t *testing.T) {
	// Setup
	name := "boogie"
	value := "woogie"
	jsonString := fmt.Sprintf("{\"%s\": \"%s\"}", name, value)
	sut := lib.NewJson(&jsonString)
	var target map[string]string

	// Test
	sut.LoadOrPanic(&target)
	res, ok := target[name]

	// Verify
	test.ExpectNonNil(target, t)
	test.ExpectBool(true, ok, t)
	test.ExpectString(value, res, t)
}

func TestThat_Json_LoadOrPanic_Panics_ForBadJsonFile(t *testing.T) {
	// Setup
	defer test.ExpectPanic(t)
	sut := lib.NewJsonFromFile("json_test.bad.json")
	var target map[string]string

	// Test
	sut.LoadOrPanic(&target)
}

func TestThat_Json_LoadOrPanic_Works_ForGoodJsonFile(t *testing.T) {
	// Setup
	sut := lib.NewJsonFromFile("json_test.good.json")
	var target map[string]string

	// Test
	sut.LoadOrPanic(&target)

	// Verify
	test.ExpectNonNil(target, t)
	for i := 1; i <= 5; i++ {
		name := fmt.Sprintf("key%d", i)
		value := fmt.Sprintf("value%d", i)
		res, ok := target[name]
		test.ExpectBool(true, ok, t)
		test.ExpectString(value, res, t)
	}
}

