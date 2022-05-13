package nullables

import(
	"fmt"
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

// Getters

func TestThat_GetInt64_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetInt64()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_GetInt64_Returns_ValuePointer_ForInt64Value(t *testing.T) {
	// Setup
	sut := NewNullable(333)

	// Test
	actual := sut.GetInt64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt64(*actual, 333, t)
	ExpectTrue(sut.IsInt64(), t)
}

func TestThat_GetInt64_Returns_ValuePointer_ForBoolValue(t *testing.T) {
	// Setup
	sut0 := NewNullable(false)
	sut1 := NewNullable(true)

	// Test
	actual0 := sut0.GetInt64()
	actual1 := sut1.GetInt64()

	// Verify
	ExpectNonNil(actual0, t)
	ExpectInt64(*actual0, 0, t)
	ExpectTrue(sut0.IsBool(), t)
	ExpectNonNil(actual1, t)
	ExpectInt64(*actual1, 1, t)
	ExpectTrue(sut1.IsBool(), t)
}

func TestThat_GetInt64_Returns_ValuePointer_ForFloatValue(t *testing.T) {
	// Setup
	sut0 := NewNullable(0.0)
	sut1 := NewNullable(1.1)

	// Test
	actual0 := sut0.GetInt64()
	actual1 := sut1.GetInt64()

	// Verify
	ExpectNonNil(actual0, t)
	ExpectInt64(*actual0, 0, t)
	ExpectTrue(sut0.IsFloat64(), t)
	ExpectNonNil(actual1, t)
	ExpectInt64(*actual1, 1, t)
	ExpectTrue(sut1.IsFloat64(), t)
}

func TestThat_GetInt64_Returns_ValuePointer_ForStringValue(t *testing.T) {
	// Setup
	sut := NewNullable("333")

	// Test
	actual := sut.GetInt64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt64(*actual, 333, t)
	ExpectTrue(sut.IsString(), t)
}

func TestThat_GetInt64_Returns_ValuePointer_ForTimeValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-12 19:50:0 (which was when this test was added)
	sut := NewNullable(time.Date(2022, 5, 12, 19, 50, 0, 0, time.UTC))

	// Test
	actual := sut.GetInt64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt64(*actual, 1652385000, t)
	ExpectTrue(sut.IsTime(), t)
}

func TestThat_GetBool_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetBool()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_GetBool_Returns_ValuePointer_ForBoolValue(t *testing.T) {
	// Setup
	sut0 := NewNullable(false)
	sut1 := NewNullable(true)

	// Test
	actual0 := sut0.GetBool()
	actual1 := sut1.GetBool()

	// Verify
	ExpectNonNil(actual0, t)
	ExpectFalse(*actual0, t)
	ExpectTrue(sut0.IsBool(), t)
	ExpectNonNil(actual1, t)
	ExpectTrue(*actual1, t)
	ExpectTrue(sut1.IsBool(), t)
}

func TestThat_GetBool_Returns_ValuePointer_ForInt64Value(t *testing.T) {
	// Setup
	sut0 := NewNullable(0)
	sut1 := NewNullable(1)

	// Test
	actual0 := sut0.GetBool()
	actual1 := sut1.GetBool()

	// Verify
	ExpectNonNil(actual0, t)
	ExpectFalse(*actual0, t)
	ExpectTrue(sut0.IsInt64(), t)
	ExpectNonNil(actual1, t)
	ExpectTrue(*actual1, t)
	ExpectTrue(sut1.IsInt64(), t)
}

func TestThat_GetBool_Returns_ValuePointer_ForFloat64Value(t *testing.T) {
	// Setup
	sut0 := NewNullable(0.0)
	sut1 := NewNullable(1.1)

	// Test
	actual0 := sut0.GetBool()
	actual1 := sut1.GetBool()

	// Verify
	ExpectNonNil(actual0, t)
	ExpectFalse(*actual0, t)
	ExpectTrue(sut0.IsFloat64(), t)
	ExpectNonNil(actual1, t)
	ExpectTrue(*actual1, t)
	ExpectTrue(sut1.IsFloat64(), t)
}

func TestThat_GetBool_Returns_ValuePointer_ForStringValue(t *testing.T) {
	// Setup
	sut0 := NewNullable("not true")
	sut1 := NewNullable("true")

	// Test
	actual0 := sut0.GetBool()
	actual1 := sut1.GetBool()

	// Verify
	ExpectNonNil(actual0, t)
	ExpectFalse(*actual0, t)
	ExpectTrue(sut0.IsString(), t)
	ExpectNonNil(actual1, t)
	ExpectTrue(*actual1, t)
	ExpectTrue(sut1.IsString(), t)
}

func TestThat_GetBool_Returns_ValuePointer_ForTimeValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-13 01:47:0 (which was when this test was added)
	sut := NewNullable(time.Date(2022, 5, 13, 01, 47, 0, 0, time.UTC))

	// Test
	actual := sut.GetBool()

	// Verify
	ExpectNonNil(actual, t)
	ExpectTrue(*actual, t)
	ExpectTrue(sut.IsTime(), t)
}

func TestThat_GetFloat64_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_GetFloat64_Returns_ValuePointer_ForIntValue(t *testing.T) {
	// Setup
	sut := NewNullable(1)

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNonNil(actual, t)
	var expected float64 = 1.0
	ExpectFloat64(expected, *actual, t)
}

func TestThat_GetFloat64_Returns_ValuePointer_ForBoolValue(t *testing.T) {
	// Setup
	sut0 := NewNullable(false)
	sut1 := NewNullable(true)

	// Test
	actual0 := sut0.GetFloat64()
	actual1 := sut1.GetFloat64()

	// Verify
	ExpectNonNil(actual0, t)
	var expected0 float64 = 0.0
	ExpectFloat64(expected0, *actual0, t)
	ExpectNonNil(actual1, t)
	var expected1 float64 = 2.0
	ExpectFloat64(expected1, *actual1, t)
}

func TestThat_GetFloat64_Returns_ValuePointer_ForFloatValue(t *testing.T) {
	// Setup
	var expected float64 = 1.1
	sut := NewNullable(expected)

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectFloat64(expected, *actual, t)
}

func TestThat_GetFloat64_Returns_ValuePointer_ForStringValue(t *testing.T) {
	// Setup
	var expected float64 = 1.1
	sut := NewNullable(fmt.Sprintf("%f", expected))

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectFloat64(expected, *actual, t)
}

func TestThat_GetFloat64_Returns_Nil_ForTimeValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-13 06:00:0 (which was when this test was added)
	sut := NewNullable(time.Date(2022, 5, 13, 06, 00, 0, 0, time.UTC))

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNil(actual, t)
}

