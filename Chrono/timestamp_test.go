package chrono

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_TimeStamp_NewTimeStamp_ReturnsNothing_ForNilTimesource(t *testing.T) {
	// Setup
	sut := NewTimeStamp(nil)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_TimeStamp_NewTimeStamp_ReturnsSomething_ForGoodTimeSource(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewTimeStamp(ts)

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_TimeStamp_NewFromUnixTimeStamp_ReturnsNothing_ForNilTimesource(t *testing.T) {
	// Setup
	sut := NewFromUnixTimeStamp(nil, 0)

	// Verify
	ExpectNil(sut, t)
}

func TestThat_TimeStamp_NewFromUnixTimeStamp_ReturnsSomething_ForGoodTimeSource(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewFromUnixTimeStamp(ts, 0)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(sut.ToUnixTimeStamp() == 0, t)
}

func TestThat_TimeStamp_NewTimeStampForever_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewTimeStampForever()
	ts := NewTimeSource()
	now := NewTimeStamp(ts)

	// Verify
	ExpectNonNil(sut, t)
	ExpectTrue(sut.Compare(now) == 1, t)
	ExpectTrue(sut.Diff(now) == 1, t)
	ExpectFalse(sut.IsPast(), t)
	ExpectTrue(sut.IsFuture(), t)
	ExpectTrue(sut.CompareToNow() == 1, t)
	ExpectTrue(sut.DiffNow() == 1, t)
	ExpectTrue(sut.ToUnixTimeStamp() == 0, t)
	ExpectTrue(sut.IsForever(), t)
	ExpectTrue(sut.Add(1000).ToUnixTimeStamp() == 0, t)
}

func TestThat_TimeStamp_Add_ReturnsGoodTimeStamp_ForPositiveValue(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewFromUnixTimeStamp(ts, 0)

	// Test
	sut.Add(1000)

	// Verify
	ExpectTrue(sut.ToUnixTimeStamp() == 1000, t)
}

func TestThat_TimeStamp_Add_ReturnsGoodTimeStamp_ForNegativeValue(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewFromUnixTimeStamp(ts, 0)

	// Test
	sut.Add(-1000)

	// Verify
	ExpectTrue(sut.ToUnixTimeStamp() == -1000, t)
}

func TestThat_TimeStamp_Comparisons_ReturnNegative_WhenTimeStampIsOlderThanSupplied(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewTimeStamp(ts)

	// Test
	sut.Add(-1000)

	// Verify
	ExpectInt(-1, sut.Compare(NewTimeStamp(ts)), t)
	ExpectTrue(sut.CompareToNow() == -1, t)
	ExpectTrue(sut.Diff(NewTimeStamp(ts)) == -1000, t)
	ExpectTrue(sut.IsPast(), t)
	ExpectFalse(sut.IsFuture(), t)
	// FIXME: run at just the wrong can cause these checks to fail; add delay to wait for 1-second boundary
	ExpectTrue(sut.DiffNow() == -1000, t)
}

func TestThat_TimeStamp_Comparisons_ReturnPositive_WhenTimeStampIsNewerThanSupplied(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewTimeStamp(ts)

	// Test
	sut.Add(1000)

	// Verify
	ExpectTrue(sut.Compare(NewTimeStamp(ts)) == 1, t)
	ExpectTrue(sut.CompareToNow() == 1, t)
	ExpectTrue(sut.Diff(NewTimeStamp(ts)) == 1000, t)
	ExpectFalse(sut.IsPast(), t)
	ExpectTrue(sut.IsFuture(), t)
	// FIXME: run at just the wrong can cause these checks to fail; add delay to wait for 1-second boundary
	ExpectTrue(sut.DiffNow() == 1000, t)
}

func TestThat_TimeStamp_Comparisons_ReturnZero_WhenTimeStampIsSameAsSupplied(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewTimeStamp(ts)
	timeStamp := NewFromUnixTimeStamp(ts, sut.ToUnixTimeStamp())

	// Verify
	ExpectTrue(sut.Compare(timeStamp) == 0, t)
	ExpectTrue(sut.Diff(timeStamp) == 0, t)
	ExpectFalse(sut.IsPast(), t)
	ExpectFalse(sut.IsFuture(), t)
	// FIXME: run at just the wrong can cause these checks to fail; add delay to wait for 1-second boundary
	ExpectTrue(sut.CompareToNow() == 0, t)
	ExpectTrue(sut.DiffNow() == 0, t)
}
