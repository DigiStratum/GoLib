# GoLib Architecture

This document describes the high-level architecture and design principles of GoLib.

## Design Philosophy

GoLib follows several core design principles:

1. **Accept Interfaces, Return Structs** - Functions accept interfaces for flexibility but return concrete structs (with intentional exceptions)
2. **Composition over Inheritance** - Go's embedding is used to compose behaviors
3. **Dependency Injection** - Components receive dependencies rather than creating them
4. **Interface Segregation** - Small, focused interfaces over large monolithic ones
5. **Cloud Native** - Abstractions that work across local and cloud environments

## Module Structure

```
GoLib/
├── Auth/                    # Authentication
│   └── oauth2/              # OAuth 2.0 implementation (RFC 6749)
├── Cache/                   # In-memory caching
├── Chrono/                  # Time utilities
├── CLI/                     # Command-line helpers
├── Cloud/                   # Cloud provider integrations
│   └── aws/                 # AWS SDK wrappers
├── Config/                  # Configuration management
├── Crypto/                  # Cryptography
│   └── asymmetric/          # Asymmetric encryption
├── Data/                    # Data structures
│   ├── hashmap/             # Key-value maps
│   ├── collection/          # Collections
│   └── sizeable/            # Size calculation
├── DB/                      # Database connectivity
│   └── MySQL/               # MySQL driver with pooling
│       └── nullables/       # Nullable type helpers
├── Dependencies/            # Dependency injection
├── Events/                  # Event system
├── FileIO/                  # File operations
├── Logger/                  # Logging
│   └── logwriter/           # Output writers
├── Net/                     # Network utilities
│   └── http/                # HTTP client/server
├── Object/                  # Object abstraction
│   ├── collection/          # Object collections
│   ├── field/               # Field definitions
│   ├── prototype/           # Object prototypes
│   └── store/               # Storage backends
│       ├── objectstoredynamo/
│       ├── objectstoremanager/
│       ├── objectstoremysql/
│       ├── objectstores3/
│       └── objectstoresecrets/
├── Process/                 # Process management
├── Starter/                 # Application lifecycle
├── Testing/                 # Test utilities
│   ├── mocks/               # Mock implementations
│   └── TestRunner/          # Test harness
├── Version/                 # Versioning
├── WebUI/                   # UI templating
└── tools/                   # Build/dev tools
```

## Core Components

### Logger

The Logger component provides structured logging with configurable levels:

```
CRAZY → TRACE → DEBUG → INFO → WARN → ERROR → FATAL
```

- **Singleton Pattern**: Default logger accessible via `GetLogger()`
- **Prefixed Loggers**: Create contextual loggers with `GetNewPrefixedLogger()`
- **Pluggable Writers**: Implement `LogWriterIfc` for custom output destinations
- **Thread-Safe**: Uses mutex for concurrent access

### Cache

In-memory LRU cache with automatic expiration:

- **Size/Count Limits**: Configure maximum entries or total size
- **LRU Eviction**: Least recently used items evicted first
- **Background Pruning**: Async goroutine removes expired entries
- **Thread-Safe**: Mutex-protected operations
- **Configurable**: Via `ConfigIfc` for runtime configuration

### Config

Configuration management with template dereferencing:

- **HashMap Extension**: Builds on Data/hashmap
- **Reference Resolution**: `%key%` syntax for value substitution
- **Subset Extraction**: Filter configs by key prefix
- **Merge Support**: Combine multiple configurations

### Object Store

Abstraction layer for object storage:

```
ObjectStoreIfc
    ├── ObjectStoreS3         (AWS S3)
    ├── ObjectStoreDynamo     (AWS DynamoDB)
    ├── ObjectStoreMySQL      (MySQL database)
    └── ObjectStoreSecrets    (AWS Secrets Manager)
```

Objects are identified by paths: `SCOPE/CONTEXT/LANGUAGE/path/filename`
- **SCOPE**: public, private, protected
- **CONTEXT**: Application-specific path components
- **LANGUAGE**: Locale code (e.g., en-US, default)

### HTTP Client

Builder-pattern HTTP client:

```
HttpRequestBuilder → HttpRequest → HttpClient → HttpResponse
```

- **Request Builder**: Fluent API for constructing requests
- **Headers Builder**: Separate builder for headers
- **Body Builder**: Request body construction
- **Standards-Based**: Follows RFC 7231

### Database (MySQL)

Connection pooling and query abstraction:

```
ConnectionPool
    └── PooledConnection
            └── LeasedConnection
                    └── Query → ResultSet → ResultRow
```

- **Connection Pooling**: Reuse connections efficiently
- **Lease Pattern**: Check out/return connections
- **Query Builder**: Parameterized query construction
- **Nullable Types**: Handle SQL NULL values

### Dependencies

Lightweight dependency injection:

- **Convention-Based**: Not a full DI container
- **Explicit Wiring**: Dependencies passed at construction
- **Interface-Driven**: Depend on interfaces, not implementations

## Data Flow Patterns

### Request Handling

```
HTTP Request
    → Router
        → Module (with Scheme)
            → ObjectStore.Get(path)
            → Template Rendering
                → Magic Tag Substitution
                    → HTTP Response
```

### Configuration Loading

```
JSON File
    → Config.LoadFromFile()
        → Merge with overrides
            → Dereference %references%
                → Apply to components
```

### Logging Flow

```
Application Code
    → Logger.Info(format, args)
        → Check MinLogLevel
            → Format with timestamp
                → LogWriter.Log(message)
                    → StdOut / File / etc.
```

## Concurrency Model

- **Mutex Locks**: Used for mutable shared state (Cache, Logger)
- **Channels**: Preferred for goroutine orchestration
- **Goroutines**: Background tasks (cache pruning, async operations)
- **Pointer Receivers**: For mutable operations
- **Value Receivers**: For read-only operations

## Error Handling

- **Return Errors**: Functions return `error` as last value
- **fmt.Errorf**: Preferred over `errors.New(fmt.Sprintf(...))`
- **No Library Logging**: Libraries don't log errors; callers decide
- **Nil Receiver Handling**: Methods handle `r == nil` gracefully

## Extension Points

1. **LogWriter**: Implement `LogWriterIfc` for custom log destinations
2. **ObjectStore**: Implement `ObjectStoreIfc` for new storage backends
3. **TimeSource**: Implement `TimeSourceIfc` for testing time-dependent code
4. **Config**: Extend `ConfigIfc` for custom configuration sources
