package memcache

import(
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

	chrono "github.com/DigiStratum/GoLib/Chrono"
)

func TestThat_NewDefaultMemcacheClient_ReturnsError_WhenNoHostsSpecified(t *testing.T) {
	// Setup
	ts := chrono.NewTimeSource()

	// Test
	sut, err := NewDefaultMemcacheClient(ts)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

