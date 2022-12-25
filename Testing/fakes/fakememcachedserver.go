package fakes

/*

A fake Memcached Server implementation

-----

ref: https://github.com/memcached/memcached/blob/master/doc/protocol.txt

-----

ref: https://docs.oracle.com/cd/E17952_01/mysql-5.6-en/ha-memcached-interfaces-protocol.html

## Storage commands to the server take the form:

>> command key [flags] [exptime] length [noreply]

Or when using compare and swap (cas):

>> cas key [flags] [exptime] length [casunique] [noreply]

Where:
 * command: The command name.
 * set: Store value against key
 * add: Store this value against key if the key does not already exist
 * replace: Store this value against key if the key already exists
 * append: Append the supplied value to the end of the value for the specified key. The flags and
   exptime arguments should not be used.
 * prepend: Append value currently in the cache to the end of the supplied value for the specified
   key. The flags and exptime arguments should not be used.
 * cas: Set the specified key to the supplied value, only if the supplied casunique matches. This is
   effectively the equivalent of change the information if nobody has updated it since I last fetched it.
 * key: The key. All data is stored using a the specific key. The key cannot contain control characters
   or whitespace, and can be up to 250 characters in size.
 * flags: The flags for the operation (as an integer). Flags in memcached are transparent. The memcached
   server ignores the contents of the flags. They can be used by the client to indicate any type of
   information. In memcached 1.2.0 and lower the value is a 16-bit integer value. In memcached 1.2.1 and
   higher the value is a 32-bit integer.
 * exptime: The expiry time, or zero for no expiry.
 * length: The length of the supplied value block in bytes, excluding the terminating \r\n characters.
 * casunique: A unique 64-bit value of an existing entry. This is used to compare against the existing
   value. Use the value returned by the gets command when issuing cas updates.
 * noreply: Tells the server not to reply to the command.

The return value from the server is one line, specifying the status or error information.

...



## Retrieval commands take the form:

>> get key1 [key2 .... keyn]
>> gets key1 [key2 ... keyn]

You can supply multiple keys to the commands, with each requested key separated by whitespace.

The server responds with an information line of the form:

>> VALUE key flags bytes [casunique]

Where:

 * key: The key name.
 * flags: The value of the flag integer supplied to the memcached server when the value was stored.
 * bytes: The size (excluding the terminating \r\n character sequence) of the stored value.
 * casunique: The unique 64-bit integer that identifies the item.

The information line is immediately followed by the value data block.

## Deletion commands take the form:

>> delete key [time] [noreply]

Where:
 * key: The key name.
 * time: The time in seconds (or a specific Unix time) for which the client wishes the server to
   refuse add or replace commands on this key. All add, replace, get, and gets commands fail during
   this period. set operations succeed. After this period, the key is deleted permanently and all
   commands are accepted.

If not supplied, the value is assumed to be zero (delete immediately).

noreply: Tells the server not to reply to the command.

Responses to the command are either DELETED to indicate that the key was successfully removed, or NOT_FOUND to indicate that the specified key could not be found.

## The increment and decrement commands change the value of a key within the server without performing a separate get/set sequence. The operations assume that the currently stored value is a 64-bit integer. If the stored value is not a 64-bit integer, then the value is assumed to be zero before the increment or decrement operation is applied.

## Increment and decrement commands take the form:

>> incr key value [noreply]
>> decr key value [noreply]

Where:
 * key: The key name.
 * value: An integer to be used as the increment or decrement value.
 * noreply: Tells the server not to reply to the command.

The response is:

>> NOT_FOUND: The specified key could not be located.
>> value: The new value associated with the specified key.

Values are assumed to be unsigned. For decr operations, the value is never decremented below 0. For incr operations, the value wraps around the 64-bit maximum.

## The stats command provides detailed statistical information about the current status of the memcached instance and the data it is storing.

Statistics commands take the form:

>> STAT [name] [value]

Where:
 * name: The optional name of the statistics to return. If not specified, the general statistics are returned.
 * value: A specific value to be used when performing certain statistics operations.

The return value is a list of statistics data, formatted as follows:

>> STAT name value

The statistics are terminated with a single line, END.

-----

Command		Command Formats
set		set key flags exptime length, set key flags exptime length noreply
add		add key flags exptime length, add key flags exptime length noreply
replace		replace key flags exptime length, replace key flags exptime length noreply
append		append key length, append key length noreply
prepend		prepend key length, prepend key length noreply
cas		cas key flags exptime length casunique, cas key flags exptime length casunique noreply
get		get key1 [key2 ... keyn]
gets
delete		delete key, delete key noreply, delete key expiry, delete key expiry noreply
incr		incr key, incr key noreply, incr key value, incr key value noreply
decr		decr key, decr key noreply, decr key value, decr key value noreply
stat		stat, stat name, stat name value

-----

## memcached Protocol Responses

String				Description
STORED				Value has successfully been stored.
NOT_STORED			The value was not stored, but not because of an error. For commands
				where you are adding a or updating a value if it exists (such as add
				and replace), or where the item has already been set to be deleted.
EXISTS				When using a cas command, the item you are trying to store already
				exists and has been modified since you last checked it.
NOT_FOUND			The item you are trying to store, update or delete does not exist or
				has already been deleted.
ERROR				You submitted a nonexistent command name.
CLIENT_ERROR errorstring	There was an error in the input line, the detail is contained in
				errorstring.
SERVER_ERROR errorstring	There was an error in the server that prevents it from returning the
				information. In extreme conditions, the server may disconnect the
				client after this error occurs.
VALUE keys flags length		The requested key has been found, and the stored key, flags and data
				block are returned, of the specified length.
DELETED				The requested key was deleted from the server.
STAT name value			A line of statistics data.
END				The end of the statistics data.

-----

TODO:
 * Simulate memcached version <= 1.2.0 for 16 bit flags on cache items vs 32 bits vs 1.2.1+

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
	Flags		uint32	// Note:  In memcached 1.2.0 and lower the value is a 16-bit integer value. In memcached 1.2.1 and higher the value is a 32-bit integer.
	Expires		int64	// 0 for non-expiring items
	Accessed	int64	// When was this item last accessed (factors into LRU algorithm)
	CASUnique	string	// Simple hash to compare and swap (md5?)
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
			// TODO: Make the version variable/configurable, simulate different behaviors expected for different versions
			response = "VERSION 0.0\n"

		case "set":
			// Most common command. Store this data, possibly overwriting any existing data. New items are at the top of the LRU.

		case "add":
			// Store this data, only if it does not already exist. New items are at the top of the LRU. If an item already exists and an add fails, it promotes the item to the front of the LRU anyway.
			if len(commandWords) >= 2 {
				key := commandWords[1]
				ci, exists := r.cache[key]
				if exists {
					// TODO: Touch the last access timestamp
					ci.Accessed = time.Now().Unix()
				} else {
					// TODO: Set normally
				}
			} else {
				// TODO: What response is expected for a get with no specified key?
			}
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
				if ci, exists := r.cache[key]; exists {
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

