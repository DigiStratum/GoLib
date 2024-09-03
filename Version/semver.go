package version

/*

Semantic Versioning (AKA SEMVER) representation

We support major.minor.patch version scheme in accordance with accepted norms around SEMVER.

ref: https://semver.npmjs.com/

TODO:
 * Support NPM style version/range specifiers: https://semver.npmjs.com/

*/

import (
	"fmt"
	"strings"
	"strconv"
)

type SemVerIfc interface {
	// Embeded Interface(s)
	VersionIfc

	// Our own interface
	GetVersionMajor() int
	GetVersionMinor() int
	GetVersionPatch() int
}

type semver struct {
	version		string
	major		int
	minor		int
	patch		int
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewSemVer(version string) *semver {
	r := semver{
		version:	version,
	}

	// Parse the version string into major.minor.patch parts
	versionParts := strings.Split(version, ".")
	if (len(versionParts) == 0) || (len(versionParts) > 3) { return nil }
	var err error
	r.major, err = strconv.Atoi(versionParts[0]);
	if nil != err { return nil }
	if len(versionParts) > 1 {
		r.minor, err = strconv.Atoi(versionParts[1]);
		if nil != err { return nil }
		if len(versionParts) > 2 {
			r.patch, err = strconv.Atoi(versionParts[2]);
			if nil != err { return nil }
		}
	}

	return &r
}

// -------------------------------------------------------------------------------------------------
// VersionIfc
// -------------------------------------------------------------------------------------------------

func (r *semver) GetVersion() string {
	if nil == r { return "" }
	return r.version
}

func (r *semver) GetScheme() string {
	return "SEMVER"
}

// Compare our version to another; return -1 if ours < other, 0 if ours ==, 1 if ours >;
// If versions are not comparable (i.e. mismatched Scheme), then return 0 + non-nil error
func (r *semver) Compare(version VersionIfc) (int, error) {
	if version.GetScheme() != r.GetScheme() {
		return 0, fmt.Errorf("Version scheme mismatch, cannot compare!")
	}

	v, ok := version.(SemVerIfc)
	if ! ok {
		return 0, fmt.Errorf("Version does not implement SemVerIfc")
	}

	// If Major.Minor.Patch are == for both, then 0!
	if r.GetVersionMajor() == v.GetVersionMajor() {
		if r.GetVersionMinor() == v.GetVersionMinor() {
			if r.GetVersionPatch() == v.GetVersionPatch() {
				return 0, nil
			} else if r.GetVersionPatch() < v.GetVersionPatch() { return -1, nil }
		} else if r.GetVersionMinor() < v.GetVersionMinor() { return -1, nil }
	} else if r.GetVersionMajor() < v.GetVersionMajor() { return -1, nil }

	return 1, nil
}

// -------------------------------------------------------------------------------------------------
// SemVerIfc
// -------------------------------------------------------------------------------------------------

func (r *semver) GetVersionMajor() int {
	return r.major
}

func (r *semver) GetVersionMinor() int {
	return r.minor
}

func (r *semver) GetVersionPatch() int {
	return r.patch
}

