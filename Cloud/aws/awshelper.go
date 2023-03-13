package aws

/*
Cloud helper library for AWS services.

AWS, of course, has a vast, sprawling API with more endpoints, capabilities, and details that one
might care to count. There are, however, a number of helpful boilerplate type operations that are
broadly applicable which we will capture here.

ref: https://docs.awssdk.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
ref: https://github.com/aws/aws-sdk-go/tree/main/aws

TODO:
 * Validate config values however we can; awsRegion, for example, can only be so many things
*/

import (
        "fmt"

        awssdk "github.com/aws/aws-sdk-go/aws"
        awssdksession "github.com/aws/aws-sdk-go/aws/session"
        awssdkcredentials "github.com/aws/aws-sdk-go/aws/credentials"

	cfg "github.com/DigiStratum/GoLib/Config"
	"github.com/DigiStratum/GoLib/Starter"
)

type AWSHelperIfc interface {
	// Embedded Interface(s)
	starter.StartableIfc
	cfg.ConfigurableIfc

	// Out own interface
	GetSession() (*awssdksession.Session, error)
}

type AWSHelper struct {
	*starter.Startable
	*cfg.Configurable
	awsSession		*awssdksession.Session
	awsRegion		string
	awsAccessKeyId		string
	awsSecretAccessKeyId	string
	awsSessionToken		string
}

// Make a new one of these
func NewAWSHelper() *AWSHelper {
	awsh := AWSHelper{ }

	// Declare Configuration
	awsh.Configurable = cfg.NewConfigurable(
		cfg.NewConfigItem("awsRegion").CaptureWith(awsh.captureConfigAwsRegion),
		cfg.NewConfigItem("awsAccessKeyId").CaptureWith(awsh.captureConfigAwsAccessKeyId),
		cfg.NewConfigItem("awsSecretAccessKeyId").CaptureWith(awsh.captureConfigAwsSecretAccessKeyId),
		cfg.NewConfigItem("awsSessionToken").CaptureWith(awsh.captureConfigAwsSessionToken),
	)

	awsh.Startable = starter.NewStartable(
		awsh.Configurable,
	)

	return &awsh
}

// -------------------------------------------------------------------------------------------------
// StartableIfc
// -------------------------------------------------------------------------------------------------

func (r *AWSHelper) Start() error {
	return r.Startable.Start()
}

// -------------------------------------------------------------------------------------------------
// ConfigurableIfc
// -------------------------------------------------------------------------------------------------

// Optionally accept overrides for defaults in configuration
func (r *AWSHelper) Configure(config cfg.ConfigIfc) error {
	// If we have already been configured, do not accept a second configuration
	if r.Startable.IsStarted() { return nil }

	return r.Configurable.Configure(config)
}

func (r *AWSHelper) captureConfigAwsRegion(value string) error {
	r.awsRegion = value
	return nil
}

func (r *AWSHelper) captureConfigAwsAccessKeyId(value string) error {
	r.awsAccessKeyId = value
	return nil
}

func (r *AWSHelper) captureConfigAwsSecretAccessKeyId(value string) error {
	r.awsSecretAccessKeyId = value
	return nil
}

func (r *AWSHelper) captureConfigAwsSessionToken(value string) error {
	r.awsSessionToken = value
	return nil
}

// -------------------------------------------------------------------------------------------------
// AWSHelperIfc
// -------------------------------------------------------------------------------------------------

// Get our AWS session
func (r *AWSHelper) GetSession() (*awssdksession.Session, error) {
	if nil == r.awsSession {
		config := awssdk.Config{}
		if len(r.awsRegion) > 0 {
			config.Region = awssdk.String(r.awsRegion)
		}
		if ((len(r.awsAccessKeyId) > 0) ||
			(len(r.awsSecretAccessKeyId) > 0) ||
			(len(r.awsSessionToken) > 0)) {
			config.Credentials = awssdkcredentials.NewStaticCredentials(
				r.awsAccessKeyId,
				r.awsSecretAccessKeyId,
				r.awsSessionToken,
			)
		}
		sess, err := awssdksession.NewSession(&config)
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
