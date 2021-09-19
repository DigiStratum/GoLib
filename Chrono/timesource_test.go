package chrono

import(
	"time"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

const TEST_MSEC_STEP = 10

func TestThat_Timesource_NowUnixTimeStamp_UpdatesOncePerSecond(t *testing.T) {
	// Setup
	sut := NewTimeSource()

	// Test
	res, t1 := waitForSecondChange(sut)
	ExpectTrue(res, t)

	msecElapsed := 0
	// Count how many milliseconds until TimeStamp changes again...
	for ; msecElapsed < 1500; msecElapsed += TEST_MSEC_STEP {
		t2 := sut.NowUnixTimeStamp()
		if t2 > t1 { break; }
		time.Sleep(TEST_MSEC_STEP * time.Millisecond)
	}

	// Verity
	// Expect 1000msec between TimeStamp changes +/- 10 msec for test imprecision
	ExpectTrue(
		(msecElapsed >= (1000 - TEST_MSEC_STEP)) && (msecElapsed <= (1000 + TEST_MSEC_STEP)),
		t,
	)
}

func TestThat_Timesource_Now_Returns_GoodTimeStamp(t *testing.T) {
	// Setup
	sut := NewTimeSource()

	// Test
	res, _ := waitForSecondChange(sut)
	ExpectTrue(res, t)
	ts := sut.Now()
	ExpectNonNil(ts, t)
	actual := ts.ToUnixTimeStamp()
	expected := time.Now().Unix()

	// Verify
	ExpectInt64(expected, actual, t)
}

// Get us to within 5msec of the next change of TimeStamp according to this TimeSource
func waitForSecondChange(ts *TimeSource) (bool, int64) {
	t1 := ts.NowUnixTimeStamp()
	t2 := t1
	maxIter := 250
	for ; maxIter > 0; maxIter-- {
		t2 = ts.NowUnixTimeStamp()
		if t2 > t1 { return true, t2 }
		time.Sleep(5 * time.Millisecond)
	}
	// If >= 250 * 5msec transition from one second to next, then time is broken
	return false, -1
}
