package version

/*

Major.Minor versioning scheme implementation

We support major.minor version scheme in accordance with commonly used schemes such as "1.0"

*/

import (
	"fmt"
	"strings"
	"strconv"
)

type MajorMinorIfc interface {
	// Embeded Interface(s)
	VersionIfc

	// Our own interface
	GetVersionMajor() int
	GetVersionMinor() int
}

type majmin struct {
	version		string
	major		int
	minor		int
}

// -------------------------------------------------------------------------------------------------
// Factory functions
// -------------------------------------------------------------------------------------------------

func NewMajorMinor(version string) *majmin {
	r := majmin{
		version:	version,
	}

	// Parse the version string into major.minor parts
	versionParts := strings.Split(version, ".")
	if (len(versionParts) == 0) || (len(versionParts) > 2) { return nil }
	var err error
	r.major, err = strconv.Atoi(versionParts[0]);
	if nil != err { return nil }
	if len(versionParts) > 1 {
		r.minor, err = strconv.Atoi(versionParts[1]);
		if nil != err { return nil }
	}

	return &r
}

// -------------------------------------------------------------------------------------------------
// VersionIfc
// -------------------------------------------------------------------------------------------------

func (r *majmin) GetVersion() string {
	if nil == r { return "" }
	return r.version
}

func (r *majmin) GetScheme() string {
	return "MAJMIN"
}

// Compare our version to another; return -1 if ours < other, 0 if ours ==, 1 if ours >;
// If versions are not comparable (i.e. mismatched Scheme), then return 0 + non-nil error
func (r *majmin) Compare(version VersionIfc) (int, error) {
	if version.GetScheme() != r.GetScheme() {
		return 0, fmt.Errorf("Version scheme mismatch, cannot compare!")
	}

	v, ok := version.(MajorMinorIfc)
	if ! ok {
		return 0, fmt.Errorf("Version does not implement MajorMinorIfc")
	}

	// If Major.Minor are == for both, then 0!
	if r.GetVersionMajor() == v.GetVersionMajor() {
		if r.GetVersionMinor() == v.GetVersionMinor() {
			return 0, nil
		} else if r.GetVersionMinor() < v.GetVersionMinor() { return -1, nil }
	} else if r.GetVersionMajor() < v.GetVersionMajor() { return -1, nil }

	return 1, nil
}

// -------------------------------------------------------------------------------------------------
// MajorMinorIfc
// -------------------------------------------------------------------------------------------------

func (r *majmin) GetVersionMajor() int {
	return r.major
}

func (r *majmin) GetVersionMinor() int {
	return r.minor
}

