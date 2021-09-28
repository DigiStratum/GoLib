package cloud

/*
Cloud helper library for AWS services.

AWS, of course, has a vast, sprawling API with more endpoints, capabilities, and details that one
might care to count. There are, however, a number of helpful boilerplate type operations that are
broadly applicable which we will capture here.

*/

import (
        "fmt"

        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
)

type AWSHelperIfc interface {
	GetSession() (*session.Session, error)
}

type AWSHelper struct {
	awsRegion	string
	awsSession	*session.Session
}

// Make a new one of these
func NewAWSHelper(awsRegion string) *AWSHelper {
	return &AWSHelper{
		awsRegion: awsRegion,
	}
}

// -------------------------------------------------------------------------------------------------
// AWSHelperIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Get our AWS session
func (r *AWSHelper) GetSession() (*session.Session, error) {
	if nil == r.awsSession {
		sess, err := session.NewSession(
			&aws.Config{ Region: aws.String(r.awsRegion) },
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
