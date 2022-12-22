package memcache

import(
	"fmt"
	"net"
	"sync"
	"strings"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

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
	fakeServer, err := newFakeMemcachedServer(true)
	ExpectNoError(err, t)
        defer fakeServer.Close()
	ts := chrono.NewTimeSource()
	hosts := []string{ fmt.Sprintf("%s:%d", FAKE_MEMCACHED_HOST, FAKE_MEMCACHED_PORT) }

	// Test
	sut, err := NewDefaultMemcacheClient(ts, hosts...)

	// Verify
	ExpectNonNil(sut, t)
	ExpectNoError(err, t)
}

type fakeMemcachedServer struct {
	host		string
	listener	net.Listener
	waitGroup	sync.WaitGroup	// ref: https://gobyexample.com/waitgroups
	verbose		bool
}

func newFakeMemcachedServer(verbose bool) (*fakeMemcachedServer, error) {
	host := fmt.Sprintf("%s:%d", FAKE_MEMCACHED_HOST, FAKE_MEMCACHED_PORT)
	listener, err := net.Listen( "tcp", host )
	if nil != err { return nil, err }
	fms := fakeMemcachedServer{
		listener:	listener,
		verbose:	verbose,
		host:		host,
	}
	// Spin off a Go Routine for the FMS Listener
	go fms.Listen()
	return &fms, nil
}

// TODO: Any value in a generalized CloseableIfc to organize different Close()able types through?
// ref: https://www.reddit.com/r/golang/comments/a4nim7/nonblocking_accept_for_tcp_connections/
func (r *fakeMemcachedServer) Close() {
	r.waitGroup.Wait()
	r.listener.Close()
}

// ref: https://www.developer.com/languages/intro-socket-programming-go/
func (r *fakeMemcachedServer) Listen() {
	if r.verbose {
		fmt.Printf("\tFakeMemcachedServer: Ready to Accept Connections on '%s'\n", r.host)
	}
	for {
		connection, err := r.listener.Accept()
		if err == nil {
			if r.verbose {
				fmt.Println("\tFakeMemcachedServer: Got new client connection!\t")
			}
			go r.handleConnection(connection)
		}
	}
}

// Now that we have a memcache client connected, let's interact
func (r *fakeMemcachedServer) handleConnection(connection net.Conn) {
	// Make sure the Server waits for us to finish before Close()ing the Listener
	r.waitGroup.Add(1)
	defer r.waitGroup.Done()

	// Read a message
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		if r.verbose {
			fmt.Println("\tFakeMemcachedServer: Error reading:", err.Error())
		}
		return
	}

	// Got a message
	// TODO: Break off just the first word as the directive
	directive := strings.Trim(string(buffer[:mLen]), " \t\n\r")
	switch directive {
		case "version":
			if r.verbose {
				fmt.Printf("\tFakeMemcachedServer: Got 'version' directive!\n")
			}
			connection.Write([]byte("VERSION 0.0\n"))

		default:
			if r.verbose {
				fmt.Printf("\tFakeMemcachedServer: Unhandled directive: '%s'\n", directive)
			}
	}

	connection.Close()
}

