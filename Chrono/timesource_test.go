package chrono

import(
	"time"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_Timesource_NowUnixTimeStamp_UpdatesOncePerSecond(t *testing.T) {
	// Setup
	sut := NewTimeSource()

	// Test
	t1 := sut.NowUnixTimeStamp()
	t2 := t1
	for maxIter := 250; maxIter > 0; maxIter-- {
		t2 = sut.NowUnixTimeStamp()
		if t2 > t1 { break; }
		time.Sleep(5 * time.Millisecond)
	}
	msec := 0
	for ; msec < 1500; msec += 10 {
		t3 := sut.NowUnixTimeStamp()
		if t3 > t2 { break; }
		time.Sleep(10 * time.Millisecond)
	}

	// Verity
	ExpectTrue((msec >= 990) && (msec <= 1010),t)
}
