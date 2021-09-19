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

func TestThat_TimeStamp_NewTimeStamp_ReturnsSomething(t *testing.T) {
	// Setup
	ts := NewTimeSource()
	sut := NewTimeStamp(ts)

	// Verify
	ExpectNonNil(sut, t)
}
