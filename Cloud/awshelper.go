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

	lib "github.com/DigiStratum/GoLib"
)

type AWSHelperIfc interface {
	GetSession() *session.Session
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

// Get our AWS session
func (r *AWSHelper) GetSession() *session.Session {
	if nil == r.awsSession {
		sess, err := session.NewSession(
			&aws.Config{ Region: aws.String(r.awsRegion) },
		)
		if nil != err {
			lib.GetLogger().Error(fmt.Sprintf(
				"Failed to establish an AWS session in region '%s': '%s'",
				r.awsRegion,
				err.Error(),
			))
			return nil
		}
		r.awsSession = sess
	}
	return r.awsSession
}

