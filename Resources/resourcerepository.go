package resources

/*
A Resource Repository to serve as a the abstraction of accessing static resources from any number of
sources.

TODO: Add support for adding resource to a given repository
*/

const (
	REPO_TYPE_MAP	int = iota
	REPO_TYPE_DIR
	REPO_TYPE_DB
	REPO_TYPE_S3
	REPO_TYPE_URL
)

type ResourceRepository struct {
	type		int,
	resourceMap	*ResourceMap,	// REPO_TYPE_MAP
	baseDir		*string,	// REPO_TYPE_DIR
	db		*string,	// REPO_TYPE_DB
	bucket		*string,	// REPO_TYPE_S3
	baseUrl		*string,	// REPO_TYPE_URL
}

func NewCompiledRepository(resourceMap *ResourceMap) *ResourceRepository {
	rr := ResourceRepository{
		type:		REPO_TYPE_MAP,
		resourceMap:	resourceMap,
	}
	return &rr
}

func NewLocalRespository(baseDir string) *ResourceRepository {
	rr := ResourceRepository {
		type:		REPO_TYPE_DIR,
		baseDir:	baseDir,
	}
	return &rr
}

func NewDBRepository(db	string) *ResourceRepository {
	rr := ResourceRepository {
		type:		REPO_TYPE_DB,
		db:		db,
	}
	return &rr
}

func NewS3Repository(bucket string) *ResourceRepository {
	rr := ResourceRepository {
		type:		REPO_TYPE_S3,
		bucket:		bucket,
	}
	return &rr
}

func NewHttpRepository(baseUrl string) *ResourceRepository {
	rr := ResourceRepository {
		type:		REPO_TYPE_URL,
		baseUrl:	baseUrl,
	}
	return &rr
}

func (rr *ResourceRepository) GetResource(path string) *Resource {
	switch (rr.type) {
		case REPO_TYPE_MAP:
			return rr.resourceMap.GetResource(path)
		// TODO: Add support for the other repository types
	}
	return nil
}

