package nullables

import(
	"time"
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

func TestThat_NewNullable_Returns_Nil_ForComplex(t *testing.T) {
	// Setup
	var c64 complex64
	var c128 complex128

	// Test
	sutc64 := NewNullable(c64)
	sutc128 := NewNullable(c128)

	// Verify
	ExpectNil(sutc64, t)
	ExpectNil(sutc128, t)
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

func TestThat_NewNullable_Returns_NullableInt64_ForInts(t *testing.T) {
	// Test
	var i int
	var i8 int8
	var i16 int16
	var i32 int32
	var i64 int64
	suti := NewNullable(i)
	suti8 := NewNullable(i8)
	suti16 := NewNullable(i16)
	suti32 := NewNullable(i32)
	suti64 := NewNullable(i64)

	// Verify
	ExpectNonNil(suti, t)
	ExpectFalse(suti.IsNil(), t)
	ExpectTrue(suti.GetType() == NULLABLE_INT64, t)
	ExpectNonNil(suti8, t)
	ExpectFalse(suti8.IsNil(), t)
	ExpectTrue(suti8.GetType() == NULLABLE_INT64, t)
	ExpectNonNil(suti16, t)
	ExpectFalse(suti16.IsNil(), t)
	ExpectTrue(suti16.GetType() == NULLABLE_INT64, t)
	ExpectNonNil(suti32, t)
	ExpectFalse(suti32.IsNil(), t)
	ExpectTrue(suti32.GetType() == NULLABLE_INT64, t)
	ExpectNonNil(suti64, t)
	ExpectFalse(suti64.IsNil(), t)
	ExpectTrue(suti64.GetType() == NULLABLE_INT64, t)
}

func TestThat_NewNullable_Returns_NullableFloat64_ForFloats(t *testing.T) {
	// Test
	var f32 float32
	var f64 float64
	sutf32 := NewNullable(f32)
	sutf64 := NewNullable(f64)

	// Verify
	ExpectNonNil(sutf32, t)
	ExpectFalse(sutf32.IsNil(), t)
	ExpectTrue(sutf32.GetType() == NULLABLE_FLOAT64, t)
	ExpectNonNil(sutf64, t)
	ExpectFalse(sutf64.IsNil(), t)
	ExpectTrue(sutf64.GetType() == NULLABLE_FLOAT64, t)
}

func TestThat_NewNullable_Returns_NullableBool_ForBool(t *testing.T) {
	// Test
	sut := NewNullable(true)

	// Verify
	ExpectNonNil(sut, t)
	ExpectFalse(sut.IsNil(), t)
	ExpectTrue(sut.GetType() == NULLABLE_BOOL, t)
}

func TestThat_NewNullable_Returns_NullableString_ForString(t *testing.T) {
	// Test
	sut := NewNullable("so stringy!")

	// Verify
	ExpectNonNil(sut, t)
	ExpectFalse(sut.IsNil(), t)
	ExpectTrue(sut.GetType() == NULLABLE_STRING, t)
}

func TestThat_NewNullable_Returns_NullableTime_ForTime(t *testing.T) {
	// Test
	sut := NewNullable(time.Now())

	// Verify
	ExpectNonNil(sut, t)
	ExpectFalse(sut.IsNil(), t)
	ExpectTrue(sut.GetType() == NULLABLE_TIME, t)
}

