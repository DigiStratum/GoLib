package version

/*

Semantic Versioning (AKA SEMVER) representation with compararison methods.

We strictly support major.minor.patch version scheme in accordance with accepted norms around
SEMVER. In theory, we could also support any number of sub-version nodes, however this complicates
the semantics of the versioning scheme unnecessarily for purposes of SEMVER; a different versioning
scheme should be a different interface altogether.

Our comparison methods offer a couple different coverages. The main Compare() compares the entire
version string for <=>. However we also have support for CompareMajor() which just checks <=> for
two versions being within the same major version (i.e. minor & patch do not matter), or
CompareMajorMinor() to check <=> for two within the same major.minor version. This is important to
be clear, in accordance with SEMVER, that a versioned interface is compatible (or not) for example,
with version N.* or version N+, or that an update is needed for version N-, etc. The same applies to
major.minor comparisons where we want to check N.M.*, N.M+, or N.M-. Use the Compare() function for
exact versio (major.minor.patch) comparisons.

TODO:
 * Consider options to standardize the interface of comparison function(s) for different versioning
   schemes. If we can come up with some sort of scheme where, say, the comparison version is
   specified as some sort of matcher string (e.g. "1.2.*") instead of some version scheme specific
   calling notation, then we can standardize a VersionIfc interface such that SemVer is but one of
   many possible versioning scheme implementations that we can implement. Note that the interface
   would also need to indicate which scheme is in use so that comparability is determinable.

*/

import (
	"fmt"
	"strings"
)

type SemVerIfc interface {
	ToString() string

	// TODO
	// Compare our version to another; return -1 if ours < other, 0 if ours ==, 1 if ours >
	// Compare(semver SemVerIfc) int
	// Compare our version MAJOR to another; return -1 if ours < other, 0 if ours ==, 1 if ours >
	// CompareMajor(semver SemVerIfc) int
	// Compare our version MAJOR.MINOR to another; return -1 if ours < other, 0 if ours ==, 1 if ours >
	// CompareMajorMinor(semver SemVerIfc) int
}

type semver struct {
	major		*int
	minor		*int
	patch		*int
}

func NewSemVer(version string) *semver {
	// TODO: Parse the version string into major.minor.patch parts
	r := semver{}
	return &r
}

func (r *semver) ToString() string {
	if nil == r { return "" }
	var sb strings.Builder
	if nil != r.major {
		sb.WriteString(fmt.Sprintf("%d", *r.major))
		if nil != r.minor {
			sb.WriteString(fmt.Sprintf(".%d", *r.minor))
			if nil != r.patch{
				sb.WriteString(fmt.Sprintf(".%d", *r.patch))
			}
		}
	}
	return sb.String()
}

