package fakes

/*

A fake Memcached Server implementation

ref: https://github.com/memcached/memcached/blob/master/doc/protocol.txt
ref: https://docs.oracle.com/cd/E17952_01/mysql-5.6-en/ha-memcached-interfaces-protocol.html

*/

import(
	"fmt"
	"net"
	"sync"
	"time"
	"strings"
	"strconv"

	cfg "github.com/DigiStratum/GoLib/Config"
)

const FAKE_MEMCACHED_DEFAULT_PORT = 21212
const FAKE_MEMCACHED_DEFAULT_HOST = "localhost"

type FakeMemcachedServerIfc interface {
	Listen() error		// Start listening on host:port (if not already)
	Close()			// Stop listening/accepting new connections
	Verbose()		// Enable verbose output for this instance
}

type fakeCacheItem struct {
	Value		[]byte
	Flags		uint32
	Expires		time.Time
}

type fakeMemcachedServer struct {
	host		string
	listening	bool
	listener	net.Listener
	waitGroup	sync.WaitGroup	// ref: https://gobyexample.com/waitgroups
	verbose		bool
	cache		map[string]fakeCacheItem
}

// ------------------------------------------------------------------------------------------------
// Factory Functions
// ------------------------------------------------------------------------------------------------

// Instantiate FakeMemcachedServer with optional config items: 'port' and 'host'
func NewFakeMemcachedServer(config ...cfg.ConfigItemIfc) (*fakeMemcachedServer, error) {
	// Configure
	var err error
	port := FAKE_MEMCACHED_DEFAULT_PORT
	host := FAKE_MEMCACHED_DEFAULT_HOST
	for _, ci := range config {
		switch ci.GetName() {
			case "port":
				vi := ci.GetValue()
				if vs, ok := vi.(string); ok {
					port, err = strconv.Atoi(vs)
				}
			case "host":
				vi := ci.GetValue()
				if vs, ok := vi.(string); ok {
					host = vs
				}
		}
	}

	// Make a new one of these
	fms := fakeMemcachedServer{
		host: fmt.Sprintf("%s:%d", host, port ),
		cache: make(map[string]fakeCacheItem),
	}

	// Start up a socket listener
	if err = fms.Listen(); nil != err { return nil, err }

	return &fms, nil
}

// Start listening on host:port (if not already)
func (r *fakeMemcachedServer) Listen() error {
	if nil == r { return fmt.Errorf("Nope: nil!") }
	if r.listening { return fmt.Errorf("Nope: already listening!") }

	// Start listening...
	listener, err := net.Listen( "tcp", r.host )
	if nil != err { return err }
	r.listener = listener
	r.vprintf("Listening on '%s'", r.host)
	r.listening = true

	// Spin off a Go Routine for the FMS Listener
	go r.accept()

	return nil
}

// Stop listening/accepting new connections
func (r *fakeMemcachedServer) Close() {
	if (nil == r) || (! r.listening) { return }

	// Stop listening
	// ref: https://www.reddit.com/r/golang/comments/a4nim7/nonblocking_accept_for_tcp_connections/
	r.listener.Close()

	// Wait for any outstanding connections to be handled
	r.waitGroup.Wait()
}

// Enable verbose output for this instance
func (r *fakeMemcachedServer) Verbose() {
	if (nil == r) || r.verbose { return }
	r.verbose = true
	r.vprintf("Listening on '%s'", r.host)
}


// ------------------------------------------------------------------------------------------------
// fakeMemcachedServer implementaiton
// ------------------------------------------------------------------------------------------------

// Acceept new connections (until listener is Close()'d)
// ref: https://www.developer.com/languages/intro-socket-programming-go/
func (r *fakeMemcachedServer) accept() {
	if nil == r { return }
	r.vprintf("Ready to Accept Connections")
	for {
		connection, err := r.listener.Accept()
		if err == nil {
			r.vprintf("Got new client connection!")
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
		r.vprintf("Error reading: %s", err.Error())
		return
	}

	// Got a message
	// TODO: Break off just the first word as the directive
	commandLine := strings.Trim(string(buffer[:mLen]), " \t\n\r")
	commandWords := strings.Fields(commandLine)
	if len(commandWords) == 0 {
		r.vprintf("Empty command, nothing to do!")
		return
	}
	var response string
	switch commandWords[0] {
		// TODO: Add other commands
		case "version":
			r.vprintf("Got 'version' directive!")
			response = "VERSION 0.0\n"

		case "set":
			// Most common command. Store this data, possibly overwriting any existing data. New items are at the top of the LRU.

		case "add":
			// Store this data, only if it does not already exist. New items are at the top of the LRU. If an item already exists and an add fails, it promotes the item to the front of the LRU anyway.
		case "replace":
			// Store this data, but only if the data already exists. Almost never used, and exists for protocol completeness (set, add, replace, etc)
		case "append":
			// Add this data after the last byte in an existing item. This does not allow you to extend past the item limit. Useful for managing lists.
		case "prepend":
			// Same as append, but adding new data before existing data.
		case "cas":
			// Check And Set (or Compare And Swap). An operation that stores data, but only if no one else has updated the data since you read it last. Useful for resolving race conditions on updating cache data.
		case "get":
			// Command for retrieving data. Takes one or more keys and returns all found items.
			// get xyzkey\r\n
			// VALUE xyzkey 0 6\r\n
			// abcdef\r\n
			if len(commandWords) >= 2 {
				key := commandWords[1]
				if ci, ok := r.cache[key]; ok {
					response = r.getValueResponse(key, &ci)
				} else {
					// TODO: What response is expected for get of invalid key (doesn't exist)?
				}
			} else {
				// TODO: What response is expected for a get with no specified key?
			}
		case "gets":
			// An alternative get command for using with CAS. Returns a CAS identifier (a unique 64bit number) with the item. Return this value with the cas command. If the item's CAS value has changed since you gets'ed it, it will not be stored.
		case "delete":
			// Removes an item from the cache, if it exists.
		case "incr":
			// Increment and Decrement. If an item stored is the string representation of an unsigned 64bit integer, you may run incr or decr commands to modify that number. You may only incr by positive values, or decr by positive values. They do not accept negative values.

			// If a value does not already exist, incr/decr will fail.
		case "decr":
		case "stats":
			// basic stats command.
			// sub-commands are "items", "slabs", and "sizes"
			// items: Returns some information, broken down by slab, about items stored in memcached.
			// slabs: Returns more information, broken down by slab, about items stored in memcached. More centered to performance of a slab rather than counts of particular items.
			// sizes: A special command that shows you how items would be distributed if slabs were broken into 32byte buckets instead of your current number of slabs. Useful for determining how efficient your slab sizing is.

			// WARNING this is a development command. As of 1.4 it is still the only command which will lock your memcached instance for some time. If you have many millions of stored items, it can become unresponsive for several minutes. Run this at your own risk. It is roadmapped to either make this feature optional or at least speed it up.
		case "flush_all":
			// Invalidate all existing cache items. Optionally takes a parameter, which means to invalidate all items after N seconds have passed.
			// This command does not pause the server, as it returns immediately. It does not free up or flush memory at all, it just causes all items to expire.
		default:
			r.vprintf("Unhandled directive: '%s'", commandWords[0])
			return
	}
	connection.Write([]byte(response))
}

func (r *fakeMemcachedServer) getValueResponse(key string, ci *fakeCacheItem) string {
	return fmt.Sprintf(
		"VALUE %s %d %d\r\n%s\r\n",
		key,
		ci.Flags,
		len(ci.Value),
		string(ci.Value),
	)
}


// Put out a verbose message, ala Printf formatting
func (r *fakeMemcachedServer) vprintf(formatter string, args ...interface{}) {
	if (nil == r) || ! r.verbose { return }
	if ! r.verbose { return }
	fmt.Printf("\tFakeMemcachedServer: " + formatter + "\n", args...)
}

