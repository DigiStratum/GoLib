# GoLib Integrations

This document describes the integration points, external dependencies, and API contracts for GoLib.

## External Dependencies

### Go Module Dependencies

From `go.mod`:

```go
require (
    github.com/DATA-DOG/go-sqlmock v1.5.2        // SQL mocking for tests
    github.com/DigiStratum/Go-cbroglie-mustache v1.0.1  // Mustache templating
    github.com/aws/aws-sdk-go v1.55.7            // AWS SDK
    github.com/go-sql-driver/mysql v1.9.3        // MySQL driver
)
```

### AWS Services

GoLib integrates with the following AWS services:

| Service | Package | Purpose |
|---------|---------|---------|
| S3 | `Cloud/aws`, `Object/store/objectstores3` | Object storage |
| DynamoDB | `Object/store/objectstoredynamo` | NoSQL object storage |
| Secrets Manager | `Object/store/objectstoresecrets` | Secure credential storage |

#### AWS Credentials

AWS integrations use the standard AWS SDK credential chain:
1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM instance profile (EC2/ECS)

### MySQL Database

The `DB/MySQL` package provides MySQL connectivity:

```go
// Connection configuration via ConfigIfc
config := {
    "host":     "localhost",
    "port":     "3306",
    "database": "mydb",
    "user":     "username",
    "password": "password",
}
```

## API Contracts

### Logger Interface

```go
type LoggerIfc interface {
    GetNewPrefixedLogger(prefix string) *Logger
    SetMinLogLevel(minLogLevel LogLevel) *Logger
    SetLogWriter(logWriter lw.LogWriterIfc) *Logger
    LogTimestamp(logTimestamp bool) *Logger
    Any(level LogLevel, format string, a ...interface{}) error
    Crazy(format string, a ...interface{}) error
    Trace(format string, a ...interface{}) error
    Debug(format string, a ...interface{}) error
    Info(format string, a ...interface{}) error
    Warn(format string, a ...interface{}) error
    Error(format string, a ...interface{}) error
    Fatal(format string, a ...interface{}) error
}
```

### LogWriter Interface

```go
type LogWriterIfc interface {
    Log(message string)
}
```

Implementations:
- `StdOutLogWriter` - Writes to standard output

### Cache Interface

```go
type CacheIfc interface {
    Configure(config cfg.ConfigIfc) error
    SetTimeSource(timeSource chrono.TimeSourceIfc)
    IsEmpty() bool
    Size() int64
    Count() int
    Set(key string, value interface{}) bool
    SetExpires(key string, expires chrono.TimeStampIfc) bool
    GetExpires(key string) chrono.TimeStampIfc
    Get(key string) interface{}
    GetKeys() []string
    Has(key string) bool
    HasAll(keys *[]string) bool
    Drop(key string) (bool, error)
    DropAll(keys *[]string) (int, error)
    Flush()
    Close() error
}
```

#### Cache Configuration

| Key | Type | Description |
|-----|------|-------------|
| `newItemExpires` | int64 | Seconds until new items expire (0 = never) |
| `totalCountLimit` | int64 | Maximum number of cached items (0 = unlimited) |
| `totalSizeLimit` | int64 | Maximum total size in bytes (0 = unlimited) |

### Config Interface

```go
type ConfigIfc interface {
    hashmap.HashMapIfc
    MergeConfig(mergeCfg ConfigIfc) *Config
    GetSubsetConfig(prefix string) *Config
    GetSubsetKeys(keys *[]string) *Config
    GetInverseSubsetConfig(prefix string) *Config
    DereferenceString(str string) *string
    Dereference(referenceConfig ConfigIfc) int
    DereferenceAll(referenceConfigs ...ConfigIfc) int
    DereferenceLoop(maxLoops int, referenceConfig ConfigIfc) bool
}
```

### HashMap Interface

```go
type HashMapIfc interface {
    IsEmpty() bool
    Size() int
    Has(key string) bool
    HasAll(keys *[]string) bool
    Get(key string) *string
    GetInt64(key string) *int64
    GetFloat64(key string) *float64
    Set(key string, value string)
    GetKeys() *[]string
    Merge(hashMap HashMapIfc)
    GetIterator() func() interface{}
    // ... additional methods
}
```

### HTTP Client Interface

```go
// HttpClient
type HttpClientIfc interface {
    Send(request HttpRequestIfc) (HttpResponseIfc, error)
}

// HttpRequest
type HttpRequestIfc interface {
    GetMethod() string
    GetUrl() string
    GetHeaders() HttpHeadersIfc
    GetBody() HttpRequestBodyIfc
}

// HttpResponse  
type HttpResponseIfc interface {
    GetStatusCode() int
    GetHeaders() HttpHeadersIfc
    GetBody() []byte
}
```

### Object Store Interface

```go
type ObjectStoreIfc interface {
    Get(path string) (interface{}, error)
    Put(path string, object interface{}) error
    Delete(path string) error
    Exists(path string) bool
    List(prefix string) ([]string, error)
}
```

#### Object Path Format

```
SCOPE/CONTEXT/LANGUAGE/relative/path/filename
```

- **SCOPE**: `public` | `private` | `protected`
- **CONTEXT**: Application-defined path segments
- **LANGUAGE**: ISO locale code (e.g., `en-US`) or `default`

### TimeSource Interface

```go
type TimeSourceIfc interface {
    Now() TimeStampIfc
}
```

Used for testing time-dependent code by injecting mock time sources.

### TimeStamp Interface

```go
type TimeStampIfc interface {
    Add(seconds int64) TimeStampIfc
    IsExpired() bool
    Unix() int64
    // ... additional methods
}
```

### Configurable Interface

```go
type ConfigurableIfc interface {
    Configure(config cfg.ConfigIfc) error
}
```

Components implementing this interface accept runtime configuration.

### Runnable Interface

```go
type RunnableIfc interface {
    Run()
    IsRunning() bool
    Stop()
}
```

Used by `Process` and `Starter` packages for lifecycle management.

### io.Closer Implementation

Many components implement `io.Closer`:

```go
type Closer interface {
    Close() error
}
```

Components managing resources (connections, files, caches) should be closed when done.

## OAuth 2.0 Integration

The `Auth/oauth2` package supports OAuth 2.0 per RFC 6749:

### Token Response Format

```json
{
    "access_token": "YOUR_ACCESS_TOKEN_STRING",
    "token_type": "Bearer",
    "expires_in": 3600,
    "refresh_token": "YOUR_REFRESH_TOKEN_STRING",
    "scope": "openid email profile",
    "id_token": "YOUR_ID_TOKEN_STRING"
}
```

### Supported Grant Types

- Authorization Code
- Client Credentials
- Refresh Token

## Integration Patterns

### Dependency Injection Pattern

```go
// Create dependencies
logger := logger.NewLogger("my-app")
cache := cache.NewCache()
config := config.NewConfig()

// Inject into component
myService := NewMyService(logger, cache, config)
```

### Configuration Pattern

```go
// Load base config
config := config.NewConfig()
config.LoadFromJSON(jsonBytes)

// Merge environment-specific config
envConfig := config.NewConfig()
envConfig.LoadFromJSON(envJsonBytes)
config.MergeConfig(envConfig)

// Resolve references
config.DereferenceLoop(10, config)

// Apply to component
cache.Configure(config.GetSubsetConfig("cache."))
```

### Connection Pool Pattern

```go
// Create pool
pool := mysql.NewConnectionPool()
pool.Configure(config)

// Lease connection
conn, err := pool.Lease()
if err != nil { return err }
defer conn.Release()

// Execute queries
result, err := conn.Query("SELECT * FROM users")
```

### Object Store Abstraction

```go
// Select store implementation based on config
var store object.ObjectStoreIfc
switch storeType {
case "s3":
    store = objectstores3.NewObjectStoreS3()
case "mysql":
    store = objectstoremysql.NewObjectStoreMySQL()
case "dynamo":
    store = objectstoredynamo.NewObjectStoreDynamo()
}

// Use uniformly
data, err := store.Get("public/default/en-US/data/config.json")
```

## Error Handling Contract

All integrations follow these error conventions:

1. **Return errors** - Don't panic; return errors for caller to handle
2. **Wrap context** - Use `fmt.Errorf("context: %w", err)` for error chains
3. **Nil checks** - Methods handle nil receivers gracefully
4. **No library logging** - Callers decide whether to log errors
