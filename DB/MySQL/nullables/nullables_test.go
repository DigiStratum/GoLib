package nullables

import(
	"fmt"
	"time"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

// Initialization with unsupported data types (mainly "derived types")

func TestThat_NewNullable_ReturnsNil_ForStruct(t *testing.T) {
	// Setup
	type miniStruct struct {
		testProp	string
	}
	testStruct := miniStruct{ testProp: "test value" }

	// Test
	var sut *Nullable = NewNullable(testStruct)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_ReturnsNil_ForArray(t *testing.T) {
	// Setup
	testArray := make([]int, 5)

	// Test
	var sut *Nullable = NewNullable(testArray)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_ReturnsNil_ForPointer(t *testing.T) {
	// Setup
	value := "test value"

	// Test
	var sut *Nullable = NewNullable(&value)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_ReturnsNil_ForFunc(t *testing.T) {
	// Setup
	localFunc := func () bool { return true }

	// Test
	var sut *Nullable = NewNullable(localFunc)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_ReturnsNil_ForMap(t *testing.T) {
	// Setup
	testMap := make(map[int]int)

	// Test
	var sut *Nullable = NewNullable(testMap)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_ReturnsNil_ForChannel(t *testing.T) {
	// Setup
	testChan := make(chan int)

	// Test
	var sut *Nullable = NewNullable(testChan)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_ReturnsNil_ForSlice(t *testing.T) {
	// Setup
	var testSlice []int

	// Test
	var sut *Nullable = NewNullable(testSlice)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_NewNullable_ReturnsNil_ForComplex(t *testing.T) {
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

// Checkers

func TestThat_Nullable_IsNil_Returns_ExpectedResult_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Now())
	suts := NewNullable("super stringy!")

	// Test
	actualn := sutn.IsNil()
	actuali := suti.IsNil()
	actualb := sutb.IsNil()
	actualf := sutf.IsNil()
	actualt := sutt.IsNil()
	actuals := suts.IsNil()

	// Verify
	ExpectTrue(actualn, t)
	ExpectFalse(actuali, t)
	ExpectFalse(actualb, t)
	ExpectFalse(actualf, t)
	ExpectFalse(actualt, t)
	ExpectFalse(actuals, t)
}

func TestThat_Nullable_IsInt64_Returns_ExpectedResult_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Now())
	suts := NewNullable("super stringy!")

	// Test
	actualn := sutn.IsInt64()
	actuali := suti.IsInt64()
	actualb := sutb.IsInt64()
	actualf := sutf.IsInt64()
	actualt := sutt.IsInt64()
	actuals := suts.IsInt64()

	// Verify
	ExpectFalse(actualn, t)
	ExpectTrue(actuali, t)
	ExpectFalse(actualb, t)
	ExpectFalse(actualf, t)
	ExpectFalse(actualt, t)
	ExpectFalse(actuals, t)
}

func TestThat_Nullable_IsBool_Returns_ExpectedResult_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Now())
	suts := NewNullable("super stringy!")

	// Test
	actualn := sutn.IsBool()
	actuali := suti.IsBool()
	actualb := sutb.IsBool()
	actualf := sutf.IsBool()
	actualt := sutt.IsBool()
	actuals := suts.IsBool()

	// Verify
	ExpectFalse(actualn, t)
	ExpectFalse(actuali, t)
	ExpectTrue(actualb, t)
	ExpectFalse(actualf, t)
	ExpectFalse(actualt, t)
	ExpectFalse(actuals, t)
}

func TestThat_Nullable_IsFloat64_Returns_ExpectedResult_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Now())
	suts := NewNullable("super stringy!")

	// Test
	actualn := sutn.IsFloat64()
	actuali := suti.IsFloat64()
	actualb := sutb.IsFloat64()
	actualf := sutf.IsFloat64()
	actualt := sutt.IsFloat64()
	actuals := suts.IsFloat64()

	// Verify
	ExpectFalse(actualn, t)
	ExpectFalse(actuali, t)
	ExpectFalse(actualb, t)
	ExpectTrue(actualf, t)
	ExpectFalse(actualt, t)
	ExpectFalse(actuals, t)
}

func TestThat_Nullable_IsTime_Returns_ExpectedResult_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Now())
	suts := NewNullable("super stringy!")

	// Test
	actualn := sutn.IsTime()
	actuali := suti.IsTime()
	actualb := sutb.IsTime()
	actualf := sutf.IsTime()
	actualt := sutt.IsTime()
	actuals := suts.IsTime()

	// Verify
	ExpectFalse(actualn, t)
	ExpectFalse(actuali, t)
	ExpectFalse(actualb, t)
	ExpectFalse(actualf, t)
	ExpectTrue(actualt, t)
	ExpectFalse(actuals, t)
}

func TestThat_Nullable_IsString_Returns_ExpectedResult_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Now())
	suts := NewNullable("super stringy!")

	// Test
	actualn := sutn.IsString()
	actuali := suti.IsString()
	actualb := sutb.IsString()
	actualf := sutf.IsString()
	actualt := sutt.IsString()
	actuals := suts.IsString()

	// Verify
	ExpectFalse(actualn, t)
	ExpectFalse(actuali, t)
	ExpectFalse(actualb, t)
	ExpectFalse(actualf, t)
	ExpectFalse(actualt, t)
	ExpectTrue(actuals, t)
}

// Getters

func TestThat_Nullable_GetInt64_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetInt64()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Nullable_GetInt64_Returns_ValuePointer_ForInt64Value(t *testing.T) {
	// Setup
	sut := NewNullable(333)

	// Test
	actual := sut.GetInt64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt64(*actual, 333, t)
	ExpectTrue(sut.IsInt64(), t)
}

func TestThat_Nullable_GetInt64_Returns_ValuePointer_ForBoolValue(t *testing.T) {
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

func TestThat_Nullable_GetInt64_Returns_ValuePointer_ForFloatValue(t *testing.T) {
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

func TestThat_Nullable_GetInt64_Returns_ValuePointer_ForStringValue(t *testing.T) {
	// Setup
	sut := NewNullable("333")

	// Test
	actual := sut.GetInt64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectInt64(*actual, 333, t)
	ExpectTrue(sut.IsString(), t)
}

func TestThat_Nullable_GetInt64_Returns_ValuePointer_ForTimeValue(t *testing.T) {
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

func TestThat_Nullable_GetBool_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetBool()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Nullable_GetBool_Returns_ValuePointer_ForBoolValue(t *testing.T) {
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

func TestThat_Nullable_GetBool_Returns_ValuePointer_ForInt64Value(t *testing.T) {
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

func TestThat_Nullable_GetBool_Returns_ValuePointer_ForFloat64Value(t *testing.T) {
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

func TestThat_Nullable_GetBool_Returns_ValuePointer_ForStringValue(t *testing.T) {
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

func TestThat_Nullable_GetBool_Returns_ValuePointer_ForTimeValue(t *testing.T) {
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

func TestThat_Nullable_GetFloat64_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Nullable_GetFloat64_Returns_ValuePointer_ForIntValue(t *testing.T) {
	// Setup
	sut := NewNullable(1)

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNonNil(actual, t)
	var expected float64 = 1.0
	ExpectFloat64(expected, *actual, t)
}

func TestThat_Nullable_GetFloat64_Returns_ValuePointer_ForBoolValue(t *testing.T) {
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

func TestThat_Nullable_GetFloat64_Returns_ValuePointer_ForFloatValue(t *testing.T) {
	// Setup
	var expected float64 = 1.1
	sut := NewNullable(expected)

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectFloat64(expected, *actual, t)
}

func TestThat_Nullable_GetFloat64_Returns_ValuePointer_ForStringValue(t *testing.T) {
	// Setup
	var expected float64 = 1.1
	sut := NewNullable(fmt.Sprintf("%f", expected))

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNonNil(actual, t)
	ExpectFloat64(expected, *actual, t)
}

func TestThat_Nullable_GetFloat64_Returns_Nil_ForTimeValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-13 06:00:0 (which was when this test was added)
	sut := NewNullable(time.Date(2022, 5, 13, 06, 00, 0, 0, time.UTC))

	// Test
	actual := sut.GetFloat64()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Nullable_GetString_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetString()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Nullable_GetString_Returns_ValuePointer_ForIntValue(t *testing.T) {
	// Setup
	expected := 333
	sut := NewNullable(expected)

	// Test
	actual := sut.GetString()

	// Verify
	ExpectNonNil(actual, t)
	ExpectString(fmt.Sprintf("%d", expected), *actual, t)
}

func TestThat_Nullable_GetString_Returns_ValuePointer_ForBoolValue(t *testing.T) {
	// Setup
	sut0 := NewNullable(false)
	sut1 := NewNullable(true)

	// Test
	actual0 := sut0.GetString()
	actual1 := sut1.GetString()

	// Verify
	ExpectNonNil(actual0, t)
	ExpectString("false", *actual0, t)
	ExpectNonNil(actual1, t)
	ExpectString("true", *actual1, t)
}

func TestThat_Nullable_GetString_Returns_ValuePointer_ForFloatValue(t *testing.T) {
	// Setup
	sut := NewNullable(1.1)

	// Test
	actual := sut.GetString()

	// Verify
	ExpectNonNil(actual, t)
	ExpectString("1.1E+00", *actual, t)
}

func TestThat_Nullable_GetString_Returns_ValuePointer_ForStringValue(t *testing.T) {
	// Setup
	expected := "super stringy!"
	sut := NewNullable(expected)

	// Test
	actual := sut.GetString()

	// Verify
	ExpectNonNil(actual, t)
	ExpectString(expected, *actual, t)
}

func TestThat_Nullable_GetString_Returns_Nil_ForTimeValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-13 06:48:0 (which was when this test was added)
	sut := NewNullable(time.Date(2022, 5, 13, 06, 48, 0, 0, time.UTC))

	// Test
	actual := sut.GetString()

	// Verify
	ExpectNonNil(actual, t)
	ExpectString("2022-05-13T06:48:00Z", *actual, t)
}

func TestThat_Nullable_GetTime_Returns_Nil_ForNilValue(t *testing.T) {
	// Setup
	sut := NewNullable(nil)

	// Test
	actual := sut.GetTime()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Nullable_GetTime_Returns_ValuePointer_ForIntValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-13 07:12:0 (which was when this test was added)
	sut := NewNullable(1652451120)

	// Test
	actual := sut.GetTime()

	// Verify
	ExpectNonNil(actual, t)
	actualStr := (*actual).Format("2006-01-02T15:04:05Z")
	ExpectString("2022-05-13T07:12:00Z", actualStr, t)
}

func TestThat_Nullable_GetTime_Returns_Nil_ForBoolValue(t *testing.T) {
	// Setup
	sut := NewNullable(true)

	// Test
	actual := sut.GetTime()

	// Verify
	ExpectNil(actual, t)
}

func TestThat_Nullable_GetTime_Returns_ValuePointer_ForFloatValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-13 07:12:0 (which was when this test was added)
	sut := NewNullable(1652451120.0)

	// Test
	actual := sut.GetTime()

	// Verify
	ExpectNonNil(actual, t)
	actualStr := (*actual).Format("2006-01-02T15:04:05Z")
	ExpectString("2022-05-13T07:12:00Z", actualStr, t)
}

func TestThat_Nullable_GetTime_Returns_ValuePointer_ForStringValue(t *testing.T) {
	// Setup
	// UTC Timestamp for 2022-05-13 07:29:0 (which was when this test was added)
	sut := NewNullable("2022-05-13T07:29:00Z")

	// Test
	actual := sut.GetTime()

	// Verify
	ExpectNonNil(actual, t)
	actualStr := (*actual).Format("2006-01-02T15:04:05Z")
	ExpectString("2022-05-13T07:29:00Z", actualStr, t)
}

func TestThat_Nullable_GetTime_Returns_ValuePointer_ForTimeValue(t *testing.T) {
	// Setup
	expected := time.Now()
	sut := NewNullable(expected)

	// Test
	actual := sut.GetTime()

	// Verify
	ExpectNonNil(actual, t)
	expectedStr := (expected).Format("2006-01-02T15:04:05Z")
	actualStr := (*actual).Format("2006-01-02T15:04:05Z")
	ExpectString(expectedStr, actualStr, t)
}

func TestThat_Nullable_MarshalJSON_Returns_JSONByteSliceWithoutError_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Date(2022, 5, 14, 17, 21, 0, 0, time.UTC))
	suts := NewNullable("super stringy!")

	// Test
	actualn, err1 := sutn.MarshalJSON()
	actuali, err2 := suti.MarshalJSON()
	actualb, err3 := sutb.MarshalJSON()
	actualf, err4 := sutf.MarshalJSON()
	actualt, err5 := sutt.MarshalJSON()
	actuals, err6 := suts.MarshalJSON()

	// Verify
	ExpectNoError(err1, t)
	ExpectString("null", string(actualn), t)
	ExpectNoError(err2, t)
	ExpectString("333", string(actuali), t)
	ExpectNoError(err3, t)
	ExpectString("true", string(actualb), t)
	ExpectNoError(err4, t)
	ExpectString("3.3", string(actualf), t)
	ExpectNoError(err5, t)
	ExpectString("\"2022-05-14T17:21:00Z\"", string(actualt), t)
	ExpectNoError(err6, t)
	ExpectString("\"super stringy!\"", string(actuals), t)
}

func TestThat_Nullable_UmmarshalJSON_Returns_JSONByteSliceWithoutError_ForEachNullableType(t *testing.T) {
	// Setup
	sutn := NewNullable(nil)
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Date(2022, 5, 14, 17, 44, 0, 0, time.UTC))
	suts := NewNullable("super stringy!")

	// Test
	err1 := sutn.UnmarshalJSON([]byte("null"))
	err2 := suti.UnmarshalJSON([]byte("444"))
	actuali := suti.GetInt64()
	err3 := sutb.UnmarshalJSON([]byte("false"))
	actualb := sutb.GetBool()
	err4 := sutf.UnmarshalJSON([]byte("4.4"))
	actualf := sutf.GetFloat64()
	err5 := sutt.UnmarshalJSON([]byte("\"2000-01-02T12:34:56Z\""))
	actualt := sutt.GetTime()
	err6 := suts.UnmarshalJSON([]byte("\"silly string!\""))
	actuals := suts.GetString()

	// Verify
	ExpectNoError(err1, t)
	ExpectTrue(sutn.IsNil(), t)

	ExpectNoError(err2, t)
	ExpectNonNil(actuali, t)
	ExpectTrue(suti.IsInt64(), t)
	ExpectInt64(444, *actuali, t)

	ExpectNoError(err3, t)
	ExpectNonNil(actualb, t)
	ExpectTrue(sutb.IsBool(), t)
	ExpectFalse(*actualb, t)

	ExpectNoError(err4, t)
	ExpectNonNil(actualf, t)
	ExpectTrue(sutf.IsFloat64(), t)
	ExpectFloat64(4.4, *actualf, t)

	ExpectNoError(err5, t)
	ExpectNonNil(actualt, t)
	ExpectTrue(sutt.IsTime(), t)
	actualStr := (*actualt).Format("2006-01-02T15:04:05Z")
	ExpectString("2000-01-02T12:34:56Z", actualStr, t)

	ExpectNoError(err6, t)
	ExpectNonNil(actuals, t)
	ExpectTrue(suts.IsString(), t)
	ExpectString("silly string!", *actuals, t)
}

func TestThat_Nullable_Scan_Returns_WithoutError_ForEachNullableType(t *testing.T) {
	// Setup
	var expectedn *string = nil
	sutn := NewNullable(nil)
	var expectedi int64 = 444
	suti := NewNullable(333)
	expectedb := false
	sutb := NewNullable(true)
	var expectedf float64 = 4.4
	sutf := NewNullable(3.3)
	expectedt := time.Date(2000, 1, 2, 12, 34, 56, 0, time.UTC)
	sutt := NewNullable(time.Date(2022, 5, 14, 17, 44, 0, 0, time.UTC))
	expecteds := "silly string!"
	suts := NewNullable("super stringy!")

	// Test
	err1 := sutn.Scan(expectedn)
	err2 := suti.Scan(expectedi)
	actuali := suti.GetInt64()
	err3 := sutb.Scan(expectedb)
	actualb := sutb.GetBool()
	err4 := sutf.Scan(expectedf)
	actualf := sutf.GetFloat64()
	err5 := sutt.Scan(expectedt)
	actualt := sutt.GetTime()
	err6 := suts.Scan(expecteds)
	actuals := suts.GetString()

	// Verify
	ExpectNoError(err1, t)
	ExpectTrue(sutn.IsNil(), t)

	ExpectNoError(err2, t)
	ExpectNonNil(actuali, t)
	ExpectTrue(suti.IsInt64(), t)
	ExpectInt64(expectedi, *actuali, t)

	ExpectNoError(err3, t)
	ExpectNonNil(actualb, t)
	ExpectTrue(sutb.IsBool(), t)
	ExpectBool(expectedb, *actualb, t)

	ExpectNoError(err4, t)
	ExpectNonNil(actualf, t)
	ExpectTrue(sutf.IsFloat64(), t)
	ExpectFloat64(expectedf, *actualf, t)

	ExpectNoError(err5, t)
	ExpectNonNil(actualt, t)
	ExpectTrue(sutt.IsTime(), t)
	actualStr := (*actualt).Format("2006-01-02T15:04:05Z")
	expectedStr := (expectedt).Format("2006-01-02T15:04:05Z")
	ExpectString(expectedStr, actualStr, t)

	ExpectNoError(err6, t)
	ExpectNonNil(actuals, t)
	ExpectTrue(suts.IsString(), t)
	ExpectString(expecteds, *actuals, t)
}

func TestThat_Nullable_Scan_Returns_Error_ForEachNullableType(t *testing.T) {
	// Setup
	// Note: NULLABLE_NIL always works without error
	suti := NewNullable(333)
	sutb := NewNullable(true)
	sutf := NewNullable(3.3)
	sutt := NewNullable(time.Date(2022, 5, 14, 17, 44, 0, 0, time.UTC))
	suts := NewNullable("super stringy!")

	// Test
	err1 := suti.Scan(nil) // <- not an int!
	actuali := suti.GetInt64()
	err2 := sutb.Scan(nil) // <- not a bool!
	actualb := sutb.GetBool()
	err3 := sutf.Scan(nil) // <- not a float!
	actualf := sutf.GetFloat64()
	err4 := sutt.Scan("not a time")
	actualt := sutt.GetTime()
	err5 := suts.Scan(nil) // <- not a string!
	actuals := suts.GetString()

	// Verify
	ExpectError(err1, t)
	ExpectNil(actuali, t)
	ExpectTrue(suti.IsInt64(), t)

	ExpectError(err2, t)
	ExpectNil(actualb, t)
	ExpectTrue(sutb.IsBool(), t)

	ExpectError(err3, t)
	ExpectNil(actualf, t)
	ExpectTrue(sutf.IsFloat64(), t)

	ExpectError(err4, t)
	ExpectNil(actualt, t)
	ExpectTrue(sutt.IsTime(), t)

	ExpectError(err5, t)
	ExpectNil(actuals, t)
	ExpectTrue(suts.IsString(), t)
}

