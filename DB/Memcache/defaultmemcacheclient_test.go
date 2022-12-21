package memcache

import(
	"net"
	"testing"

	. "github.com/DigiStratum/GoLib/Testing"

	chrono "github.com/DigiStratum/GoLib/Chrono"
)

const FAKE_MEMCACHED_PORT = 21212
const FAKE_MEMCACHED_HOST = "localhost:" + FAKE_MEMCACHED_PORT

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
	fakeServer, err := fakeMemcached()
	ExpectNoError(err, t)
        defer fakeServer.Close()
	hosts := []string{ FAKE_MEMCACHED_HOST }

	// Test
	sut, err := NewDefaultMemcacheClient(nil, hosts...)

	// Verify
	ExpectNil(sut, t)
	ExpectError(err, t)
}

type fakeMemcachedServer struct {
	listener	net.Listener
	terminate	bool
}

func newFakeMemcachedServer() (*fakeMemcachedServer, error) {
	listener, err := net.Listen("tcp", FAKE_MEMCACHED_HOST)
	if nil != err { return nil, err }
	fms := fakeMemcachedServer{
		listener:	listener,
		terminate:	false,
	}
	// Spin off a Go Routine for the FMS Listener
	go fms.Listen()
	return &fms
}

// ref: https://www.reddit.com/r/golang/comments/a4nim7/nonblocking_accept_for_tcp_connections/
func (r *fakeMemcachedServer) Stop() {
	r.listener.Close()
}

// ref: https://www.developer.com/languages/intro-socket-programming-go/
func (r *fakeMemcachedServer) Listen() {
	for {
		connection, err := r.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			break
		}
		fmt.Println("client connected")
		go r.handleConnection(connection)
	}
}

// Now that we have a memcache client connected, let's interact
func (r *fakeMemcachedServer) handleConnection(connection net.Conn) {
	buffer := make([]byte, 1024)
	mLen, err := connection.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received: ", string(buffer[:mLen]))
	_, err = connection.Write([]byte("Thanks! Got your message:" + string(buffer[:mLen])))
	connection.Close()
}

