package version

/*

Versioning interface

*/

type VersionIfc interface {
	GetVersion() string
	GetScheme() string
	Compare(version VersionIfc) (int, error)
}

