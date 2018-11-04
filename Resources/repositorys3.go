package resources

/*

Resource Repository for AWS S3 service

Ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html

S3 Resource storage mode adheres to our normal path model (see README.md), but with the additional
conditions that:

a) There exists an AWS / S3 service account
b) There exists a bucket within that S3 account
c) There exists a folder within that bucket within which the Resource paths are organized

This enables us to maintain any number of collections of Resources within a given S3 bucket by
separating the collections into different folders.

Configuration:
	* awsregion	- AWS Region identifier e.g. "us-west-1"
	* s3bucket	- AWS S3 Bucket to retrieve content from
	* s3folder	- AWS S3 Folder to prepend to any path (no trailing slash)

*/

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	lib "github.com/DigiStratum/GoLib"
)

type RepositoryS3 struct {
	repoConfig	*lib.Config
	awsSession	*session.Session
	awsS3		*s3.S3
	awsS3Downloader	*s3manager.Downloader
	readCache	ResourceMap
}

// Make a new one of these!
func NewRepositoryS3() *RepositoryS3 {
	r := RepositoryS3{
		readCache: make(ResourceMap),
	}
	return &r
}

// Satisfies RespositoryIfc
func (r *RepositoryS3) Configure(config *lib.Config) error {

	// Validate that the config has what we need for S3!
	requiredConfig := []string{ "awsregion", "s3bucket", "s3folder" }
	if ! (config.HasAll(&requiredConfig)) {
		return errors.New("Incomplete RepositoryS3 configuration provided")
	}

	r.repoConfig = config

	return nil
}

// Satisfies RepositoryIfc
// Ref: https://stackoverflow.com/questions/41645377/golang-s3-download-to-buffer-using-s3manager-downloader
func (r *RepositoryS3) GetResource(path string) *Resource {
	// If it's not yet in the cache
	if _, ok := r.readCache[path]; !ok {
		// Read the Resource from our S3 bucket into cache
		buff := &aws.WriteAtBuffer{}
		downloader := r.getS3Downloader()
		_, err := downloader.Download(
			buff,
			&s3.GetObjectInput{
				Bucket:	aws.String(r.repoConfig.Get("s3bucket")),
				Key:	aws.String(r.repoConfig.Get("s3folder") + "/" + path),
			},
		)
		// Error = no Resource!
		if nil != err { return nil }
		r.readCache[path] = NewResourceFromString(string(buff.Bytes()))
	}
	return r.readCache[path]
}

// Satisfies RepositoryIfc
func (r *RepositoryS3) HasResource(path string) bool {
	// If it's already in the cache, then we know we have it!
	if _, ok := r.readCache[path]; ok { return true }

	// If there's S3 metadata with no error, then there's a Resource!
	// ref: github.com/aws/aws-sdk-go/service/s3/examples_test.go ("HeadObject")
	awsS3 := r.getS3()
	_, err := awsS3.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(r.repoConfig.Get("s3bucket")),
			Key:	aws.String(r.repoConfig.Get("s3folder") + "/" + path),
		},
	)
	return nil == err
}

// Satisfies WritableRepositoryIfc
func (r *RepositoryS3) PutResource(resource *Resource, path string) error {
	// TODO: Actually implement WRITE operation to S3 here
	return errors.New("Not Yet Implemented!")
}

// Get our AWS session
func (r *RepositoryS3) getSession() *session.Session {
	if nil == r.awsSession {
		sess, err := session.NewSession(
			&aws.Config{ Region: aws.String(r.repoConfig.Get("awsregion")) },
		)
		if nil != err {
			l := lib.GetLogger()
			l.Error("Failed to establish and AWS session")
			return nil
		}
		r.awsSession = sess
	}
	return r.awsSession
}

// Get our S3 connection
func (r *RepositoryS3) getS3() *s3.S3 {
	if nil == r.awsS3 {
		sess := r.getSession();
		if nil == sess { return nil }
		r.awsS3 = s3.New(sess)
	}
	return r.awsS3
}

// Get our S3 Downloader
func (r *RepositoryS3) getS3Downloader() *s3manager.Downloader {
	if nil == r.awsS3Downloader {
		sess := r.getSession();
		if nil == sess { return nil }
		r.awsS3Downloader = s3manager.NewDownloader(sess)
	}
	return r.awsS3Downloader
}

