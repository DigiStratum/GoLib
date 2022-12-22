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
	fms, err := newFakeMemcachedServer()
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

type fakeMemcachedServer struct {
	host		string
	listening	bool
	listener	net.Listener
	waitGroup	sync.WaitGroup	// ref: https://gobyexample.com/waitgroups
	verbose		bool
}

func newFakeMemcachedServer() (*fakeMemcachedServer, error) {
	// Make a new one of these
	fms := fakeMemcachedServer{
		host:		fmt.Sprintf("%s:%d", FAKE_MEMCACHED_HOST, FAKE_MEMCACHED_PORT),
	}

	// Start up a socket listener
	err := fms.Listen()
	if nil != err { return nil, err }

	// Spin off a Go Routine for the FMS Listener
	go fms.Accept()
	return &fms, nil
}
func (r *fakeMemcachedServer) Listen() error {
	if nil == r { return fmt.Errorf("Nope: nil!") }
	if r.listening { return fmt.Errorf("Nope: already listening!") }
	listener, err := net.Listen( "tcp", r.host )
	if nil != err { return err }
	r.listener = listener
	r.VPrintf("Listening on '%s'", r.host)
	r.listening = true
	return nil
}

// ref: https://www.reddit.com/r/golang/comments/a4nim7/nonblocking_accept_for_tcp_connections/
func (r *fakeMemcachedServer) Close() {
	if (nil == r) || (! r.listening) { return }
	r.waitGroup.Wait()
	r.listener.Close()
}

// ref: https://www.developer.com/languages/intro-socket-programming-go/
func (r *fakeMemcachedServer) Accept() {
	if nil == r { return }
	r.VPrintf("Ready to Accept Connections")
	for {
		connection, err := r.listener.Accept()
		if err == nil {
			r.VPrintf("Got new client connection!")
			go r.handleConnection(connection)
		}
	}
}

// Now that we have a memcache client connected, let's interact
func (r *fakeMemcachedServer) handleConnection(connection net.Conn) {
	if nil == r { return }

	// No matter what happens, close this connection before we return to caller
	defer connection.Close()

	// Make sure the Server waits for us to finish before Close()ing the Listener
	r.waitGroup.Add(1)
	defer r.waitGroup.Done()

	// Read a message
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		r.VPrintf("Error reading: %s", err.Error())
		return
	}

	// Got a message
	// TODO: Break off just the first word as the directive
	directive := strings.Trim(string(buffer[:mLen]), " \t\n\r")
	switch directive {
		case "version":
			r.VPrintf("Got 'version' directive!")
			connection.Write([]byte("VERSION 0.0\n"))

		default:
			r.VPrintf("Unhandled directive: '%s'", directive)
	}
}

// Enable verbose output for this instance
func (r *fakeMemcachedServer) Verbose() {
	if (nil == r) || r.verbose { return }
	r.verbose = true
	r.VPrintf("Listening on '%s'", r.host)
}

// Put out a verbose message, ala Printf formatting
func (r *fakeMemcachedServer) VPrintf(formatter string, args ...interface{}) {
	if (nil == r) || ! r.verbose { return }
	if ! r.verbose { return }
	fmt.Printf("\tFakeMemcachedServer: " + formatter + "\n", args...)
}

