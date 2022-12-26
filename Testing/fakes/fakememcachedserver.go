package fakes

/*

A fake Memcached Server implementation

ref: https://github.com/memcached/memcached/blob/master/doc/protocol.txt
ref: https://docs.oracle.com/cd/E17952_01/mysql-5.6-en/ha-memcached-interfaces-protocol.html
ref: https://github.com/memcached/memcached/wiki/MetaCommands

TODO:
 * Simulate memcached version <= 1.2.0 for 16 bit flags on cache items vs 32 bits vs 1.2.1+
 * Simulate max item size and rejection of over-sized values

*/

import(
	"fmt"
	"net"
	"sync"
	"strconv"
	"regexp"

	cfg "github.com/DigiStratum/GoLib/Config"
	chrono "github.com/DigiStratum/GoLib/Chrono"
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
	timeSource	chrono.TimeSourceIfc
}

// ------------------------------------------------------------------------------------------------
// Factory Functions
// ------------------------------------------------------------------------------------------------

// Instantiate FakeMemcachedServer with optional config items: 'port' and 'host'
func NewFakeMemcachedServer(config ...cfg.ConfigItemIfc) (*fakeMemcachedServer, error) {
	var err error

	// Configure with sane/working defaults
	port := FAKE_MEMCACHED_DEFAULT_PORT
	host := FAKE_MEMCACHED_DEFAULT_HOST
	var timeSource chrono.TimeSourceIfc = chrono.NewTimeSource()

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
			case "TimeSource":
				vi := ci.GetValue()
				if vts, ok := vi.(chrono.TimeSourceIfc); ok {
					timeSource = vts
				}
		}
	}

	// Make a new one of these
	fms := fakeMemcachedServer{
		host: fmt.Sprintf("%s:%d", host, port ),
		cache: make(map[string]fakeCacheItem),
		timeSource: timeSource,
	}

	// Start up a socket listener
	if err = fms.Listen(); nil != err { return nil, err }

	return &fms, nil
}

// ------------------------------------------------------------------------------------------------
// fakeMemcachedServer implementation
// ------------------------------------------------------------------------------------------------

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

// -----------------------------------------------
// SOCKET HANDLERS
// -----------------------------------------------

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
	message := make([]byte, 1024)
	if _, err := connection.Read(message); nil != err {
		r.vprintf("Error reading: %s", err.Error())
		return
	}

	// Send response to the message based on the handler's result
	connection.Write([]byte(r.handleCommand(newCommandTokenizer(&message))))
}

// -----------------------------------------------
// COMMAND HANDLERS
// -----------------------------------------------

func (r *fakeMemcachedServer) handleCommand(cmd *commandTokenizer) string {
	command := cmd.GetTokenString()
	if nil == command {
		r.vprintf("Empty command, nothing to do!")
		return "" // TODO: What is the expected response for an empty/blank command?
	}
	switch *command {
		case "version": return r.handleVersionCommand(cmd)
		case "set": return r.handleSetCommand(cmd)
		case "add": return r.handleAddCommand(cmd)
		case "replace": return r.handleReplaceCommand(cmd)
		case "append": return r.handleAppendCommand(cmd)
		case "prepend": return r.handlePrependCommand(cmd)
		case "cas":
			// Check And Set (or Compare And Swap). An operation that stores data, but only if no one else has updated the data since you read it last. Useful for resolving race conditions on updating cache data.
		case "get":
			// Command for retrieving data. Takes one or more keys and returns all found items.
			// get xyzkey\r\n
			// VALUE xyzkey 0 6\r\n
			// abcdef\r\n
			/*
			if len(commandWords) >= 2 {
				key := commandWords[1]
				if ci := r.readCacheItem(key); nil != ci {
					response = r.getValueResponse(key, ci)
				} else {
					// TODO: What response is expected for get of invalid key (doesn't exist)?
				}
			} else {
				// TODO: What response is expected for a get with no specified key?
			}
			*/
		case "gets":
			// An alternative get command for using with CAS. Returns a CAS identifier (a unique 64bit number) with the item. Return this value with the cas command. If the item's CAS value has changed since you gets'ed it, it will not be stored.
		case "delete":
			// Removes an item from the cache, if it exists.
		case "incr": return r.handleIncrementCommand(cmd, 1)
		case "decr": return r.handleIncrementCommand(cmd, -1)
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
			/*
			r.cache = make(map[string]fakeCacheItem)
			if len(commandWords) >= 2 {
				next := commandWords[1]
				if "noreply" != next {
					response = r.getOkResponse()
				}
			}
			*/
		default:
	}
	return r.getErrorResponse(fmt.Sprintf("Unhandled command: '%s'", *command))
}

// Get the fake memcached server version
func (r *fakeMemcachedServer) handleVersionCommand(cmd *commandTokenizer) string {
	return r.getVersionResponse()
}

// Most common command. Store this data, possibly overwriting any existing data. New items are at the top of the LRU.
func (r *fakeMemcachedServer) handleSetCommand(cmd *commandTokenizer) string {
	key, flags, expires, noReply, value, err := r.parseStorageCmd(cmd)
	if nil != err { return r. getErrorResponse(err.Error()) }
	ci := fakeCacheItem{
		Value: value,
		Expires: int64(expires),
		Flags: uint32(flags),
	}
	r.writeCacheItem(key, &ci)
	if noReply { return "" }
	return r.getStoredResponse()
}

// Store this data, only if it does not already exist. New items are at the top of the LRU. If an item already exists and an add fails, it promotes the item to the front of the LRU anyway.
func (r *fakeMemcachedServer) handleAddCommand(cmd *commandTokenizer) string {
	// Parse: add <key> [...] just to peek at the key and put it back
	cmd.SetRewindPoint()
	key := cmd.GetTokenString()
	cmd.Rewind()
	if nil == key {
		return r.getErrorResponse("requires `command <key> <flags> <exptime> <bytes> [noreply]\\r\\n<data>\\r\\n`")
	}
	// If this item doesn't exist in the cache, then set it normally
	if ! r.existsCacheItem(*key) { return r.handleSetCommand(cmd) }
	// Otherwise just update the last accessed timestamp
	r.touchCacheItem(*key)
	return r.getOkResponse()
}

// Store this data, but only if the data already exists. Almost never used, and exists for protocol completeness (set, add, replace, etc)
func (r *fakeMemcachedServer) handleReplaceCommand(cmd *commandTokenizer) string {
	// Parse: replace <key> <flags> <exptime> <bytes> [noreply]\r\n<data>\r\n
	cmd.SetRewindPoint()
	key := cmd.GetTokenString()
	cmd.Rewind()
	if nil == key {
		return r.getErrorResponse("requires `command <key> <flags> <exptime> <bytes> [noreply]\\r\\n<data>\\r\\n`")
	}
	if r.existsCacheItem(*key) { return r.getNotStoredResponse() }
	return r.handleSetCommand(cmd)
}

// Add this data after the last byte in an existing item. This does not allow you to extend past the item limit. Useful for managing lists.
func (r *fakeMemcachedServer) handleAppendCommand(cmd *commandTokenizer) string {
	key, _, _, noReply, value, err := r.parseStorageCmd(cmd)
	if nil != err { return r. getErrorResponse(err.Error()) }
	if ci := r.readCacheItem(key); nil != ci {
		// Append supplied value to the existing ci Value
		ci.Value = append(ci.Value, value...)
		// TODO: Validate max length of Value
		// TODO: What is expected of provided expires/flags? Are we supposed to replace these values or ignore?
		r.writeCacheItem(key, ci)
		if noReply { return "" }
		return r.getStoredResponse()
	}
	if noReply { return "" }
	return r.getNotFoundResponse()
}

// Same as append, but adding new data before existing data.
func (r *fakeMemcachedServer) handlePrependCommand(cmd *commandTokenizer) string {
	key, _, _, noReply, value, err := r.parseStorageCmd(cmd)
	if nil != err { return r. getErrorResponse(err.Error()) }
	if ci := r.readCacheItem(key); nil != ci {
		// Append the existing ci value to the supplied value
		ci.Value = append(value, ci.Value...)
		// TODO: Validate max length of Value
		// TODO: What is expected of provided expires/flags? Are we supposed to replace these values or ignore?
		r.writeCacheItem(key, ci)
		if noReply { return "" }
		return r.getStoredResponse()
	}
	if noReply { return "" }
	return r.getNotFoundResponse()
}

// Increment and Decrement. If an item stored is the string representation of an unsigned 64bit integer, you may run incr or decr commands to modify that number. You may only incr by positive values, or decr by positive values. They do not accept negative values.
func (r *fakeMemcachedServer) handleIncrementCommand(cmd *commandTokenizer, direction int) string {
	key := cmd.GetTokenString()
	delta := cmd.GetTokenInt()
	if (nil == key) || (nil == delta) { return r.getErrorResponse("incr/decr command: must supply <key> <value>") }
	noReply := cmd.IsNoReply()
	if ci := r.readCacheItem(*key); nil != ci {
		intval, err := strconv.Atoi(string(ci.Value))
		if nil != err { return r.getErrorResponse("inc/dec command: existing entry does not appear to be numeric") }
		newintval := intval+(*delta*direction)
		if newintval < 0 { newintval = 0 }
		newValue := fmt.Sprintf("%d", newintval)
		ci.Value = []byte(newValue)
		r.writeCacheItem(*key, ci)
		if noReply { return "" }
		return r.getValueOnlyResponse(newValue)
	}
	if noReply { return "" }
	return r.getNotFoundResponse()
}

/*
func (r *fakeMemcachedServer) handleCommand(cmd *commandTokenizer) string {
}
*/

// -----------------------------------------------
// HELPERS
// -----------------------------------------------

// Parse: command <key> <flags> <exptime> <bytes> [noreply]\r\n<data>\r\n
func (r *fakeMemcachedServer) parseStorageCmd(cmd *commandTokenizer) (string, int, int, bool, []byte, error) {
	key := cmd.GetTokenString()
	flags := cmd.GetTokenInt()
	expires := cmd.GetTokenInt()
	bytelen := cmd.GetTokenInt()
	if (nil == key) || (nil == flags) || (nil == expires) || (nil == bytelen) {
		return "", 0, 0, false, []byte{}, fmt.Errorf("requires `command <key> <flags> <exptime> <bytes> [noreply]\\r\\n<data>\\r\\n`")
	}
	noReply := cmd.IsNoReply()
	value := cmd.GetTokenBytes(*bytelen)
	if nil == value {
		return "", 0, 0, false, []byte{}, fmt.Errorf("storage command: failed to read %d bytes for value", *bytelen)
	}
	return *key, *flags, *expires, noReply, *value, nil
}

func (r *fakeMemcachedServer) readCacheItem(key string) *fakeCacheItem {
	ci, _ := r.cache[key]
	// Make sure it either never expires or not expired (Accessed + Expires < now() )
	if (0 == ci.Expires) || (ci.Accessed + ci.Expires < r.timeSource.NowUnixTimeStamp()) {
		return &ci
	}
	return nil
}

func (r *fakeMemcachedServer) existsCacheItem(key string) bool {
	return nil != r.readCacheItem(key)
}

func (r *fakeMemcachedServer) writeCacheItem(key string, ci *fakeCacheItem) {
	ci.Accessed = r.timeSource.NowUnixTimeStamp()
	r.cache[key] = *ci
}

func (r *fakeMemcachedServer) touchCacheItem(key string) {
	if ci :=r.readCacheItem(key); nil != ci {
		ci.Accessed = r.timeSource.NowUnixTimeStamp()
		r.writeCacheItem(key, ci)
	}
}

// -----------------------------------------------
// RESPONSES
// -----------------------------------------------

func (r *fakeMemcachedServer) getValueOnlyResponse(value string) string {
	return fmt.Sprintf(
		"VALUE %s\r\n",
		value,
	)
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

func (r *fakeMemcachedServer) getOkResponse() string {
	return fmt.Sprintf("OK\r\n")
}

func (r *fakeMemcachedServer) getErrorResponse(msg string) string {
	return fmt.Sprintf("ERROR %s\r\n", msg)
}

func (r *fakeMemcachedServer) getStoredResponse() string {
	return fmt.Sprintf("STORED\r\n")
}

func (r *fakeMemcachedServer) getNotStoredResponse() string {
	return fmt.Sprintf("NOT_STORED\r\n")
}

func (r *fakeMemcachedServer) getNotFoundResponse() string {
	return fmt.Sprintf("NOT_FOUND\r\n")
}

func (r *fakeMemcachedServer) getVersionResponse() string {
	// TODO: Make version configurable; we want to also be able to alter certain behaviors based on version differences
	vmajor := 0
	vminor := 0
	vpatch := 0
	return fmt.Sprintf("VERSION %d.%d.%d\r\n", vmajor, vminor, vpatch)
}

// Put out a verbose message, ala Printf formatting
func (r *fakeMemcachedServer) vprintf(formatter string, args ...interface{}) {
	if (nil == r) || ! r.verbose { return }
	if ! r.verbose { return }
	fmt.Printf("\tFakeMemcachedServer: " + formatter + "\n", args...)
}

// -------------------------------------------------------------------------------------------------
// Lexical Command Tokenizer Implementation
// -------------------------------------------------------------------------------------------------

type commandTokenizer struct {
	cursor		int		// Current cursor location
	rewind		int		// Rewind location
	command		[]byte		// Command "text" that we are lexing
}

func newCommandTokenizer(command *[]byte) *commandTokenizer {
	return &commandTokenizer{
		cursor:		0,
		rewind:		0,
		command:	*command,
	}
}

func (r *commandTokenizer) SetRewindPoint() {
	r.rewind = r.cursor
}

func (r *commandTokenizer) Rewind() {
	r.cursor = r.rewind
}

// Absorb white space/separators, then get any non-white space as the token up to next white-space/separator/EOL
func (r *commandTokenizer) GetTokenString() *string {
	inToken := false
	start := r.cursor
	for ; r.cursor < len(r.command); r.cursor++ {
		ch := r.command[r.cursor]
		// Any char that's white space or separator is a token boundary:
		if (ch == ' ') || (ch == '\r') || (ch == '\n') || (ch == '\t') {
			if inToken { break }
			r.cursor++ // Absorb leading white-space/separators until we get to the beginning of the token
		} else { inToken = true }
	}
	// Seems we found no token, only white-space
	if ! inToken { return nil }
	token := r.command[start:r.cursor]
	tokenstr := string(token)
	return &tokenstr
}

// Return optional patemeter (such as '[noreply]') matching pattern in the next Token, or don't consume it
func (r *commandTokenizer) GetOptionalTokenString(pattern string) *string {
	// Use this instead of Rewind so that outside caller can still Rewind() back past OptionalTokenStrings...
	initialCursor := r.cursor
	token := r.GetTokenString()
	if nil == token { return nil }
	if matched, _ := regexp.MatchString(pattern, *token); ! matched {
		// No Match! Rewind the cursor
		r.cursor = initialCursor
		return nil
	}
	return token
}

// Check whether we have an optional noreply token at the current cursor position
func (r *commandTokenizer) IsNoReply() bool {
	return nil != r.GetOptionalTokenString("^noreply$")
}

// Get the next token, and convert to int
// TODO: is int32 good enough, or do we need another type? how big can these numbers get?
func (r *commandTokenizer) GetTokenInt() *int {
	tokenstr := r.GetTokenString()
	if nil == tokenstr { return nil }
	token, err := strconv.Atoi(*tokenstr)
	if nil != err { return nil }
	return &token
}

// Get tokenlen worth of bytes from the command at the cursor (useful when we know byte length of the value)
func (r *commandTokenizer) GetTokenBytes(tokenlen int) *[]byte {
	// If there aren't that many characters left, then decline
	if (len(r.command) - r.cursor) < tokenlen { return nil }
	token := r.command[r.cursor:r.cursor+tokenlen]
	r.cursor += tokenlen
	return &token
}

