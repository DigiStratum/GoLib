# GoLib Documentation

GoLib is a comprehensive Go library providing reusable, modular components for building robust applications. It offers a wide range of functionality from caching and configuration management to database connectivity, HTTP clients, and cloud integrations.

## Overview

GoLib serves as the foundational library for DigiStratum projects, providing:

- **Cache** - In-memory LRU cache with configurable expiration and size limits
- **Chrono** - Time manipulation and timestamp utilities
- **CLI** - Command-line interface helpers
- **Cloud** - AWS integrations (S3, Secrets Manager, DynamoDB)
- **Config** - Configuration management with template dereferencing
- **Crypto** - Asymmetric cryptography utilities
- **Data** - Data structures (hashmaps, collections, sizeables)
- **DB** - MySQL database connectivity with connection pooling
- **Dependencies** - Lightweight dependency injection framework
- **Events** - Event handling and dispatch
- **FileIO** - File system operations and resource management
- **Logger** - Structured logging with configurable levels and writers
- **Net/HTTP** - HTTP client with builder patterns for requests/responses
- **Object** - Object storage abstraction layer (S3, MySQL, DynamoDB, Secrets)
- **Process** - Process management and runnable interfaces
- **Starter** - Application lifecycle management
- **Testing** - Mock utilities and test runners
- **Version** - Version management utilities
- **WebUI** - UI modeling and templating system

## Installation

```go
go get github.com/DigiStratum/GoLib
```

## Requirements

- Go 1.24.2 or later

## Quick Start

### Logging

```go
import "github.com/DigiStratum/GoLib/Logger"

// Use the singleton logger
log := logger.GetLogger()
log.SetMinLogLevel(logger.DEBUG)
log.Info("Application started")
log.Debug("Debug message: %s", someVar)
```

### Caching

```go
import "github.com/DigiStratum/GoLib/Cache"

c := cache.NewCache()
c.Configure(config) // Optional: set limits and expiration
c.Set("key", value)
result := c.Get("key")
```

### HTTP Client

```go
import "github.com/DigiStratum/GoLib/Net/http"

client := http.NewHttpClient()
request := http.NewHttpRequestBuilder().
    SetMethod(http.GET).
    SetUrl("https://api.example.com/data").
    Build()
response, err := client.Send(request)
```

### Database (MySQL)

```go
import "github.com/DigiStratum/GoLib/DB/MySQL"

pool := mysql.NewConnectionPool()
conn, err := pool.Lease()
defer conn.Release()
result, err := conn.Query("SELECT * FROM users WHERE id = ?", userId)
```

## Documentation Index

- [ARCHITECTURE.md](ARCHITECTURE.md) - High-level design and module structure
- [INTEGRATIONS.md](INTEGRATIONS.md) - Integration points and API contracts
- [TEST_PLAN.md](TEST_PLAN.md) - Test coverage and testing approach
- [STANDARDS.md](STANDARDS.md) - Coding standards and conventions

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test -count=1 ./Cache/...

# Generate coverage report
go test -coverprofile=coverage.txt ./...
go tool cover -html=coverage.txt
```

## License

MIT License - see [LICENSE](../LICENSE) for details.
