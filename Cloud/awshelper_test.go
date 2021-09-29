package cloud

/*

Unit Tests for Cloud

*/

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
)

func TestThat_AWSHelper_NewAWSHelper_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewAWSHelper()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_AWSHelper_GetSession_ReturnsSessionNoError(t *testing.T) {
	// Setup
	sut := NewAWSHelper()

	// Test
	session, err := sut.GetSession()

	// Verify
	ExpectNil(err, t)
	ExpectNonNil(session, t)
}
