package dsaws

/*

Unit Tests for Cloud

*/

import(
	//"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

	cfg "github.com/DigiStratum/GoLib/Config"
)

func TestThat_AWSHelper_NewAWSHelper_ReturnsSomething(t *testing.T) {
	// Setup
	sut := NewAWSHelper()

	// Verify
	ExpectNonNil(sut, t)
}

func TestThat_AWSHelper_ImplementsConfigurableIfc(t *testing.T) {
	// Setup
	var sut interface{} = NewAWSHelper()

	// Test
	_, ok := sut.(cfg.ConfigurableIfc)

	// Verify
	ExpectTrue(ok, t)
}

func TestThat_AWSHelper_Configure_AppliesConfigSettings(t *testing.T) {
	// Setup
	sut := NewAWSHelper()
	config := cfg.NewConfig()
	config.Set("awsRegion", "awsRegion")
	config.Set("awsAccessKeyId", "awsAccessKeyId")
	config.Set("awsSecretAccessKeyId", "awsSecretAccessKeyId")
	config.Set("awsSessionToken", "awsSessionToken")

	// Test
	err := sut.Configure(config)

	// Verify
	ExpectNil(err, t)
	ExpectString("awsRegion", sut.awsRegion, t)
	ExpectString("awsAccessKeyId", sut.awsAccessKeyId, t)
	ExpectString("awsSecretAccessKeyId", sut.awsSecretAccessKeyId, t)
	ExpectString("awsSessionToken", sut.awsSessionToken, t)
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
