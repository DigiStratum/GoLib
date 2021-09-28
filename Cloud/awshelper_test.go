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
	awsRegion := "us-west-2"
	sut := NewAWSHelper(awsRegion)

	// Verify
	ExpectNonNil(sut, t)
}
