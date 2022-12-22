package memcache

import(
	"fmt"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"
	"github.com/DigiStratum/GoLib/Testing/fakes"

	chrono "github.com/DigiStratum/GoLib/Chrono"
)

const FAKE_MEMCACHED_PORT = 21212
const FAKE_MEMCACHED_HOST = "localhost"

func TestThat_NewDefaultMemcacheClient_ReturnsError_WhenNoHostsSpecified(t *testing.T) {
	// Setup
	ts := chrono.NewTimeSource()

	// Test
	sut, err := NewDefaultMemcacheClient(ts)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

func TestThat_NewDefaultMemcacheClient_ReturnsError_WhenTimeSourceNil(t *testing.T) {
	// Setup
	hosts := []string{ fmt.Sprintf("%s:%d", FAKE_MEMCACHED_HOST, FAKE_MEMCACHED_PORT) }

	// Test
	sut, err := NewDefaultMemcacheClient(nil, hosts...)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

func TestThat_NewDefaultMemcacheClient_ReturnsNonNil_NoError(t *testing.T) {
	// Setup
	fms, err := fakes.NewFakeMemcachedServer()
	ExpectNoError(err, t)
        defer fms.Close()
	//fms.Verbose()

	ts := chrono.NewTimeSource()
	hosts := []string{ fmt.Sprintf("%s:%d", FAKE_MEMCACHED_HOST, FAKE_MEMCACHED_PORT) }

	// Test
	sut, err := NewDefaultMemcacheClient(ts, hosts...)

	// Verify
	ExpectNonNil(sut, t)
	ExpectNoError(err, t)
}

