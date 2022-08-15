package objectstores3

/*

ObjectStore for AWS S3 service

Ref: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/s3-example-basic-bucket-operations.html

S3 Object storage mode adheres to our normal path model (see README.md), but with the additional
conditions that:

a) There exists an AWS / S3 service account
b) There exists a bucket within that S3 account
c) There exists a folder within that bucket within which the Object paths are organized

This enables us to maintain any number of collections of Objects within a given S3 bucket by
separating the collections into different folders.

Configuration:
	* awsregion	- AWS Region identifier e.g. "us-west-1"
	* s3bucket	- AWS S3 Bucket to retrieve content from
	* s3folder	- AWS S3 Folder to prepend to any path (no trailing slash)

*/

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	cfg "github.com/DigiStratum/GoLib/Config"
	obj "github.com/DigiStratum/GoLib/Object"
	"github.com/DigiStratum/GoLib/Cloud"
)

type ObjectStoreS3 struct {
	storeConfig	cfg.ConfigIfc
	awsS3		*s3.S3
	awsS3Downloader	*s3manager.Downloader
	readCache	*MutableObjectStore
	awsHelper	*cloud.AWSHelper
}

// -------------------------------------------------------------------------------------------------
// Factory Functions
// -------------------------------------------------------------------------------------------------

// Make a new one of these!
func NewObjectStoreS3() *ObjectStoreS3 {
	r := ObjectStoreS3{
		readCache: NewMutableObjectStore(),
	}
	return &r
}

// -------------------------------------------------------------------------------------------------
// Satisfies ConfigurableIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectStoreS3) Configure(config cfg.ConfigIfc) error {

	// Validate that the config has what we need for S3!
	requiredConfig := []string{ "awsregion", "s3bucket", "s3folder" }
	if ! (config.HasAll(&requiredConfig)) {
		return fmt.Errorf("Incomplete ObjectStoreS3 configuration provided")
	}
	r.storeConfig = config

	// Light up our AWS Helper with the region from our configuration data
	r.awsHelper = cloud.NewAWSHelper(config.Get("awsregion"))
	return nil
}

// -------------------------------------------------------------------------------------------------
// Satisfies ObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

// Ref: https://stackoverflow.com/questions/41645377/golang-s3-download-to-buffer-using-s3manager-downloader
func (r *ObjectStoreS3) GetObject(path string) (*obj.Object, error) {
	// Require configuration
	if nil == r.storeConfig { return nil, fmt.Errorf("Not Configured!") }

	// If it's not yet in the cache
	if ! r.readCache.HasObject(path) {
		// Read the Object from our S3 bucket into cache
		buff := &aws.WriteAtBuffer{}
		downloader := r.getS3Downloader()

		// The S3 key is the path prefixed with our configured folder for this store, if any
		s3Folder := r.storeConfig.Get("s3folder")
		key := path
		if len(s3Folder) > 0 { key = fmt.Sprintf("%s/%s", s3Folder, path) }

		// Now try to download the object from S3
		_, err := downloader.Download(
			buff,
			&s3.GetObjectInput{
				Bucket:	aws.String(os.storeConfig.Get("s3bucket")),
				Key:	aws.String(key),
			},
		)
		// Error = no Object!
		if nil != err { return nil, fmt.Errorf(
			"ObjectStoreS3.GetObject(%s) Error : '%s'",
			path,
			err.Error(),
		)}
		r.readCache.PutObject(path, NewObjectFromString(string(buff.Bytes())))
	}
	return r.readCache.GetObject(path), nil
}

func (r ObjectStoreS3) HasObject(path string) (bool, error) {
	// Require configuration
	if nil == r.storeConfig { return false, fmt.Errorf("Not Configured!") }

	// If it's already in the cache, then we know we have it!
	if r.readCache.HasObject(path) { return true, nil }

	// If there's S3 metadata with no error, then there's a Object!
	// ref: github.com/aws/aws-sdk-go/service/s3/examples_test.go ("HeadObject")
	awsS3 := r.getS3()
	_, err := awsS3.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(r.storeConfig.Get("s3bucket")),
			Key:	aws.String(r.storeConfig.Get("s3folder") + "/" + path),
		},
	)
	return nil == err, err
}

// -------------------------------------------------------------------------------------------------
// Satisfies MutableObjectStoreIfc Public Interface
// -------------------------------------------------------------------------------------------------

func (r *ObjectStoreS3) PutObject(path string, object *obj.Object) error {
	// Require configuration
	if nil == r.storeConfig { return fmt.Errorf("Not Configured!") }

	// TODO: Actually implement WRITE operation to S3 here
	return fmt.Errorf("Not Yet Implemented!")
}

// -------------------------------------------------------------------------------------------------
// ObjectStoreS3 Private Interface
// -------------------------------------------------------------------------------------------------

// Get our S3 connection
func (r ObjectStoreS3) getS3() *s3.S3 {
	if nil == r.awsS3 {
		sess := r.awsHelper.GetSession();
		if nil == sess { return nil }
		r.awsS3 = s3.New(sess)
	}
	return r.awsS3
}

// Get our S3 Downloader
func (r ObjectStoreS3) getS3Downloader() *s3manager.Downloader {
	if nil == r.awsS3Downloader {
		sess := r.awsHelper.GetSession();
		if nil == sess { return nil }
		r.awsS3Downloader = s3manager.NewDownloader(sess)
	}
	return r.awsS3Downloader
}
