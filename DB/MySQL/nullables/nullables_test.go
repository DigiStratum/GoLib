package nullables

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Initialization with unsupported data types (mainly "derived types")

func TestThat_NewNullable_Returns_Nil_ForStruct(t *testing.T) {
	// Setup
	type miniStruct struct {
		testProp	string
	}
	testStruct := miniStruct{ testProp: "test value" }

	// Test
	sut := NewNullable(testStruct)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_Returns_Nil_ForArray(t *testing.T) {
	// Setup
	testArray := make([]int, 5)

	// Test
	sut := NewNullable(testArray)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_Returns_Nil_ForPointer(t *testing.T) {
	// Setup
	value := "test value"

	// Test
	sut := NewNullable(&value)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_Returns_Nil_ForFunc(t *testing.T) {
	// Setup
	localFunc := func () bool { return true }

	// Test
	sut := NewNullable(localFunc)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_Returns_Nil_ForMap(t *testing.T) {
	// Setup
	testMap := make(map[int]int)

	// Test
	sut := NewNullable(testMap)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_Returns_Nil_ForChannel(t *testing.T) {
	// Setup
	testChan := make(chan int)

	// Test
	sut := NewNullable(testChan)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_Returns_Nil_ForSlice(t *testing.T) {
	// Setup
	var testSlice []int

	// Test
	sut := NewNullable(testSlice)

	// Verify
	ExpectNil(sut, t)
}

// Initialization with supported data types

func TestThat_NewNullable_Returns_NullableNil_ForNil(t *testing.T) {
	// Test
	sut := NewNullable(nil)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(sut.IsNil(), t)
	ExpectTrue(sut.GetType() == NULLABLE_NIL, t)
}

