package cloud

/*
Cloud helper library for AWS services.

AWS, of course, has a vast, sprawling API with more endpoints, capabilities, and details that one
might care to count. There are, however, a number of helpful boilerplate type operations that are
broadly applicable which we will capture here.

ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
ref: https://github.com/aws/aws-sdk-go/tree/main/aws

*/

import (
        "fmt"

        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/aws/credentials"

	cfg "github.com/DigiStratum/GoLib/Config"
)

type AWSHelperIfc interface {
	GetSession() (*session.Session, error)
}

type AWSHelper struct {
	awsSession		*session.Session
	awsRegion		string
	awsAccessKeyId		string
	awsSecretAccessKeyId	string
	awsSessionToken		string
}

// Make a new one of these
func NewAWSHelper() *AWSHelper {
	return &AWSHelper{}
}

// -------------------------------------------------------------------------------------------------
// cfg.ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *AWSHelper) Configure(config cfg.ConfigIfc) error {
	if nil == config { return fmt.Errorf("AWSHelper.Configure() - Configuration was nil") }

	if config.Has("awsRegion") {
		awsRegion := config.Get("awsRegion")
		if nil != awsRegion { r.awsRegion = *awsRegion }
	}

	if config.Has("awsAccessKeyId") {
		awsAccessKeyId := config.Get("awsAccessKeyId")
		if nil != awsAccessKeyId { r.awsAccessKeyId = *awsAccessKeyId }
	}

	if config.Has("awsSecretAccessKeyId") {
		awsSecretAccessKeyId := config.Get("awsSecretAccessKeyId")
		if nil != awsSecretAccessKeyId { r.awsSecretAccessKeyId = *awsSecretAccessKeyId }
	}

	if config.Has("awsSessionToken") {
		awsSessionToken := config.Get("awsSessionToken")
		if nil != awsSessionToken { r.awsSessionToken = *awsSessionToken }
	}

	return nil
}

// -------------------------------------------------------------------------------------------------
// AWSHelperIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get our AWS session
func (r *AWSHelper) GetSession() (*session.Session, error) {
	if nil == r.awsSession {
		sess, err := session.NewSession(
			&aws.Config{
				Region: aws.String(r.awsRegion),
				Credentials: credentials.NewStaticCredentials(
					r.awsAccessKeyId,
					r.awsSecretAccessKeyId,
					r.awsSessionToken,
				),
			},
		)
		if nil != err {
			return nil, fmt.Errorf(
				"Failed to establish an AWS session in region '%s': '%s'",
				r.awsRegion,
				err.Error(),
			)
		}
		r.awsSession = sess
	}
	return r.awsSession, nil
}
