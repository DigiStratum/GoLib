package chrono

import(
	"time"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

const TEST_MSEC_STEP = 10

func TestThat_TimeSource_NewTimeSource_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewTimeSource()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_Timesource_NowUnixTimeStamp_UpdatesOncePerSecond(t *testing.T) {
	// Setup
	sut := NewTimeSource()

	// Test
	t1 := sut.NowUnixTimeStamp()
	time.Sleep(1 * time.Second)
	t2 := sut.NowUnixTimeStamp()

	// Verify
	ExpectInt64(1, (t2 - t1), t)
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
	maxIter := 600
	for ; maxIter > 0; maxIter-- {
		t2 := ts.NowUnixTimeStampMilli()
		msec :=  t2 % 1000
		if msec <= 20 {
//fmt.Printf("Time approached second boundary with %d msec over in %d cycles!\n", msec, maxIter)
			return true, t2 / 1000
		}
		time.Sleep(2 * time.Millisecond)
	}
	// If >= 250 * 5msec transition from one second to next, then time is broken
//fmt.Printf("Time is broken!\n")
	return false, -1
}
