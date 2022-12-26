# A fake Memcached Server implementation

-----

ref: https://github.com/memcached/memcached/blob/master/doc/protocol.txt

-----

ref: https://docs.oracle.com/cd/E17952_01/mysql-5.6-en/ha-memcached-interfaces-protocol.html

## Storage commands to the server take the form:

>> <command name> <key> <flags> <exptime> <bytes> [noreply]\r\n

such that <command name> = set|add|replace|append|prepend

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

Stats:
|-----------------------+---------+-------------------------------------------|
| Name                  | Type    | Meaning                                   |
|-----------------------+---------+-------------------------------------------|
| pid                   | 32u     | Process id of this server process         |
| uptime                | 32u     | Number of secs since the server started   |
| time                  | 32u     | current UNIX time according to the server |
| version               | string  | Version string of this server             |
| pointer_size          | 32      | Default size of pointers on the host OS   |
|                       |         | (generally 32 or 64)                      |
| rusage_user           | 32u.32u | Accumulated user time for this process    |
|                       |         | (seconds:microseconds)                    |
| rusage_system         | 32u.32u | Accumulated system time for this process  |
|                       |         | (seconds:microseconds)                    |
| curr_items            | 64u     | Current number of items stored            |
| total_items           | 64u     | Total number of items stored since        |
|                       |         | the server started                        |
| bytes                 | 64u     | Current number of bytes used              |
|                       |         | to store items                            |
| max_connections       | 32u     | Max number of simultaneous connections    |
| curr_connections      | 32u     | Number of open connections                |
| total_connections     | 32u     | Total number of connections opened since  |
|                       |         | the server started running                |
| rejected_connections  | 64u     | Conns rejected in maxconns_fast mode      |
| connection_structures | 32u     | Number of connection structures allocated |
|                       |         | by the server                             |
| response_obj_oom      | 64u     | Connections closed by lack of memory      |
| response_obj_count    | 64u     | Total response objects in use             |
| response_obj_bytes    | 64u     | Total bytes used for resp. objects. is a  |
|                       |         | subset of bytes from read_buf_bytes.      |
| read_buf_count        | 64u     | Total read/resp buffers allocated         |
| read_buf_bytes        | 64u     | Total read/resp buffer bytes allocated    |
| read_buf_bytes_free   | 64u     | Total read/resp buffer bytes cached       |
| read_buf_oom          | 64u     | Connections closed by lack of memory      |
| reserved_fds          | 32u     | Number of misc fds used internally        |
| proxy_conn_requests   | 64u     | Number of requests received by the proxy  |
| proxy_conn_errors     | 64u     | Number of internal errors from proxy      |
| proxy_conn_oom        | 64u     | Number of out of memory errors while      |
|                       |         | serving proxy requests                    |
| proxy_req_active      | 64u     | Number of in-flight proxy requests        |
| cmd_get               | 64u     | Cumulative number of retrieval reqs       |
| cmd_set               | 64u     | Cumulative number of storage reqs         |
| cmd_flush             | 64u     | Cumulative number of flush reqs           |
| cmd_touch             | 64u     | Cumulative number of touch reqs           |
| get_hits              | 64u     | Number of keys that have been requested   |
|                       |         | and found present                         |
| get_misses            | 64u     | Number of items that have been requested  |
|                       |         | and not found                             |
| get_expired           | 64u     | Number of items that have been requested  |
|                       |         | but had already expired.                  |
| get_flushed           | 64u     | Number of items that have been requested  |
|                       |         | but have been flushed via flush_all       |
| delete_misses         | 64u     | Number of deletions reqs for missing keys |
| delete_hits           | 64u     | Number of deletion reqs resulting in      |
|                       |         | an item being removed.                    |
| incr_misses           | 64u     | Number of incr reqs against missing keys. |
| incr_hits             | 64u     | Number of successful incr reqs.           |
| decr_misses           | 64u     | Number of decr reqs against missing keys. |
| decr_hits             | 64u     | Number of successful decr reqs.           |
| cas_misses            | 64u     | Number of CAS reqs against missing keys.  |
| cas_hits              | 64u     | Number of successful CAS reqs.            |
| cas_badval            | 64u     | Number of CAS reqs for which a key was    |
|                       |         | found, but the CAS value did not match.   |
| touch_hits            | 64u     | Number of keys that have been touched     |
|                       |         | with a new expiration time                |
| touch_misses          | 64u     | Number of items that have been touched    |
|                       |         | and not found                             |
| store_too_large       | 64u     | Number of rejected storage requests       |
|                       |         | caused by attempting to write a value     |
|                       |         | larger than the -I limit                  |
| store_no_memory       | 64u     | Number of rejected storage requests       |
|                       |         | caused by exhaustion of the -m memory     |
|                       |         | limit (relevant when -M is used)          |
| auth_cmds             | 64u     | Number of authentication commands         |
|                       |         | handled, success or failure.              |
| auth_errors           | 64u     | Number of failed authentications.         |
| idle_kicks            | 64u     | Number of connections closed due to       |
|                       |         | reaching their idle timeout.              |
| evictions             | 64u     | Number of valid items removed from cache  |
|                       |         | to free memory for new items              |
| reclaimed             | 64u     | Number of times an entry was stored using |
|                       |         | memory from an expired entry              |
| bytes_read            | 64u     | Total number of bytes read by this server |
|                       |         | from network                              |
| bytes_written         | 64u     | Total number of bytes sent by this server |
|                       |         | to network                                |
| limit_maxbytes        | size_t  | Number of bytes this server is allowed to |
|                       |         | use for storage.                          |
| accepting_conns       | bool    | Whether or not server is accepting conns  |
| listen_disabled_num   | 64u     | Number of times server has stopped        |
|                       |         | accepting new connections (maxconns).     |
| time_in_listen_disabled_us                                                  |
|                       | 64u     | Number of microseconds in maxconns.       |
| threads               | 32u     | Number of worker threads requested.       |
|                       |         | (see doc/threads.txt)                     |
| conn_yields           | 64u     | Number of times any connection yielded to |
|                       |         | another due to hitting the -R limit.      |
| hash_power_level      | 32u     | Current size multiplier for hash table    |
| hash_bytes            | 64u     | Bytes currently used by hash tables       |
| hash_is_expanding     | bool    | Indicates if the hash table is being      |
|                       |         | grown to a new size                       |
| expired_unfetched     | 64u     | Items pulled from LRU that were never     |
|                       |         | touched by get/incr/append/etc before     |
|                       |         | expiring                                  |
| evicted_unfetched     | 64u     | Items evicted from LRU that were never    |
|                       |         | touched by get/incr/append/etc.           |
| evicted_active        | 64u     | Items evicted from LRU that had been hit  |
|                       |         | recently but did not jump to top of LRU   |
| slab_reassign_running | bool    | If a slab page is being moved             |
| slabs_moved           | 64u     | Total slab pages moved                    |
| crawler_reclaimed     | 64u     | Total items freed by LRU Crawler          |
| crawler_items_checked | 64u     | Total items examined by LRU Crawler       |
| lrutail_reflocked     | 64u     | Times LRU tail was found with active ref. |
|                       |         | Items can be evicted to avoid OOM errors. |
| moves_to_cold         | 64u     | Items moved from HOT/WARM to COLD LRU's   |
| moves_to_warm         | 64u     | Items moved from COLD to WARM LRU         |
| moves_within_lru      | 64u     | Items reshuffled within HOT or WARM LRU's |
| direct_reclaims       | 64u     | Times worker threads had to directly      |
|                       |         | reclaim or evict items.                   |
| lru_crawler_starts    | 64u     | Times an LRU crawler was started          |
| lru_maintainer_juggles                                                      |
|                       | 64u     | Number of times the LRU bg thread woke up |
| slab_global_page_pool | 32u     | Slab pages returned to global pool for    |
|                       |         | reassignment to other slab classes.       |
| slab_reassign_rescues | 64u     | Items rescued from eviction in page move  |
| slab_reassign_evictions_nomem                                               |
|                       | 64u     | Valid items evicted during a page move    |
|                       |         | (due to no free memory in slab)           |
| slab_reassign_chunk_rescues                                                 |
|                       | 64u     | Individual sections of an item rescued    |
|                       |         | during a page move.                       |
| slab_reassign_inline_reclaim                                                |
|                       | 64u     | Internal stat counter for when the page   |
|                       |         | mover clears memory from the chunk        |
|                       |         | freelist when it wasn't expecting to.     |
| slab_reassign_busy_items                                                    |
|                       | 64u     | Items busy during page move, requiring a  |
|                       |         | retry before page can be moved.           |
| slab_reassign_busy_deletes                                                  |
|                       | 64u     | Items busy during page move, requiring    |
|                       |         | deletion before page can be moved.        |
| log_worker_dropped    | 64u     | Logs a worker never wrote due to full buf |
| log_worker_written    | 64u     | Logs written by a worker, to be picked up |
| log_watcher_skipped   | 64u     | Logs not sent to slow watchers.           |
| log_watcher_sent      | 64u     | Logs written to watchers.                 |
| log_watchers          | 64u     | Number of currently active watchers.      |
| unexpected_napi_ids   | 64u     | Number of times an unexpected napi id is  |
|                       |         | is received. See doc/napi_ids.txt         |
| round_robin_fallback  | 64u     | Number of times napi id of 0 is received  |
|                       |         | resulting in fallback to round robin      |
|                       |         | thread selection. See doc/napi_ids.txt    |
|-----------------------+---------+-------------------------------------------|

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

## Other Commands
ref: https://github.com/memcached/memcached/blob/master/doc/protocol.txt

"cache_memlimit" is a command with a numeric argument. This allows runtime
adjustments of the cache memory limit. It returns "OK\r\n" or an error (unless
"noreply" is given as the last parameter).

"shutdown" is a command with an optional argument used to stop memcached with
a kill signal. By default, "shutdown" alone raises SIGINT, though "graceful"
may be specified as the single argument to instead trigger a graceful shutdown
with SIGUSR1. The shutdown command is disabled by default, and can be enabled
with the -A/--enable-shutdown flag.

"version" is a command with no arguments:

version\r\n

In response, the server sends

"VERSION <version>\r\n", where <version> is the version string for the
server.

"verbosity" is a command with a numeric argument. It always succeeds,
and the server sends "OK\r\n" in response (unless "noreply" is given
as the last parameter). Its effect is to set the verbosity level of
the logging output.

"quit" is a command with no arguments:

quit\r\n

Upon receiving this command, the server closes the
connection. However, the client may also simply close the connection
when it no longer needs it, without issuing this command.

-----

TODO:
 * Simulate memcached version <= 1.2.0 for 16 bit flags on cache items vs 32 bits vs 1.2.1+
