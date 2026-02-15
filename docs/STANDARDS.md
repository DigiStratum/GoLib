# GoLib Coding Standards

This document defines the coding standards, conventions, and best practices for GoLib.

## Go Style Guidelines

GoLib follows the standard Go style guidelines with project-specific conventions.

### Code Formatting

- Use `gofmt` or `goimports` for all code
- Maximum line length: 100 characters (soft limit)
- Use tabs for indentation (Go standard)

### Naming Conventions

#### Packages

```go
// Good - lowercase, single word when possible
package cache
package logger
package http

// Acceptable - lowercase, no underscores
package logwriter
package nullables
```

#### Interfaces

```go
// Interface names end with "Ifc" suffix
type CacheIfc interface {}
type LoggerIfc interface {}
type ConfigIfc interface {}
type HttpClientIfc interface {}
```

#### Structs

```go
// PascalCase for exported types
type Cache struct {}
type Logger struct {}
type HttpClient struct {}
```

#### Methods

```go
// PascalCase for exported methods
func (r *Cache) Set(key string, value interface{}) bool
func (r *Cache) Get(key string) interface{}

// camelCase for private methods
func (r *Cache) pruneExpired() error
```

#### Receivers

```go
// Use 'r' for receiver everywhere for consistent copy/paste
func (r *Cache) Set(key string, value interface{}) bool
func (r Cache) Get(key string) interface{}

// Use pointer receiver (*r) only when necessary
func (r *Cache) Set(...) // mutates state - use pointer
func (r Cache) Get(...)  // read-only - use value
```

### File Organization

#### One Exported Type Per File

```
cache/
├── cache.go           # Cache struct and CacheIfc
├── cacheitem.go       # CacheItem struct
└── cache_test.go      # Tests
```

Exception: Small, tightly coupled types may share a file.

#### Import Ordering

```go
import (
    // Standard library
    "fmt"
    "time"
    "errors"
    
    // Third-party packages
    "github.com/aws/aws-sdk-go/aws"
    
    // Internal packages (with aliases as needed)
    cfg "github.com/DigiStratum/GoLib/Config"
    chrono "github.com/DigiStratum/GoLib/Chrono"
)
```

### Interface Design

#### Accept Interfaces, Return Structs

```go
// Good - accept interface
func NewService(logger logger.LoggerIfc) *Service

// Good - return concrete struct
func NewCache() *Cache

// Exception: Factory functions may return interfaces when needed
func NewObjectStore(storeType string) ObjectStoreIfc
```

#### Small, Focused Interfaces

```go
// Good - small interface
type LogWriterIfc interface {
    Log(message string)
}

// Avoid - large interfaces that are hard to implement/mock
type EverythingIfc interface {
    MethodA()
    MethodB()
    MethodC()
    // ... many more
}
```

### Error Handling

#### Return Errors

```go
// Good - return error for caller to handle
func (r *Cache) Drop(key string) (bool, error) {
    // ...
    if desync {
        return true, fmt.Errorf("cache desync detected for key '%s'", key)
    }
    return true, nil
}
```

#### Use fmt.Errorf

```go
// Good
return fmt.Errorf("failed to connect: %w", err)

// Avoid
return errors.New(fmt.Sprintf("failed to connect: %v", err))
```

#### Don't Log from Libraries

```go
// Good - return error, let caller decide whether to log
func (r *ObjectStore) Get(path string) (interface{}, error) {
    data, err := r.fetch(path)
    if err != nil {
        return nil, fmt.Errorf("ObjectStore.Get(%s): %w", path, err)
    }
    return data, nil
}

// Avoid - logging in library code
func (r *ObjectStore) Get(path string) (interface{}, error) {
    data, err := r.fetch(path)
    if err != nil {
        log.Printf("ERROR: failed to get %s: %v", path, err)  // Don't do this
        return nil, err
    }
    return data, nil
}
```

### Concurrency

#### Mutex for Mutable State

```go
type Cache struct {
    cache map[string]*cacheItem
    mutex sync.Mutex
}

func (r *Cache) Set(key string, value interface{}) bool {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    return r.set(key, value)
}
```

#### Channels for Orchestration

```go
// Prefer channels for goroutine coordination
done := make(chan bool)
go func() {
    // do work
    done <- true
}()
<-done
```

#### Pointer vs Value Receivers

```go
// Use pointer receiver for mutable (write) operations
func (r *Cache) Set(key string, value interface{}) bool

// Use value receiver for immutable (read) operations
func (r Cache) Get(key string) interface{}
func (r Cache) Has(key string) bool
```

### Nil Handling

#### Handle Nil Receivers

```go
func (r *Cache) Set(key string, value interface{}) bool {
    if r == nil {
        return false
    }
    // ...
}
```

Reference: https://tour.golang.org/methods/12

### Resource Management

#### Implement io.Closer

Components managing resources should implement `io.Closer`:

```go
type Cache struct {
    closed bool
    // ...
}

func (r *Cache) Close() error {
    r.closed = true
    r.flush()
    return nil
}
```

### Documentation

#### Package Comments

```go
// Package cache provides an in-memory LRU cache with configurable
// size limits and automatic expiration.
//
// The cache is safe for concurrent use.
package cache
```

#### Function Comments

```go
// Set stores a value in the cache with the given key.
// If the key already exists, the value is updated and the
// item is moved to the front of the LRU list.
// Returns true if the value was stored, false if it was too
// large to fit in the cache.
func (r *Cache) Set(key string, value interface{}) bool
```

#### TODO Comments

```go
// TODO: Add support for multiple LogWriter's
// TODO: Replace expiresList with Binary Tree for O(log n) operations
// FIXME: This logic assumes sorted iteration which is incorrect
```

### Testing

#### Test File Naming

```go
// Source file: cache.go
// Test file:   cache_test.go
```

#### Table-Driven Tests

```go
func TestCache_Set(t *testing.T) {
    tests := []struct {
        name     string
        key      string
        value    interface{}
        expected bool
    }{
        {"simple string", "key", "value", true},
        {"empty key", "", "value", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewCache()
            result := c.Set(tt.key, tt.value)
            if result != tt.expected {
                t.Errorf("Set() = %v, want %v", result, tt.expected)
            }
        })
    }
}
```

### Code Patterns

#### Builder Pattern

```go
// Use builder pattern for complex object construction
request := http.NewHttpRequestBuilder().
    SetMethod(http.GET).
    SetUrl("https://example.com").
    AddHeader("Authorization", "Bearer token").
    Build()
```

#### Factory Functions

```go
// Use New* functions for construction
func NewCache() *Cache {
    return &Cache{
        cache: make(map[string]*cacheItem),
    }
}

// Support singleton pattern where appropriate
var instance *Logger

func GetLogger() *Logger {
    return instance
}
```

#### Fluent Interface

```go
// Enable method chaining for configuration
func (r *Logger) SetMinLogLevel(level LogLevel) *Logger {
    r.minLogLevel = level
    return r
}

// Usage:
logger.SetMinLogLevel(DEBUG).SetLogWriter(writer).LogTimestamp(true)
```

## Refactoring Guidelines

From the project README:

1. Accept Interfaces, return structs (except by exception)
2. Use pointer receiver for mutable operations, value for immutable
3. Use mutex/semaphore/channels for mutable operations
4. Prefer go-routine+channel over mutex for concurrency orchestration
5. Use `r` for receiver everywhere for better copy/paste
6. Clean up TODOs/FIXMEs as reasonably able
7. Add test coverage as reasonably able
8. Add documentation (godoc, README) as reasonably able
9. Add working examples as reasonably able
10. Use `fmt.Errorf()` instead of `errors.New()`
11. Don't produce error log output from library functions
12. ONE exported struct+interface per source file (with exceptions)
13. Log Trace() messages to track entry into library functions
14. Implement `io.Closer` for classes managing external resources
15. Handle `r == nil` receivers

## Code Review Checklist

- [ ] Follows naming conventions
- [ ] Uses `r` as receiver name
- [ ] Returns errors (doesn't panic)
- [ ] Uses `fmt.Errorf()` for error creation
- [ ] Handles nil receivers where appropriate
- [ ] Has tests for new functionality
- [ ] Documentation for exported types/functions
- [ ] No logging in library code
- [ ] Uses appropriate mutex/channel for concurrency
- [ ] Implements `io.Closer` if managing resources
