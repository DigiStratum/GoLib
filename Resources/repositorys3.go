package resources

/*

Resource Repository for AWS S3 service

*/

import (
	lib "github.com/digiStratum/GoLib"
)

type RepositoryS3 {
	repoConfig	*lib.Config
}

// Make a new one of these!
func NewRepositoryS3() *RepositoryS3 {
	r := repositoryS3{ }
	return &r
}

// Satisfies RespositoryIfc
func (r *RepositoryS3) Configure(repoConfig *lib.Config) error {
	r.repoConfig = repoConfig
	// TODO: Validate that the config has what we need for S3!
	return nil
}

// Satisfies RepositoryIfc
func (r *RepositoryS3) GetResource(path string) *Resource {
	// TODO: Actually implement READ operation to S3 here
	return nil
}

// Satisfies RepositoryIfc
func (r *RepositoryS3) HasResource(path string) bool {
	// TODO: Actually implement CHECK operation to S3 here
	return false
}

// Satisfies WritableRepositoryIfc
func (r *RepositoryS3) PutResource(resource *Resource, path string) error {
	// TODO: Actually implement WRITE operation to S3 here
	return nil
}

